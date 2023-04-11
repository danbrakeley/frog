package frog

import "testing"

func Test_TeeLoggerInterfaces(t *testing.T) {
	var _ Logger = &TeeLogger{}
	var _ ChildLogger = &TeeLogger{}
}
