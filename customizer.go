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

func (l *CustomizerLogger) Log(level Level, opts []PrinterOption, msg string, fields []Fielder) Logger {
	// static fields should come first, then any line-specific fields
	l.parent.Log(level, append(l.opts, opts...), msg, append(l.fields, fields...))
	return l
}

func (l *CustomizerLogger) Transient(msg string, fields ...Fielder) Logger {
	l.Log(Transient, nil, msg, fields)
	return l
}

func (l *CustomizerLogger) Verbose(msg string, fields ...Fielder) Logger {
	l.Log(Verbose, nil, msg, fields)
	return l
}

func (l *CustomizerLogger) Info(msg string, fields ...Fielder) Logger {
	l.Log(Info, nil, msg, fields)
	return l
}

func (l *CustomizerLogger) Warning(msg string, fields ...Fielder) Logger {
	l.Log(Warning, nil, msg, fields)
	return l
}

func (l *CustomizerLogger) Error(msg string, fields ...Fielder) Logger {
	l.Log(Error, nil, msg, fields)
	return l
}
