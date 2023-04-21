package frog

import "testing"

func Test_NoAnchorLoggerInterfaces(t *testing.T) {
	var _ Logger = &NoAnchorLogger{}
	var _ ChildLogger = &NoAnchorLogger{}
	var _ AnchorRemover = &NoAnchorLogger{}
}
