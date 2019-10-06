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
	Progress Level = iota // real-time progress bars, more verbose than a log file would need
	Verbose               // debugging info
	Info                  // normal message
	Warning               // something unusual happened
	Error                 // something bad happened
	Fatal                 // it's all over

	levelMax
	levelMin Level = 0
)

type Logger interface {
	// Close should be called before the process ends, to ensure everything is flushed.
	Close()

	// RootLogger returns the parent Logger, or itself if it has no parent.
	RootLogger() Logger

	// AddFixedLine creates a logger that always overwrites the same terminal line,
	// and always writes line level Progress.
	// Returned Logger should have its Close() called before its parent.
	AddFixedLine() Logger

	// MinLevel returns the lowest Level that will be logged.
	MinLevel() Level
	// SetMinLevel sets the lowest Level that will be logged.
	SetMinLevel(level Level) Logger

	// Printf is how log lines are added.
	Printf(level Level, format string, a ...interface{}) Logger

	// the remaining calls are just shortcuts for calling Printf with a specific level

	Progressf(format string, a ...interface{}) Logger
	Verbosef(format string, a ...interface{}) Logger
	Infof(format string, a ...interface{}) Logger
	Warningf(format string, a ...interface{}) Logger
	Errorf(format string, a ...interface{}) Logger
	Fatalf(format string, a ...interface{})
}

// New creates a Buffered logger that writes to os.Stdout, and autodetects
// any attached Terminal on stdout to decide if ANSI should be used.
// The caller is responsible for calling Close() before the process ends.
func New() Logger {
	cfg := Config{
		Writer:   os.Stdout,
		UseAnsi:  isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()),
		UseColor: !IsNoColorSet,
	}
	prn := Printer{
		PrintTime:  true,
		PrintLevel: true,
	}
	return NewBuffered(cfg, prn)
}
