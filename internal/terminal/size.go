package terminal

import (
	"bufio"
	"fmt"

	"github.com/danbrakeley/frog/ansi"
	"github.com/mattn/go-tty"
)

// created with assistance of ChatGPT v4
func GetSize() (col, row int, err error) {
	t, err := tty.Open()
	if err != nil {
		return 0, 0, fmt.Errorf("Error creating tty: %v", err)
	}
	defer t.Close()

	// Switch to raw mode
	close, err := t.Raw()
	if err != nil {
		return 0, 0, fmt.Errorf("Error setting tty to raw mode: %v", err)
	}
	defer close()

	// Save position, move curser to bottom right corner, then ask for the position
	fmt.Fprint(t.Output(), ansi.PosSave+ansi.BottomRight+ansi.GetCursorPos)

	// Read response from stdin
	reader := bufio.NewReader(t.Input())
	response, err := reader.ReadSlice('R')
	if err != nil {
		return 0, 0, err
	}

	// Restore cursor position
	fmt.Fprint(t.Output(), ansi.PosRestore)

	// Parse response
	n, err := fmt.Sscanf(string(response), ansi.CSI+"%d;%dR", &row, &col)
	if err != nil || n != 2 {
		return 0, 0, fmt.Errorf("Error parsing cursor position: %s", err)
	}

	return col, row, nil
}
