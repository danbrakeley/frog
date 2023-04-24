# frog

## Overview

Frog is a package for structured logging that is fast, easy to use, good looking, and customizable.

Frog includes:

- Built-in support for plain text or JSON output.
- Plain text output optinally supports ANSI colors per log level, including user-defined palettes.
  - Respects [NO_COLOR](https://no-color.org) env var
- [Anchoring](#anchoring) of log lines to the bottom of the terminal output, for progress bars and real-time status updates.
- Detection of terminal/tty and disabling of ANSI/anchoring when none is found.
- Nesting of Loggers to add fields, anchored lines, custom line rendering settings, and other custom behavior.
	- Each additional nested layer adds context without altering its parent(s).
- User-customizable line rendering via the `Printer` interface.
- Five log levels:

level | description
--- | ---
`Transient` | Output that is safe to ignore (like progress bars and estimated time remaining).
`Verbose` | Output for debugging (disabled by default).
`Info` | Normal events.
`Warning` | Unusual events.
`Error` | Something went wrong.

## Anchoring

![animated gif of anchoring in action](images/anchor-demo-0.9.2.gif)

Anchors require an ANSI-compatible terminal connected to the output.

To add an anchored Logger, call `frog.AddAnchor()` and pass in an existing Logger, and it will return a Logger whose Transient log lines will target a newly created anchored line.

Behind the scenes, `AddAnchor` is looking to see if the given Logger or any of its ancestors implement the AnchorAdder interface. Currently `Buffered` is the only included Logger that supports anchors, and its AddAnchor returns an instance of `AnchoredLogger` that wraps the given Logger. If their is no Buffered/AnchorAdder in the ancestry, then it will use `NoAnchorLogger` instead.

Anchored lines are drawn to the terminal using ANSI escape codes for manipulating the cursor. Because the Buffered logger serializes logging from any number of goroutines, it provides a safe environment in which to manipulate cursor position temporarily, and to re-draw anchored lines as needed when they get blown away by non-Transient log lines.

Note that sending a Verbose, Info, Warning, or Error log line via an AnchoredLogger will not target the anchored line, but will log the results as if you had sent the log lines through the parent that was passed into `frog.AddAnchor` in the first place.

When you are done with an anchored line and wish to have it stop redrawing itself at the bottom of the output, just call `frog.RemoveAnchor()` on the logger that was returned from `frog.AddAnchor`. At that point you can keep using the logger if you wish, or discard it.

Calling `RemoveAnchor` is optional.

You are free to `AddAnchor` and `RemoveAnchor` at any time and in any order.

The code that generated the example shown output is here: [cmd/anchors/main.go](cmd/anchors/main.go), but the core of what it is doing is:

```go
func main() {
	log := frog.New(frog.Auto)
	defer log.Close()

	wg := new(sync.WaitGroup)
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(log frog.Logger, n int) {
			defer wg.Done()
			defer frog.RemoveAnchor(log)
			for j := 0; j <= 100; j++ {
				log.Transient(" + Status", frog.Int("thread", n), frog.Int("percent", j))
				time.Sleep(time.Duration(50) * time.Millisecond)
			}
		}(frog.AddAnchor(log), i)
	}
	wg.Wait()
}
```

Note that `frog.New(frog.Auto)` automatically detects if a terminal is connected to the output, and if not, it turns off anchors. To see this in action, you can try piping the output of running the previous demo into a file. For example:

```txt
$ go run ./cmd/anchors/ -> out.txt && cat out.txt
2023.04.21-03:49:13 [nfo] Spawning example threads...   count=3
2023.04.21-03:49:15 [nfo] waited for one second...
2023.04.21-03:49:16 [WRN] waited for two seconds...
2023.04.21-03:49:17 [ERR] BORED OF WAITING
2023.04.21-03:49:19 [nfo] All threads done!
```

You can see that no ANSI escape sequences or transient log lines ended up in the resulting file.

## Nesting

TODO: build nested loggers, then draw graph to illustrate the parent/child relationships that are formed

## Windows compatibility

Frog uses ANSI/VT-100 commands to change colors and move the cursor, and for this to display properly, you must be using Windows 10 build 1511 (circa 2016) or newer, or be using a third-party terminal application like ConEmu/Cmdr or mintty. There's no planned supoprt for the native command prompts of earlier versions of Windows.

Windows Terminal also works great, but has problems with ANSI before the [1.17.1023 Preview build](https://github.com/microsoft/terminal/releases/tag/v1.17.1023) (released on Jan 1, 2023).

## Usage

The quickest way to get started is to create one of the default `Logger`s via a call to `frog.New`. The parameter `frog.Auto` tells `New` to autodetect if there's a terminal on stdout, and if so, to enable support for colors and anchored lines. There are other default styles you can pass to `New` as well, like `frog.Basic` and `frog.JSON`. See the implementation of the `New` function in [frog.go](https://github.com/danbrakeley/frog/blob/main/frog.go) for details.

The JSON output from using `frog.JSON` will output each log line as a single JSON object. This allows structured data to be easily consumed by a log parser that supports it (e.g. [filebeat](https://www.elastic.co/products/beats/filebeat)).

## TODO

- handle dicritics in uicode on long transient lines
- go doc pass
- test on linux and mac
- handle terminal width size changing while an app is running (for anchored lines)

## Known Issues

- When using anchored lines, if you resize the terminal to be narrowing than when frog was initialized, lines won't be properly cropped, and a long enough line could cause extra wrapping that would break the anchored line's ability to redraw itself. The result would be slightly garbled output. See the TODO in the previous section about this.
- A single log line will print out all given fields, even if multiple fields use the same name. When outputting JSON, this can result in a JSON object that has multiple fields with the same name. This is not necessarily considered invalid, but it can result in ambiguous behavior.
  - Frog will output the field names in the same order as they are passed to Log/Transient/Verbose/Info/Warning/Error (even when outputting JSON).
  - When there are parent/child relationships, the fields are printed starting with the parent, and then each child's static fields (if any) are added in order as you traverse down, child to child. Any fields passed with the log line itself are added last.

## Release Notes

### 0.9.3

- Improved JSONPrinter performance by switching to StringBuilder

### 0.9.2

- **API BREAKING CHANGES**
- Changed frog.Palette from an enum to an array of frog.Color, which allows customizing colors used for each log level.
- TextPrinter no longer exports any of its fields, and instead users should use the printer options, e.g. POPalette(...)
- `Logger` interface changes:
	- `Log()` re-added, is just a passthrough to LogImpl
	- `LogImpl()` arguments re-ordered, and anchored line moved to new ImplData, which also handles min levels
- Added `NoAnchorLogger` to ensure consistent nesting behavior when the RootLogger does not support anchors.
- Removed the "buffered log closing" Debug log line that was previously sent when a Buffered Logger was closed.

### 0.9.0

- **API BREAKING CHANGES**
- `Logger` interface: removed `Log()`, added `LogImpl()`.
  - It is likely the signature of LogImpl will change in the near future.
- This is fallout from a large refactor to fix a bug when an `AnchoredLogger` wraps a `CustomizerLogger`
  - `AnchoredLogger` was written back when it was the only child logger, and it did dumb things like traverse up the parent chain until it found a Buffered, then set that as its parent (ignoring anything between it and the Buffered). Also it kept a copy of the Buffered's Printer.
  - Now all Loggers that wrap other loggers are expected to know their parent and pass through requests, making modifications as needed, until the request hits the root Logger. This is why LogImpl currently has the anchored line as a parameter. I hope to obfuscate this in the future.
- Added an extra Customizer to the relevant tests, so going forward this case will be tested.

### 0.8.4

- Handle anchored lines that are longer than the terminal width more gracefully
  - New behavior is that we detect the terminal width when frog comes up, then crops transient lines to that width
  - Can be manually set with printer option `POTransientLineLength(len)`

### 0.8.0

- Re-worked printer options to be algebraic data types (see printeroptions.go and printer.go)
- Allow overriding printer options per log line (see changes to the Printer's `Render` method and Logger's `Log` method).
- `frog.WithFields(parent, fields...)` allows creating a logger that is a pass-through to `parent`, but always adds the specified fields to each log line.
- `frog.WithPrinterOptions(parent, opts...)` allows creating a logger is a pass-through to `parent`, but always sets the specified PrinterOptions on each log line.
