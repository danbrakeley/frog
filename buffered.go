package frog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"

	"github.com/danbrakeley/frog/ansi"
)

type Buffered struct {
	minLevel Level
	writer   io.Writer
	p        Printer
	ch       chan Msg
	wg       sync.WaitGroup

	lastLine int32 // to keep thread safe, use atomic reads/writes/math
}

type MsgType byte

const (
	MsgPrint MsgType = iota
	MsgFatal
	MsgAddLine
)

type Msg struct {
	Type MsgType
	Line int32
	Msg  string
}

func NewBuffered(w io.Writer, p Printer) *Buffered {
	l := &Buffered{
		minLevel: Info,
		writer:   w,
		p:        p,
		ch:       make(chan Msg),
		wg:       sync.WaitGroup{},
	}

	l.wg.Add(1)
	go func() {
		l.processor()
		l.wg.Done()
	}()

	return l
}

// Close should be called before the app exits, to ensure any last line is done writing.
func (l *Buffered) Close() {
	l.Verbosef("buffered log closing")
	close(l.ch)
	l.wg.Wait()
}

// AddFixedLine creates a Logger that overwrites the same line.
// Thread safe.
func (l *Buffered) AddFixedLine() Logger {
	lastLine := atomic.AddInt32(&l.lastLine, 1)
	l.ch <- Msg{Type: MsgAddLine}
	return newFixedLine(l, lastLine)
}

func (l *Buffered) processor() {
	for {
		msg, ok := <-l.ch
		if !ok || msg.Type == MsgFatal {
			return
		}

		lastLine := atomic.LoadInt32(&l.lastLine)

		// if we can't use ansi, or there have been no fixed lines added, then just skip fancy formatting
		if !l.p.CanUseAnsi || lastLine == 0 {
			if msg.Type == MsgPrint {
				fmt.Fprintf(l.writer, "%s\n", msg.Msg)
			}
			continue
		}

		switch msg.Type {
		case MsgAddLine:
			// A new line has been added, so use a newline to scroll the console (if needed).
			// The cursor is always left at the bottom, so it is already in the right spot for this.
			fmt.Fprintf(l.writer, "\n")
			continue
		case MsgPrint:
			// Something tried to log without specifying a line, which isn't currently supported.
			// For now, just print it at the bottom (will overwrite any previous line).
			if msg.Line <= 0 {
				fmt.Fprintf(l.writer, "%s%s\n", msg.Msg, ansi.EraseEOL)
				continue
			}
			// fall out of this switch (without continuing)
		default:
			continue
		}

		// The cursor is kept under the bottom-most fixed line, so we'll move to the correct line,
		// print, then move back.
		offset := int(1 + lastLine - msg.Line)
		fmt.Fprint(l.writer, ansi.PrevLine(offset))
		fmt.Fprint(l.writer, msg.Msg)
		fmt.Fprint(l.writer, ansi.EraseEOL)
		fmt.Fprint(l.writer, ansi.NextLine(offset))
	}
}

func (l *Buffered) MinLevel() Level {
	return l.minLevel
}

func (l *Buffered) SetMinLevel(level Level) Logger {
	if level < levelMin || level >= levelMax {
		panic(fmt.Errorf("level %v is not in valid range [%v,%v)", level, levelMin, levelMax))
	}
	l.minLevel = level
	return l
}

func (l *Buffered) Printf(level Level, format string, a ...interface{}) Logger {
	if level < l.minLevel {
		return l
	}
	l.printfImpl(0, level, format, a...)
	return l
}

func (l *Buffered) printfImpl(fixedLine int32, level Level, format string, a ...interface{}) {
	l.ch <- Msg{Line: fixedLine, Msg: l.p.Sprintf(level, format, a...)}
	if level == Fatal {
		// we can't just close this channel, because another thread may still be trying to write to it
		l.ch <- Msg{Type: MsgFatal}
		l.wg.Wait()
		os.Exit(-1)
	}
}

func (l *Buffered) Progressf(format string, a ...interface{}) Logger {
	l.Printf(Progress, format, a...)
	return l
}

func (l *Buffered) Verbosef(format string, a ...interface{}) Logger {
	l.Printf(Verbose, format, a...)
	return l
}

func (l *Buffered) Infof(format string, a ...interface{}) Logger {
	l.Printf(Info, format, a...)
	return l
}

func (l *Buffered) Warningf(format string, a ...interface{}) Logger {
	l.Printf(Warning, format, a...)
	return l
}

func (l *Buffered) Errorf(format string, a ...interface{}) Logger {
	l.Printf(Error, format, a...)
	return l
}

func (l *Buffered) Fatalf(format string, a ...interface{}) {
	l.Printf(Fatal, format, a...)
}
