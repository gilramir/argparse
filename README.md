# argparse

Argparse is a Go module that makes it easy to write the command-line parsing
part of your program. It loosely follows the conceptual model of the Python
argparse module.

[![Go Reference](https://pkg.go.dev/badge/github.com/gilramir/argparse/v2.svg)](https://pkg.go.dev/github.com/gilramir/argparse/v2)

Highlights:

* You can have nested subcommands.
* The values for the command-line options are stored in a struct of your
  creation.
* Argparse can deduce the name of the value field in the struct by looking
  at the name of the option. Or, you can tell it exactly which field to use.
* Argparse will tell you if a particular option was present on the command-line
  or not present, in case you need that information.
* Options can be inherited by sub-commands, so you need only define them once.
* The user-facing text (help output and error messages) is translatable, with
  English and Korean built in.

# Installation

```
go get github.com/gilramir/argparse/v2
```

Argparse requires Go 1.14 or later, and is imported as package `argparse`:

```go
import "github.com/gilramir/argparse/v2"
```

Note the `/v2` in the module path: that is the import path, and the package
name inside it is `argparse`.

# Contents

* [Quick start](#quick-start)
* [Usage](#usage)
  * [Single-level commands](#single-level-commands)
  * [Sub-commands](#sub-commands)
  * [Default values and "Seen" arguments](#default-values-and-seen-arguments)
* [Parse vs. ParseAndExit vs. ParseArgs](#parse-vs-parseandexit-vs-parseargs)
* [Version](#version)
* [Accepted command-line syntax](#accepted-command-line-syntax)
* [Capturing leftover arguments](#capturing-leftover-arguments)
* [Values struct and field names](#values-struct-and-field-names)
  * [Custom value types](#custom-value-types)
* [Command](#command)
* [Argument](#argument)
* [NumArgs and NumArgsGlob](#numargs-and-numargsglob)
* [Mutually-exclusive groups](#mutually-exclusive-groups)
* [Inheritance by sub-commands](#inheritance-by-sub-commands)
* [Translation](#translation)
* [Examples](#examples)
* [License](#license)

# Quick start

A complete, runnable single-level program:

```go
package main

import (
	"fmt"
	"time"

	"github.com/gilramir/argparse/v2"
)

// The struct that will hold the values parsed from the command-line.
type MyOptions struct {
	Count      int
	Expiration time.Duration
	Verbose    bool
	Names      []string
}

func main() {
	opts := &MyOptions{}

	ap := argparse.New(&argparse.Command{
		Name:        "example1",
		Description: "This is an example program",
		Values:      opts,
	})

	// Switch arguments
	ap.Add(&argparse.Argument{
		Switches: []string{"--count"},
		MetaVar:  "N",
		Help:     "How many items",
	})
	ap.Add(&argparse.Argument{
		Switches: []string{"--expiration", "-x"},
		Help:     "How long: #(h|m|s|ms|us|ns)",
	})
	ap.Add(&argparse.Argument{
		Switches: []string{"-v", "--verbose"},
		Help:     "Set verbose mode",
	})

	// A positional argument that requires one or more values
	ap.Add(&argparse.Argument{
		Name:        "names",
		Help:        "Some names passed into the program",
		NumArgsGlob: "+",
	})

	// On -h/--help this prints the help and exits; on a bad command-line it
	// prints the error and exits; otherwise it fills in opts and returns.
	ap.Parse()

	if opts.Verbose {
		fmt.Printf("Verbose mode is on!\n")
	}
	fmt.Printf("Count is %d\n", opts.Count)
	for i, name := range opts.Names {
		fmt.Printf("%d. %s\n", i+1, name)
	}
}
```

Running it with `-h` produces:

```
$ example1 -h
example1

This is an example program

usage: example1 [--count N] [--expiration EXPIRATION] [-v] [-h] names [...]

        --count=N                      How many items
        --expiration,-x=EXPIRATION     How long: #(h|m|s|ms|us|ns)
        -v,--verbose                   Set verbose mode
        -h,--help                      See this list of options
        names[ ... ]                   Some names passed into the program
```

(The exact column spacing adapts to your terminal width.)

The rest of this document explains each piece in detail. For more working
programs, see the [`examples/`](v2/examples) directory.

# Usage

## Single-level commands

Building the parser from the [Quick start](#quick-start) above, step by step:

1. Define the struct that will hold the values from the parse of the
   command-line.

```go
type MyOptions struct {
	Count      int
	Expiration time.Duration
	Verbose    bool
	Names      []string
}
```

2. Instantiate an Argparse object with a root Command object. This lets you
   describe your program, and points argparse to the value struct object.

```go
opts := &MyOptions{}
ap := argparse.New(&argparse.Command{
	Name:        "example1",
	Description: "This is an example program",
	Values:      opts,
})
```

3. Add options to the root Command object (via the Argparse object):

```go
// These are switch arguments
ap.Add(&argparse.Argument{
	Switches: []string{"--count"},
	MetaVar:  "N",
	Help:     "How many items",
})

ap.Add(&argparse.Argument{
	Switches: []string{"--expiration", "-x"},
	Help:     "How long: #(h|m|s|ms|us|ns)",
})

ap.Add(&argparse.Argument{
	Switches: []string{"-v", "--verbose"},
	Help:     "Set verbose mode",
})

// This is a positional argument
ap.Add(&argparse.Argument{
	Name: "names",
	Help: "Some names passed into the program",
	// We require one or more names
	NumArgsGlob: "+",
})
```

When each Argument is added, the argparse logic looks in the Command's Values
struct for a field name that matches either a "Switches" or "Name" value of the
argument. If it fails to find a matching field, the code will `panic()`. (This
is a programming error, caught the first time you run your program, not a
command-line error.)

4. Perform the parse.

```go
ap.Parse()
```

If the user requests help, the help text is shown and the program exits. If the
user gives an illegal command-line, the error message is shown and the program
exits. Otherwise, on success, the program continues to the next statement. See
[Parse vs. ParseAndExit](#parse-vs-parseandexit) for the details.

5. Use the values.

```go
if opts.Verbose {
	fmt.Printf("Verbose mode is on!\n")
}
fmt.Printf("Count is %d\n", opts.Count)
```

## Sub-commands

Sub-commands are also supported. You add a Command as a child of its parent
Command, and then can add arguments to that new Command.

```go
// Add "open" as a sub-command of the root parser "ap".
open_ap := ap.New(&argparse.Command{
	Name:        "open",
	Description: "Open something",
	Function:    DoOpen,
	Values:      &OpenOptions{},
})

open_ap.Add(&argparse.Argument{
	Switches: []string{"-r", "--reason"},
	Help:     "Why you are opening this",
})
```

With sub-commands, the `Function` field is a callback to your code, run when the
sub-command is chosen. The callback receives two arguments: a pointer to the
leaf `argparse.Command` object that was triggered by the command-line, and the
`Values` associated with that Command. Its type is:

```go
type ParserCallback func(*Command, Values) error
```

Since sub-commands have callback functions, it's usually better to perform the
parse with `ParseAndExit`. In that way, the program exits after completing the
sub-command's callback function.

```go
ap.ParseAndExit()
```

To use the Values, you will need to coerce them from the `argparse.Values`
interface to the actual struct-pointer that they are:

```go
func DoOpen(cmd *argparse.Command, values argparse.Values) error {
	opts := values.(*OpenOptions)
	// Now you can use opts
	return nil
}
```

## Default values and "Seen" arguments

Because you supply the struct that holds the values seen on the command-line,
you can initialize values to anything you want. Those are thus the default
values.

Sometimes you may also wish to know if the user actually set the value on the
command-line, even if the value is the same as the default value. You can learn
whether the user actually provided an option by checking the
`argparse.Command.Seen` map, which is filled in after parsing happens.

The function callbacks used for sub-commands provide the `argparse.Command` to
your code. But if you're using a single-level parser with no callback, the
Command object is the **Root** field of the ArgumentParser object.

The **Seen** map uses the name of the struct field as keys. For example, this
tells you if `--count` was given, by using the **Root** field of the
ArgumentParser object:

```go
if ap.Root.Seen["Count"] {
	// do something
}
```

But if this were code in a Function callback for a subcommand, it would simply
use the **Command** object passed to the function:

```go
func DoOpen(cmd *argparse.Command, values argparse.Values) error {
	opts := values.(*OpenOptions)
	// Now you can use opts

	if cmd.Seen["Count"] {
		// do something
	}
	return nil
}
```

# Parse vs. ParseAndExit vs. ParseArgs

`Parse()` and `ParseAndExit()` parse `os.Args`, handle `-h`/`--help`, and report
command-line errors. They differ in what happens after a successful parse:

| Situation | `Parse()` | `ParseAndExit()` |
|---|---|---|
| `-h` / `--help` given | print help to Stdout, exit 0 | print help to Stdout, exit 0 |
| invalid command-line | print error to Stderr, exit 1 | print error to Stderr, exit 1 |
| matched command has a `Function` | run it; if it returns an error, print it and exit 1; otherwise **return to the caller** | run it; if it returns an error, print it and exit 1; otherwise **exit 0** |
| matched command has no `Function` | **return to the caller** | print usage to Stderr, exit 1 |

Use `Parse()` for a single-level program where you inspect the values struct
after parsing. Use `ParseAndExit()` when your commands do their work in
`Function` callbacks; it guarantees that running without selecting a command
that does something is reported as an error.

If you don't want argparse to call `os.Exit` at all — for example to embed it in
a REPL or sub-shell, or to test your own command-line handling — use
`ParseArgs`, which parses an explicit argument slice and returns an error
instead of exiting:

```go
cmd, err := ap.ParseArgs([]string{"open", "--reason", "maintenance"})
if err == argparse.ErrHelp {
	// The user asked for help; the help text was already written to Stdout.
	return
} else if err != nil {
	fmt.Fprintln(os.Stderr, err)
	return
}
// cmd is the triggered Command, with its Values filled in. ParseArgs does not
// run any Function callback.
```

`ParseArgs` returns `argparse.ErrHelp` when help was requested (after writing the
help text to the parser's `Stdout`), the parse error on bad input, or `nil` on
success.

# Version

Set the `Version` field on the `ArgumentParser` to enable a version switch. When
it is set, `--version` prints the string and stops (the same way `-h` prints the
help): `Parse`/`ParseAndExit` exit 0, and `ParseArgs` returns `argparse.ErrVersion`
after writing the version to `Stdout`. If `Version` is empty, `--version` is not
treated specially. The version switch also appears in the `--help` output.

```go
ap := argparse.New(&argparse.Command{Name: "mytool", Values: opts})
ap.Version = "mytool 1.4.0"
// Optionally change which switches request it (default is just "--version"):
ap.VersionSwitches = []string{"-V", "--version"}
ap.ParseAndExit()
```

# Accepted command-line syntax

* A boolean switch takes no value: `-v`, `--verbose`.
* A switch that takes a value accepts either a space or an `=`:
  `--count 5`, `--count=5`, `-x 5m`, `-x=5m`.
* A switch with `NumArgs` greater than 1 consumes that many following values:
  `--coords 1 2 3`.
* Short-flag bundling (such as `-vx` for `-v -x`) is **not** supported; give
  each short switch separately.
* `--` marks the end of switches: every token after it is treated as a
  positional value, even if it begins with a hyphen.

Because a bare token that begins with `-` is treated as a switch, a *positional*
value that starts with a hyphen (a negative number, a path like `-foo`) must
come after `--`:

```
myprog -- -5 -file.txt
```

The value of a switch may begin with a hyphen without `--` (for example
`--count -5` or `--count=-5`), because argparse knows a value is expected there.

# Capturing leftover arguments

For wrapper programs that forward arguments to another command (think
`mytool run -- some-cmd --its-own-flags`), set `PassThrough: true` on a trailing
positional argument with a `[]string` destination. It captures the remaining
tokens **verbatim** — including ones beginning with `-` and any literal `--` —
without trying to interpret them as this program's switches.

```go
type Options struct {
	Verbose bool
	Cmd     []string // everything to forward
}

ap.Add(&argparse.Argument{Switches: []string{"-v", "--verbose"}})
ap.Add(&argparse.Argument{Name: "Cmd", PassThrough: true})
```

Capture begins at the first of:

* a `--` token (which is consumed as the separator) — this works even when there
  are no other positional arguments; or
* the first token that would be assigned to the pass-through argument (after any
  preceding switches and positionals have been handled).

Once capture begins, nothing further is interpreted: switches that happen to
match this program's own (and additional `--` tokens) are stored as-is. The
program's switches and any preceding positionals before capture begins are
parsed normally.

Using the `Options` above (a `-v`/`--verbose` switch and a `Cmd` pass-through),
here is how different numbers of `--` behave:

```
# No "--": mytool's own switches are parsed; capture starts at the first
# non-switch token ("build"), and everything from there is taken raw.
$ mytool -v build --release
Verbose=true  Cmd=["build" "--release"]

# One "--": needed here because the first forwarded token begins with "-".
# Without the "--", argparse would treat "--release" as one of mytool's
# switches and fail. The "--" is consumed as the separator.
$ mytool -- --release build
Verbose=false  Cmd=["--release" "build"]

# One "--": a forwarded flag that collides with mytool's own (--verbose) is
# still captured raw, because it comes after capture has begun.
$ mytool --verbose -- inner --verbose
Verbose=true  Cmd=["inner" "--verbose"]

# Two "--": the first is the separator; the second is inside the captured
# remainder and is kept literally.
$ mytool -- build -- --release
Verbose=false  Cmd=["build" "--" "--release"]
```

So the *first* `--` (before capture begins) is a separator and is dropped; any
later `--` is part of the forwarded arguments and is kept. Use a leading `--`
when the first forwarded token begins with a hyphen and there is no preceding
positional. A pass-through argument must be the last positional and have a
`[]string` destination.

# Values struct and field names

The Values struct needs to have field names that match either the short name or
the long name of the Arguments added to a Command. Because the field name needs
to be inspected by argparse, it must start with an upper-case character (so that
Go exports the field). Also, any embedded dashes are removed and the field name
is expected to be in CamelCase. For example:

* `-s` : the field name is `S`
* `--input` : the field name is `Input`
* `--no-verify` : the field name is `NoVerify`

If the deduced name doesn't match your field, set the Argument's `Dest` field to
the exact field name instead.

The fields for switch or positional arguments can be of these scalar types:

* **bool** — For a switch, if the switch is present, the value is set to true.
* **string**
* **float64**, **float32**
* **int**, **int64**, **int8**, **int16**, **int32**
* **uint**, **uint64**, **uint8**, **uint16**, **uint32** — same decimal/hex/octal
  forms as the signed integers, but negative values are rejected
* **time.Duration** — parsed by `time.ParseDuration()`

Or they can be the following slice types. A slice value for a switch argument
means the switch may be given more than once on the command-line. A slice value
for a positional argument means the positional argument can appear more than
once. In that case, **NumArgs** or **NumArgsGlob** should be used to define how
many times it must appear; otherwise argparse still accepts it only once.

* **[]bool**
* **[]string**
* **[]float64**, **[]float32**
* **[]int**, **[]int64**, **[]int8**, **[]int16**, **[]int32** — these can be
  given in decimal, in hex if they start with `0x` (as in `0xff`), or in octal
  if they start with `0o` or just `0`.
* **[]uint**, **[]uint64**, **[]uint8**, **[]uint16**, **[]uint32**
* **[]time.Duration** — each value is parsed by `time.ParseDuration()`

## Custom value types

A field of any type that implements
[`encoding.TextUnmarshaler`](https://pkg.go.dev/encoding#TextUnmarshaler) is
parsed by calling its `UnmarshalText` method, so you can accept IP addresses,
URLs, enums, `uint`s, or any type of your own. A slice of such a type works too
(one element parsed per value). Many standard-library types already qualify
(for example `net.IP` and `time.Time`).

```go
type Options struct {
	Listen net.IP   // --listen 10.0.0.1
	Hosts  []net.IP // --host 10.0.0.1 --host 10.0.0.2
}
```

For your own types, define `UnmarshalText` (a pointer receiver) and return an
error to reject invalid input:

```go
type Level int

func (l *Level) UnmarshalText(text []byte) error {
	switch string(text) {
	case "low":
		*l = 0
	case "high":
		*l = 1
	default:
		return fmt.Errorf("unknown level %q", text)
	}
	return nil
}
```

`Choices` is not supported for custom types (setting it panics); enforce the
allowed set inside `UnmarshalText` instead.

# Command

A `Command` describes the program (the root command) or a sub-command. The
fields you set are:

* **Name**: The name of the program or sub-command. For the root command, if
  left empty it defaults to `os.Args[0]`. For a sub-command, this is the word
  the user types to select it.
* **Description**: A description shown after the command name and before the
  options in `--help`. Can be multi-line.
* **Epilog**: Text shown after all the options in `--help`. Can be multi-line.
* **Values**: A pointer to the struct that will receive the parsed values.
* **Function**: The callback to run when this command is selected. See
  [Sub-commands](#sub-commands).

After parsing, **Seen** and **CommandSeen** maps on the Command record which
arguments and which sub-commands were present on the command-line.

# Argument

The following fields can be set in an `Argument`:

* **Switches**: (optional) All the accepted versions for this switch. Each one
  must start with at least one hyphen. Give this for switch arguments.

* **Name**: (optional) For positional arguments (those that don't start with a
  hyphen), the name of the argument. While the user does not type this name on
  the command-line, it is used in the help statement produced by argparse. Give
  either `Switches` or `Name`, not both.

* **Dest**: The name of the field in the Values struct where the value will be
  stored. This is only needed if you wish to override the
  [deduced field name](#values-struct-and-field-names).

* **Help**: The help string to display to the user. Can be multi-line.

* **MetaVar**: The text to use as the name of the value in the `--help` output.

* **NumArgs**: (optional) For positional and switch arguments, the exact number
  of values that must follow. See [NumArgs and NumArgsGlob](#numargs-and-numargsglob).

* **NumArgsGlob**: (optional) For positional arguments only, a pattern
  (`"*"`, `"+"`, or `"?"`) describing how many values may or must be given. See
  [NumArgs and NumArgsGlob](#numargs-and-numargsglob).

* **PassThrough**: (positional arguments only) If true, this argument captures
  the rest of the command-line verbatim. It must be the last positional and have
  a `[]string` destination. See
  [Capturing leftover arguments](#capturing-leftover-arguments).

* **Required**: (switch arguments only) If true, the switch must be given on the
  command-line; otherwise a parse error is returned. The requiredness of
  positional arguments is controlled by `NumArgs`/`NumArgsGlob` instead, so
  setting `Required` on a positional argument panics.

* **Inherit**: If true, then all sub-commands of this Command will automatically
  inherit a copy of this Argument. This also means the Values struct of each
  sub-command must have a field whose name and type work for this Argument; if
  not, the `New()` that adds the sub-command will panic. If you `Add()` an
  inherited argument after already adding a sub-command with `New()`, that
  `Add()` will panic.

* **Choices**: (optional) A slice (even when the field value is a scalar) listing
  the only possible values for the argument. If a user gives a value not in this
  list, an error is returned to the user. The slice type must match the value
  type for this Argument: `[]bool`, `[]string`, `[]int`, or `[]float64`. The
  choices are listed in the `--help` output.

The `--help` output also shows an argument's default value (the value held by its
Values-struct field when the argument is added), unless that default is the
type's zero value.

# NumArgs and NumArgsGlob

These two fields control how many values an argument consumes.

**NumArgs** is an exact count. It applies to both switch and positional
arguments:

* If neither `NumArgs` nor `NumArgsGlob` is given, `NumArgs` defaults to 1 (0 for
  a boolean switch, which takes no value).
* `NumArgs` greater than 1 means exactly that many values must follow; giving
  fewer is a command-line error. The destination field must be a slice.

**NumArgsGlob** is a pattern, allowed for **positional arguments only**:

* `"*"` — zero or more
* `"+"` — one or more
* `"?"` — zero or one

Rules:

* `NumArgsGlob` is not allowed for switch arguments.
* A `"?"` positional argument may be followed by other positional arguments.
* A `"*"` or `"+"` positional argument may not be followed by any other
  positional argument, because it consumes an unlimited number of values.
* If `"*"` or `"+"` is used, the destination field must be a slice.

# Mutually-exclusive groups

Use `AddMutuallyExclusive` to declare a set of switch arguments of which at most
one may be given. It adds the arguments to the Command and records the group; if
`required` is true, exactly one of them must be given.

```go
ap.AddMutuallyExclusive(false,
	&argparse.Argument{Switches: []string{"--json"}, Help: "JSON output"},
	&argparse.Argument{Switches: []string{"--yaml"}, Help: "YAML output"},
)
```

Giving more than one member is a parse error ("Only one of ... may be given");
if the group is `required` and none is given, that is also an error ("One of ...
is required"). Only switch arguments may be grouped, and the members must not
set `Required` individually (the group's `required` flag covers that). On a
sub-command, call the method on the sub-command's `*Command`.

# Inheritance by sub-commands

It's often the case that some arguments can be given at any level in the command
hierarchy. For example, a `-v`/`--verbose` option could be given for a root
command or a sub-command. Instead of defining the same argument at each level,
argparse lets you define it at the root level and have the sub-commands inherit
it.

You will probably also want your Values structs to inherit the values through
Go's composition. The
[`examples/twolevels_with_defaults`](v2/examples/twolevels_with_defaults)
example shows this. The root options have Verbose and Debug, and the
OpenOptions and CloseOptions inherit them through composition.

```go
type RootOptions struct {
	Debug   bool
	Verbose bool
	Reason  string
}

type OpenOptions struct {
	RootOptions

	Name string
}

type CloseOptions struct {
	RootOptions

	Name string
}
```

The root command parser defines the arguments and sets their `Inherit` flag to
true:

```go
ap := argparse.New(&argparse.Command{
	Description: "This is an example program",
	Values:      opts,
})

ap.Add(&argparse.Argument{
	Switches: []string{"--debug"},
	Help:     "Set debug mode",
	Inherit:  true,
})

ap.Add(&argparse.Argument{
	Switches: []string{"-v", "--verbose"},
	Help:     "Set verbose mode",
	Inherit:  true,
})
```

After that, the parsers for the open and close sub-commands do not need to
define the arguments. Because the parent command (the root) had these arguments
with `Inherit = true`, and because the Values struct for open and close have
destination fields for Verbose and Debug, argparse copies the argument
definitions from the root command to the open and close commands.

# Translation

The `argparse.New()` function returns an `ArgumentParser` whose following fields
can be changed (for example, for internationalization):

* **Stdout** — an `io.Writer` to send the output of `--help` to instead of
  `os.Stdout`.

* **Stderr** — an `io.Writer` to send error messages to instead of `os.Stderr`.

* **HelpSwitches** — the switches that argparse interprets as a request for help.
  The default is `[]string{"-h", "--help"}`.

* **Messages** — a `Messages` struct holding every string that argparse prints to
  the user: the help-output headers and all command-line error messages. Override
  it to translate that text.

```go
ap := argparse.New(&argparse.Command{
	Description: "This is an example program",
	Values:      opts,
})
ap.HelpSwitches = []string{"-h", "--help", "--aide"}
```

Argparse ships with built-in `Messages` values you can assign directly:

* `argparse.DefaultMessages_en` — English (the default).
* `argparse.DefaultMessages_ko` — Korean.

```go
ap := argparse.New(&argparse.Command{
	Description: "이것은 예제 프로그램입니다",
	Values:      opts,
})
ap.Messages = argparse.DefaultMessages_ko
```

To provide your own translation, copy `DefaultMessages_en`, then change the
fields you want:

```go
msgs := argparse.DefaultMessages_en
msgs.NoSuchSwitchFmt = "No existe la opción: %s"
msgs.HelpDescription = "Ver esta lista de opciones"
ap.Messages = msgs
```

Note: fields whose name ends in `Fmt` are format strings for `fmt.Sprintf` /
`fmt.Errorf`. A translation must keep the same verbs (`%s`, `%v`, `%w`) so the
argument values are still substituted correctly.

Only text shown to the *user* of your program is translatable. Errors caused by
misusing the argparse API itself (such as adding an argument with no matching
struct field) are reported via `panic()` and are intentionally not part of
`Messages`.

# Examples

For working examples, see the [`examples/`](v2/examples) directory in the source
code:

* [`example1`](v2/examples/example1) — a single-level parser.
* [`onelevel`](v2/examples/onelevel) — another single-level parser.
* [`twolevels`](v2/examples/twolevels) — a parser with sub-commands.
* [`twolevels_with_defaults`](v2/examples/twolevels_with_defaults) — sub-commands
  with inherited arguments and composed option structs.

# License

Argparse is distributed under the MIT License. See the [LICENSE](LICENSE) file
for details.
