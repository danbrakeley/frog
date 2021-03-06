package frog

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/danbrakeley/frog/ansi"
)

type Buffered struct {
	minLevel Level
	cfg      Config
	prn      Printer
	ch       chan bufmsg
	wg       sync.WaitGroup

	isClosed    int32 // to keep thread safe, use atomic reads/writes/math
	openAnchors int32 // to keep thread safe, use atomic reads/writes/math
}

type msgType byte

const (
	mtPrint      msgType = iota // string to print
	mtFatal                     // stop processor without closing channel
	mtAddLine                   // add an anchored line
	mtRemoveLine                // remove an acnhored line
)

type bufmsg struct {
	Type  msgType
	Line  int32
	Level Level
	Msg   string
}

func NewBuffered(cfg Config, prn Printer) *Buffered {
	l := &Buffered{
		minLevel: Info,
		cfg:      cfg,
		prn:      prn,
		ch:       make(chan bufmsg),
		wg:       sync.WaitGroup{},
	}

	l.wg.Add(1)
	if l.cfg.Writer == nil {
		return nil
	}
	go func() {
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
		return
	}

	if l.minLevel <= Verbose {
		_, file, line, ok := runtime.Caller(1)
		l.Verbose("buffered log closing", String("file", file), Int("line", line), Bool("ok", ok))
	}
	close(l.ch)
	l.wg.Wait()
}

// AddAnchor creates a Logger that is "achored" to the bottom of the output.
// This "anchoring" is achieved by using ANSI to re-draw the anchored line at
// the bottom as the output scrolls up.
// Thread safe.
func (l *Buffered) AddAnchor() Logger {
	lineNum := atomic.AddInt32(&l.openAnchors, 1)
	l.ch <- bufmsg{Type: mtAddLine, Line: lineNum}
	onClose := func() {
		atomic.AddInt32(&l.openAnchors, -1)
		l.ch <- bufmsg{Type: mtRemoveLine, Line: lineNum}
	}
	return newAnchor(l, lineNum, onClose)
}

func (l *Buffered) processor() {
	type anchoredLine struct {
		lineNum int32
		str     string
	}
	var anchoredLines []anchoredLine

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
		if !ok || msg.Type == mtFatal {
			return
		}

		switch msg.Type {
		case mtAddLine:
			// ensure terminal scrolls down if needed to add a new line
			fmt.Fprint(l.cfg.Writer, "\n")

			// figure out where this new line goes
			idx := 0
			for _, v := range anchoredLines {
				if msg.Line < v.lineNum {
					break
				}
				idx++
			}

			// and then insert it there
			anchoredLines = append(anchoredLines, anchoredLine{})
			copy(anchoredLines[idx+1:], anchoredLines[idx:])
			anchoredLines[idx] = anchoredLine{lineNum: msg.Line}

			// if inserting has bumped any lines down, re-draw those lines in their new home
			if idx < len(anchoredLines)-1 {
				fmt.Fprint(l.cfg.Writer, ansi.PrevLine(len(anchoredLines)-(idx+1)))
				for i := idx + 1; i < len(anchoredLines); i++ {
					fmt.Fprint(l.cfg.Writer, anchoredLines[i].str)
					fmt.Fprint(l.cfg.Writer, ansi.EraseEOL)
					fmt.Fprint(l.cfg.Writer, ansi.NextLine(1))
				}
			}

		case mtRemoveLine:
			// find the line we are removing
			idx := fnMustFindIdx(msg.Line)

			// remove element
			copy(anchoredLines[idx:], anchoredLines[idx+1:])
			anchoredLines = anchoredLines[:len(anchoredLines)-1]

			// redraw/erase bottom lines as needed
			fmt.Fprint(l.cfg.Writer, ansi.PrevLine(1+len(anchoredLines)-idx))
			for i := idx; i < len(anchoredLines); i++ {
				fmt.Fprint(l.cfg.Writer, anchoredLines[i].str)
				fmt.Fprint(l.cfg.Writer, ansi.EraseEOL)
				fmt.Fprint(l.cfg.Writer, ansi.NextLine(1))
			}
			fmt.Fprint(l.cfg.Writer, ansi.EraseEOL)

		case mtPrint:
			// if we aren't using anchored lines, then just print normally
			if len(anchoredLines) == 0 {
				fmt.Fprintf(l.cfg.Writer, "%s\n", msg.Msg)
				continue
			}

			// If we are using anchored lines, but this msg doesn't have one specified, then move all
			// the anchored lines down, and draw this line above them.
			// If this does have an anchored line, but it is not Transient level, then also print it above.
			if msg.Line <= 0 || msg.Level > Transient {
				fmt.Fprint(l.cfg.Writer, "\n")
				fmt.Fprint(l.cfg.Writer, ansi.PrevLine(1+len(anchoredLines)))
				fmt.Fprintf(l.cfg.Writer, "%s%s\n", msg.Msg, ansi.EraseEOL)

				for _, v := range anchoredLines {
					fmt.Fprintf(l.cfg.Writer, "%s%s\n", v.str, ansi.EraseEOL)
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
			fmt.Fprint(l.cfg.Writer, ansi.PrevLine(offset))
			fmt.Fprint(l.cfg.Writer, msg.Msg)
			fmt.Fprint(l.cfg.Writer, ansi.EraseEOL)
			fmt.Fprint(l.cfg.Writer, ansi.NextLine(offset))

		default:
		}
	}
}

func (l *Buffered) SetMinLevel(level Level) Logger {
	if level < levelMin || level >= levelMax {
		panic(fmt.Errorf("level %v is not in valid range [%v,%v)", level, levelMin, levelMax))
	}
	l.minLevel = level
	return l
}

func (l *Buffered) Log(level Level, format string, fields ...Fielder) Logger {
	if level < l.minLevel {
		return l
	}
	l.logImpl(l.prn, 0, level, format, fields...)
	return l
}

// logImpl is only called if the line will be shown, regardless of level, line, etc
func (l *Buffered) logImpl(prn Printer, anchoredLine int32, level Level, format string, fields ...Fielder) {
	l.ch <- bufmsg{
		Line:  anchoredLine,
		Level: level,
		Msg:   prn.Render(l.cfg.UseColor, level, format, fields...),
	}
	if level == Fatal {
		// we can't just close this channel, because another thread may still be trying to write to it
		l.ch <- bufmsg{Type: mtFatal}
		l.wg.Wait()
		os.Exit(-1)
	}
}

func (l *Buffered) Transient(format string, a ...Fielder) Logger {
	l.Log(Transient, format, a...)
	return l
}

func (l *Buffered) Verbose(format string, a ...Fielder) Logger {
	l.Log(Verbose, format, a...)
	return l
}

func (l *Buffered) Info(format string, a ...Fielder) Logger {
	l.Log(Info, format, a...)
	return l
}

func (l *Buffered) Warning(format string, a ...Fielder) Logger {
	l.Log(Warning, format, a...)
	return l
}

func (l *Buffered) Error(format string, a ...Fielder) Logger {
	l.Log(Error, format, a...)
	return l
}

func (l *Buffered) Fatal(format string, a ...Fielder) {
	l.Log(Fatal, format, a...)
}
