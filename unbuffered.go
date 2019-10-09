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

func (l *Unbuffered) Logf(level Level, format string, a ...interface{}) Logger {
	if level < l.minLevel {
		return l
	}
	msg := l.prn.Render(false, false, level, format, a...)
	fmt.Fprintf(l.w, "%s\n", msg)
	if level == Fatal {
		os.Exit(-1)
	}
	return l
}

func (l *Unbuffered) Transientf(format string, a ...interface{}) Logger {
	l.Logf(Transient, format, a...)
	return l
}

func (l *Unbuffered) Verbosef(format string, a ...interface{}) Logger {
	l.Logf(Verbose, format, a...)
	return l
}

func (l *Unbuffered) Infof(format string, a ...interface{}) Logger {
	l.Logf(Info, format, a...)
	return l
}

func (l *Unbuffered) Warningf(format string, a ...interface{}) Logger {
	l.Logf(Warning, format, a...)
	return l
}

func (l *Unbuffered) Errorf(format string, a ...interface{}) Logger {
	l.Logf(Error, format, a...)
	return l
}

func (l *Unbuffered) Fatalf(format string, a ...interface{}) {
	l.Logf(Fatal, format, a...)
}
