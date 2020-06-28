package frog

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

var isNoColorSet = false

func init() {
	_, exists := os.LookupEnv("NO_COLOR")
	isNoColorSet = exists
}

type NewLogger byte

const (
	Auto NewLogger = iota
	Basic
	JSON
)

type Option byte

const (
	// UseStdout sends all output to stdout (default).
	UseStdout Option = iota
	// UseStderr sends all output to stderr.
	UseStderr
	// Color enables use of ANSI commands to add color (default). Note that colors cannot be enabled when
	// using the Basic or JSON logger types, or if a NO_COLOR environment variable is present.
	Color
	// NoColor disables use of ANSI commands to add color.
	NoColor
	// ShowTimestamps enables the inclusion of time/date (default). Note that JSON type always adds timestamps.
	ShowTimestamps
	// HideTimestamps disables the inclusion of time/date.
	HideTimestamps
	// ShowLevel enables the inclusion of the log level of each log line (ie "[nfo]", "[WRN]", "[ERR]", etc).
	ShowLevel
	// HideLevel disables the inclusion of the log level with each log line. Note that JSON type always adds log level.
	HideLevel
)

// HasTerminal returns true if the passed writer is connected to a terminal
func HasTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fd := f.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

// New creates a Logger that writes to os.Stdout, depending on the NewLogger type passed to it:
// - Auto - if terminal detected on stdout, then colors and fixed lines are supported (else, uses Basic)
// - Basic - no colors or fixed lines, no buffering
// - JSON - no colors or fixed lines, no buffering, and each line is a valid JSON object
// Resulting Logger can be modified by including 1 or more NewOpts after the NewLogger type.
// The caller is responsible for calling Close() when done with the returned Logger.
func New(t NewLogger, opts ...Option) Logger {
	if t == Auto && !HasTerminal(os.Stdout) {
		t = Basic
	}

	// process options
	var useColor bool = true
	var showTime bool = true
	var showLevel bool = true
	var writer io.Writer = os.Stdout
	for _, opt := range opts {
		switch opt {
		case UseStdout:
			writer = os.Stdout
		case UseStderr:
			writer = os.Stderr
		case Color:
			useColor = true
		case NoColor:
			useColor = false
		case ShowTimestamps:
			showTime = true
		case HideTimestamps:
			showTime = false
		case ShowLevel:
			showLevel = true
		case HideLevel:
			showLevel = false
		}
	}

	if isNoColorSet {
		useColor = false
	}

	switch t {
	case Auto:
		return NewBuffered(
			Config{
				Writer:   writer,
				UseColor: useColor,
			},
			&TextPrinter{
				PrintTime:  showTime,
				PrintLevel: showLevel,
			},
		)
	case Basic:
		return NewUnbuffered(
			Config{
				Writer:   writer,
				UseColor: false,
			},
			&TextPrinter{
				PrintTime:  showTime,
				PrintLevel: showLevel,
			},
		)
	case JSON:
		return NewUnbuffered(
			Config{
				Writer:   writer,
				UseColor: false,
			},
			&JSONPrinter{},
		)
	}

	return nil
}

// AddFixedLine adds a new logger on a fixed line, if supported.
// Else, returns passed in Logger.
func AddFixedLine(log Logger) Logger {
	fla, ok := log.(FixedLineAdder)
	if ok {
		fl := fla.AddFixedLine()
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

func RemoveFixedLine(log Logger) {
	flr, ok := log.(FixedLineRemover)
	if ok {
		flr.RemoveFixedLine()
	}
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
