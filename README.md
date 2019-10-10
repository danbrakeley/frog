# frog

## Overview

Frog is a package for logging output.

When connected to a terminal, Frog allows lines to be fixed to the bottom of the output and updated quickly and efficiently. This allows for display of data that would be surperfluous in a log file, but is really nice to have when watching the output in a terminal.

Frog uses ANSI/VT-100 commands to change colors and move the cursor. Most terminals support these commands, although Windows is a special case. Windows consoles (cmd, powershell) only support ANSI in the more recent Windows 10 builds. Older versions of Windows require a third-party terminal like ConEmu/Cmdr or mintty to display Frog's colors and fixed lines properly.

![animated gif of frog in action](https://the.real.danbrakeley.com/github/frog-0.2.0-demo.gif)

## Features

- Multiple logging levels
- Fixing one or more lines in place for progress bars and real-time status displays (via ANSI commands)
- Colorized output (via ANSI commands)
- Autodetection of connected terminal, and disabling of ANSI-based output as needed.
- Loggers target the io.Writer interface, allowing your logs to be easily sent to files, buffers, etc.
- Customizable line style, currently supporting a simple plain text format, as well as JSON
  - Alternatively, you can write your own `Printer` that can style each line however you like
- Loggers can easily be written to target other loggers (see `TeeLogger`)

## Usage

The quickest way to get started is to create one of the default `Logger`s via a call to `New`:

```go
  log := frog.New(frog.Auto)
  defer log.Close()

  log.Infof("My first boring log line")
  log.Warningf("Also boring")

  status := frog.AddFixedLine(log)
  for i := 0; i <= 100; i++ {
    status.Transientf(" + %d%% complete", n)
    time.Sleep(time.Duration(10) * time.Millisecond)
  }
  status.Infof("Still boring!", n)
```

`frog.Auto` will use a Buffered Logger, but will only use ANSI if there's a detected terminal, and only displays color if ANSI is supported and if there's no `NO_COLOR` environment variable.

Another option is `frog.JSON`, which uses an Unbuffered Logger that formats each log line as a single JSON object. Switching the previous code example is as simple as changing the first line:

```go
  log := frog.New(frog.JSON)
  ⋮
```

Notice that when you run this second example, all the Transient lines are not output. Transient lines are by default only output when sent to a `Logger` responsible for a fixed line.

Also notice that when we want to create a fixed line, we make a call on the `frog` package, and not directly to our `Logger`. `Logger` is an interface, and most implementations of that interface will probably not care about fixed lines, so it was decided not to make it a mandatory part of the `Logger` interface. Instead, any `Logger` that wants to support fixed liens will need to implement not only the `Logger` interface, but also the `FixedLineLogger` interface, which is what `AddFixedLine` looks for on the passed in `Logger`.

Here's a complete list of Frog's log levels, and their intended uses:

level | description
--- | ---
`Transient` | Data that is safe to ignore (like progress bars and estimated time remaining).
`Verbose` | Data for debugging (disabled most of the time).
`Info` | Normal log lines (usually enabled).
`Warning` | Unusual events.
`Error` | Something went wrong.
`Fatal` | Display this text then halt the entire app.

## TODO

- ✓ ~~Write to two Loggers simultaneously~~
- ✓ ~~JSON Printer~~
- Close should not be how you remove a fixed log line (so that Close isn't called early when fixed lines aren't supported and AddFixedLine returns whatever log you passed it)
- hide cursor while it moves around during fixed line redraws
- Structured logging
- Only re-draw the parts of fixed lines that have changed
