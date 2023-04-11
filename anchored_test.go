package frog

import "testing"

func Test_AnchoredLoggerInterfaces(t *testing.T) {
	var _ Logger = &AnchoredLogger{}
	var _ ChildLogger = &AnchoredLogger{}
	var _ AnchorRemover = &AnchoredLogger{}
}
