package frog

// TeeLogger directs all traffic to both a primary and secondary logger
// Note that of if one of your loggers supports anchor, make sure that is the Primary anchor.
type TeeLogger struct {
	Primary   Logger // Anchors are only supported though this logger
	Secondary Logger
}

func (l *TeeLogger) Close() {
	l.Primary.Close()
	l.Secondary.Close()
}

func (l *TeeLogger) Parent() Logger {
	// Anchoring relies on there being a single root parent.
	return l.Primary
}

func (l *TeeLogger) SetMinLevel(level Level) Logger {
	l.Primary.SetMinLevel(level)
	l.Secondary.SetMinLevel(level)
	return l
}

func (l *TeeLogger) LogImpl(anchoredLine int32, opts []PrinterOption, level Level, msg string, fields []Fielder) {
	l.Primary.LogImpl(anchoredLine, opts, level, msg, fields)
	l.Secondary.LogImpl(0, opts, level, msg, fields)
}

func (l *TeeLogger) Transient(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Transient, msg, fields)
	return l
}

func (l *TeeLogger) Verbose(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Verbose, msg, fields)
	return l
}

func (l *TeeLogger) Info(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Info, msg, fields)
	return l
}

func (l *TeeLogger) Warning(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Warning, msg, fields)
	return l
}

func (l *TeeLogger) Error(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Error, msg, fields)
	return l
}
