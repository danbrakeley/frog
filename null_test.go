package frog

import "testing"

func Test_NullLoggerInterfaces(t *testing.T) {
	var _ RootLogger = &NullLogger{}
}
