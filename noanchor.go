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

func (l *NoAnchorLogger) LogImpl(level Level, msg string, fielders []Fielder, opts []PrinterOption, d ImplData) {
	d.MergeMinLevel(l.minLevel)
	l.parent.LogImpl(level, msg, fielders, opts, d)
}

func (l *NoAnchorLogger) Transient(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Transient, msg, fielders, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Verbose(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Verbose, msg, fielders, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Info(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Info, msg, fielders, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Warning(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Warning, msg, fielders, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Error(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Error, msg, fielders, nil, ImplData{})
	return l
}

func (l *NoAnchorLogger) Log(level Level, msg string, fielders ...Fielder) Logger {
	l.LogImpl(level, msg, fielders, nil, ImplData{})
	return l
}
