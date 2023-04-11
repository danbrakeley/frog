package frog

import "testing"

func Test_BufferedInterfaces(t *testing.T) {
	var _ Logger = &Buffered{}
	var _ AnchorAdder = &Buffered{}
}
