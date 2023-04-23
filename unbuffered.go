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
	// this space intentionally left blank (nothing to cleanup or flush)
}

func (l *Unbuffered) MinLevel() Level {
	return l.minLevel
}

func (l *Unbuffered) SetMinLevel(level Level) Logger {
	l.minLevel = level
	return l
}

func (l *Unbuffered) LogImpl(level Level, msg string, fielders []Fielder, opts []PrinterOption, d ImplData) {
	d.MergeMinLevel(l.minLevel)
	if level < d.MinLevel {
		return
	}

	fmt.Fprintf(l.writer, "%s\n", l.prn.Render(level, opts, msg, FieldifyAndAppend(d.Fields, fielders)))
}

func (l *Unbuffered) Transient(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Transient, msg, fielders, nil, ImplData{})
	return l
}

func (l *Unbuffered) Verbose(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Verbose, msg, fielders, nil, ImplData{})
	return l
}

func (l *Unbuffered) Info(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Info, msg, fielders, nil, ImplData{})
	return l
}

func (l *Unbuffered) Warning(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Warning, msg, fielders, nil, ImplData{})
	return l
}

func (l *Unbuffered) Error(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Error, msg, fielders, nil, ImplData{})
	return l
}

func (l *Unbuffered) Log(level Level, msg string, fielders ...Fielder) Logger {
	l.LogImpl(level, msg, fielders, nil, ImplData{})
	return l
}
