package frog

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

type UseAnsiType bool
type UseColorType bool

const (
	NoAnsi  UseAnsiType = false
	UseAnsi UseAnsiType = true

	NoColor  UseColorType = false
	UseColor UseColorType = true
)

func AssertGolden(t *testing.T, testName string, actual []byte) {
	t.Helper()
	golden := filepath.Join("test-fixtures", testName+".golden")
	if *update {
		ioutil.WriteFile(golden, actual, 0644)
	}
	expected, _ := ioutil.ReadFile(golden)
	if !bytes.Equal(actual, expected) {
		t.Errorf(
			"golden file %s does not match output:\nGolden File:\n%s\nActual:\n%s",
			golden, string(expected), string(actual),
		)
	}
}

func Test_Golden(t *testing.T) {
	cases := []struct {
		Name     string
		DoWork   func(Logger)
		UseAnsi  UseAnsiType
		UseColor UseColorType
		Printer  Printer
	}{
		{
			"min-level", minLevel, NoAnsi, UseColor,
			&TextPrinter{PrintTime: false, PrintLevel: true},
		},
		{
			"min-level-ansi", minLevel, UseAnsi, UseColor,
			&TextPrinter{PrintTime: false, PrintLevel: true},
		},
		{
			"trims-newlines", newlineVariations, NoAnsi, UseColor,
			&TextPrinter{PrintTime: false, PrintLevel: true},
		},
		{
			"fixed-lines-ansi", fixedLines, UseAnsi, UseColor,
			&TextPrinter{PrintTime: false, PrintLevel: true},
		},
		{
			"fixed-lines-ansi-no-color", fixedLines, UseAnsi, NoColor,
			&TextPrinter{PrintTime: false, PrintLevel: true},
		},
		{
			"fixed-lines-no-ansi", fixedLines, NoAnsi, UseColor,
			&TextPrinter{PrintTime: false, PrintLevel: true},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// run against a Buffered Logger
			{
				buf := &bytes.Buffer{}
				cfg := Config{
					Writer:   buf,
					UseAnsi:  tc.UseAnsi == UseAnsi,
					UseColor: tc.UseColor == UseColor,
				}
				l := NewBuffered(cfg, tc.Printer)
				tc.DoWork(l)
				l.Close()
				AssertGolden(t, tc.Name+".buffered", buf.Bytes())
			}

			// run against an Unbuffered Logger
			{
				buf := &bytes.Buffer{}
				l := NewUnbuffered(buf, tc.Printer)
				tc.DoWork(l)
				l.Close()
				AssertGolden(t, tc.Name+".unbuffered", buf.Bytes())
			}

			// run against a TeeLogger, with a Buffered as Primary and Unbuffered as Secondary
			{
				buf1 := &bytes.Buffer{}
				cfg := Config{
					Writer:   buf1,
					UseAnsi:  tc.UseAnsi == UseAnsi,
					UseColor: tc.UseColor == UseColor,
				}
				bl := NewBuffered(cfg, tc.Printer)

				buf2 := &bytes.Buffer{}
				ul := NewUnbuffered(buf2, tc.Printer)

				tee := &TeeLogger{Primary: bl, Secondary: ul}
				tc.DoWork(tee)
				tee.Close()
				AssertGolden(t, tc.Name+".buffered", buf1.Bytes())
				AssertGolden(t, tc.Name+".unbuffered", buf2.Bytes())
			}
		})
	}
}

func minLevel(l Logger) {
	for _, level := range []Level{Transient, Verbose, Info, Warning, Error} {
		l.SetMinLevel(level)
		l.Transientf("this is a transient line")
		l.Verbosef("this is a verbose line")
		l.Infof("this is an info line")
		l.Warningf("this is a warning line")
		l.Errorf("this is an error line")
	}
}

func newlineVariations(l Logger) {
	l.Infof("most of these lines will end up the same")
	l.Infof("\nmost of these lines will end up the same")
	l.Infof("\n\nmost of these lines will end up the same")
	l.Infof("\n\n\nmost of these lines will end up the same")
	l.Infof("most of these lines will end up the same\n")
	l.Infof("most of these lines will end up the same\n\n")
	l.Infof("most of these lines will end up the same\n\n\n")
	l.Infof("\nmost of these lines will end up the same\n")
	l.Infof("except\nthese last couple of lines, which have newline breaks")
	l.Infof("\nexcept these\nlast couple of lines, which\nhave newline breaks")
	l.Infof("\n\nexcept these last\ncouple of lines,\n\nwhich have newline breaks")
	l.Infof("except these last couple\n\nof lines, which have newline breaks\n")
	l.Infof("except these last couple of lines, which\nhave\nnewline breaks\n\n")
	l.Infof("\nexcept these last couple\nof lines, which have\n\n\nnewline breaks\n")
}

func fixedLines(l Logger) {
	f := []Logger{
		AddFixedLine(l),
		AddFixedLine(l),
		AddFixedLine(l),
	}
	f[0].Transientf("write to first line")
	f[1].Transientf("write to second line")
	f[2].Transientf("write to third line")
	f[0].Transientf("write back to first line")
	f[2].Transientf("now we're on the third line again")
	f[0].Warningf("something unexpected happened on the first line")
	f[1].Transientf("done")
	f[2].Transientf("done")
	f[0].Transientf("done")
}
