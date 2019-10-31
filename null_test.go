package frog

import "testing"

func Test_NullIsLogger(t *testing.T) {
	var log Logger
	nl := NullLogger{}

	// The following is a compile error if NullLogger doesn't implement the Logger interface.
	log = &nl

	// We don't actually want to use log for anything, so do this to placate the compiler.
	_ = log
}
