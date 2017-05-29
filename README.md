# argparse
CLI argument parsing for Go, that follows the model of the Python argparse module.

You create a parser object and add arguments to it. The parser object is associated
with a struct of your own design, and each argument is associated with a field
in that struct. The struct holds the values from the command-line after the
parser object is done parsing them.

A parser object chan have children parser objects. This is how sub-commands are
implemented.

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
        err := dosomething(self.Input, self.Output)
        return err
    }
       
    var options MyProgramOptions

    func example() (error) {
        p := &argparse.ArgumentParser{
            Name: "my_program",
            ShortDescription: "A utility program",
            Destination: &options,
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

* Short: (optional) The sort (one-dash) version of this argument. You must supply the "-".

* Long: (optional) The long (two-dash) version of this argument. You must supply the "--".

* Name: (optional) For positional arguments (after all "dash" arguments), the name of
    the argument. While this name is not used by the user when giving the CLI string,
    it is used in the help statement produced by argparse.

* Description: A description of the argument. Can be multi-line.

* Metavar: The text to use as the name of the value in the --help output.

* NumArgs: (optional) For positional arguments, a rune that specifies how many values can or must be
    provided. If not given, then only one value can be given:

    * '\*' - zero or more

    * '+' - one or more

    * '?' - one or zero

    This is not used for "dash" arguments. If the type of the value of a dash argument is a slice,
    then argparse allows the argument to occur more than once, and saves each value.


# Tutorial and Examples

## Create a CLI that accepts no options

See [ex1.go](examples/ex1.go)

    func main() () {
        p := &argparse.ArgumentParser{
            Name: "my_program",
            ShortDescription: "A utility program",
        }
        p.ParseArgs()
    }

## Create a CLI with an option and a positional argument

See [ex2.go](examples/ex2.go)

    type Options struct {
        Pattern     string
        Filenames   []string
    }

    func (self *Options) Run([]argparse.Destination) (error) {
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
