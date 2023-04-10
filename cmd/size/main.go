package main

import (
	"fmt"
	"os"
	"time"

	"github.com/danbrakeley/frog/internal/terminal"
)

func main() {
	start := time.Now()
	width, height, err := terminal.GetSize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting terminal size: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Width: %d, Height: %d\nTime: %v\n", width, height, time.Now().Sub(start))
}
