package frog

type Logger interface {
	// Close ensures any buffers are flushed and any resources released.
	// It is safe to call Close more than once (but consecutive calls do nothing).
	Close()

	// SetMinLevel sets the lowest Level that will be logged.
	SetMinLevel(level Level) Logger

	// Log is how all log lines are added.
	Log(level Level, format string, a ...Fielder) Logger

	// Transient et al are just shortcuts for calling Log with specific levels.
	// Note that Fatal will never return, as it flushes any buffers then calls os.Exit(-1).
	Transient(format string, a ...Fielder) Logger
	Verbose(format string, a ...Fielder) Logger
	Info(format string, a ...Fielder) Logger
	Warning(format string, a ...Fielder) Logger
	Error(format string, a ...Fielder) Logger
	Fatal(format string, a ...Fielder)
}

// ChildLogger is the interface for loggers that feed back to a parent.
type ChildLogger interface {
	// Parent returns the parent Logger, or nil if it has no parent.
	Parent() Logger
}

// FixedLineAdder is the interface for loggers that support fixing a line in place,
// for progress bars or other transient status messages.
type FixedLineAdder interface {
	AddFixedLine() Logger
}

// FixedLineRemover is the interface that a fixed line logger must implement
// in order for the fixed line to be removed before app end.
type FixedLineRemover interface {
	RemoveFixedLine()
}
