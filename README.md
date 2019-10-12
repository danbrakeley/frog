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

The quickest way to get started is to create one of the default `Logger`s via a call to `New`:

```go
  log := frog.New(frog.Auto)
  defer log.Close()

  log.Infof("Example log line")
  log.Warningf("Example warning")

  status := frog.AddFixedLine(log)
  for i := 0; i <= 100; i++ {
    status.Transientf(" + %d%% complete", n)
    time.Sleep(time.Duration(10) * time.Millisecond)
  }
  status.Infof("Done", n)
```

`frog.Auto` will use a Buffered Logger, but will only use ANSI if there's a detected terminal, and only displays color if ANSI is supported and if there's no `NO_COLOR` environment variable.

The implementation of the `New` func shows a couple of different examlpes of how to create customized versions of the built-in Loggers, and is a good reference. See [frog.go](https://github.com/danbrakeley/frog/blob/master/frog.go#L20-L46).

Another value you can pass to `New` is `frog.JSON`, which uses an Unbuffered Logger that formats each log line as a single JSON object. Switching the previous code example is as simple as changing the first line:

```go
  log := frog.New(frog.JSON)
  ⋮
```

Notice that when you run this second example, all the `Transient` level log lines are not output. `Transient` lines are by default only output when sent to a fixed line, and if called on a `Logger` that does not currently represent a fixed line, then those lines are dropped.

Also notice that when we want to create a fixed line, we make a call on the `frog` package, and not directly to our `Logger`. Similarly, if you want to remove a fixed line, you call `frog.RemoveFixedLine` on the fixed line logger. Both calls are safe to call on `Loggers` that do not support fixed lines, and in that case will just return the passed `Logger`. The intent with this API is that you can write your code to always assume fixed line support, and it will just work when there isn't fixed line support present at runtime.

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
- Structured logging
- hide cursor while it moves around during fixed line redraws
