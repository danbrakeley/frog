package frog

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/danbrakeley/frog/ansi"
)

var isNoColorSet = false

func init() {
	_, exists := os.LookupEnv("NO_COLOR")
	isNoColorSet = exists
}

type Palette byte

const (
	PalNone Palette = iota
	PalColor
	PalDark
)

type Printer interface {
	Render(Level, string, ...Fielder) string
	SetOptions(...PrinterOption) Printer
}

type TextPrinter struct {
	Palette    Palette
	PrintTime  bool
	PrintLevel bool
	// FieldIndent controls where the first field begins rendering, compared to the message.
	// Note that the first field will always be at least 3 spaces from the end of the message,
	// and always be aligned with an offset that is a multiple of 5
	// For example:
	//   [nfo] a long message that overflows the first field indent   fieldIndent=40
	//   [nfo] short message                           fieldIndent=40
	//   [nfo] short message                 fieldIndent=30
	//   [nfo] short message       fieldIndent=20
	//   [nfo] short message       fieldIndent=10
	//   [nfo] short     fieldIndent=10
	//   [nfo] short     fieldIndent=0
	//   [nfo] sh   fieldIndent=0
	FieldIndent int
	// PrintMessageLast will cause the message to display after the fields, instead of before
	// For example:
	//   [nfo] fieldIndent=40                          short message
	//   [nfo] fieldIndent=10      something  a long message that overflows the first field indent
	PrintMessageLast bool
}

func (p *TextPrinter) SetOptions(opts ...PrinterOption) Printer {
	for _, opt := range opts {
		switch opt.Value() {
		case poPalette:
			p.Palette = Palette(opt.AsInt())
		case poShowTime:
			p.PrintTime = opt.AsBool()
		case poShowLevel:
			p.PrintLevel = opt.AsBool()
		case poFieldIndent:
			p.FieldIndent = opt.AsInt()
		case poMessageLast:
			p.PrintMessageLast = opt.AsBool()
		}
	}
	return p
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

var colorsMain = [levelMax][2]string{
	{ansi.FgDarkGreen, ansi.FgDarkGray}, // Transient
	{ansi.FgCyan, ansi.FgDarkCyan},      // Verbose
	{ansi.FgWhite, ansi.FgLightGray},    // Info
	{ansi.FgYellow, ansi.FgDarkYellow},  // Warning
	{ansi.FgRed, ansi.FgDarkRed},        // Error
	{ansi.FgRed, ansi.FgDarkRed},        // Fatal
}

var colorsDark = [levelMax][2]string{
	{ansi.FgDarkGray, ansi.FgDarkGray}, // Transient
	{ansi.FgDarkGray, ansi.FgDarkGray}, // Verbose
	{ansi.FgDarkGray, ansi.FgDarkGray}, // Info
	{ansi.FgDarkGray, ansi.FgDarkGray}, // Warning
	{ansi.FgDarkGray, ansi.FgDarkGray}, // Error
	{ansi.FgDarkGray, ansi.FgDarkGray}, // Fatal
}

func (p *TextPrinter) Render(level Level, msg string, fields ...Fielder) string {
	var useColor bool
	var colorPrimary, colorSecondary string

	if !isNoColorSet {
		switch p.Palette {
		case PalNone:
			useColor = false
		case PalColor:
			useColor = true
			colorPrimary = ansi.CSI + colorsMain[level][0] + "m"
			colorSecondary = ansi.CSI + colorsMain[level][1] + "m"
		case PalDark:
			useColor = true
			colorPrimary = ansi.CSI + colorsDark[level][0] + "m"
			colorSecondary = ansi.CSI + colorsDark[level][1] + "m"
		}
	}

	msg = escapeStringForTerminal(trimNewlines(msg))

	var sb strings.Builder
	sb.Grow(256)

	if useColor {
		sb.WriteString(colorSecondary)
	}

	if p.PrintTime {
		sb.WriteString(fmt.Sprintf("%s ", time.Now().Format("2006.01.02-15:04:05")))
	}

	if p.PrintLevel {
		switch level {
		case Transient:
			sb.WriteString("[==>] ")
		case Verbose:
			sb.WriteString("[dbg] ")
		case Info:
			sb.WriteString("[nfo] ")
		case Warning:
			sb.WriteString("[WRN] ")
		case Error:
			sb.WriteString("[ERR] ")
		case Fatal:
			sb.WriteString("[!!!] ")
		default:
			sb.WriteString("[???] ")
		}
	}

	fnWriteMsg := func() int {
		if useColor {
			sb.WriteString(colorPrimary)
		}
		sb.WriteString(msg)
		return utf8.RuneCountInString(msg)
	}

	fnWriteFields := func() int {
		count := 0
		for i, f := range fields {
			if i != 0 {
				sb.WriteByte(' ')
				count++
			}
			field := f.Field()
			v := field.Value
			if field.IsJSONString {
				if !field.IsJSONSafe {
					v = escapeStringForTerminal(v)
				}
				if len(v) == 0 || strings.ContainsAny(v, " \\") {
					v = "\"" + v + "\""
				}
			}

			if useColor {
				sb.WriteString(colorSecondary)
			}
			sb.WriteString(field.Name)
			count += utf8.RuneCountInString(field.Name)
			sb.WriteByte('=')
			count += 1
			if useColor {
				sb.WriteString(colorPrimary)
			}
			sb.WriteString(v)
			count += utf8.RuneCountInString(v)
		}
		return count
	}

	// write left side
	var visibleRuneCount int
	var hasRightSide bool
	if p.PrintMessageLast {
		visibleRuneCount = fnWriteFields()
		hasRightSide = len(msg) > 0
	} else {
		visibleRuneCount = fnWriteMsg()
		hasRightSide = len(fields) > 0
	}

	// write indentation
	if visibleRuneCount > 0 && hasRightSide {
		minLen := p.FieldIndent
		const minSpace = 3
		const tabWidth = 5

		space := minSpace
		if visibleRuneCount+space < minLen {
			space = minLen - visibleRuneCount
		}

		offset := (((visibleRuneCount + space - 1) / tabWidth) + 1) * tabWidth
		for i := 0; i < offset-visibleRuneCount; i++ {
			sb.WriteByte(' ')
		}
	}

	if hasRightSide {
		// write right side
		if p.PrintMessageLast {
			fnWriteMsg()
		} else {
			fnWriteFields()
		}
	}

	if useColor {
		sb.WriteString(ansi.CSI + ansi.Reset + "m")
	}

	return sb.String()
}

type JSONPrinter struct {
	TimeOverride time.Time
}

func (p *JSONPrinter) SetOptions(opts ...PrinterOption) Printer {
	// JSONPrinter doesn't currently respect any printer options
	return p
}

func (p *JSONPrinter) Render(level Level, msg string, fields ...Fielder) string {
	var stamp time.Time
	if !p.TimeOverride.IsZero() {
		stamp = p.TimeOverride
	} else {
		stamp = time.Now()
	}

	out := fmt.Sprintf(`{"timestamp":"%s","level":"%s","msg":"%s"`,
		stamp.Format(time.RFC3339),
		level.String(),
		escapeStringForJSON(trimNewlines(msg)),
	)

	for _, f := range fields {
		field := f.Field()
		if field.IsJSONString {
			if field.IsJSONSafe {
				out += fmt.Sprintf(`,"%s":"%s"`, field.Name, field.Value)
			} else {
				out += fmt.Sprintf(`,"%s":"%s"`, field.Name, escapeStringForJSON(field.Value))
			}
		} else {
			out += fmt.Sprintf(`,"%s":%s`, field.Name, field.Value)
		}
	}

	out += "}"
	return out
}

func escapeStringForTerminal(s string) string {
	sb := strings.Builder{}
	sb.Grow(len(s) * 2) // worst case
	for _, r := range s {
		switch r {
		case '\t': // TAB
			sb.WriteString(`\t`)
		case '\n': // LF
			sb.WriteString(`\n`)
		case '\r': // CR
			sb.WriteString(`\r`)
		case '"': // 0x22
			sb.WriteString(`\"`)
		case '\\': // 0x5c
			sb.WriteString(`\\`)
		default:
			// the rest can represent themselves safely
			sb.WriteRune(r)
		}
	}
	return sb.String()
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
