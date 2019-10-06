package frog

import "sync"

// FixedLine is a Logger that attempts to overwrite the same line in the terminal,
// allowing progress bars and other simple UI for the human to consume.
//
// If the Printer's CanUseAnsi is false, then it simply redirects to the normal
// behavior of the parent Buffered Logger.
//
// If the Printer's CanUseAnsi is true, then the Progress level is always
// printed, regardless of the MinLevel. This allows progress bars that do
// not pollute logs with garbage when not connected to a terminal.
type FixedLine struct {
	parent *Buffered
	prn    Printer
	line   int32

	mutex     sync.RWMutex
	fnOnClose func()
}

func newFixedLine(b *Buffered, line int32, fnOnClose func()) *FixedLine {
	return &FixedLine{
		parent:    b,
		prn:       b.prn,
		line:      line,
		fnOnClose: fnOnClose,
	}
}

func (l *FixedLine) Close() {
	l.mutex.Lock()
	if l.fnOnClose != nil {
		l.fnOnClose()
		l.fnOnClose = nil
	}
	l.mutex.Unlock()
}

func (l *FixedLine) RootLogger() Logger {
	return l.parent.RootLogger()
}

func (l *FixedLine) AddFixedLine() Logger {
	return l.parent.AddFixedLine()
}

func (l *FixedLine) MinLevel() Level {
	return l.parent.MinLevel()
}

func (l *FixedLine) SetMinLevel(level Level) Logger {
	l.parent.SetMinLevel(level)
	return l
}

func (l *FixedLine) Printf(level Level, format string, a ...interface{}) Logger {
	l.mutex.RLock()
	isClosed := l.fnOnClose == nil
	l.mutex.RUnlock()

	// if this line is closed, or if the parent can't use ansi, then divert request to parent
	if isClosed || !l.parent.cfg.UseAnsi {
		l.parent.Printf(level, format, a...)
		return l
	}

	// fixed lines allow Progress to print, regardless of MinLevel
	if level < l.parent.MinLevel() && level != Progress {
		return l
	}

	l.parent.printfImpl(l.prn, l.line, level, format, a...)
	return l
}

func (l *FixedLine) Progressf(format string, a ...interface{}) Logger {
	l.Printf(Progress, format, a...)
	return l
}

func (l *FixedLine) Verbosef(format string, a ...interface{}) Logger {
	l.Printf(Verbose, format, a...)
	return l
}

func (l *FixedLine) Infof(format string, a ...interface{}) Logger {
	l.Printf(Info, format, a...)
	return l
}

func (l *FixedLine) Warningf(format string, a ...interface{}) Logger {
	l.Printf(Warning, format, a...)
	return l
}

func (l *FixedLine) Errorf(format string, a ...interface{}) Logger {
	l.Printf(Error, format, a...)
	return l
}

func (l *FixedLine) Fatalf(format string, a ...interface{}) {
	l.Printf(Fatal, format, a...)
}
