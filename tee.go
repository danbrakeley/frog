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

func (l *TeeLogger) AddFixedLine() Logger {
	return &TeeLogger{
		Primary:   AddFixedLine(l.Primary),
		Secondary: AddFixedLine(l.Secondary),
	}
}

func (l *TeeLogger) RemoveFixedLine() {
	RemoveFixedLine(l.Primary)
	RemoveFixedLine(l.Secondary)
}

func (l *TeeLogger) SetMinLevel(level Level) Logger {
	l.Primary.SetMinLevel(level)
	l.Secondary.SetMinLevel(level)
	return l
}

func (l *TeeLogger) Log(level Level, msg string, fields ...Fielder) Logger {
	l.Primary.Log(level, msg, fields...)
	l.Secondary.Log(level, msg, fields...)
	return l
}

func (l *TeeLogger) Transient(msg string, fields ...Fielder) Logger {
	l.Log(Transient, msg, fields...)
	return l
}

func (l *TeeLogger) Verbose(msg string, fields ...Fielder) Logger {
	l.Log(Verbose, msg, fields...)
	return l
}

func (l *TeeLogger) Info(msg string, fields ...Fielder) Logger {
	l.Log(Info, msg, fields...)
	return l
}

func (l *TeeLogger) Warning(msg string, fields ...Fielder) Logger {
	l.Log(Warning, msg, fields...)
	return l
}

func (l *TeeLogger) Error(msg string, fields ...Fielder) Logger {
	l.Log(Error, msg, fields...)
	return l
}

func (l *TeeLogger) Fatal(msg string, fields ...Fielder) {
	l.Log(Fatal, msg, fields...)
}
