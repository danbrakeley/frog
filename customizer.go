package frog

// CustomizerLogger is a Logger that adds specific fields and/or sets specific printer options for each line it logs
type CustomizerLogger struct {
	parent   Logger
	opts     []PrinterOption
	fields   []Fielder
	minLevel Level // defaults to Transient
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

func (l *CustomizerLogger) MinLevel() Level {
	return l.minLevel
}

func (l *CustomizerLogger) SetMinLevel(level Level) Logger {
	l.minLevel = level
	return l
}

func (l *CustomizerLogger) LogImpl(level Level, msg string, fields []Fielder, opts []PrinterOption, d ImplData) {
	d.MergeMinLevel(l.minLevel) // ensure our minLevel is taken into account

	// static fields should come first, then any line-specific fields
	l.parent.LogImpl(level, msg, append(l.fields, fields...), append(l.opts, opts...), d)
}

func (l *CustomizerLogger) Transient(msg string, fields ...Fielder) Logger {
	l.LogImpl(Transient, msg, fields, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Verbose(msg string, fields ...Fielder) Logger {
	l.LogImpl(Verbose, msg, fields, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Info(msg string, fields ...Fielder) Logger {
	l.LogImpl(Info, msg, fields, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Warning(msg string, fields ...Fielder) Logger {
	l.LogImpl(Warning, msg, fields, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Error(msg string, fields ...Fielder) Logger {
	l.LogImpl(Error, msg, fields, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Log(level Level, msg string, fields ...Fielder) Logger {
	l.LogImpl(level, msg, fields, nil, ImplData{})
	return l
}
