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

func moveBetweenFixedLines(l Logger) {
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

func addAndRemoveFixedLines(l Logger) {
	l.Infof("before adding fixed logger 1")
	fl := AddFixedLine(l)
	fl.Transientf("first fixed line")
	l.Warningf("main logger should still log properly")
	fl.Transientf("first fixed line again")
	RemoveFixedLine(fl)
	l.Infof("after removing fixed logger 1")

	l.Infof("before adding fixed logger 2")
	fl2 := AddFixedLine(l)
	l.Infof("before adding fixed logger 3")
	fl3 := AddFixedLine(l)
	fl2.Transientf("logger 2 status update A")
	l.Infof("regular log")
	fl3.Transientf("logger 3 status update A")
	fl2.Infof("another regular log (via fl2)")
	fl3.Transientf("logger 3 status update B")
	fl2.Transientf("logger 2 status update B")
	RemoveFixedLine(fl3)
	fl3.Transientf("THIS SHOULD NOT BE OUTPUT")
	fl3.Infof("this should be redirected to parent logger")

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
