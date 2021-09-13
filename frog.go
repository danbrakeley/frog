package frog

import (
	"io"
	"os"

	"github.com/mattn/go-isatty"
)

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
// - Auto - if terminal detected on stdout, then colors and anchored lines are supported (else, uses Basic)
// - Basic - no colors or anchored lines, no buffering
// - JSON - no colors or anchored lines, no buffering, and each line is a valid JSON object
// Resulting Logger can be modified by including 1 or more NewOpts after the NewLogger type.
// The caller is responsible for calling Close() when done with the returned Logger.
func New(t NewLogger, opts ...PrinterOption) Logger {
	if t == Auto && !HasTerminal(os.Stdout) {
		t = Basic
	}

	switch t {
	case Auto:
		prn := TextPrinter{Palette: PalColor, PrintTime: true, PrintLevel: true, FieldIndent: 20, PrintMessageLast: false}
		return NewBuffered(os.Stdout, prn.SetOptions(opts...))
	case Basic:
		prn := TextPrinter{Palette: PalColor, PrintTime: true, PrintLevel: true, FieldIndent: 20, PrintMessageLast: false}
		return NewUnbuffered(os.Stdout, prn.SetOptions(opts...))
	case JSON:
		return NewUnbuffered(os.Stdout, &JSONPrinter{})
	}

	return nil
}

// AddAnchor adds a new logger on an anchored line, if supported.
// Else, returns passed in Logger.
func AddAnchor(log Logger) Logger {
	aa, ok := log.(AnchorAdder)
	if ok {
		a := aa.AddAnchor()
		if a == nil {
			return log
		}
		return a
	}
	// if we are a child that doesn't create its own anchored lines, then pass up to parent
	parent := Parent(log)
	if parent != nil {
		return AddAnchor(parent)
	}
	return log
}

func RemoveAnchor(log Logger) {
	ar, ok := log.(AnchorRemover)
	if ok {
		ar.RemoveAnchor()
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
