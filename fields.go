package frog

import (
	"strconv"
	"time"
)

type Field struct {
	Name         string
	Value        string
	IsJSONString bool // if true, the string in Value should be bookended by double quotes to be valid JSON
	IsJSONSafe   bool // if true, this string only contains alpha-numerics, spaces, and safe punctuation
}

// Fielder is an interface used to add structured logging to calls to Logger methods
type Fielder interface {
	Field() Field
}

// Fieldify takes a slice of Fielders and calls Field() on each, resulting in a slice of Fields
func Fieldify(f []Fielder) []Field {
	fields := make([]Field, len(f))
	for i := range f {
		fields[i] = f[i].Field()
	}
	return fields
}

// FieldifyAndAppend returns a slice of Fields that starts with a copy the passed in []Field,
// and then appends the passed in []Fielder, after first rendering them to Fields.
func FieldifyAndAppend(fields []Field, fielders []Fielder) []Field {
	var out []Field
	if len(fielders)+len(fields) > 0 {
		out = make([]Field, len(fields), len(fielders)+len(fields))
		copy(out, fields)
		for _, fielder := range fielders {
			out = append(out, fielder.Field())
		}
	}
	return out
}

// Bool adds a field whose value will be true or false
func Bool(name string, value bool) FieldBool {
	return FieldBool{Name: name, Value: value}
}

// Byte adds an 8-bit unsigned integer field
func Byte(name string, value byte) FieldUint64 {
	return FieldUint64{Name: name, Value: uint64(value)}
}

// Dur adds a time.Duration field
func Dur(name string, value time.Duration) FieldDuration {
	return Duration(name, value)
}

// Duration adds a time.Duration field
func Duration(name string, value time.Duration) FieldDuration {
	return FieldDuration{Name: name, Value: value}
}

// Err adds an error field named "error"
func Err(value error) FieldError {
	return FieldError{Name: "error", Value: value}
}

// Float32 adds a 32-bit floating point number field
func Float32(name string, value float32) FieldFloat32 {
	return FieldFloat32{Name: name, Value: value}
}

// Float64 adds a 64-bit floating point number field
func Float64(name string, value float64) FieldFloat64 {
	return FieldFloat64{Name: name, Value: value}
}

// Int adds a signed integer field
func Int(name string, value int) FieldInt64 {
	return FieldInt64{Name: name, Value: int64(value)}
}

// Int8 adds an 8-bit signed integer field
func Int8(name string, value int8) FieldInt64 {
	return FieldInt64{Name: name, Value: int64(value)}
}

// Int16 adds a 16-bit signed integer field
func Int16(name string, value int16) FieldInt64 {
	return FieldInt64{Name: name, Value: int64(value)}
}

// Int32 adds a 32-bit signed integer field
func Int32(name string, value int32) FieldInt64 {
	return FieldInt64{Name: name, Value: int64(value)}
}

// Int64 adds a 64-bit signed integer field
func Int64(name string, value int64) FieldInt64 {
	return FieldInt64{Name: name, Value: value}
}

// String adds an escaped and quoted string field
func String(name string, value string) FieldString {
	return FieldString{Name: name, Value: value}
}

// Time adds a time.Time field that will output a string formatted using RFC 3339 (ISO 8601)
func Time(name string, value time.Time) FieldTimeFormat {
	return FieldTimeFormat{Name: name, Value: value, Format: time.RFC3339}
}

// TimeNano adds a time.Time field that will output a string formatted using RFC 3339 with nanosecond precision
func TimeNano(name string, value time.Time) FieldTimeFormat {
	return FieldTimeFormat{Name: name, Value: value, Format: time.RFC3339Nano}
}

// TimeUnix adds a time.Time field that outputs as a unix epoch (unsigned integer)
func TimeUnix(name string, value time.Time) FieldTimeUnix {
	return FieldTimeUnix{Name: name, Value: value}
}

// TimeUnixNano adds a time.Time field that outputs as a unix epoch with nanosecond precision (unsigned integer)
func TimeUnixNano(name string, value time.Time) FieldTimeUnixNano {
	return FieldTimeUnixNano{Name: name, Value: value}
}

// Uint adds an unsigned integer field
func Uint(name string, value uint) FieldUint64 {
	return FieldUint64{Name: name, Value: uint64(value)}
}

// Uint8 adds an 8-bit unsigned integer field
func Uint8(name string, value uint8) FieldUint64 {
	return FieldUint64{Name: name, Value: uint64(value)}
}

// Uint16 adds a 16-bit unsigned integer field
func Uint16(name string, value uint16) FieldUint64 {
	return FieldUint64{Name: name, Value: uint64(value)}
}

// Uint32 adds a 32-bit unsigned integer field
func Uint32(name string, value uint32) FieldUint64 {
	return FieldUint64{Name: name, Value: uint64(value)}
}

// Uint64 adds a 64-bit unsigned integer field
func Uint64(name string, value uint64) FieldUint64 {
	return FieldUint64{Name: name, Value: value}
}

// Bool

type FieldBool struct {
	Name  string
	Value bool
}

func (f FieldBool) Field() Field {
	if f.Value {
		return Field{Name: f.Name, Value: "true"}
	}
	return Field{Name: f.Name, Value: "false"}
}

// Duration

type FieldDuration struct {
	Name  string
	Value time.Duration
}

func (f FieldDuration) Field() Field {
	return Field{Name: f.Name, Value: f.Value.String(), IsJSONString: true, IsJSONSafe: true}
}

// Error

type FieldError struct {
	Name  string
	Value error
}

func (f FieldError) Field() Field {
	if f.Value == nil {
		return Field{Name: f.Name, Value: "null"}
	}
	return Field{Name: f.Name, Value: f.Value.Error(), IsJSONString: true}
}

// Float32

type FieldFloat32 struct {
	Name  string
	Value float32
}

func (f FieldFloat32) Field() Field {
	return Field{Name: f.Name, Value: strconv.FormatFloat(float64(f.Value), 'g', -1, 32)}
}

// Float64

type FieldFloat64 struct {
	Name  string
	Value float64
}

func (f FieldFloat64) Field() Field {
	return Field{Name: f.Name, Value: strconv.FormatFloat(f.Value, 'g', -1, 64)}
}

// Int, Int8, Int16, Int32, Int64

type FieldInt64 struct {
	Name  string
	Value int64
}

func (f FieldInt64) Field() Field {
	return Field{Name: f.Name, Value: strconv.FormatInt(f.Value, 10)}
}

// String

type FieldString struct {
	Name  string
	Value string
}

func (f FieldString) Field() Field {
	return Field{Name: f.Name, Value: f.Value, IsJSONString: true, IsJSONSafe: false}
}

// Time

type FieldTimeFormat struct {
	Name   string
	Value  time.Time
	Format string
}

func (f FieldTimeFormat) Field() Field {
	return Field{Name: f.Name, Value: f.Value.Format(f.Format), IsJSONString: true, IsJSONSafe: true}
}

type FieldTimeUnix struct {
	Name  string
	Value time.Time
}

func (f FieldTimeUnix) Field() Field {
	return Field{Name: f.Name, Value: strconv.FormatInt(f.Value.Unix(), 10)}
}

type FieldTimeUnixNano struct {
	Name  string
	Value time.Time
}

func (f FieldTimeUnixNano) Field() Field {
	return Field{Name: f.Name, Value: strconv.FormatInt(f.Value.UnixNano(), 10)}
}

// Uint, Uint8, Uit16, Uint32, Uint64, Byte

type FieldUint64 struct {
	Name  string
	Value uint64
}

func (f FieldUint64) Field() Field {
	return Field{Name: f.Name, Value: strconv.FormatUint(f.Value, 10)}
}
