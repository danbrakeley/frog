package main

import (
	"fmt"
	"strings"

	"github.com/danbrakeley/frog/ansi"
)

func main() {
	fmt.Println(ansi.CSI + ansi.FgBlack + "m" + "FgBlack" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgDarkGray + "m" + "FgDarkGray" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgLightGray + "m" + "FgLightGray" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgWhite + "m" + "FgWhite" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgDarkRed + "m" + "FgDarkRed" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgRed + "m" + "FgRed" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgDarkGreen + "m" + "FgDarkGreen" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgGreen + "m" + "FgGreen" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgDarkYellow + "m" + "FgDarkYellow" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgYellow + "m" + "FgYellow" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgDarkBlue + "m" + "FgDarkBlue" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgBlue + "m" + "FgBlue" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgDarkMagenta + "m" + "FgDarkMagenta" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgMagenta + "m" + "FgMagenta" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgDarkCyan + "m" + "FgDarkCyan" + ansi.CSI + ansi.Reset + "m")
	fmt.Println(ansi.CSI + ansi.FgCyan + "m" + "FgCyan" + ansi.CSI + ansi.Reset + "m")
	fmt.Println("")

	fmt.Println(strings.Join([]string{
		ansi.CSI + ansi.FgDarkGray + "m",
		"2020.06.28-18:25:17 [==>] ",
		ansi.CSI + ansi.FgDarkGreen + "m",
		"test message        ",
		ansi.CSI + ansi.FgDarkGray + "m",
		"field=",
		ansi.CSI + ansi.FgDarkGreen + "m",
		"value",
		ansi.CSI + ansi.Reset + "m",
	}, ""))

	fmt.Println(strings.Join([]string{
		ansi.CSI + ansi.FgDarkCyan + "m",
		"2020.06.28-18:25:17 [dbg] ",
		ansi.CSI + ansi.FgCyan + "m",
		"test message        ",
		ansi.CSI + ansi.FgDarkCyan + "m",
		"field=",
		ansi.CSI + ansi.FgCyan + "m",
		"value",
		ansi.CSI + ansi.Reset + "m",
	}, ""))

	fmt.Println(strings.Join([]string{
		ansi.CSI + ansi.FgLightGray + "m",
		"2020.06.28-18:25:17 [nfo] ",
		ansi.CSI + ansi.FgWhite + "m",
		"test message        ",
		ansi.CSI + ansi.FgLightGray + "m",
		"field=",
		ansi.CSI + ansi.FgWhite + "m",
		"value",
		ansi.CSI + ansi.Reset + "m",
	}, ""))

	fmt.Println(strings.Join([]string{
		ansi.CSI + ansi.FgDarkYellow + "m",
		"2020.06.28-18:25:17 [WRN] ",
		ansi.CSI + ansi.FgYellow + "m",
		"test message        ",
		ansi.CSI + ansi.FgDarkYellow + "m",
		"field=",
		ansi.CSI + ansi.FgYellow + "m",
		"value",
		ansi.CSI + ansi.Reset + "m",
	}, ""))

	fmt.Println(strings.Join([]string{
		ansi.CSI + ansi.FgDarkRed + "m",
		"2020.06.28-18:25:17 [ERR] ",
		ansi.CSI + ansi.FgRed + "m",
		"test message        ",
		ansi.CSI + ansi.FgDarkRed + "m",
		"field=",
		ansi.CSI + ansi.FgRed + "m",
		"value",
		ansi.CSI + ansi.Reset + "m",
	}, ""))
}
