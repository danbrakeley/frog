package frog

import (
	"fmt"
	"sync"
)

// FixedLine is a Logger that treats Transient level log data in a special way:
// - Transient level is never ignored, and always overwrites the same output line.
// - Non-Transient level is sent to the parent logger.
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

	// the parent can handle non-Transient lines
	if isClosed || level != Transient {
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
