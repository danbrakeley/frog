package frog

type RootLogger interface {
	Logger

	// Close ensures any buffers are flushed and any resources released.
	// It is safe to call Close more than once (but consecutive calls do nothing).
	Close()
}

type Logger interface {
	// MinLevel gets the minimum level that is filtered by this Logger instance.
	// If this Logger is part of a chain of nested Loggers, note that that this only returns the min
	// level of this link in the chain. This Logger's parents may have more restrictive min levels
	// that prevent log lines from being displayed.
	MinLevel() Level
	// SetMinLevel sets the lowest Level that will be accepted by this Logger.
	// If this Logger has parent(s), the effective MinLevel will be the max of each logger's min level.
	SetMinLevel(level Level) Logger

	// Transient logs a string (with optional fields) with the log level set to Transient.
	Transient(msg string, fields ...Fielder) Logger
	// Verbose logs a string (with optional fields) with the log level set to Verbose.
	Verbose(msg string, fields ...Fielder) Logger
	// Info logs a string (with optional fields) with the log level set to Info.
	Info(msg string, fields ...Fielder) Logger
	// Warning logs a string (with optional fields) with the log level set to Warning.
	Warning(msg string, fields ...Fielder) Logger
	// Error logs a string (with optional fields) with the log level set to Error.
	Error(msg string, fields ...Fielder) Logger
	// Log logs a string (with optional fields) with the log level set to the passed in value.
	Log(level Level, msg string, fields ...Fielder) Logger

	// LogImpl is called by children to pass up log events to the root Logger.
	LogImpl(level Level, msg string, fields []Fielder, opts []PrinterOption, d ImplData)
}

// ChildLogger is the interface for loggers that feed back to a parent.
type ChildLogger interface {
	// Parent returns the parent Logger, or nil if it has no parent.
	Parent() Logger
}

// AnchorAdder is the interface for loggers that support anchoring a line to the
// bottom of the output, for progress bars or other transient status messages.
type AnchorAdder interface {
	AddAnchor(parent Logger) Logger
}

// AnchorRemover is the interface that an anchor logger must implement
// in order for the anchor to be removed before app end.
type AnchorRemover interface {
	RemoveAnchor()
}
