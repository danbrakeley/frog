package frog

import (
	"os"
)

var IsNoColorSet = false

func init() {
	_, exists := os.LookupEnv("NO_COLOR")
	IsNoColorSet = exists
}

type Logger interface {
	// Close ensures any buffers are flushed and any resources released.
	// It is safe to call Close more than once (but consecutive calls do nothing).
	Close()

	// SetMinLevel sets the lowest Level that will be logged.
	SetMinLevel(level Level) Logger

	// Logf is how log lines are added.
	Logf(level Level, format string, a ...interface{}) Logger

	// Transientf et al are just shortcuts for calling Logf with specific levels
	// Note that Fatalf doesn't return itself like the others do because it isn't expected to return at all.
	Transientf(format string, a ...interface{}) Logger
	Verbosef(format string, a ...interface{}) Logger
	Infof(format string, a ...interface{}) Logger
	Warningf(format string, a ...interface{}) Logger
	Errorf(format string, a ...interface{}) Logger
	Fatalf(format string, a ...interface{})
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
