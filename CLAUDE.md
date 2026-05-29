# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Module layout

This is a Go library (`github.com/gilramir/argparse/v2`), a command-line argument
parser loosely modeled on Python's `argparse`. **All source lives in the `v2/`
subdirectory** — that is where `go.mod` is. Run all Go commands from inside `v2/`.

## Commands

```bash
cd v2

go test ./...                              # run all tests
go test -check.f TestChooseHelpLongAtRoot  # run a single test by name (gocheck filter, regex)
go test -check.vv                          # verbose gocheck output
go vet ./...
go build ./...

go run ./examples/example1 -h              # run an example program
```

Tests use the **gocheck** framework (`gopkg.in/check.v1`), not the stdlib `testing`
package directly. `init_test.go` wires gocheck into `go test`; all tests are methods
on `MySuite` with a `*C` receiver, and `-run`/`go test` filtering by test name does
**not** work — use `-check.f <regex>` to select individual tests.

## Architecture

The parse flow is: build a `Command` tree → `New()` finalizes it → `Parse()` /
`ParseAndExit()` walks `os.Args` through a state machine → reflection writes results
into the user's struct.

- **`ArgumentParser`** (`argparse.go`) — top-level handle from `New()`. Holds the
  `Root` *Command, the help switches, i18n `Messages`, and output writers. `Parse()`
  prints help and `os.Exit(0)` on `-h`, prints the error and `os.Exit(1)` on a parse
  error. `ParseAndExit()` additionally runs the triggered command's `Function` and
  exits with its return status.

- **`Command`** (`command.go`) — one node in the (possibly nested) subcommand tree.
  Carries the user's `Values` struct, an optional `Function` callback, and after
  parsing the `Seen` / `CommandSeen` maps (keyed by destination field name). When a
  `Command` is added, `init()` runs and copies any parent arguments marked `Inherit`
  down into it.

- **`Argument`** (`argument.go`) — a switch (starts with `-`) or positional. At
  add-time, argparse reflects over the `Command.Values` struct to find the matching
  exported field (deduced from `Switches`/`Name`, or overridden by `Dest`), and
  **panics** if no compatible field exists. Field-name deduction (`toSafeCamelCase`,
  `firstRuneUpper`): strip leading dashes, CamelCase, e.g. `--no-verify` → `NoVerify`.

- **`value.go`** — a `valueT` implementation per supported field type (bool, string,
  int, int64, float64, time.Duration, and their slice forms). Each knows how to
  `parse()` a string into the reflected field, validate `Choices`, and report its
  default switch arity. **Adding a new supported field type means adding a new
  `valueT` here** and registering it in the type-dispatch in `sanityCheckValueType`.

- **`parse.go`** — the core: a lexer/state-machine (`parserState`, `stateFunc`)
  that tokenizes argv and transitions between `stateArgument`,
  `stateSwitchArgument`, `statePositionalArgument`, `statePassThrough` (after `--`),
  etc. Returns a `parseResults` with the triggered command, ancestor chain, and any
  parse error or help request. This is the most intricate file — read it fully
  before changing parsing semantics.

- **`help.go` / `help_formatter.go`** — render usage and `--help` output. The
  formatter does column layout against the terminal width (via `consolesize` and
  `unicodemonowidth` for wide-character handling).

- **`messages.go`** — all user-facing strings (`DefaultMessages_en`), overridable on
  the parser for i18n.

## Key behaviors to preserve

- **Misconfiguration panics; bad user input returns an error.** Programmer mistakes
  (no matching struct field, wrong value type, `Inherit` with no destination field,
  invalid `NumArgs`/`NumArgsGlob` combinations) panic at setup time. Invalid
  command-line input from the end user produces a `parseError` surfaced via `Stderr`.

- **`NumArgsGlob`** (`"*"`, `"+"`, `"?"`) is positional-only and requires a slice
  field for `*`/`+`. A `*`/`+` positional must be last (unbounded); a `?` positional
  may be followed by others. Switch arities come from `NumArgs` + the field type.

- **Inheritance** is resolved by copying argument definitions parent→child, so the
  child's `Values` struct must independently contain compatible fields — composition
  (embedding the parent options struct) is the intended pattern. See
  `examples/twolevels_with_defaults`.

The `examples/` directory (`example1`, `onelevel`, `twolevels`,
`twolevels_with_defaults`) are runnable references for the public API.
