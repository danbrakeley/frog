package main

import (
	"fmt"
	"strings"

	"github.com/danbrakeley/ansi"
)

func main() {
	fmt.Println(ansi.FgBlack + "FgBlack" + ansi.Reset)
	fmt.Println(ansi.FgDarkGray + "FgDarkGray" + ansi.Reset)
	fmt.Println(ansi.FgLightGray + "FgLightGray" + ansi.Reset)
	fmt.Println(ansi.FgWhite + "FgWhite" + ansi.Reset)
	fmt.Println(ansi.FgDarkRed + "FgDarkRed" + ansi.Reset)
	fmt.Println(ansi.FgRed + "FgRed" + ansi.Reset)
	fmt.Println(ansi.FgDarkGreen + "FgDarkGreen" + ansi.Reset)
	fmt.Println(ansi.FgGreen + "FgGreen" + ansi.Reset)
	fmt.Println(ansi.FgDarkYellow + "FgDarkYellow" + ansi.Reset)
	fmt.Println(ansi.FgYellow + "FgYellow" + ansi.Reset)
	fmt.Println(ansi.FgDarkBlue + "FgDarkBlue" + ansi.Reset)
	fmt.Println(ansi.FgBlue + "FgBlue" + ansi.Reset)
	fmt.Println(ansi.FgDarkMagenta + "FgDarkMagenta" + ansi.Reset)
	fmt.Println(ansi.FgMagenta + "FgMagenta" + ansi.Reset)
	fmt.Println(ansi.FgDarkCyan + "FgDarkCyan" + ansi.Reset)
	fmt.Println(ansi.FgCyan + "FgCyan" + ansi.Reset)
	fmt.Println("")

	fmt.Println(strings.Join([]string{
		ansi.FgDarkGray,
		"2020.06.28-18:25:17 [==>] ",
		ansi.FgDarkGreen,
		"test message        ",
		ansi.FgDarkGray,
		"field=",
		ansi.FgDarkGreen,
		"value",
		ansi.Reset,
	}, ""))

	fmt.Println(strings.Join([]string{
		ansi.FgDarkCyan,
		"2020.06.28-18:25:17 [dbg] ",
		ansi.FgCyan,
		"test message        ",
		ansi.FgDarkCyan,
		"field=",
		ansi.FgCyan,
		"value",
		ansi.Reset,
	}, ""))

	fmt.Println(strings.Join([]string{
		ansi.FgLightGray,
		"2020.06.28-18:25:17 [nfo] ",
		ansi.FgWhite,
		"test message        ",
		ansi.FgLightGray,
		"field=",
		ansi.FgWhite,
		"value",
		ansi.Reset,
	}, ""))

	fmt.Println(strings.Join([]string{
		ansi.FgDarkYellow,
		"2020.06.28-18:25:17 [WRN] ",
		ansi.FgYellow,
		"test message        ",
		ansi.FgDarkYellow,
		"field=",
		ansi.FgYellow,
		"value",
		ansi.Reset,
	}, ""))

	fmt.Println(strings.Join([]string{
		ansi.FgDarkRed,
		"2020.06.28-18:25:17 [ERR] ",
		ansi.FgRed,
		"test message        ",
		ansi.FgDarkRed,
		"field=",
		ansi.FgRed,
		"value",
		ansi.Reset,
	}, ""))
}
