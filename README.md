# frog

## Overview

Frog is a package for logging output.

When connected to a terminal, Frog allows lines to be anchored to the bottom of the output. This allows progress bars and other UX that is only meant for human consumption. When not connected to a terminal, by default these "transient" lines are skipped (can be overridden).

Windows Users: Frog uses ANSI/VT-100 commands to change colors and move the cursor, and for this to display properly, you must be using Windows 10 build 1511 (circa 2016) or newer, or be using a third-party terminal application like ConEmu/Cmdr or mintty. There's no planned supoprt for the native command prompts of earlier versions of Windows.

![animated gif of frog in action](https://the.real.danbrakeley.com/github/frog-0.2.0-demo.gif)

## Features

- Multiple logging levels
- Anchoring of one or more lines for progress bars and real-time status displays
  - disabled if not connected terminal
- Colorized output
  - disabled if not connected terminal
  - disabled if [NO_COLOR](https://no-color.org) is set
- Use of standard io.Writer interface allows easy retargetting of logs to files, buffers, etc
- Customizable line style via `Printer` interface
  - built in styles for lines in plain text or JSON
- Loggers can be written to target other loggers (see `TeeLogger`)

## Usage

The quickest way to get started is to create one of the default `Logger`s via a call to `frog.New`:

```go
  log := frog.New(frog.Auto)
  defer log.Close()

  log.Info("Example log line")
  log.Warning("Example warning")

  status := frog.AddAnchor(log)
  defer frog.RemoveAnchor(status)
  for i := 0; i <= 100; i++ {
    status.Transient(" + complete", frog.Int("percent", n))
    time.Sleep(time.Duration(10) * time.Millisecond)
  }
  status.Info("Done", n)
```

The parameter `frog.Auto` tells `New` to autodetect if there's a terminal on stdout, and if so, to enable support for anchored lines and colors. There are other default styles you can pass to `New` as well, like `frog.Basic` and `frog.JSON`.

`frog.JSON` will output each log line as a single JSON object. This allows structured data to be easily consumed by a log parser that supports it (ie [filebeat](https://www.elastic.co/products/beats/filebeat)). A JSON logger created in this way doesn't support anchored lines, and by default will not output Transient level lines. Note that you can still call AddAnchor and Transient on such a logger, as the API remains consistent and valid, but nothing changes in the output as a result of these calls.

You can also build a custom Logger, if you prefer. See the implementation of the `New` function in [frog.go](https://github.com/danbrakeley/frog/blob/main/frog.go#L28-L55) for examples.

Here's a complete list of Frog's log levels, and their intended uses:

level | description
--- | ---
`Transient` | Output that is safe to ignore (like progress bars and estimated time remaining).
`Verbose` | Output for debugging (disabled most of the time).
`Info` | Normal log lines (usually enabled).
`Warning` | Unusual events.
`Error` | Something went wrong.

## TODO

- go doc pass
- test on linux and mac
- make some benchmarks, maybe do a pass on performance
- handle terminal width size changing while an app is running (for anchored lines)

## Known Issues

- A single log line will print out all given fields, even if multiple fields use the same name. When outputting JSON, this can result in a JSON object that has multiple fields with the same name. This is not necessarily considered invalid, but it can result in ambiguous behavior.
  - Frog will output the field names in the same order as they are passed to Log/Transient/Verbose/Info/Warning/Error (even when outputting JSON).
  - When there are parent/child relationships, the outermost parent adds their fields first, then the next outermost parent, etc.

## Release Notes

### 0.9.0

- API BREAKING CHANGES
- `Logger` interface: removed `Log()`, added `LogImpl()`.
  - It is likely the signature of LogImpl will change in the near future.
- This is fallout from a large refactor to fix a bug when an `AnchoredLogger` wraps a `CustomizerLogger`
  - `Anchored` was written back when it was the only child logger, and it did dumb things like traverse up the parent chain until it found a Buffered, then set that as its parent (ignoring anything between it and the Buffered). Also it kept a copy of the Buffered's Printer.
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
