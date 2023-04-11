package frog

type NullLogger struct {
}

func (n *NullLogger) Close() {
}

func (n *NullLogger) SetMinLevel(level Level) Logger {
	return n
}

func (n *NullLogger) LogImpl(anchoredLine int32, opts []PrinterOption, level Level, msg string, fields []Fielder) {
}

func (n *NullLogger) Transient(format string, a ...Fielder) Logger {
	return n
}

func (n *NullLogger) Verbose(format string, a ...Fielder) Logger {
	return n
}

func (n *NullLogger) Info(format string, a ...Fielder) Logger {
	return n
}

func (n *NullLogger) Warning(format string, a ...Fielder) Logger {
	return n
}

func (n *NullLogger) Error(format string, a ...Fielder) Logger {
	return n
}
