package frog

type Logger interface {
	// Close ensures any buffers are flushed and any resources released.
	// It is safe to call Close more than once (but consecutive calls do nothing).
	Close()

	// SetMinLevel sets the lowest Level that will be logged.
	SetMinLevel(level Level) Logger

	// Log is how all log lines are added.
	Log(level Level, opts []PrinterOption, msg string, fields []Fielder) Logger

	// Transient et al are just shortcuts for calling Log with specific levels.
	Transient(msg string, fields ...Fielder) Logger
	Verbose(msg string, fields ...Fielder) Logger
	Info(msg string, fields ...Fielder) Logger
	Warning(msg string, fields ...Fielder) Logger
	Error(msg string, fields ...Fielder) Logger
}

// ChildLogger is the interface for loggers that feed back to a parent.
type ChildLogger interface {
	// Parent returns the parent Logger, or nil if it has no parent.
	Parent() Logger
}

// AnchorAdder is the interface for loggers that support anchoring a line to the
// bottom of the output, for progress bars or other transient status messages.
type AnchorAdder interface {
	AddAnchor() Logger
}

// AnchorRemover is the interface that an anchor logger must implement
// in order for the anchor to be removed before app end.
type AnchorRemover interface {
	RemoveAnchor()
}

// Printerer is an interface that exposes one's Printer
type Printerer interface {
	Printer() Printer
}
