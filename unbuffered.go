package frog

import (
	"fmt"
	"io"
	"os"
)

type Unbuffered struct {
	w        io.Writer
	prn      Printer
	minLevel Level
}

func NewUnbuffered(w io.Writer, prn Printer) *Unbuffered {
	return &Unbuffered{
		w:        w,
		prn:      prn,
		minLevel: Info,
	}
}

func (l *Unbuffered) Close() {
}

func (l *Unbuffered) SetMinLevel(level Level) Logger {
	l.minLevel = level
	return l
}

func (l *Unbuffered) Log(level Level, msg string, fields ...Fielder) Logger {
	if level < l.minLevel {
		return l
	}
	fmt.Fprintf(l.w, "%s\n", l.prn.Render(false, false, level, msg, fields...))
	if level == Fatal {
		os.Exit(-1)
	}
	return l
}

func (l *Unbuffered) Transient(msg string, fields ...Fielder) Logger {
	l.Log(Transient, msg, fields...)
	return l
}

func (l *Unbuffered) Verbose(msg string, fields ...Fielder) Logger {
	l.Log(Verbose, msg, fields...)
	return l
}

func (l *Unbuffered) Info(msg string, fields ...Fielder) Logger {
	l.Log(Info, msg, fields...)
	return l
}

func (l *Unbuffered) Warning(msg string, fields ...Fielder) Logger {
	l.Log(Warning, msg, fields...)
	return l
}

func (l *Unbuffered) Error(msg string, fields ...Fielder) Logger {
	l.Log(Error, msg, fields...)
	return l
}

func (l *Unbuffered) Fatal(msg string, fields ...Fielder) {
	l.Log(Fatal, msg, fields...)
}
