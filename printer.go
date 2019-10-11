package frog

import (
	"fmt"
	"strings"
	"time"

	"github.com/danbrakeley/frog/ansi"
)

type Printer interface {
	Render(useAnsi bool, useColor bool, level Level, format string, a ...interface{}) string
}

type TextPrinter struct {
	PrintTime  bool
	PrintLevel bool
}

func trimNewlines(str string) string {
	for strings.HasPrefix(str, "\n") {
		str = str[1:]
	}
	for strings.HasSuffix(str, "\n") {
		str = str[:len(str)-len("\n")]
	}
	return str
}

func (p *TextPrinter) Render(useAnsi bool, useColor bool, level Level, format string, a ...interface{}) string {
	format = trimNewlines(format)

	// and replace any newlines that remain
	format = strings.ReplaceAll(format, "\n", "¶")

	var out []string

	if useAnsi && useColor {
		var str string
		switch level {
		case Transient:
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
		case Transient:
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

	if useAnsi && useColor {
		out = append(out, ansi.Esc+ansi.Reset+"m")
	}

	return strings.Join(out, "")
}

type JSONPrinter struct {
	TimeOverride time.Time
}

func (p *JSONPrinter) Render(useAnsi bool, useColor bool, level Level, format string, a ...interface{}) string {
	var stamp time.Time
	if !p.TimeOverride.IsZero() {
		stamp = p.TimeOverride
	} else {
		stamp = time.Now()
	}

	return fmt.Sprintf(`{"timestamp":"%s","level":"%s","msg":"%s"}`,
		stamp.Format(time.RFC3339),
		level.String(),
		escapeStringForJSON(fmt.Sprintf(trimNewlines(format), a...)),
	)
}

func escapeStringForJSON(s string) string {
	sb := strings.Builder{}
	sb.Grow(len(s) * 6) // worst case
	for _, r := range s {
		switch r {
		case 0x00: // NUL
			sb.WriteString(`\u0000`)
		case 0x01: // SOH
			sb.WriteString(`\u0001`)
		case 0x02: // STX
			sb.WriteString(`\u0002`)
		case 0x03: // ETX
			sb.WriteString(`\u0003`)
		case 0x04: // EOT
			sb.WriteString(`\u0004`)
		case 0x05: // ENQ
			sb.WriteString(`\u0005`)
		case 0x06: // ACK
			sb.WriteString(`\u0006`)
		case 0x07: // BEL
			sb.WriteString(`\u0007`)
		case 0x08: // BS
			sb.WriteString(`\u0008`)
		case '\t': // TAB
			sb.WriteString(`\t`)
		case '\n': // LF
			sb.WriteString(`\n`)
		case 0x0b: // VT
			sb.WriteString(`\u000b`)
		case 0x0c: // FF
			sb.WriteString(`\u000c`)
		case '\r': // CR
			sb.WriteString(`\r`)
		case 0x0e: // SO
			sb.WriteString(`\u000e`)
		case 0x0f: // SI
			sb.WriteString(`\u000f`)
		case 0x10: // DLE
			sb.WriteString(`\u0010`)
		case 0x11: // DC1
			sb.WriteString(`\u0011`)
		case 0x12: // DC2
			sb.WriteString(`\u0012`)
		case 0x13: // DC3
			sb.WriteString(`\u0013`)
		case 0x14: // DC4
			sb.WriteString(`\u0014`)
		case 0x15: // NAK
			sb.WriteString(`\u0015`)
		case 0x16: // SYN
			sb.WriteString(`\u0016`)
		case 0x17: // ETB
			sb.WriteString(`\u0017`)
		case 0x18: // CAN
			sb.WriteString(`\u0018`)
		case 0x19: // EM
			sb.WriteString(`\u0019`)
		case 0x1a: // SUB
			sb.WriteString(`\u001a`)
		case 0x1b: // ESC
			sb.WriteString(`\u001b`)
		case 0x1c: // FS
			sb.WriteString(`\u001c`)
		case 0x1d: // GS
			sb.WriteString(`\u001d`)
		case 0x1e: // RS
			sb.WriteString(`\u001e`)
		case 0x1f: // US
			sb.WriteString(`\u001f`)
		case '"': // 0x22
			sb.WriteString(`\"`)
		case '&': // 0x26
			sb.WriteString(`\u0026`)
		case '<': // 0x3c
			sb.WriteString(`\u003c`)
		case '>': // 0x3e
			sb.WriteString(`\u003e`)
		case '\\': // 0x5c
			sb.WriteString(`\\`)
		case '\u2028': // Line separator (considered a line terminator in JS)
			sb.WriteString(`\u2028`)
		case '\u2029': // Paragraph separator (considered a line terminator in JS)
			sb.WriteString(`\u2029`)
		default:
			// the rest can represent themselves safely
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
