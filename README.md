# frog

## Overview

Frog is a package for logging output.

When connected to a terminal, Frog allows lines to be anchored to the bottom of the output. This allows progress bars and other UX that is only meant for human consumption, which are not included when logging to a file or pipe.

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
  for i := 0; i <= 100; i++ {
    status.Transient(" + complete", frog.Int("percent", n))
    time.Sleep(time.Duration(10) * time.Millisecond)
  }
  status.Info("Done", n)
```

The parameter `frog.Auto` tells `New` to autodetect if there's a terminal on stdout, and if so, to enable support for anchored lines and colors. There are other default styles you can pass to `New` as well, like `frog.Basic` and `frog.JSON`.

`frog.JSON` will output each log line as a single JSON object. This allows structured data to be easily consumed by a log parser that supports it (ie [filebeat](https://www.elastic.co/products/beats/filebeat)). A JSON logger created in this way doesn't support anchored lines, and by default will not output Transient level lines. Note that you can still call AddAnchor and Transient on such a logger, as the API remains consistent and valid, but nothing changes in the output as a result of these calls.

You can also build a custom Logger, if you prefer. See the implementation of the `New` function in [frog.go](https://github.com/danbrakeley/frog/blob/master/frog.go#L40-L79) for examples.

Here's a complete list of Frog's log levels, and their intended uses:

level | description
--- | ---
`Transient` | Output that is safe to ignore (like progress bars and estimated time remaining).
`Verbose` | Output for debugging (disabled most of the time).
`Info` | Normal log lines (usually enabled).
`Warning` | Unusual events.
`Error` | Something went wrong.

## TODO

- create a Logger with static fields
- go doc pass
- test on linux and mac
- make some benchmarks, maybe do a pass on performance
