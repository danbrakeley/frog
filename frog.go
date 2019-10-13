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
// - Auto - if terminal detected on stdout, then colors and fixed lines are supported
// - Basic - fixed lines are not supported, but colors are, and output never needs flushing
// - JSON - no colors or fixed lines, and each line is a valid JSON object
// The caller is responsible for calling Close() when done with the returned Logger.
func New(t NewLogger) Logger {
	if t == Auto && !HasTerminal(os.Stdout) {
		t = Basic
	}

	switch t {
	case Auto:
		cfg := Config{
			Writer:   os.Stdout,
			UseColor: !isNoColorSet,
		}
		prn := &TextPrinter{
			PrintTime:  true,
			PrintLevel: true,
		}
		return NewBuffered(cfg, prn)

	case Basic:
		cfg := Config{
			Writer:   os.Stdout,
			UseColor: !isNoColorSet,
		}
		prn := &TextPrinter{
			PrintTime:  true,
			PrintLevel: true,
		}
		return NewUnbuffered(cfg, prn)

	case JSON:
		return NewUnbuffered(Config{Writer: os.Stdout, UseColor: false}, &JSONPrinter{})
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
