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

To add an anchored Logger, call `frog.AddAnchor()` and pass in an existing Logger. Frog will walk up the chain of parent/child Loggers until it finds a Logger that also implements the AnchorAdder inteface. If it finds one, it calls that object's AddLogger on that interface, and if it doesn't, it wraps the passed in Logger in a NoAnchorLogger and returns that. The NoAnchorLogger ensures that the behavior is identical to the caller, so that the resulting Logger chain is of a consistent depth. NoAnchorLogger supports SetMinLevel, just like AnchorLogger does.

Once you have a Logger returned from `frog.AddAnchor()`, sending a Transient log line will have special behavior that involves moving the cursor down to a previously established line in the output, overwriting any existing line, then returning the cursor back up, ready to handle any new log lines like normal. Sending a Verbose, Info, Warning, or Error log line will not alter the anchored line, and instead will output the line as if the line had been sent directly to the anchor's parent Logger.

Anchored lines are remembered and redraw as more output is logged above them, so that they always show up at the bottom of the output. If you are done with an anchored line and wish to delete it and stop re-drawing it, then call `frog.RemoveAnchor()`, passing in the Logger that was originally returned from `AddAnchor`. This is optional, and calling `Close()` on a RootLogger that has active anchored lines will remove those lines and clean up the output before `Close()` returns.

You are free to `AddAnchor` and `RemoveAnchor` at any time and in any order.

The code that generated the example output above is in [cmd/anchors/main.go](cmd/anchors/main.go).

Here's a TL;DR on that code:

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

## Windows compatibility

Frog uses ANSI/VT-100 commands to change colors and move the cursor, and for this to display properly, you must be using Windows 10 build 1511 (circa 2016) or newer, or be using a third-party terminal application like ConEmu/Cmdr or mintty. There's no planned supoprt for the native command prompts of earlier versions of Windows.

Windows Terminal also works great, but has problems with ANSI before the [1.17.1023 Preview build](https://github.com/microsoft/terminal/releases/tag/v1.17.1023) (released on 1 Jan 2023, but as of 20 April 2023, it still isn't in the mainline builds).

## Usage

The quickest way to get started is to create one of the default `Logger`s via a call to `frog.New`. The parameter `frog.Auto` tells `New` to autodetect if there's a terminal on stdout, and if so, to enable support for colors and anchored lines. There are other default styles you can pass to `New` as well, like `frog.Basic` and `frog.JSON`. See the implementation of the `New` function in [frog.go](https://github.com/danbrakeley/frog/blob/main/frog.go) for details.

The JSON output from using `frog.JSON` will output each log line as a single JSON object. This allows structured data to be easily consumed by a log parser that supports it (e.g. [filebeat](https://www.elastic.co/products/beats/filebeat)).

## TODO

- handle dicritics in uicode on long transient lines
- make some benchmarks, maybe do a pass on performance
- go doc pass
- test on linux and mac
- handle terminal width size changing while an app is running (for anchored lines)

## Known Issues

- A single log line will print out all given fields, even if multiple fields use the same name. When outputting JSON, this can result in a JSON object that has multiple fields with the same name. This is not necessarily considered invalid, but it can result in ambiguous behavior.
  - Frog will output the field names in the same order as they are passed to Log/Transient/Verbose/Info/Warning/Error (even when outputting JSON).
  - When there are parent/child relationships, the fields are printed starting with the parent, and then each child's static fields (if any) are added in order as you traverse down, child to child. Any fields passed directly to the Logger are added last (and thus will print last in the output).

## Release Notes

### 0.9.2

- **API BREAKING CHANGES**
- Changed frog.Palette from an enum to an array of frog.Color, which allows customizing colors used for each log level.
- TextPrinter no longer exports any of its fields, and instead users should use the printer options, e.g. POPalette(...)
- `Logger` interface changes:
	- `Log()` re-added, is just a passthrough to LogImpl
	- `LogImpl()` arguments re-ordered, and anchored line moved to new ImplData, which also handles min levels
- Added `NoAnchorLogger` to ensure consistent nesting behavior when the RootLogger does not support anchors.

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
