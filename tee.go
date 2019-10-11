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

func (l *TeeLogger) Logf(level Level, format string, a ...interface{}) Logger {
	l.Primary.Logf(level, format, a...)
	l.Secondary.Logf(level, format, a...)
	return l
}

func (l *TeeLogger) Transientf(format string, a ...interface{}) Logger {
	l.Logf(Transient, format, a...)
	return l
}

func (l *TeeLogger) Verbosef(format string, a ...interface{}) Logger {
	l.Logf(Verbose, format, a...)
	return l
}

func (l *TeeLogger) Infof(format string, a ...interface{}) Logger {
	l.Logf(Info, format, a...)
	return l
}

func (l *TeeLogger) Warningf(format string, a ...interface{}) Logger {
	l.Logf(Warning, format, a...)
	return l
}

func (l *TeeLogger) Errorf(format string, a ...interface{}) Logger {
	l.Logf(Error, format, a...)
	return l
}

func (l *TeeLogger) Fatalf(format string, a ...interface{}) {
	l.Logf(Fatal, format, a...)
}
