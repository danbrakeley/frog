# frog

## Overview

Frog is a package for logging output. When connected to a terminal, it supports marking one or more lines as bring fixed to the bottom of the output, with new sequential log lines appearing above the fixed lines. The fixed lines can be redrawn individually, allowing for one or more progress bars and other real-time status displays.

When not connected to a terminal, any superfluous fixed-line logs can be skipped.

Frog uses ANSI/VT-100 commands to change colors and move the cursor. Most terminals support this, although Windows cmd prompt has only supported it recently, and even then apps have to specifically enable the support via a call to SetConsoleMode (which frog does). For older versions of Windows, you can use ConEmu, Cmdr, and other third-party terminal apps that parse ANSI/VT-100 commands.

## Features

- Multiple logging levels:
  - `Progress`: Intended for status animations (skipped if stdout not connected to a terminal).
  - `Verbose`: Intended for debugging of your app's behavior.
  - `Info`: Intended for typical log lines.
  - `Warning`: Intended for calling attention to something non-critical, but unusual or important.
  - `Error`: Intended for calling attention to something that went wrong.
  - `Fatal`: Causes app to immediately exit with non-zero exit code.
- Optional use of ANSI colors
  - Windows: cmd.exe: ANSI support requires one of the more recent releases of Windows 10.
  - Windows: ConEmu/Cmdr: Should work regardless of Windows version.
  - Any other environment: relies on your terminal supporting ANSI/VT-100 support
- Fixing one or more lines in place
  - Intended for real-time display of progress/status/percentage/time/etc.
  - Non-fixed log lines continue to output above the fixed lines.

## TODO

- Example demo app
- Only re-draw the parts of fixed lines that have changed
- Ensure all lines are written even if app crashes (maybe use unbuffered logger when no terminal is attached?)
- Optionally caching errors, so a summary can be generated of all the errors so far.
- JSON output
