package frog

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/danbrakeley/ansi"
)

var hasEnvVarNoColor = false

func init() {
	_, exists := os.LookupEnv("NO_COLOR")
	hasEnvVarNoColor = exists
}

type Printer interface {
	Render(Level, []PrinterOption, string, []Field) string
	SetOptions(...PrinterOption) Printer
}

type TextPrinter struct {
	palette    ansicolors
	printTime  bool
	printLevel bool

	// fieldIndent controls where the first field begins rendering, compared to the message.
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
	fieldIndent int

	// printMessageLast will cause the message to display after the fields, instead of before
	// For example:
	//   [nfo] fieldIndent=40                          short message
	//   [nfo] fieldIndent=10      something  a long message that overflows the first field indent
	printMessageLast bool

	// transientLineLength is used to crop anchored lines, to avoid having them overflow into a second line.
	// A value of 0 means no cropping will occur. If an anchored line ends up being wider than the terminal,
	// it will wrap, which will throw off the formatting and scramble the output.
	transientLineLength int
}

func (p *TextPrinter) SetOptions(opts ...PrinterOption) Printer {
	for _, o := range opts {
		switch ot := o.(type) {
		case poPalette:
			p.palette = ot.ANSIColors
		case poTime:
			p.printTime = ot.Visible
		case poLevel:
			p.printLevel = ot.Visible
		case poFieldIndent:
			p.fieldIndent = ot.Indent
		case poMsgLeftFieldsRight:
			p.printMessageLast = false
		case poFieldsLeftMsgRight:
			p.printMessageLast = true
		case poTransientLineLength:
			p.transientLineLength = ot.Cols
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

func (p *TextPrinter) Render(level Level, opts []PrinterOption, msg string, fields []Field) string {
	// To override printer options just for this one line, make a copy of the existing printer settings,
	// then set our one-off options on that copy, then finally call Render on the altered copy.
	if len(opts) > 0 {
		tmp := *p
		tmp.SetOptions(opts...)
		return tmp.Render(level, nil, msg, fields)
	}

	var useColor bool
	var colorPrimary, colorSecondary string

	if !hasEnvVarNoColor {
		colorPrimary = p.palette[level][0]
		colorSecondary = p.palette[level][1]
		useColor = len(colorPrimary) > 0 && len(colorSecondary) > 0
	}

	msg = escapeMessageForTerminal(trimNewlines(msg))

	var sb strings.Builder
	sb.Grow(256)

	if useColor {
		sb.WriteString(colorSecondary)
	}

	if p.printTime {
		sb.WriteString(fmt.Sprintf("%s ", time.Now().Format("2006.01.02-15:04:05")))
	}

	if p.printLevel {
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
		for i, field := range fields {
			if i != 0 {
				sb.WriteByte(' ')
				count++
			}
			v := field.Value
			if field.IsJSONString {
				if !field.IsJSONSafe {
					v = escapeStringFieldForTerminal(v)
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
	if p.printMessageLast {
		visibleRuneCount = fnWriteFields()
		hasRightSide = len(msg) > 0
	} else {
		visibleRuneCount = fnWriteMsg()
		hasRightSide = len(fields) > 0
	}

	// write indentation
	if visibleRuneCount > 0 && hasRightSide {
		minLen := p.fieldIndent
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
		if p.printMessageLast {
			fnWriteMsg()
		} else {
			fnWriteFields()
		}
	}

	if useColor {
		sb.WriteString(ansi.CSI + ansi.Reset + "m")
	}

	out := sb.String()

	if level == Transient && p.transientLineLength > 0 {
		runeCount := len([]rune(out))
		if runeCount > p.transientLineLength {
			out = ansi.CropPreservingANSI(out, p.transientLineLength)
		}
	}

	return out
}

type JSONPrinter struct {
	TimeOverride time.Time // TODO: only tests use this currently, can we instead support POTime for tests?
}

func (p *JSONPrinter) SetOptions(opts ...PrinterOption) Printer {
	// JSONPrinter doesn't currently respect any printer options
	return p
}

func (p *JSONPrinter) Render(level Level, opts []PrinterOption, msg string, fields []Field) string {
	var stamp time.Time
	if !p.TimeOverride.IsZero() {
		stamp = p.TimeOverride
	} else {
		stamp = time.Now()
	}

	var sb strings.Builder
	sb.Grow(70 + len(msg) + len(fields)*50)

	sb.WriteString(`{"timestamp":"`)
	sb.WriteString(stamp.Format(time.RFC3339))
	sb.WriteString(`","level":"`)
	sb.WriteString(level.String())
	sb.WriteString(`","msg":"`)
	sb.WriteString(escapeStringForJSON(trimNewlines(msg)))
	sb.WriteString(`"`)

	for _, field := range fields {
		if field.IsJSONString {
			sb.WriteString(`,"`)
			sb.WriteString(field.Name)
			sb.WriteString(`":"`)
			if field.IsJSONSafe {
				sb.WriteString(field.Value)
			} else {
				sb.WriteString(escapeStringForJSON(field.Value))
			}
			sb.WriteString(`"`)
		} else {
			sb.WriteString(`,"`)
			sb.WriteString(field.Name)
			sb.WriteString(`":`)
			sb.WriteString(field.Value)
		}
	}

	sb.WriteString(`}`)
	return sb.String()
}

func escapeMessageForTerminal(s string) string {
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
		case '\\': // 0x5c
			sb.WriteString(`\\`)
		default:
			// the rest can represent themselves safely
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func escapeStringFieldForTerminal(s string) string {
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
