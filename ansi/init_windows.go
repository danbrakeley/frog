package ansi

import (
	"os"
	"syscall"
	"unsafe"

	isatty "github.com/mattn/go-isatty"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

const (
	enableProcessedInput       uintptr = 0x0001 // ENABLE_PROCESSED_INPUT
	enableLineInput                    = 0x0002 // ENABLE_LINE_INPUT
	enableEchoInput                    = 0x0004 // ENABLE_ECHO_INPUT
	enableWindowInput                  = 0x0008 // ENABLE_WINDOW_INPUT
	enableMouseInput                   = 0x0010 // ENABLE_MOUSE_INPUT
	enableInsertMode                   = 0x0020 // ENABLE_INSERT_MODE
	enableQuickEditMode                = 0x0040 // ENABLE_QUICK_EDIT_MODE
	enableExtendedFlags                = 0x0080 // ENABLE_EXTENDED_FLAGS
	enableVirtualTerminalInput         = 0x0200 // ENABLE_VIRTUAL_TERMINAL_INPUT

	enableProcessedOutput           uintptr = 0x0001 // ENABLE_PROCESSED_OUTPUT
	enableWrapAtEolOutput                   = 0x0002 // ENABLE_WRAP_AT_EOL_OUTPUT
	enableVirtualTerminalProcessing         = 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING
	disableNewlineAutoReturn                = 0x0008 // DISABLE_NEWLINE_AUTO_RETURN
	enableLvbGridWorldwide                  = 0x0010 // ENABLE_LVB_GRID_WORLDWIDE
)

func init() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		var mode uint32
		procGetConsoleMode.Call(os.Stdout.Fd(), uintptr(unsafe.Pointer(&mode)))
		if (mode & enableVirtualTerminalProcessing) != enableVirtualTerminalProcessing {
			procSetConsoleMode.Call(os.Stdout.Fd(), uintptr(mode|enableVirtualTerminalProcessing))
		}
	}
}
