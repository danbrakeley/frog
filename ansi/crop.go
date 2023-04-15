package ansi

import "strings"

// CropPreservingANSI crops the unicode runes down to given length, but preserves all ANSI/VT100 escape sequences.
func CropPreservingANSI(str string, max int) string {
	inEscape := false
	visibleCount := 0

	var sb strings.Builder
	sb.Grow(len(str))

	for _, r := range str {
		if inEscape {
			sb.WriteRune(r)
			if r == '[' || (r >= '0' && r <= '9') || r == ';' || r == '?' {
				continue
			}
			inEscape = false
			continue
		}
		if r == EscRune {
			sb.WriteRune(r)
			inEscape = true
			continue
		}

		if visibleCount < max {
			sb.WriteRune(r)
			visibleCount++
			continue
		}
	}

	return sb.String()
}
