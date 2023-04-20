package frog

import "testing"

func Test_BufferedInterfaces(t *testing.T) {
	var _ RootLogger = &Buffered{}
	var _ AnchorAdder = &Buffered{}
}
