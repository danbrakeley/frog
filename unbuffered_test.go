package frog

import "testing"

func Test_UnbufferedInterfaces(t *testing.T) {
	var _ Logger = &Unbuffered{}
}
