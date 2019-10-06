package frog

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func Test_Golden(t *testing.T) {
	cases := []struct {
		Name    string
		DoWork  func(Logger)
		Printer Printer
	}{
		{
			"min-level", minLevel,
			Printer{CanUseAnsi: false, PrintTime: false, PrintLevel: true},
		},
		{
			"min-level-ansi", minLevel,
			Printer{CanUseAnsi: true, PrintTime: false, PrintLevel: true},
		},
		{
			"trims-newlines", newlineVariations,
			Printer{CanUseAnsi: false, PrintTime: false, PrintLevel: true},
		},
		{
			"fixed-lines-ansi", fixedLines,
			Printer{CanUseAnsi: true, PrintTime: false, PrintLevel: true},
		},
		{
			"fixed-lines-no-ansi", fixedLines,
			Printer{CanUseAnsi: false, PrintTime: false, PrintLevel: true},
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			l := NewBuffered(buf, tc.Printer)
			tc.DoWork(l)
			l.Close()
			actual := buf.Bytes()
			golden := filepath.Join("test-fixtures", tc.Name+".golden")
			if *update {
				ioutil.WriteFile(golden, actual, 0644)
			}

			expected, _ := ioutil.ReadFile(golden)
			if !bytes.Equal(actual, expected) {
				t.Errorf(
					"golden file does not match output:\nGolden File:\n%s\nActual:\n%s",
					string(expected), string(actual),
				)
			}
		})
	}
}

func minLevel(l Logger) {
	for _, level := range []Level{Progress, Verbose, Info, Warning, Error} {
		l.SetMinLevel(level)
		l.Progressf("this is a progress line")
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
		l.AddFixedLine(),
		l.AddFixedLine(),
		l.AddFixedLine(),
	}
	f[0].Progressf("write to first line")
	f[1].Progressf("write to second line")
	f[2].Progressf("write to third line")
	f[0].Progressf("write back to first line")
	f[2].Progressf("now we're on the third line again")
	f[0].Warningf("something unexpected happened on the first line")
	f[1].Progressf("done")
	f[2].Progressf("done")
	f[0].Progressf("done")
}
