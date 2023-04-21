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

// HasTerminal returns true if the passed writer is connected to a terminal.
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
func New(t NewLogger, opts ...PrinterOption) RootLogger {
	hasTerminal := false
	if t == Auto {
		hasTerminal = HasTerminal(os.Stdout)
		if !hasTerminal {
			t = Basic
		}
	}

	switch t {
	case Auto:
		prn := TextPrinter{palette: DefaultPalette.toANSI(), printTime: true, printLevel: true, fieldIndent: 20}
		return NewBuffered(os.Stdout, hasTerminal, prn.SetOptions(opts...))
	case Basic:
		prn := TextPrinter{printTime: true, printLevel: true, fieldIndent: 20}
		return NewUnbuffered(os.Stdout, prn.SetOptions(opts...))
	case JSON:
		return NewUnbuffered(os.Stdout, &JSONPrinter{})
	}

	return nil
}

// AddAnchor adds a new logger on an anchored line (if supported).
// Else, returns passed in Logger.
func AddAnchor(log Logger) Logger {
	var aa AnchorAdder
	var ok bool

	// search up the chain of parents for something that supports adding anchors
	tmp := log
	for {
		if tmp == nil {
			break
		}
		aa, ok = tmp.(AnchorAdder)
		if ok {
			break
		}
		tmp = Parent(tmp)
	}

	// If there's no AnchorAdder available, then just wrap the given Logger in a NoAnchorLogger.
	// Note that we can't just return the original Logger here, as the caller is expecting an extra
	// level of indirection.
	if aa == nil {
		return newNoAnchor(log)
	}

	alog := aa.AddAnchor(log)
	if alog == nil {
		return log
	}
	return alog
}

// RemoveAnchor needs to be passed the logger that was returned by AddAnchor.
// It can also work by being passed a child of that Logger.
func RemoveAnchor(log Logger) {
	var ar AnchorRemover
	var ok bool

	// search up the chain of parents for something that supports removing anchors
	tmp := log
	for {
		if tmp == nil {
			break
		}
		ar, ok = tmp.(AnchorRemover)
		if ok {
			ar.RemoveAnchor()
			break
		}
		tmp = Parent(tmp)
	}
}

func Parent(log Logger) Logger {
	child, ok := log.(ChildLogger)
	if !ok {
		return nil
	}
	return child.Parent()
}

// WithFields creates a new Logger that wraps the parent logger, adding the specified
// fields when the returned Logger is used, but leaving the parent unmodified.
func WithFields(log Logger, fields ...Fielder) Logger {
	return newCustomizerLogger(log, nil, fields)
}

// WithOptions creates a new Logger that wraps the passed Logger, adding the specified
// PrinterOptions when the returned Logger is used, but leaving the parent unmodified.
func WithOptions(log Logger, opts ...PrinterOption) Logger {
	return newCustomizerLogger(log, opts, nil)
}

// WithOptionsAndFields creates a new Logger that wraps the passed Logger, adding the
// specified fields and PrinterOptions when the return Logger is used, but leaving
// the parent unmodified.
func WithOptionsAndFields(log Logger, opts []PrinterOption, fields []Fielder) Logger {
	return newCustomizerLogger(log, opts, fields)
}
