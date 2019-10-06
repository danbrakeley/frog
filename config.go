package frog

import "io"

type Config struct {
	Writer   io.Writer
	UseAnsi  bool
	UseColor bool
}
