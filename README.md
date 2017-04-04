# argparse
CLI argument parsing for Go, that follows the model of the Python argparse module.

Your root parser can have child parsers, for sub-commands. Each parser is associated
with a struct of your design which implements the argparse.Destination interface.
The "destination" must have field names which correspond to the short or long
versions of arguments added to a parser, and, because it implements the
argparse.Destination interface, it must have a Run() method, which is executed
when that root or sub-command parser is the one triggered by the input.

# Usage

Instantiate the root argument parser with argparse.ArgumentParser. Every ArgumentParser
needs a "destination struct" which has field names which are similar to the
argument names that will be added to that parser.

    include "github.com/gilramir/argparse"

    type MyProgramOptions struct {
        Input   string
        Output  string
    }

    func (self *MyProgramOptions) Run(parents []Destination) (cliOk bool, err error) {
        cliOk = true
        err = dosomething(self.Input, self.Output)
    }
        

    func example() (error) {
        p := &argparse.ArgumentParser{
            Name: "my_program",
            ShortDescription: "A utility program",
        }

        p.AddArgument(&argparse.Argument{
            Type: "",
            Short: "-i",
            Long: "--input",
            Description: "The input file",
            Metvar: "FILE",
        })

        p.AddArgument(&argparse.Argument{
            Type: "",
            Short: "-o",
            Long: "--output",
            Description: "The output file",
            Metvar: "FILE",
        })

        err := p.ParseArgs()
        return err
    }

When each Argument is added, the argparse logic looks for a name in the Destination
struct associated with the ArgumentParser that matches either the short name or
the long name of the argument.  If it fails to find a matching field, the code
will panic().
