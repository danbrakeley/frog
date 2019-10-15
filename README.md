# frog

## Overview

Frog is a package for logging output.

When connected to a terminal, Frog allows lines to be fixed to the bottom of the output, allowing for display of data that would be surperfluous in a log file, but is really nice to have when watching the output in a terminal.

Frog uses ANSI/VT-100 commands to change colors and move the cursor. For this to work in Windows, you must either be using Windows 10 build 1511 (circa 2016) or newer, or be using a third-party terminal application like ConEmu/Cmdr or mintty. There's no planned supoprt for the native command prompts of earlier versions of Windows.

![animated gif of frog in action](https://the.real.danbrakeley.com/github/frog-0.2.0-demo.gif)

## Features

- Multiple logging levels
- Fixing one or more lines in place for progress bars and real-time status displays (via ANSI commands)
- Colorized output (via ANSI commands)
- Detects if there's a connected terminal, and disables ANSI-based output as needed
- Loggers target the io.Writer interface, allowing your logs to be easily sent to files, buffers, etc
- Customizable line style, with built in styles for lines in plain text or JSON
  - Alternatively, you can write your own `Printer` that can generate your own custom style
- Loggers can be written to target other loggers (see `TeeLogger`), allowing a single logger to target multiple outputs and/or styles simultaneously (ie plain text w/ fixed lines on stdout, while also sending json to a file on disk).

## Usage

The quickest way to get started is to create one of the default `Logger`s via a call to `frog.New`:

```go
  log := frog.New(frog.Auto)
  defer log.Close()

  log.Info("Example log line")
  log.Warning("Example warning")

  status := frog.AddFixedLine(log)
  for i := 0; i <= 100; i++ {
    status.Transient(" + complete", frog.Int("percent", n))
    time.Sleep(time.Duration(10) * time.Millisecond)
  }
  status.Info("Done", n)
```

The parameter `frog.Auto` tells `New` to autodetect if there's a terminal on stdout, and if so, to enable support for fixed lines and colors. There are other default styles you can pass to `New` as well, like `frog.Basic` and `frog.JSON`.

`frog.JSON` will output each log line as a single JSON object. This allows structured data to be easily consumed by a log parser that supports it (ie [filebeat](https://www.elastic.co/products/beats/filebeat)). A JSON logger created in this way doesn't support fixed lines, and by default will not output Transient level lines, but from the code's perspective, the API is identical, and calls to AddFixedLine and Transient continue to behave the same as before.

You can also build a custom Logger, with its own custom line printing style, and using sources other than stdout. See the implementation of the `New` function in [frog.go](https://github.com/danbrakeley/frog/blob/master/frog.go#L40-L79) for examples.

Here's a complete list of Frog's log levels, and their intended uses:

level | description
--- | ---
`Transient` | Output that is safe to ignore (like progress bars and estimated time remaining).
`Verbose` | Output for debugging (disabled most of the time).
`Info` | Normal log lines (usually enabled).
`Warning` | Unusual events.
`Error` | Something went wrong.
`Fatal` | After displaying this log line, halt the entire app.

## TODO

- ✓ ~~Write to two Loggers simultaneously~~
- ✓ ~~JSON Printer~~
- ✓ ~~rework how fixed lines are released~~
- ✓ ~~structured logging~~
- use colors to make difference betwen fields and msg more obvious
- go doc pass
- run it on other platforms
- make some benchmarks, maybe do a pass on performance
