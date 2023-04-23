package frog

import (
	"sync"
)

// AnchoredLogger is a Logger that treats Transient level log data in a special way:
// - Transient level is never ignored, and always overwrites the same output line.
// - Non-Transient level is sent to the parent logger.
type AnchoredLogger struct {
	parent   Logger
	line     int32
	minLevel Level // defaults to Transient

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

func (l *AnchoredLogger) Parent() Logger {
	return l.parent
}

func (l *AnchoredLogger) MinLevel() Level {
	return l.minLevel
}

func (l *AnchoredLogger) SetMinLevel(level Level) Logger {
	l.minLevel = level
	return l
}

func (l *AnchoredLogger) LogImpl(level Level, msg string, fielders []Fielder, opts []PrinterOption, d ImplData) {
	d.MergeMinLevel(l.minLevel)

	var line int32 // 0 if not targetting an anchored line

	// only transient lines are anchorable
	if level == Transient {
		l.mutex.RLock()
		line = l.line
		l.mutex.RUnlock()
	}

	// set our target anchor line, then pass up to the parent to render
	d.AnchoredLine = line
	l.parent.LogImpl(level, msg, fielders, opts, d)
}

func (l *AnchoredLogger) Transient(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Transient, msg, fielders, nil, ImplData{})
	return l
}

func (l *AnchoredLogger) Verbose(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Verbose, msg, fielders, nil, ImplData{})
	return l
}

func (l *AnchoredLogger) Info(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Info, msg, fielders, nil, ImplData{})
	return l
}

func (l *AnchoredLogger) Warning(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Warning, msg, fielders, nil, ImplData{})
	return l
}

func (l *AnchoredLogger) Error(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Error, msg, fielders, nil, ImplData{})
	return l
}

func (l *AnchoredLogger) Log(level Level, msg string, fielders ...Fielder) Logger {
	l.LogImpl(level, msg, fielders, nil, ImplData{})
	return l
}
