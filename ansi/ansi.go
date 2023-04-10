package ansi

import (
	"fmt"
	"strings"
)

// Useful references:
// https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
// https://en.wikipedia.org/wiki/ANSI_escape_code
// http://www.termsys.demon.co.uk/vtansi.htm

const (
	EscRune = '\u001b'
	Esc     = string(EscRune)
	CSI     = Esc + "[" // Control Sequence Introducer

	// Move cursor to screen-relative locations

	TopLeft     = CSI + "H"
	LeftMost    = CSI + "1G" // stays on current line
	BottomRight = CSI + "32767;32767H"

	// Save/Restore cursor position

	PosSave    = CSI + "s" // saves the current cursor position
	PosRestore = CSI + "u" // restores cursor position to the last save

	// Erase commands do not move cursor, escept where noted

	EraseEOL    = CSI + "K"  // erase to end of current line
	EraseSOL    = CSI + "1K" // erase to start of current line
	EraseLine   = CSI + "2K" // erase from start to end of current line
	EraseDown   = CSI + "J"  // erase to end of current line, then everything down to bottom of screen
	EraseUp     = CSI + "1J" // erase to start of current line, then everything up to top of screen
	EraseScreen = CSI + "2J" // erases everything to background colour, then moves cursor home

	// DECTCEM commands

	ShowCursor = CSI + "?25h"
	HideCursor = CSI + "?25l"

	// DECXCPR commands

	GetCursorPos = CSI + "6n" // responds with `ESC [ <r> ; <c> R`, where <r> is row and <c> is column
)

func Up(n int) string {
	return fmt.Sprintf(CSI+"%dA", n)
}

func Down(n int) string {
	return fmt.Sprintf(CSI+"%dB", n)
}

func Right(n int) string {
	return fmt.Sprintf(CSI+"%dC", n)
}

func Left(n int) string {
	return fmt.Sprintf(CSI+"%dD", n)
}

// NextLine moves cursor to start of line n below current
// (stops at bottom of viewable area, does not cause a scroll)
func NextLine(n int) string {
	return fmt.Sprintf(CSI+"%dE", n)
}

// PrevLine moves cursor to start of line n above current
// (stops at top of viewable area)
func PrevLine(n int) string {
	return fmt.Sprintf(CSI+"%dF", n)
}

const (
	Reset         = "0"
	Bold          = "1"
	Underline     = "4"
	Reverse       = "7"
	UnderlineOff  = "24"
	ReverseOff    = "27"
	FgBlack       = "30"
	FgDarkRed     = "31"
	FgDarkGreen   = "32"
	FgDarkYellow  = "33"
	FgDarkBlue    = "34"
	FgDarkMagenta = "35"
	FgDarkCyan    = "36"
	FgLightGray   = "37"
	FgReset       = "39" // sets foreground color to default
	BgBlack       = "40"
	BgDarkRed     = "41"
	BgDarkGreen   = "42"
	BgDarkYellow  = "43"
	BgDarkBlue    = "44"
	BgDarkMagenta = "45"
	BgDarkCyan    = "46"
	BgLightGray   = "47"
	BgReset       = "49" // sets background color to default
	FgDarkGray    = "90"
	FgRed         = "91"
	FgGreen       = "92"
	FgYellow      = "93"
	FgBlue        = "94"
	FgMagenta     = "95"
	FgCyan        = "96"
	FgWhite       = "97"
	BgDarkGray    = "100"
	BgRed         = "101"
	BgGreen       = "102"
	BgYellow      = "103"
	BgBlue        = "104"
	BgMagenta     = "105"
	BgCyan        = "106"
	BgWhite       = "107"
)

// SGR applies the above sgr params in the order specified (later commands may override earlier commands)
func SGR(params ...string) string {
	var sb strings.Builder
	sb.Grow(len(params)*4 + 2) // worst case: n 3-char params, n-1 semicolons, <esc>, '[', and 'm'
	sb.WriteString(CSI)
	for i, v := range params {
		if i != 0 {
			sb.WriteRune(';')
		}
		sb.WriteString(string(v))
	}
	sb.WriteRune('m')
	return sb.String()
}
