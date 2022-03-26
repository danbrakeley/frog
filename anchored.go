package frog

import (
	"fmt"
	"sync"
)

// AnchoredLogger is a Logger that treats Transient level log data in a special way:
// - Transient level is never ignored, and always overwrites the same output line.
// - Non-Transient level is sent to the parent logger.
type AnchoredLogger struct {
	parent *Buffered
	prn    Printer
	line   int32

	mutex     sync.RWMutex
	fnOnClose func()
}

func newAnchor(b *Buffered, line int32, fnOnClose func()) *AnchoredLogger {
	return &AnchoredLogger{
		parent:    b,
		prn:       b.prn,
		line:      line,
		fnOnClose: fnOnClose,
	}
}

func (l *AnchoredLogger) RemoveAnchor() {
	l.mutex.Lock()
	if l.fnOnClose != nil {
		l.fnOnClose()
		l.fnOnClose = nil
	}
	l.mutex.Unlock()
}

func (l *AnchoredLogger) Close() {
	panic(fmt.Errorf("called Close on a AnchoredLogger logger"))
}

func (l *AnchoredLogger) Parent() Logger {
	return l.parent
}

func (l *AnchoredLogger) SetMinLevel(level Level) Logger {
	l.parent.SetMinLevel(level)
	return l
}

func (l *AnchoredLogger) Log(level Level, opts []PrinterOption, msg string, fields []Fielder) Logger {
	l.mutex.RLock()
	isClosed := l.fnOnClose == nil
	l.mutex.RUnlock()

	// the parent can handle non-Transient lines
	if isClosed || level != Transient {
		l.parent.Log(level, opts, msg, fields)
		return l
	}

	// if we really do have an anchored line we want to print, then go straight to the source
	l.parent.logImpl(l.prn, opts, l.line, level, msg, fields)
	return l
}

func (l *AnchoredLogger) Transient(msg string, fields ...Fielder) Logger {
	l.Log(Transient, nil, msg, fields)
	return l
}

func (l *AnchoredLogger) Verbose(msg string, fields ...Fielder) Logger {
	l.Log(Verbose, nil, msg, fields)
	return l
}

func (l *AnchoredLogger) Info(msg string, fields ...Fielder) Logger {
	l.Log(Info, nil, msg, fields)
	return l
}

func (l *AnchoredLogger) Warning(msg string, fields ...Fielder) Logger {
	l.Log(Warning, nil, msg, fields)
	return l
}

func (l *AnchoredLogger) Error(msg string, fields ...Fielder) Logger {
	l.Log(Error, nil, msg, fields)
	return l
}
