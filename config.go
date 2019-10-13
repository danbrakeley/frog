package frog

import "io"

type Config struct {
	Writer   io.Writer
	UseColor bool
}
