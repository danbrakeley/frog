package frog

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "update golden files")

func AssertGolden(t *testing.T, testName string, actual []byte) {
	t.Helper()
	golden := filepath.Join("test-fixtures", testName+".golden")
	if *update {
		os.Mkdir("test-fixtures", 0644)
		ioutil.WriteFile(golden, actual, 0644)
	}
	expected, _ := ioutil.ReadFile(golden)
	if !bytes.Equal(actual, expected) {
		t.Fatalf(
			"golden file %s does not match output:\nGolden File:\n%s\nActual:\n%s",
			golden, string(expected), string(actual),
		)
	}
}

func Test_Golden(t *testing.T) {
	cases := []struct {
		Name    string
		DoWork  func(Logger)
		Printer Printer
	}{
		{"min-level", minLevel, &TextPrinter{PrintTime: false, PrintLevel: true}},
		{"trims-newlines", newlineVariations, &TextPrinter{PrintTime: false, PrintLevel: true}},
		{"fixed-lines-movement", moveBetweenFixedLines, &TextPrinter{PrintTime: false, PrintLevel: true}},
		{"fixed-lines-add-remove", addAndRemoveFixedLines, &TextPrinter{PrintTime: false, PrintLevel: true}},
	}

	for _, tc := range cases {
		modes := []string{"ansi-color", "ansi", "plain"}
		// run against a Buffered Logger
		for _, mode := range modes {
			t.Run(tc.Name+"."+mode+".buffered", func(t *testing.T) {
				buf := &bytes.Buffer{}
				cfg := Config{
					Writer:   buf,
					UseAnsi:  strings.HasPrefix(mode, "ansi"),
					UseColor: strings.HasSuffix(mode, "color"),
				}
				l := NewBuffered(cfg, tc.Printer)
				tc.DoWork(l)
				l.Close()
				AssertGolden(t, tc.Name+"."+mode, buf.Bytes())
			})
		}

		// run against an Unbuffered Logger
		t.Run(tc.Name+".unbuffered", func(t *testing.T) {
			buf := &bytes.Buffer{}
			l := NewUnbuffered(buf, tc.Printer)
			tc.DoWork(l)
			l.Close()
			AssertGolden(t, tc.Name+".plain", buf.Bytes())
		})

		// run against the JSON printer
		t.Run(tc.Name+".json", func(t *testing.T) {
			buf := &bytes.Buffer{}
			l := NewUnbuffered(buf, &JSONPrinter{TimeOverride: time.Date(2019, 9, 10, 21, 44, 00, 00, time.UTC)})
			tc.DoWork(l)
			l.Close()
			AssertGolden(t, tc.Name+".json", buf.Bytes())
		})

		// run against a TeeLogger, with Buffered as Primary and Unbuffered as Secondary
		for _, mode := range modes {
			t.Run(tc.Name+"."+mode+".tee", func(t *testing.T) {
				buf1 := &bytes.Buffer{}
				cfg := Config{
					Writer:   buf1,
					UseAnsi:  strings.HasPrefix(mode, "ansi"),
					UseColor: strings.HasSuffix(mode, "color"),
				}
				buf2 := &bytes.Buffer{}
				tee := &TeeLogger{
					Primary:   NewBuffered(cfg, tc.Printer),
					Secondary: NewUnbuffered(buf2, tc.Printer),
				}
				tc.DoWork(tee)
				tee.Close()
				AssertGolden(t, tc.Name+"."+mode, buf1.Bytes())
				AssertGolden(t, tc.Name+".plain", buf2.Bytes())
			})
		}
	}
}

func minLevel(l Logger) {
	for _, level := range []Level{Transient, Verbose, Info, Warning, Error} {
		l.SetMinLevel(level)
		l.Transient("this is a transient line")
		l.Verbose("this is a verbose line")
		l.Info("this is an info line")
		l.Warning("this is a warning line")
		l.Error("this is an error line")
	}
}

func newlineVariations(l Logger) {
	l.Info("most of these lines will end up the same")
	l.Info("\nmost of these lines will end up the same")
	l.Info("\n\nmost of these lines will end up the same")
	l.Info("\n\n\nmost of these lines will end up the same")
	l.Info("most of these lines will end up the same\n")
	l.Info("most of these lines will end up the same\n\n")
	l.Info("most of these lines will end up the same\n\n\n")
	l.Info("\nmost of these lines will end up the same\n")
	l.Info("except\nthese last couple of lines, which have newline breaks")
	l.Info("\nexcept these\nlast couple of lines, which\nhave newline breaks")
	l.Info("\n\nexcept these last\ncouple of lines,\n\nwhich have newline breaks")
	l.Info("except these last couple\n\nof lines, which have newline breaks\n")
	l.Info("except these last couple of lines, which\nhave\nnewline breaks\n\n")
	l.Info("\nexcept these last couple\nof lines, which have\n\n\nnewline breaks\n")
}

func moveBetweenFixedLines(l Logger) {
	f := []Logger{
		AddFixedLine(l),
		AddFixedLine(l),
		AddFixedLine(l),
	}
	f[0].Transient("write to first line")
	f[1].Transient("write to second line")
	f[2].Transient("write to third line")
	f[0].Transient("write back to first line")
	f[2].Transient("now we're on the third line again")
	f[0].Warning("something unexpected happened on the first line")
	f[1].Transient("done")
	f[2].Transient("done")
	f[0].Transient("done")
}

func addAndRemoveFixedLines(l Logger) {
	l.Info("before adding fixed logger 1")
	fl := AddFixedLine(l)
	fl.Transient("first fixed line")
	l.Warning("main logger should still log properly")
	fl.Transient("first fixed line again")
	RemoveFixedLine(fl)
	l.Info("after removing fixed logger 1")

	l.Info("before adding fixed logger 2")
	fl2 := AddFixedLine(l)
	l.Info("before adding fixed logger 3")
	fl3 := AddFixedLine(l)
	fl2.Transient("logger 2 status update A")
	l.Info("regular log")
	fl3.Transient("logger 3 status update A")
	fl2.Info("another regular log (via fl2)")
	fl3.Transient("logger 3 status update B")
	fl2.Transient("logger 2 status update B")
	RemoveFixedLine(fl3)
	fl3.Transient("THIS SHOULD NOT BE OUTPUT")
	fl3.Info("this should be redirected to parent logger")

	// End now, before removing fl2, because it should be safe to close a logger
	// without first removing all fixed lines.
}

func Test_FixedLine_Close(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		}
	}()

	// The following is the code under test
	buf := &bytes.Buffer{}
	cfg := Config{
		Writer:   buf,
		UseAnsi:  true,
		UseColor: false,
	}
	log := NewBuffered(cfg, &TextPrinter{})
	fl := AddFixedLine(log)
	// FixedLine panics if you call close (you should only call the parent's close)
	fl.Close()
	log.Close()
}

func Test_AssertInterfaces(t *testing.T) {
	fnAssertAdder := func(t *testing.T, log Logger) {
		t.Helper()
		_, ok := log.(FixedLineAdder)
		if !ok {
			t.Errorf("TeeLogger is not a FixedLineAdder")
		}
	}

	fnAssertRemover := func(t *testing.T, log Logger) {
		t.Helper()
		_, ok := log.(FixedLineRemover)
		if !ok {
			t.Errorf("TeeLogger is not a FixedLineRemover")
		}
	}

	bl := NewBuffered(Config{}, &TextPrinter{})
	fnAssertAdder(t, bl)
	blfl := AddFixedLine(bl)
	fnAssertRemover(t, blfl)

	tl := &TeeLogger{}
	fnAssertAdder(t, tl)
	tlfl := AddFixedLine(tl)
	fnAssertRemover(t, tlfl)
}
