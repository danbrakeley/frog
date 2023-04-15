package ansi

import (
	"strings"
	"testing"
)

func Test_CropPreservingANSI(t *testing.T) {
	cases := []struct {
		Name     string
		In       string
		Max      int
		Expected string
	}{
		{"no ansi", "Hello, World!", 10, "Hello, Wor"},
		{"unicode", "༼ つ ◕_◕ ༽つ", 7, "༼ つ ◕_◕"},
		{
			"colors",
			SGR(FgRed, BgDarkGreen) + "Hello, World!" + SGR(Reset),
			10,
			SGR(FgRed, BgDarkGreen) + "Hello, Wor" + SGR(Reset)},
		{
			"everything",
			TopLeft + LeftMost + BottomRight + PosSave + PosRestore + EraseEOL + EraseSOL + EraseLine +
				"Some legit text" +
				EraseDown + EraseUp + EraseScreen + ShowCursor + HideCursor + GetCursorPos +
				"©🦀💨✔" +
				Up(100) + Down(2) + Left(9999) + Right(51) + NextLine(10) + PrevLine(800) +
				SGR(FgDarkRed, BgWhite, Bold, Underline, Reverse) + "nice" + SGR(Reset),
			17,
			TopLeft + LeftMost + BottomRight + PosSave + PosRestore + EraseEOL + EraseSOL + EraseLine +
				"Some legit text" +
				EraseDown + EraseUp + EraseScreen + ShowCursor + HideCursor + GetCursorPos +
				"©🦀" +
				Up(100) + Down(2) + Left(9999) + Right(51) + NextLine(10) + PrevLine(800) +
				SGR(FgDarkRed, BgWhite, Bold, Underline, Reverse) + SGR(Reset),
		},
		// TODO: The following case is not handled properly because it uses diacritics that combine
		// multiple code points into a single grapheme.
		// Maybe [rivo/uniseg](https://github.com/rivo/uniseg) can help? Should get benchmarking in
		// first to understand performance overhead?
		// If performance is a concern, maybe add a printer option so user can choose faster if desired?
		// {"diacritics", "deͤmͫͫoͦͦͦ", 3, "deͤmͫͫ"},
	}

	escape := func(str string) string {
		var sb strings.Builder
		sb.Grow(len(str))
		for _, r := range str {
			if r == EscRune {
				sb.WriteRune('→')
				continue
			}
			sb.WriteRune(r)
		}
		return sb.String()
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := CropPreservingANSI(tc.In, tc.Max)

			if actual != tc.Expected {
				t.Errorf("\nexpected: %s\n  actual: %s", escape(tc.Expected), escape(actual))
			}
		})
	}
}
