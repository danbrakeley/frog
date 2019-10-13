package frog

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
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
		{"fields", fields, &TextPrinter{PrintTime: false, PrintLevel: true}},
	}

	modesBasic := []string{"color", "plain"}
	modesAll := []string{}
	for _, v := range modesBasic {
		modesAll = append(modesAll, v+".fixedline")
		modesAll = append(modesAll, v)
	}

	for _, tc := range cases {
		for _, mode := range modesAll {
			t.Run(tc.Name+"."+mode, func(t *testing.T) {
				buf := &bytes.Buffer{}
				cfg := Config{
					Writer:   buf,
					UseColor: strings.HasPrefix(mode, "color"),
				}
				var log Logger
				if strings.HasSuffix(mode, ".fixedline") {
					log = NewBuffered(cfg, tc.Printer)
				} else {
					log = NewUnbuffered(cfg, tc.Printer)
				}
				tc.DoWork(log)
				log.Close()
				AssertGolden(t, tc.Name+"."+mode, buf.Bytes())
			})
		}

		// run against the JSON printer
		t.Run(tc.Name+".json", func(t *testing.T) {
			buf := &bytes.Buffer{}
			cfg := Config{
				Writer:   buf,
				UseColor: false,
			}
			l := NewUnbuffered(cfg, &JSONPrinter{TimeOverride: time.Date(2019, 9, 10, 21, 44, 00, 00, time.UTC)})
			tc.DoWork(l)
			l.Close()
			AssertGolden(t, tc.Name+".json", buf.Bytes())

			// parse each line as a JSON object to ensure only valid JSON is produced
			lastLineWasEmpty := false
			for i, line := range strings.Split(string(buf.Bytes()), "\n") {
				if lastLineWasEmpty {
					// the only empty lines should be the last line
					t.Errorf("empty line in json output at line %d", i)
				}
				if len(strings.TrimSpace(line)) == 0 {
					lastLineWasEmpty = true
					continue
				}
				target := make(map[string]interface{})
				err := json.Unmarshal([]byte(line), &target)
				if err != nil {
					t.Errorf("error parsing logged json: %v\n\n%s\n\n", err, line)
				}
			}
		})

		// run against a TeeLogger, with Buffered as Primary and Unbuffered as Secondary
		for _, mode := range modesBasic {
			t.Run(tc.Name+"."+mode+".tee", func(t *testing.T) {
				useColor := strings.HasPrefix(mode, "color")
				buf1 := &bytes.Buffer{}
				buf2 := &bytes.Buffer{}
				tee := &TeeLogger{
					Primary:   NewBuffered(Config{Writer: buf1, UseColor: useColor}, tc.Printer),
					Secondary: NewUnbuffered(Config{Writer: buf2, UseColor: useColor}, tc.Printer),
				}
				tc.DoWork(tee)
				tee.Close()
				AssertGolden(t, tc.Name+"."+mode+".fixedline", buf1.Bytes())
				AssertGolden(t, tc.Name+"."+mode, buf2.Bytes())
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

func fields(l Logger) {
	// bool
	l.Info("bool", Bool("true", true))
	l.Warning("bool", Bool("false", false))

	// byte
	l.Info("byte", Byte("min", byte(0)))
	l.Warning("byte", Byte("max", byte(255)))

	// dur/duration
	l.Info("time.Duration", Dur("how_long", time.Duration(125)*time.Second))
	d, _ := time.ParseDuration("4h48m1s")
	l.Warning("time.Duration", Duration("this_long", d))

	// err
	l.Error("error", Err(fmt.Errorf("this is the error")))
	l.Warning("error", Err(nil))

	// float32
	l.Info("float32", Float32("floatymc", float32(3.3333433)))
	l.Warning("float32", Float32("floatface", float32(-0.000000000000002)))

	// float64
	l.Info("float64", Float64("flargen", float64(0)))
	l.Warning("float64", Float64("blargen", float64(-1.234456e+78)))

	// int
	l.Info("int", Int("zero", int(0)))
	l.Warning("int", Int("negative", int(-1)))

	// int8
	l.Info("int8", Int8("max", int8(127)))
	l.Warning("int8", Int8("min", int8(-128)))

	// int16
	l.Info("int16", Int16("max", int16(32767)))
	l.Warning("int16", Int16("min", int16(-32768)))

	// int32
	l.Info("int32", Int32("max", int32(2147483647)))
	l.Warning("int32", Int32("min", int32(-2147483648)))

	// int64
	l.Info("int64", Int64("max", int64(9223372036854775807)))
	l.Warning("int64", Int64("min", int64(-9223372036854775808)))

	// string
	l.Info("string", String("empty", ""))
	l.Info("string", String("space", " "))
	l.Info("string", String("quotes", "\""))
	l.Info("string", String("newline", "\n"))
	l.Info("string", String("newline", "a"))
	l.Info("string", String("punctuation", "!@#$%^&*()_+-=[]{}|;':,.<>?"))
	l.Warning("string", String("long", "this is a relatively long sentence with ʎzɐɹɔ cha\rac\ters i\n it \u0001 \"<<&&>>\""))

	// time
	l.Info("time.Time", Time("party", time.Date(1999, 01, 01, 00, 00, 00, 00, time.UTC)))
	l.Warning("time.Time", Time("future", time.Date(2038, 07, 13, 2, 55, 13, 12398456, time.UTC)))

	// timenano
	l.Info("time.Time (nano)", TimeNano("party", time.Date(1999, 01, 01, 00, 00, 00, 00, time.UTC)))
	l.Warning("time.Time (nano)", TimeNano("future", time.Date(2038, 07, 13, 2, 55, 13, 12398456, time.UTC)))

	// timeunix
	l.Info("time.Time (unix)", TimeUnix("party", time.Date(1999, 01, 01, 00, 00, 00, 00, time.UTC)))
	l.Warning("time.Time (unix)", TimeUnix("future", time.Date(2038, 07, 13, 2, 55, 13, 12398456, time.UTC)))

	// timeunixnano
	l.Info("time.Time (unix,nano)", TimeUnixNano("party", time.Date(1999, 01, 01, 00, 00, 00, 00, time.UTC)))
	l.Warning("time.Time (unix,nano)", TimeUnixNano("future", time.Date(2038, 07, 13, 2, 55, 13, 12398456, time.UTC)))

	// uint
	l.Info("uint", Uint("zero", uint(0)))
	l.Warning("uint", Uint("one", uint(1)))

	// uint8
	l.Info("uint8", Uint8("max", uint8(255)))
	l.Warning("uint8", Uint8("min", uint8(0)))

	// uint16
	l.Info("uint16", Uint16("max", uint16(65535)))
	l.Warning("uint16", Uint16("min", uint16(0)))

	// uint32
	l.Info("uint32", Uint32("max", uint32(4294967295)))
	l.Warning("uint32", Uint32("min", uint32(0)))

	// uint64
	l.Info("uint64", Uint64("max", uint64(18446744073709551615)))
	l.Warning("uint64", Uint64("min", uint64(0)))
}

func Test_FixedLine_Close(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		}
	}()

	// The following code just needs to not panic in order to pass
	buf := &bytes.Buffer{}
	cfg := Config{
		Writer:   buf,
		UseColor: false,
	}
	log := NewBuffered(cfg, &TextPrinter{})
	fl := AddFixedLine(log)
	// FixedLine panics if you call close (you should only call the parent's close)
	fl.Close()
}

func Test_AssertInterfaces(t *testing.T) {
	fnAssertAdder := func(t *testing.T, log Logger) {
		t.Helper()
		_, ok := log.(FixedLineAdder)
		if !ok {
			t.Errorf("logger is not a FixedLineAdder")
		}
	}

	fnAssertRemover := func(t *testing.T, log Logger) {
		t.Helper()
		_, ok := log.(FixedLineRemover)
		if !ok {
			t.Errorf("logger is not a FixedLineRemover")
		}
	}

	bl := NewBuffered(Config{Writer: &bytes.Buffer{}}, &TextPrinter{})
	fnAssertAdder(t, bl)
	blfl := AddFixedLine(bl)
	fnAssertRemover(t, blfl)

	tl := &TeeLogger{}
	fnAssertAdder(t, tl)
	tlfl := AddFixedLine(tl)
	fnAssertRemover(t, tlfl)
}
