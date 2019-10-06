package frog

import (
	"fmt"
	"strings"
	"time"

	"github.com/danbrakeley/frog/ansi"
)

type Printer struct {
	CanUseAnsi bool
	PrintTime  bool
	PrintLevel bool
}

func (p *Printer) Sprintf(level Level, format string, a ...interface{}) string {
	// trim newlines from the ends
	for strings.HasPrefix(format, "\n") {
		format = format[1:]
	}
	for strings.HasSuffix(format, "\n") {
		format = format[:len(format)-len("\n")]
	}

	// and replace any newlines that remain
	format = strings.ReplaceAll(format, "\n", "¶")

	var out []string

	if p.CanUseAnsi {
		var str string
		switch level {
		case Progress:
			str = ansi.Esc + ansi.FgDarkGray + "m"
		case Verbose:
			str = ansi.Esc + ansi.FgDarkCyan + "m"
		case Info:
			str = ansi.Esc + ansi.FgLightGray + "m"
		case Warning:
			str = ansi.Esc + ansi.FgYellow + "m"
		case Error, Fatal:
			str = ansi.Esc + ansi.FgRed + "m"
		default:
			str = ansi.Esc + ansi.BgWhite + ";" + ansi.FgBlack + "m"
		}
		out = append(out, str)
	}

	if p.PrintTime {
		out = append(out, fmt.Sprintf("%s ", time.Now().Format("2006.01.02-15:04:05")))
	}

	if p.PrintLevel {
		var str string
		switch level {
		case Progress:
			str = "[==>] "
		case Verbose:
			str = "[dbg] "
		case Info:
			str = "[nfo] "
		case Warning:
			str = "[WRN] "
		case Error:
			str = "[ERR] "
		case Fatal:
			str = "[!!!] "
		default:
			str = "[???] "
		}
		out = append(out, str)
	}

	out = append(out, fmt.Sprintf(format, a...))

	if p.CanUseAnsi {
		out = append(out, ansi.Esc+ansi.Reset+"m")
	}

	return strings.Join(out, "")
}
