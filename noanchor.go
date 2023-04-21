package frog

// NoAnchorLogger is a Logger that is returned from AddAnchor when the RootLogger doesn't
// support anchors.
// NoAnchorLogger allows changing the min level, just like an AnchorLogger, but otherwise
// does nothing but pass through to the parent.
type NoAnchorLogger struct {
	parent   Logger
	minLevel Level // defaults to Transient
}

func newNoAnchor(parent Logger) *NoAnchorLogger {
	return &NoAnchorLogger{parent: parent}
}

func (l *NoAnchorLogger) RemoveAnchor() {
	// This space intentially left blank.
}

func (l *NoAnchorLogger) Parent() Logger {
	return l.parent
}

func (l *NoAnchorLogger) MinLevel() Level {
	return l.minLevel
}

func (l *NoAnchorLogger) SetMinLevel(level Level) Logger {
	l.minLevel = level
	return l
}

func (l *NoAnchorLogger) LogImpl(level Level, msg string, fields []Fielder, opts []PrinterOption, d ImplData) {
	d.MergeMinLevel(l.minLevel) // ensure our minLevel is taken into account
	l.parent.LogImpl(level, msg, fields, opts, d)
}

func (l *NoAnchorLogger) Transient(msg string, fields ...Fielder) Logger {
	l.LogImpl(Transient, msg, fields, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Verbose(msg string, fields ...Fielder) Logger {
	l.LogImpl(Verbose, msg, fields, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Info(msg string, fields ...Fielder) Logger {
	l.LogImpl(Info, msg, fields, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Warning(msg string, fields ...Fielder) Logger {
	l.LogImpl(Warning, msg, fields, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Error(msg string, fields ...Fielder) Logger {
	l.LogImpl(Error, msg, fields, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Log(level Level, msg string, fields ...Fielder) Logger {
	l.LogImpl(level, msg, fields, nil, ImplData{})
	return l
}
