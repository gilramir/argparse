# argparse
CLI argument parsing for Go, that follows the model of the Python argparse module.

You create a parser object and add arguments to it. The parser object is associated
with a struct of your own design, and each argument is associated with a field
in that struct. The struct holds the values from the command-line after the
parser object is done parsing them.

A parser object chan have children parser objects. This is how sub-commands are
implemented.

See [the GoDoc documentation for argparse](https://godoc.org/github.com/gilramir/argparse)

# Usage

Instantiate the parser object with argparse.ArgumentParser. Every ArgumentParser
needs a "destination struct" which has field names which are similar to the
argument names that will be added to that parser. The argparse.Destination
interface requires a "Run" method, which is triggered when the parse is finished.

    include "github.com/gilramir/argparse"

    type MyProgramOptions struct {
        Input   string
        Output  string
    }

    func (self *MyProgramOptions) Run(parents []Destination) (error) {
        return doSomething(self.Input, self.Output)
    }

    func example() (error) {
        p := &argparse.ArgumentParser{
            Name: "my_program",
            ShortDescription: "A utility program",
            Destination: &MyProgramOptions{},
        }

        p.AddArgument(&argparse.Argument{
            Short: "-i",
            Long: "--input",
            Description: "The input file",
            Metavar: "FILE",
        })

        p.AddArgument(&argparse.Argument{
            Short: "-o",
            Long: "--output",
            Description: "The output file",
            Metavar: "FILE",
        })

        err := p.ParseArgs()
        return err
    }

When each Argument is added, the argparse logic looks for a name in the Destination
struct associated with the ArgumentParser that matches either the short name or
the long name of the argument.  If it fails to find a matching field, the code
will panic().

# ArgumentParser

The following fields can be set in ArgumentParser:

* Name - the name of the program

* ShortDescription - a one-line description of the program

* LongDescription - (optional) this can be a multi-line description of the program

* Epilog - (optional) this can be a multi-line string. It is shown after all
    the options in the "--help" output

* Destination - the struct that receives the values after parsing

* Stdout - an io.Writer to send the output of "--help" to instead of os.Stdout

# Destination interface

The Destination interface needs to have field names that match either the short
name or the long name of the Arguments added to an ArgumentParser.  Because the
field name needs to be inspected by "argparse", it must start with an upper case
character (so that Go exports those field names to other modules). Also, any embedded
dashes are removed and the field name is expected to be in CamelCase. For example

* "-s" : the field name is S

* "--input": the field name is Input

* "--no-verify": the field name is NoVerify

The fields for switch or positional arguments can be of the following types:

* bool - For a switch, if the switch is present, the value is set to true.

* string

* int

* []string - indicates a switch is accepted more than one, or a positional argument can be
    appear more than once.

The Destination interface requires a Run() method. It receives a slice of
Destination objects, which are the Destinations for any parent ArgumentParsers. If there
is only one ArgumentParser object (no sub-commands), then this slice will be empty.

    func (self *MyProgramOptions) Run(parents []Destination) (error) {
        err := dosomething(self.Input, self.Output)
        return err
    }

The return value is an error. The error is simply passed back as the return value of
Argument.ParseArgs() or Argument.ParseArgv(). There is a special error type named
argparse.ParseErr, which if returned, causes ArgumentParser to print the usage statement to
stderr before returning the error back to the caller. You can create a ParseErr
with the helper functions argparse.ParseError() or argparse.ParseErrorf().

# Argument

The following fields can be set in Argument:

* Short: (optional) The short (one-dash) version of this argument. You must supply the "-".

* Long: (optional) The long (two-dash) version of this argument. You must supply the "--".

* Name: (optional) For positional arguments (after all switch arguments, which start with dashes), the name of
    the argument. While this name is not used by the user when giving the CLI string,
    it is used in the help statement produced by argparse.

* Dest: The name of the field in the Destination where the value will be stored.
    This is only needed if you wish to override the default.

* Description: A description of the argument. Can be multi-line.

* Metavar: The text to use as the name of the value in the --help output.

* NumArgs: (optional) For positional arguments, a rune that specifies how many values can or must be
    provided. If not given, then only one value can be given:

    * '\*' - zero or more

    * '+' - one or more

    * '?' - one or zero

    This is not used for switch arguments. If the type of the value of a switch argument is a slice,
    then argparse allows the argument to occur more than once, and saves each value.

# ParseCommands

ParseCommands are special types of Arguments that tell the ArgumentParser to change behavior while it's
parsing. If an Argument is a ParseCommand, then instead of Short, Long, or Name, the String field is used
to denote what to expect on the command-line.

The only ParseCommand available right now is *PassThrough*. This tells argparse that all the remaining
arguments on the command-line should be added to the slice that "Dest" names.

You can AddArgument the ParseCommand Argument in any order. It does not have to be added via AddArgument
in any order in relation to any switch arguments or positional arguments.

See [ex5.go](examples/ex5.go)

    type Options struct {
            Filename        []string
            OtherArguments  []string
    }

    func (self *Options) Run([]argparse.Destination) error {
            fmt.Printf("Filenames: %v\n", self.Filename)
            fmt.Printf("Other Arguments: %v\n", self.OtherArguments)
            return nil
    }

    func main() {
            p := &argparse.ArgumentParser{
                    Name:             "my_program",
                    ShortDescription: "This program takes positional arguments",
                    Destination:      &Options{},
            }

            p.AddArgument(&argparse.Argument{
                    Name: "filename",
                    Help: "The file to look at. Can be given more than once.",
            })

            p.AddArgument(&argparse.Argument{
                    ParseCommand: argparse.PassThrough,
                    String:        "--",
                    Dest:         "OtherArguments",
            })

            err := p.ParseArgs()
            if err != nil {
                    fmt.Fprint(os.Stderr, err)
            }
    }



# Tutorial and Examples

## How to structure your CLI program

A good way to stay organized is to have a very simple main() function, and
a subdirectory of all the source files related to the command-line, then a sibling
subdirectory for non-CLI (non-UI) logic that the CLI code should call.

    PROJECT/
        main.go
        cmd/
            root.go
            subcommand1.go
            subcommand2.go
        lib/
            real_logic1.go
            real_logic2.go

The main.go is as simple as:

    package main

    import "MY_URL/PROJECT/cmd"

    func main() {
        cmd.Execute()
    }

If you have sub-commands, can define your ArgumentParsers with global-scope variables:

    import "github.com/gilramir/argparse"

    var rootParser = &argparse.ArgumentParser{
        .....
    }

The Execute() function in "cmd/" is as simple as:

    func Execute() {
        err := rootParser.ParseArgs()
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
    }

## Create a CLI that accepts no options

See [ex1.go](examples/ex1.go)

    func main() () {
        p := &argparse.ArgumentParser{
            Name: "my_program",
            ShortDescription: "A utility program",
        }
        p.ParseArgs()
    }

## Create a CLI with two positional arguments

See [ex2.go](examples/ex2.go)

    type Options struct {
        Pattern     string
        Filenames   []string
    }

    func (self *Options) Run([]argparse.Destination) (error) {
            fmt.Printf("Pattern: %s\n", self.Pattern)
            fmt.Printf("Filenames: %v\n", self.Filenames)
            return nil
    }

    func main() () {
        p := &argparse.ArgumentParser{
            Name: "my_program",
            ShortDescription: "This program takes positional arguments",
            Destination: &Options{},
        }

        p.AddArgument(&argparse.Argument{
            Name: "pattern",
            Help: "The pattern to look for",
        })

        p.AddArgument(&argparse.Argument{
            Name: "filenames",
            Help: "The file(s) to look at",
            NumArgs: '+',
        })

        err := p.ParseArgs()
        if err != nil {
            fmt.Fprint(os.Stderr, err)
        }
    }

## Return a ParseErr, indicating a CLI problem

See [ex3.go](examples/ex3.go)

    type Options struct {
            Filenames []string
    }

    func (self *Options) Run([]argparse.Destination) error {
            return argparse.ParseError("The CLI syntax is bad")
    }

    func main() {
            p := &argparse.ArgumentParser{
                    Name:             "my_program",
                    ShortDescription: "This program takes positional arguments",
                    Destination: &Options{},
            }

            p.AddArgument(&argparse.Argument{
                    Name:    "filenames",
                    Help:    "The file(s) to look at",
                    NumArgs: '+',
            })

            err := p.ParseArgs()
            if err != nil {
                    fmt.Fprint(os.Stderr, err)
            }
            fmt.Printf("Done.\n")
    }

## Create a CLI with two switch arguments

See [ex4.go](examples/ex4.go)

    type Options struct {
            Pattern     string
            Filename    []string
    }

    func (self *Options) Run([]argparse.Destination) error {
            fmt.Printf("Pattern: %s\n", self.Pattern)
            fmt.Printf("Filenames: %v\n", self.Filename)
            return nil
    }

    func main() {
            p := &argparse.ArgumentParser{
                    Name:             "my_program",
                    ShortDescription: "This program takes positional arguments",
                    Destination:      &Options{},
            }

            p.AddArgument(&argparse.Argument{
                    Short: "-p",
                    Long: "--pattern",
                    Help: "The pattern to look for",
            })

            p.AddArgument(&argparse.Argument{
                    Short: "-f",
                    Long: "--filename",
                    Help: "The file to look at. Can be given more than once.",
            })

            err := p.ParseArgs()
            if err != nil {
                    fmt.Fprint(os.Stderr, err)
            }
    }
