package frog

import (
	"io"
	"math"
	"testing"
	"time"
)

func bench(b *testing.B, newLogger func(fields ...Fielder) (Logger, func())) {
	b.Run("fields=none/min=error/threads=1", func(b *testing.B) {
		log, close := newLogger()
		defer close()
		log.SetMinLevel(Error)
		runInfoMsg(b, log, BMessage)
	})
	b.Run("fields=none/min=info/threads=1", func(b *testing.B) {
		log, close := newLogger()
		defer close()
		log.SetMinLevel(Info)
		runInfoMsg(b, log, BMessage)
	})
	b.Run("fields=static/min=error/threads=1", func(b *testing.B) {
		log, close := newLogger(fieldersSample()...)
		defer close()
		log.SetMinLevel(Error)
		runInfoMsg(b, log, BMessage)
	})
	b.Run("fields=static/min=info/threads=1", func(b *testing.B) {
		log, close := newLogger(fieldersSample()...)
		defer close()
		log.SetMinLevel(Info)
		runInfoMsg(b, log, BMessage)
	})
	b.Run("fields=dynamic/min=error/threads=1", func(b *testing.B) {
		log, close := newLogger()
		defer close()
		log.SetMinLevel(Error)
		runInfoMsgWithFields(b, log, BMessage, fieldersSample)
	})
	b.Run("fields=dynamic/min=info/threads=1", func(b *testing.B) {
		log, close := newLogger()
		defer close()
		log.SetMinLevel(Info)
		runInfoMsgWithFields(b, log, BMessage, fieldersSample)
	})
	b.Run("fields=halfandhalf/min=error/threads=1", func(b *testing.B) {
		log, close := newLogger(fieldersHalf1()...)
		defer close()
		log.SetMinLevel(Error)
		runInfoMsgWithFields(b, log, BMessage, fieldersHalf2)
	})
	b.Run("fields=halfandhalf/min=info/threads=1", func(b *testing.B) {
		log, close := newLogger(fieldersHalf1()...)
		defer close()
		log.SetMinLevel(Info)
		runInfoMsgWithFields(b, log, BMessage, fieldersHalf2)
	})

	// threaded
	b.Run("fields=static/min=error/threads=8", func(b *testing.B) {
		log, close := newLogger(fieldersSample()...)
		defer close()
		log.SetMinLevel(Error)
		b.SetParallelism(8)
		runParallelInfoMsg(b, log, BMessage)
	})
	b.Run("fields=static/min=info/threads=8", func(b *testing.B) {
		log, close := newLogger(fieldersSample()...)
		defer close()
		log.SetMinLevel(Info)
		b.SetParallelism(8)
		runParallelInfoMsg(b, log, BMessage)
	})
}

func Benchmark_UnbufferedJSON(b *testing.B) {
	bench(b, newUnbufJSON)
}

func Benchmark_UnbufferedText_NoColor(b *testing.B) {
	bench(b, newUnbufTextNoColor)
}

func Benchmark_UnbufferedText_Color(b *testing.B) {
	bench(b, newUnbufTextColor)
}

func Benchmark_BufferedTextColor(b *testing.B) {
	bench(b, newBufTextColor)
}

// helpers

func runInfoMsg(b *testing.B, log Logger, msg string) {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info(msg)
	}
}

func runInfoMsgWithFields(b *testing.B, log Logger, msg string, fnFielders func() []Fielder) {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info(msg, fnFielders()...)
	}
}

func runParallelInfoMsg(b *testing.B, log Logger, msg string) {
	b.Helper()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info(msg)
		}
	})
}

func newUnbufJSON(fields ...Fielder) (Logger, func()) {
	b := NewUnbuffered(io.Discard, &JSONPrinter{})
	var f Logger = b
	if len(fields) > 0 {
		f = WithFields(f, fields...)
	}
	return f, b.Close
}

func newUnbufTextNoColor(fields ...Fielder) (Logger, func()) {
	u := NewUnbuffered(io.Discard, (&TextPrinter{}).SetOptions(
		POTime(true), POLevel(true), POFieldIndent(20),
	))
	var f Logger = u
	if len(fields) > 0 {
		f = WithFields(f, fields...)
	}
	return f, u.Close
}

func newUnbufTextColor(fields ...Fielder) (Logger, func()) {
	u := NewUnbuffered(io.Discard, (&TextPrinter{}).SetOptions(
		POPalette(DefaultPalette), POTime(true), POLevel(true), POFieldIndent(20),
	))
	var f Logger = u
	if len(fields) > 0 {
		f = WithFields(f, fields...)
	}
	return f, u.Close
}

func newBufTextColor(fields ...Fielder) (Logger, func()) {
	b := NewBuffered(io.Discard, false, (&TextPrinter{}).SetOptions(
		POPalette(DefaultPalette), POTime(true), POLevel(true), POFieldIndent(20),
	))
	var f Logger = b
	if len(fields) > 0 {
		f = WithFields(f, fields...)
	}
	return f, b.Close
}

// test data

const (
	BMessage = "This is an example log line, medium length"
)

func fieldersSample() []Fielder {
	return []Fielder{
		Bool("bool", false),
		Dur("dur", time.Duration(2903458)*time.Microsecond),
		Err(io.EOF),
		Float32("float32", math.Pi),
		Int64("int64", math.MinInt64),
		String("string", "flargenblargen"),
		TimeNano("time", time.Date(1999, 11, 31, 23, 59, 59, 2391, time.UTC)),
	}
}

func fieldersHalf1() []Fielder {
	return []Fielder{
		Bool("bool", true),
		Dur("dur", time.Duration(2903458)*time.Microsecond),
		Err(io.EOF),
		Float64("float64", math.Pi),
		Int8("int8", math.MinInt8),
		Int32("int32", math.MinInt32),
		String("string", "flargenblargen"),
		TimeNano("time", time.Date(1999, 11, 31, 23, 59, 59, 2391, time.UTC)),
		TimeUnixNano("time", time.Date(1999, 11, 31, 23, 59, 59, 2391, time.UTC)),
		Uint8("uint8", math.MaxUint8),
		Uint32("uint32", math.MaxUint32),
	}
}

func fieldersHalf2() []Fielder {
	return []Fielder{
		Byte("byte", 200),
		Duration("duration", time.Minute+time.Second),
		Float32("float32", math.Pi),
		Int("int", math.MinInt),
		Int16("int16", math.MinInt16),
		Int64("int64", math.MinInt64),
		Time("time", time.Date(1999, 11, 31, 23, 59, 59, 2391, time.UTC)),
		TimeUnix("time", time.Date(1999, 11, 31, 23, 59, 59, 2391, time.UTC)),
		Uint("uint", math.MaxUint),
		Uint16("uint16", math.MaxUint16),
		Uint64("uint64", math.MaxUint64),
	}
}

func fieldersAll() []Fielder {
	return []Fielder{
		Bool("bool", true),
		Byte("byte", 200),
		Dur("dur", time.Duration(2903458)*time.Microsecond),
		Duration("duration", time.Minute+time.Second),
		Err(io.EOF),
		Float32("float32", math.Pi),
		Float64("float64", math.Pi),
		Int("int", math.MinInt),
		Int8("int8", math.MinInt8),
		Int16("int16", math.MinInt16),
		Int32("int32", math.MinInt32),
		Int64("int64", math.MinInt64),
		String("string", "flargenblargen"),
		Time("time", time.Date(1999, 11, 31, 23, 59, 59, 2391, time.UTC)),
		TimeNano("time", time.Date(1999, 11, 31, 23, 59, 59, 2391, time.UTC)),
		TimeUnix("time", time.Date(1999, 11, 31, 23, 59, 59, 2391, time.UTC)),
		TimeUnixNano("time", time.Date(1999, 11, 31, 23, 59, 59, 2391, time.UTC)),
		Uint("uint", math.MaxUint),
		Uint8("uint8", math.MaxUint8),
		Uint16("uint16", math.MaxUint16),
		Uint32("uint32", math.MaxUint32),
		Uint64("uint64", math.MaxUint64),
	}
}
