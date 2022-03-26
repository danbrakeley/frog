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

func (l *Unbuffered) Log(level Level, opts []PrinterOption, msg string, fields []Fielder) Logger {
	if level < l.minLevel {
		return l
	}
	fmt.Fprintf(l.writer, "%s\n", l.prn.Render(level, opts, msg, fields))
	return l
}

func (l *Unbuffered) Transient(msg string, fields ...Fielder) Logger {
	l.Log(Transient, nil, msg, fields)
	return l
}

func (l *Unbuffered) Verbose(msg string, fields ...Fielder) Logger {
	l.Log(Verbose, nil, msg, fields)
	return l
}

func (l *Unbuffered) Info(msg string, fields ...Fielder) Logger {
	l.Log(Info, nil, msg, fields)
	return l
}

func (l *Unbuffered) Warning(msg string, fields ...Fielder) Logger {
	l.Log(Warning, nil, msg, fields)
	return l
}

func (l *Unbuffered) Error(msg string, fields ...Fielder) Logger {
	l.Log(Error, nil, msg, fields)
	return l
}
