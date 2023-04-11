package frog

import (
	"fmt"
	"io"
)

type Unbuffered struct {
	writer   io.Writer
	prn      Printer
	minLevel Level
}

func NewUnbuffered(writer io.Writer, prn Printer) *Unbuffered {
	return &Unbuffered{
		writer:   writer,
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

func (l *Unbuffered) LogImpl(anchoredLine int32, opts []PrinterOption, level Level, msg string, fields []Fielder) {
	if level < l.minLevel {
		return
	}
	fmt.Fprintf(l.writer, "%s\n", l.prn.Render(level, opts, msg, fields))
}

func (l *Unbuffered) Transient(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Transient, msg, fields)
	return l
}

func (l *Unbuffered) Verbose(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Verbose, msg, fields)
	return l
}

func (l *Unbuffered) Info(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Info, msg, fields)
	return l
}

func (l *Unbuffered) Warning(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Warning, msg, fields)
	return l
}

func (l *Unbuffered) Error(msg string, fields ...Fielder) Logger {
	l.LogImpl(0, nil, Error, msg, fields)
	return l
}
