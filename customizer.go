package frog

// CustomizerLogger is a Logger that adds specific fields and/or sets specific printer options for each line it logs
type CustomizerLogger struct {
	parent   Logger
	opts     []PrinterOption
	fields   []Field
	minLevel Level // defaults to Transient
}

func newCustomizerLogger(l Logger, opts []PrinterOption, fielders []Fielder) *CustomizerLogger {
	return &CustomizerLogger{
		parent: l,
		opts:   opts,
		fields: Fieldify(fielders),
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

func (l *CustomizerLogger) LogImpl(level Level, msg string, fielders []Fielder, opts []PrinterOption, d ImplData) {
	d.MergeMinLevel(l.minLevel)
	d.MergeFields(l.fields)
	l.parent.LogImpl(level, msg, fielders, append(l.opts, opts...), d)
}

func (l *CustomizerLogger) Transient(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Transient, msg, fielders, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Verbose(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Verbose, msg, fielders, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Info(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Info, msg, fielders, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Warning(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Warning, msg, fielders, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Error(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Error, msg, fielders, nil, ImplData{})
	return l
}

func (l *CustomizerLogger) Log(level Level, msg string, fielders ...Fielder) Logger {
	l.LogImpl(level, msg, fielders, nil, ImplData{})
	return l
}
