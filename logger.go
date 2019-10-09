package frog

import (
	"os"

	"github.com/mattn/go-isatty"
)

var IsNoColorSet = false

func init() {
	_, exists := os.LookupEnv("NO_COLOR")
	IsNoColorSet = exists
}

type Level byte

const (
	Transient Level = iota // strictly unimportant, ie progress bars, real-time byte counts, estimated time remaining, etc
	Verbose                // debugging info
	Info                   // normal message
	Warning                // something unusual happened
	Error                  // something bad happened
	Fatal                  // stop everything right now

	levelMax
	levelMin Level = 0
)

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

// FixedLineLogger is the interface for loggers that support fixing a line in place,
// for progress bars or other transient status messages.
type FixedLineLogger interface {
	// AddFixedLine creates a logger that always overwrites the same terminal line,
	// and always writes line level Progress.
	// Returned Logger should have its Close() called before its parent.
	AddFixedLine() Logger
}

type NewLogger byte

const (
	Auto NewLogger = iota
	Basic
)

// New creates a Buffered logger that writes to os.Stdout, and autodetects
// any attached Terminal on stdout to decide if ANSI should be used.
// The caller is responsible for calling Close() before the process ends.
func New(t NewLogger) Logger {
	prn := Printer{
		PrintTime:  true,
		PrintLevel: true,
	}

	switch t {
	case Auto:
		cfg := Config{
			Writer:   os.Stdout,
			UseAnsi:  isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()),
			UseColor: !IsNoColorSet,
		}
		return NewBuffered(cfg, prn)

	case Basic:
		return NewUnbuffered(os.Stdout, prn)
	}

	return nil
}

// AddFixedLine adds a new logger on a fixed line, if supported.
// Else, returns passed in Logger.
func AddFixedLine(log Logger) Logger {
	fll, ok := log.(FixedLineLogger)
	if ok {
		fl := fll.AddFixedLine()
		if fl == nil {
			return log
		}
		return fl
	}
	// if we are a child that doesn't create its own fixed lines, then pass up to parent
	parent := Parent(log)
	if parent != nil {
		return AddFixedLine(parent)
	}
	return log
}

func Parent(log Logger) Logger {
	child, ok := log.(ChildLogger)
	if !ok {
		return nil
	}
	return child.Parent()
}

func ParentOrSelf(log Logger) Logger {
	child, ok := log.(ChildLogger)
	if !ok {
		return log
	}
	parent := child.Parent()
	if parent == nil {
		return log
	}
	return parent
}
