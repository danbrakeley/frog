package frog

// TeeLogger directs all traffic to both a primary and secondary logger
// Note that of if one of your loggers supports anchor, make sure that is the Primary anchor.
type TeeLogger struct {
	Primary   Logger // Anchors are only supported though this logger
	Secondary Logger

	minLevel Level // defaults to Transient
}

func NewRootTee(a RootLogger, b RootLogger) (*TeeLogger, func()) {
	close := func() {
		a.Close()
		b.Close()
	}

	l := &TeeLogger{
		Primary:   a,
		Secondary: b,
	}

	return l, close
}

func (l *TeeLogger) Parent() Logger {
	// Anchoring relies on there being a single root parent.
	return l.Primary
}

func (n *TeeLogger) MinLevel() Level {
	return n.minLevel
}

func (n *TeeLogger) SetMinLevel(level Level) Logger {
	n.minLevel = level
	return n
}

func (l *TeeLogger) LogImpl(level Level, msg string, fielders []Fielder, opts []PrinterOption, d ImplData) {
	d.MergeMinLevel(l.minLevel) // ensure our minLevel is taken into account

	l.Primary.LogImpl(level, msg, fielders, opts, d)

	d.AnchoredLine = 0 // secondary loggers don't support anchored lines
	l.Secondary.LogImpl(level, msg, fielders, opts, d)
}

func (l *TeeLogger) Transient(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Transient, msg, fielders, nil, ImplData{})
	return l
}

func (l *TeeLogger) Verbose(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Verbose, msg, fielders, nil, ImplData{})
	return l
}

func (l *TeeLogger) Info(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Info, msg, fielders, nil, ImplData{})
	return l
}

func (l *TeeLogger) Warning(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Warning, msg, fielders, nil, ImplData{})
	return l
}

func (l *TeeLogger) Error(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Error, msg, fielders, nil, ImplData{})
	return l
}

func (l *TeeLogger) Log(level Level, msg string, fielders ...Fielder) Logger {
	l.LogImpl(level, msg, fielders, nil, ImplData{})
	return l
}
