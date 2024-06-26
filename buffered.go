package frog

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/danbrakeley/ansi"
	"github.com/danbrakeley/frog/internal/terminal"
)

type Buffered struct {
	minLevel Level
	writer   io.Writer
	prn      Printer
	ch       chan bufmsg
	wg       sync.WaitGroup

	isClosed    int32 // to keep thread safe, use atomic reads/writes/math
	openAnchors int32 // to keep thread safe, use atomic reads/writes/math
}

type msgType byte

const (
	mtPrint      msgType = iota // string to print
	mtAddLine                   // add an anchored line
	mtRemoveLine                // remove an acnhored line
)

type bufmsg struct {
	Type  msgType
	Line  int32
	Level Level
	Msg   string
}

func NewBuffered(writer io.Writer, requestTerminalSize bool, prn Printer) *Buffered {
	l := &Buffered{
		minLevel: Info,
		writer:   writer,
		prn:      prn,
		ch:       make(chan bufmsg),
		wg:       sync.WaitGroup{},
	}

	l.wg.Add(1)
	if l.writer == nil {
		return nil
	}
	go func() {
		if requestTerminalSize {
			c, _, err := terminal.GetSize()
			if err != nil {
				c = -1
			}
			l.prn = l.prn.SetOptions(POTransientLineLength(c))
		}
		l.processor()
		l.wg.Done()
	}()

	return l
}

// Close should be called before the app exits, to ensure any buffered output is flushed.
// Thread safe.
func (l *Buffered) Close() {
	isClosed := atomic.AddInt32(&l.isClosed, 1)
	// protect against multiple calls to close
	if isClosed != 1 {
		l.wg.Wait()
		return
	}

	// // keep around for debug
	// if l.minLevel <= Verbose {
	// 	_, file, line, ok := runtime.Caller(1)
	// 	l.Verbose("buffered log closing", String("file", file), Int("line", line), Bool("ok", ok))
	// }
	close(l.ch)
	l.wg.Wait()
}

// AddAnchor creates a Logger that is "achored" to the bottom of the output.
// This "anchoring" is achieved by using ANSI to re-draw the anchored line at
// the bottom as the output scrolls up.
// Thread safe.
func (l *Buffered) AddAnchor(parent Logger) Logger {
	lineNum := atomic.AddInt32(&l.openAnchors, 1)
	l.ch <- bufmsg{Type: mtAddLine, Line: lineNum}
	onClose := func() {
		l.ch <- bufmsg{Type: mtRemoveLine, Line: lineNum}
	}
	return newAnchor(parent, lineNum, onClose)
}

func (l *Buffered) processor() {
	type anchoredLine struct {
		lineNum int32
		str     string
	}

	// avoid allocations later on by reserving space up front
	anchoredLines := make([]anchoredLine, 0, 32)

	fnMustFindIdx := func(line int32) int {
		idx := -1
		for i, v := range anchoredLines {
			if line == v.lineNum {
				idx = i
				break
			}
		}
		if idx == -1 {
			panic(fmt.Errorf("buffered logger cannot find line %d", line))
		}
		return idx
	}

	for {
		msg, ok := <-l.ch
		if !ok {
			return
		}

		switch msg.Type {
		case mtAddLine:
			// ensure terminal scrolls down if needed to add a new line
			fmt.Fprint(l.writer, "\n")
			anchoredLines = append(anchoredLines, anchoredLine{lineNum: msg.Line})

		case mtRemoveLine:
			// find the line we are removing
			idx := fnMustFindIdx(msg.Line)

			// remove element
			copy(anchoredLines[idx:], anchoredLines[idx+1:])
			anchoredLines = anchoredLines[:len(anchoredLines)-1]

			// redraw/erase bottom lines as needed
			fmt.Fprint(l.writer, ansi.PrevLine(1+len(anchoredLines)-idx))
			for i := idx; i < len(anchoredLines); i++ {
				fmt.Fprint(l.writer, anchoredLines[i].str)
				fmt.Fprint(l.writer, ansi.EraseEOL)
				fmt.Fprint(l.writer, ansi.NextLine(1))
			}
			fmt.Fprint(l.writer, ansi.EraseEOL)

		case mtPrint:
			// if we aren't using anchored lines, then just print normally
			if len(anchoredLines) == 0 {
				fmt.Fprintf(l.writer, "%s\n", msg.Msg)
				continue
			}

			// If we are using anchored lines, but this msg doesn't have one specified, then move all
			// the anchored lines down, and draw this line above them.
			// If this does have an anchored line, but it is not Transient level, then also print it above.
			if msg.Line <= 0 || msg.Level > Transient {
				fmt.Fprint(l.writer, "\n")
				fmt.Fprint(l.writer, ansi.PrevLine(1+len(anchoredLines)))
				fmt.Fprintf(l.writer, "%s%s\n", msg.Msg, ansi.EraseEOL)

				for _, v := range anchoredLines {
					fmt.Fprintf(l.writer, "%s%s\n", v.str, ansi.EraseEOL)
				}

				// if we aren't using anchored lines, then we're done here...
				if msg.Line <= 0 {
					continue
				}
			}

			// The cursor is kept under the bottom-most anchored line, so we'll move to the correct line,
			// print, then move back.
			idx := fnMustFindIdx(msg.Line)
			anchoredLines[idx].str = msg.Msg
			offset := int(len(anchoredLines) - idx)
			fmt.Fprint(l.writer, ansi.PrevLine(offset))
			fmt.Fprint(l.writer, msg.Msg)
			fmt.Fprint(l.writer, ansi.EraseEOL)
			fmt.Fprint(l.writer, ansi.NextLine(offset))

		default:
		}
	}
}

func (l *Buffered) MinLevel() Level {
	return l.minLevel
}

func (l *Buffered) SetMinLevel(level Level) Logger {
	l.minLevel = level
	return l
}

func (l *Buffered) LogImpl(level Level, msg string, fielders []Fielder, opts []PrinterOption, d ImplData) {
	d.MergeMinLevel(l.minLevel)

	keep := false
	if d.AnchoredLine == 0 {
		keep = level >= d.MinLevel
	} else {
		keep = level >= d.MinLevel || level == Transient
	}
	if !keep {
		return
	}

	l.ch <- bufmsg{
		Line:  d.AnchoredLine,
		Level: level,
		Msg:   l.prn.Render(level, opts, msg, FieldifyAndAppend(d.Fields, fielders)),
	}
}

func (l *Buffered) Transient(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Transient, msg, fielders, nil, ImplData{})
	return l
}

func (l *Buffered) Verbose(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Verbose, msg, fielders, nil, ImplData{})
	return l
}

func (l *Buffered) Info(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Info, msg, fielders, nil, ImplData{})
	return l
}

func (l *Buffered) Warning(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Warning, msg, fielders, nil, ImplData{})
	return l
}

func (l *Buffered) Error(msg string, fielders ...Fielder) Logger {
	l.LogImpl(Error, msg, fielders, nil, ImplData{})
	return l
}

func (l *Buffered) Log(level Level, msg string, fielders ...Fielder) Logger {
	l.LogImpl(level, msg, fielders, nil, ImplData{})
	return l
}
