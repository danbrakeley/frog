package frog

import (
	"fmt"
	"sync"
)

// AnchoredLogger is a Logger that treats Transient level log data in a special way:
// - Transient level is never ignored, and always overwrites the same output line.
// - Non-Transient level is sent to the parent logger.
type AnchoredLogger struct {
	parent Logger
	line   int32

	mutex     sync.RWMutex
	fnOnClose func()
}

func newAnchor(parent Logger, line int32, fnOnClose func()) *AnchoredLogger {
	return &AnchoredLogger{
		parent:    parent,
		line:      line,
		fnOnClose: fnOnClose,
	}
}

func (l *AnchoredLogger) RemoveAnchor() {
	l.mutex.Lock()
	if l.fnOnClose != nil {
		l.fnOnClose()
		l.fnOnClose = nil
		l.line = 0
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

// LogImpl has to handle log requests from this package and from any children
// The value of the anchoredLine argument is ignored in favor of the line we store internally.
func (l *AnchoredLogger) LogImpl(anchoredLine int32, opts []PrinterOption, level Level, msg string, fields []Fielder) {
	// if we are requesting to log on an anchored line, then
	var line int32

	// only transient lines are anchorable
	if level == Transient {
		if l.line != 0 {
			l.mutex.RLock()
			line = l.line
			l.mutex.RUnlock()
		}
	}

	l.parent.LogImpl(line, opts, level, msg, fields)
}

func (l *AnchoredLogger) Transient(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Transient, msg, fields)
	return l
}

func (l *AnchoredLogger) Verbose(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Verbose, msg, fields)
	return l
}

func (l *AnchoredLogger) Info(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Info, msg, fields)
	return l
}

func (l *AnchoredLogger) Warning(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Warning, msg, fields)
	return l
}

func (l *AnchoredLogger) Error(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Error, msg, fields)
	return l
}
