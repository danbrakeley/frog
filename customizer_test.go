package frog

import "testing"

func Test_CustomizerLoggerInterfaces(t *testing.T) {
	var _ Logger = &CustomizerLogger{}
	var _ ChildLogger = &CustomizerLogger{}
}
