package frog

// TeeLogger directs all traffic to both a primary and secondary logger
type TeeLogger struct {
	Primary   Logger
	Secondary Logger
}

func (l *TeeLogger) Close() {
	l.Primary.Close()
	l.Secondary.Close()
}

func (l *TeeLogger) AddAnchor() Logger {
	return &TeeLogger{
		Primary:   AddAnchor(l.Primary),
		Secondary: AddAnchor(l.Secondary),
	}
}

func (l *TeeLogger) RemoveAnchor() {
	RemoveAnchor(l.Primary)
	RemoveAnchor(l.Secondary)
}

func (l *TeeLogger) SetMinLevel(level Level) Logger {
	l.Primary.SetMinLevel(level)
	l.Secondary.SetMinLevel(level)
	return l
}

func (l *TeeLogger) Log(level Level, opts []PrinterOption, msg string, fields []Fielder) Logger {
	l.Primary.Log(level, opts, msg, fields)
	l.Secondary.Log(level, opts, msg, fields)
	return l
}

func (l *TeeLogger) Transient(msg string, fields ...Fielder) Logger {
	l.Log(Transient, nil, msg, fields)
	return l
}

func (l *TeeLogger) Verbose(msg string, fields ...Fielder) Logger {
	l.Log(Verbose, nil, msg, fields)
	return l
}

func (l *TeeLogger) Info(msg string, fields ...Fielder) Logger {
	l.Log(Info, nil, msg, fields)
	return l
}

func (l *TeeLogger) Warning(msg string, fields ...Fielder) Logger {
	l.Log(Warning, nil, msg, fields)
	return l
}

func (l *TeeLogger) Error(msg string, fields ...Fielder) Logger {
	l.Log(Error, nil, msg, fields)
	return l
}
