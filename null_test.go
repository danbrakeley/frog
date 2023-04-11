package frog

import "testing"

func Test_NullLoggerInterfaces(t *testing.T) {
	var _ Logger = &NullLogger{}
}
