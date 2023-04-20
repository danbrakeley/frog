package frog

type NullLogger struct {
	minLevel Level
}

func (n *NullLogger) Close() {
}

func (n *NullLogger) MinLevel() Level {
	return n.minLevel
}

func (n *NullLogger) SetMinLevel(level Level) Logger {
	n.minLevel = level
	return n
}

func (n *NullLogger) LogImpl(level Level, msg string, fields []Fielder, opts []PrinterOption, d ImplData) {
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

func (n *NullLogger) Log(level Level, msg string, fields ...Fielder) Logger {
	return n
}
