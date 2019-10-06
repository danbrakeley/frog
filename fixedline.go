package frog

// FixedLine is a Logger that attempts to overwrite the same line in the terminal,
// allowing progress bars and other simple UI for the human to consume.
//
// If the Printer's CanUseAnsi is false, then it simply redirects to the normal
// behavior of the parent Buffered Logger.
//
// If the Printer's CanUseAnsi is true, then the Progress level is always
// printed, regardless of the MinLevel. This allows progress bars that do
// not pollute logs with garbage when not connected to a terminal.
type FixedLine struct {
	parent *Buffered
	line   int32
}

func newFixedLine(b *Buffered, line int32) *FixedLine {
	return &FixedLine{parent: b, line: line}
}

func (l *FixedLine) Close() {
	// Currently there's nothing line-specific to clean up.
	// Unfortunately this call can't easily wait until all prints from this line are flushed.
}

func (l *FixedLine) AddFixedLine() Logger {
	// Currently just returns self
	return l
}

func (l *FixedLine) MinLevel() Level {
	return l.parent.MinLevel()
}

func (l *FixedLine) SetMinLevel(level Level) Logger {
	l.parent.SetMinLevel(level)
	return l
}

func (l *FixedLine) Printf(level Level, format string, a ...interface{}) Logger {
	if !l.parent.p.CanUseAnsi {
		l.parent.Printf(level, format, a...)
		return l
	}

	if level < l.parent.MinLevel() && level != Progress {
		return l
	}

	l.parent.printfImpl(l.line, level, format, a...)
	return l
}

func (l *FixedLine) Progressf(format string, a ...interface{}) Logger {
	l.Printf(Progress, format, a...)
	return l
}

func (l *FixedLine) Verbosef(format string, a ...interface{}) Logger {
	l.Printf(Verbose, format, a...)
	return l
}

func (l *FixedLine) Infof(format string, a ...interface{}) Logger {
	l.Printf(Info, format, a...)
	return l
}

func (l *FixedLine) Warningf(format string, a ...interface{}) Logger {
	l.Printf(Warning, format, a...)
	return l
}

func (l *FixedLine) Errorf(format string, a ...interface{}) Logger {
	l.Printf(Error, format, a...)
	return l
}

func (l *FixedLine) Fatalf(format string, a ...interface{}) {
	l.Printf(Fatal, format, a...)
}
