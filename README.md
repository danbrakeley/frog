# frog

## Overview

Frog is a package for logging output.

When connected to a terminal, Frog allows lines to be fixed to the bottom of the output and updated quickly and efficiently. This allows for display of data that would be surperfluous in a log file, but is really nice to have when watching the output in a terminal.

Frog uses ANSI/VT-100 commands to change colors and move the cursor. Most terminals support these commands, although Windows is a special case. The Windows cmd prompt only supports ANSI in the more recent Windows 10 builds. Older versions of Windows requires a third-party terminal like ConEmu/Cmdr or mintty.

![frog demo](https://the.real.danbrakeley.com/github/frog-0.2.0-demo.gif)

## Features

- Multiple logging levels
  - `Transient`: Data that is safe to ignore (like progress bars and estimated time remaining).
  - `Verbose`: Data for debugging (disabled most of the time).
  - `Info`: Normal log lines (usually enabled).
  - `Warning`: Unusual events.
  - `Error`: Something went wrong.
  - `Fatal`: Display this text then halt the entire app.
- Fixing one or more lines in place
  - Creates a sub-Logger that re-writes Transient log lines.
  - non-Transient lines get sent up to the parent logger for normal processing and display.
- Optional use of ANSI colors
  - Automatic disabling of ANSI if no terminal detected.
  - If your terminal is a Windows cmd prompt, attempt to enable native ANSI support (recent updates of Windows 10 only).
  - Support `NO_COLOR` env var to disable colorized output (does not disable all ANSI commands, just ANSI colors).

## TODO

- ✓ ~~Write to two Loggers simultaneously~~
- JSON Printer
- Structured logging
- Only re-draw the parts of fixed lines that have changed
