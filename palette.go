package frog

import "github.com/danbrakeley/frog/ansi"

// Color is the public interface to the underlying ANSI colors.
type Color byte

const (
	Black Color = iota
	DarkRed
	DarkGreen
	DarkYellow
	DarkBlue
	DarkMagenta
	DarkCyan
	LightGray
	DarkGray
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Palette is the primary (0) and secondary (1) colors for each log level
type Palette [levelMax][2]Color

var DefaultPalette = Palette{
	{DarkGreen, DarkGray}, // Transient
	{Cyan, DarkCyan},      // Verbose
	{White, LightGray},    // Info
	{Yellow, DarkYellow},  // Warning
	{Red, DarkRed},        // Error
}

var DarkPalette = Palette{
	{DarkGray, DarkGray}, // Transient
	{DarkGray, DarkGray}, // Verbose
	{DarkGray, DarkGray}, // Info
	{DarkGray, DarkGray}, // Warning
	{DarkGray, DarkGray}, // Error
}

// internally, a palette is an ANSI escape sequence (string) for each pair of colors at each log level

type ansicolors [levelMax][2]string

func (p *Palette) toANSI() ansicolors {
	var out ansicolors
	for i := levelMin; i < levelMax; i++ {
		out[i][0] = ansiFgColor(p[i][0])
		out[i][1] = ansiFgColor(p[i][1])
	}
	return out
}

func ansiFgColor(c Color) string {
	switch c {
	case Black:
		return ansi.SGR(ansi.FgBlack)
	case DarkRed:
		return ansi.SGR(ansi.FgDarkRed)
	case DarkGreen:
		return ansi.SGR(ansi.FgDarkGreen)
	case DarkYellow:
		return ansi.SGR(ansi.FgDarkYellow)
	case DarkBlue:
		return ansi.SGR(ansi.FgDarkBlue)
	case DarkMagenta:
		return ansi.SGR(ansi.FgDarkMagenta)
	case DarkCyan:
		return ansi.SGR(ansi.FgDarkCyan)
	case LightGray:
		return ansi.SGR(ansi.FgLightGray)
	case DarkGray:
		return ansi.SGR(ansi.FgDarkGray)
	case Red:
		return ansi.SGR(ansi.FgRed)
	case Green:
		return ansi.SGR(ansi.FgGreen)
	case Yellow:
		return ansi.SGR(ansi.FgYellow)
	case Blue:
		return ansi.SGR(ansi.FgBlue)
	case Magenta:
		return ansi.SGR(ansi.FgMagenta)
	case Cyan:
		return ansi.SGR(ansi.FgCyan)
	case White:
		return ansi.SGR(ansi.FgWhite)
	}
	return ""
}
