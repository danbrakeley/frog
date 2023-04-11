package frog

import (
	"fmt"
)

// CustomizerLogger is a Logger that adds specific fields and/or sets specific printer options for each line it logs
type CustomizerLogger struct {
	parent Logger
	opts   []PrinterOption
	fields []Fielder
}

func newCustomizerLogger(l Logger, opts []PrinterOption, f []Fielder) *CustomizerLogger {
	return &CustomizerLogger{
		parent: l,
		opts:   opts,
		fields: f,
	}
}

func (l *CustomizerLogger) Parent() Logger {
	return l.parent
}

func (l *CustomizerLogger) Close() {
	panic(fmt.Errorf("called Close on a CustomizerLogger logger"))
}

func (l *CustomizerLogger) SetMinLevel(level Level) Logger {
	l.parent.SetMinLevel(level)
	return l
}

func (l *CustomizerLogger) LogImpl(anchoredLine int32, opts []PrinterOption, level Level, msg string, fields []Fielder) {
	// static fields should come first, then any line-specific fields
	l.parent.LogImpl(anchoredLine, append(l.opts, opts...), level, msg, append(l.fields, fields...))
}

func (l *CustomizerLogger) Transient(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Transient, msg, fields)
	return l
}

func (l *CustomizerLogger) Verbose(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Verbose, msg, fields)
	return l
}

func (l *CustomizerLogger) Info(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Info, msg, fields)
	return l
}

func (l *CustomizerLogger) Warning(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Warning, msg, fields)
	return l
}

func (l *CustomizerLogger) Error(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Error, msg, fields)
	return l
}
