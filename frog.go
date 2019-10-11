package frog

import (
	"os"

	"github.com/mattn/go-isatty"
)

type NewLogger byte

const (
	Auto NewLogger = iota
	Basic
	JSON
)

// New creates a Buffered logger that writes to os.Stdout, and autodetects
// any attached Terminal on stdout to decide if ANSI should be used.
// The caller is responsible for calling Close() before the process ends.
func New(t NewLogger) Logger {
	switch t {
	case Auto:
		cfg := Config{
			Writer:   os.Stdout,
			UseAnsi:  isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()),
			UseColor: !IsNoColorSet,
		}
		prn := &TextPrinter{
			PrintTime:  true,
			PrintLevel: true,
		}
		return NewBuffered(cfg, prn)

	case Basic:
		prn := &TextPrinter{
			PrintTime:  true,
			PrintLevel: true,
		}
		return NewUnbuffered(os.Stdout, prn)

	case JSON:
		return NewUnbuffered(os.Stdout, &JSONPrinter{})
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
