package frog

import (
	"fmt"
	"sync"
)

// FixedLine is a Logger that attempts to overwrite the same line in the terminal,
// allowing progress bars and other simple UI for the human to consume.
//
// If the Printer's CanUseAnsi is false, then it simply redirects to the normal
// behavior of the parent Buffered Logger.
//
// If the Printer's CanUseAnsi is true, then the Transient level is always
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

func (l *FixedLine) RemoveFixedLine() {
	l.mutex.Lock()
	if l.fnOnClose != nil {
		l.fnOnClose()
		l.fnOnClose = nil
	}
	l.mutex.Unlock()
}

func (l *FixedLine) Close() {
	panic(fmt.Errorf("called Close on a FixedLine logger"))
}

func (l *FixedLine) Parent() Logger {
	return l.parent
}

func (l *FixedLine) SetMinLevel(level Level) Logger {
	l.parent.SetMinLevel(level)
	return l
}

func (l *FixedLine) Log(level Level, msg string, fields ...Fielder) Logger {
	l.mutex.RLock()
	isClosed := l.fnOnClose == nil
	l.mutex.RUnlock()

	// most of the time, we want to our default parent behavior
	if isClosed || !l.parent.cfg.UseAnsi || level != Transient {
		l.parent.Log(level, msg, fields...)
		return l
	}

	// if we really do have a fixed line we want to print, then go straight to the source
	l.parent.logImpl(l.prn, l.line, level, msg, fields...)
	return l
}

func (l *FixedLine) Transient(msg string, fields ...Fielder) Logger {
	l.Log(Transient, msg, fields...)
	return l
}

func (l *FixedLine) Verbose(msg string, fields ...Fielder) Logger {
	l.Log(Verbose, msg, fields...)
	return l
}

func (l *FixedLine) Info(msg string, fields ...Fielder) Logger {
	l.Log(Info, msg, fields...)
	return l
}

func (l *FixedLine) Warning(msg string, fields ...Fielder) Logger {
	l.Log(Warning, msg, fields...)
	return l
}

func (l *FixedLine) Error(msg string, fields ...Fielder) Logger {
	l.Log(Error, msg, fields...)
	return l
}

func (l *FixedLine) Fatal(msg string, fields ...Fielder) {
	l.Log(Fatal, msg, fields...)
}
