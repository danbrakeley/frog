# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run all tests
go test ./...

# Run a single test by name
go test -run Test_BufferedLogger/anchors-add-remove ./...

# Update golden test files (when output changes intentionally)
go test -run Test_BufferedLogger -update ./...

# Run benchmarks (requires a unique name for the output files)
./bench.sh <name>

# Run example app that demonstrates the basic features of Frog, and accepts optional flags for customizing the output format
go run ./cmd/demo

# Run examples that focus on specific edge cases
go run ./cmd/anchors
go run ./cmd/lengthtest
```

## Architecture

This is a Go structured logging library (`github.com/danbrakeley/frog`). The core design is a **chain of Logger nodes** that pass log events upward to a root logger for rendering and output.

### Logger hierarchy

Every logger implements the `Logger` interface (`logger.go`). There are two root loggers that own I/O:

- **`Buffered`** (`buffered.go`) — goroutine-safe, uses a channel to serialize all writes. The only logger that supports anchored lines. Spawns a background goroutine (`processor()`) that handles ANSI cursor manipulation for anchors.
- **`Unbuffered`** (`unbuffered.go`) — simple synchronous writer, no anchor support.

Child loggers wrap a parent and pass calls up via `LogImpl`, adding their contribution to `ImplData` along the way:

- **`AnchoredLogger`** (`anchored.go`) — sets `ImplData.AnchoredLine` so that `Transient` calls target a specific anchored line in the Buffered processor.
- **`CustomizerLogger`** (`customizer.go`) — merges static fields and/or `PrinterOption`s into every log call. Created by `WithFields`, `WithOptions`, and `WithOptionsAndFields`.
- **`NoAnchorLogger`** (`noanchor.go`) — a no-op wrapper returned by `AddAnchor` when no `Buffered` is in the ancestry.

### `LogImpl` and `ImplData`

`LogImpl(level, msg, fielders, opts, ImplData)` is the internal protocol. As the call traverses up the chain, each node:

1. Calls `d.MergeMinLevel(l.minLevel)` to propagate its min level constraint.
2. Adds its own fields via `d.MergeFields(...)` (CustomizerLogger only).
3. Passes the modified `ImplData` to `parent.LogImpl(...)`.

The root logger (Buffered/Unbuffered) is where level filtering and rendering finally happen.

### Rendering

The `Printer` interface (`printer.go`) has two implementations:

- **`TextPrinter`** — colored/plain text with configurable timestamp, level prefix, field indentation, and transient line cropping.
- **`JSONPrinter`** — one JSON object per line.

`PrinterOption` values (`printeroptions.go`) are algebraic types (empty structs or structs with data) that modify a Printer's behavior. They can be attached to an entire logger chain via `WithOptions`, or applied per-log-line via `Log(level, msg, opts..., fields...)`.

### Fields

`Fielder` (`fields.go`) is an interface with a single `Field() Field` method. The package provides typed constructors like `frog.String`, `frog.Int`, `frog.Err`, `frog.Bool`, `frog.Dur`, `frog.Path`, etc. Fields carry `IsJSONString` and `IsJSONSafe` flags so the printer knows how to quote/escape them.

### Anchoring internals

`Buffered.AddAnchor` assigns a monotonically increasing line number and sends an `mtAddLine` message to the channel. The processor maintains an ordered slice of `anchoredLine` structs; on each `mtPrint` with an anchored target, it uses ANSI escape sequences (`ansi.PrevLine`, `ansi.NextLine`, `ansi.EraseEOL`) to overwrite the correct line without disturbing others. `frog.AddAnchor` (in `frog.go`) walks the parent chain to find an `AnchorAdder`, falling back to `NoAnchorLogger` if none is found.

### Test golden files

Tests in `frog_test.go` compare rendered output against `.golden` files in `testdata/`. The `-update` flag regenerates them. When changing rendering behavior, run tests with `-update` first, then review the diff.
