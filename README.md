# argparse
A Go module to parse CLI argument, that loosely follows the model of the Python
argparse module.

Highlights:

* You can have nested subcommands.
* The values for a parser are stored in a struct.
* Argparse can deduce the name of the value field in the struct by looking
	at the name of the option.
* Argparse fills in the struct with values from the command-line, and also
	tells you all the options that were seen.

See [the GoDoc documentation for argparse](https://godoc.org/github.com/gilramir/argparse)

# Usage

1. Define the struct that will hold the values from the parse of the command-line.

	type MyOptions struct {
		Debug bool
		Verbose	bool
		Names []string
	}

2. Instantiate an Argparse object with a root Command object. This lets
you describe your program, and points argparse to the value struct object.

	opts := &MyOptions{}
	ap := argparse.New(&argparse.Command{
		Description:	"This is an example program",
		Values:		opts,
	})

3. Add options to root Command object (via the Argparse object):

	ap.Add(&argparse.Argument{
		Switches:	[]string{"--debug"},
		Help:		"Set debug mode",
	})

	ap.Add(&argparse.Argument{
		Switches:	[]string{"-v", "--verbose"},
		Help:		"Set verbose mode",
	})

	ap.Add(&argparse.String{
		Name:		"names",
		Help:		"Some names passed into the program",
		NumArgsGlob:	'+',
	})

When each Argument is added, the argparse logic looks in the Command's
Value struct for a field name that matches either "switch" or "name" of the
argument.  If it fails to find a matching field, the code will panic().

4. Finally, perform the parse.

	ap.Parse()

If the user requests help, the help text will be given, and the program exits.
If the users gives an illegal command-line, the error message is shown, and the
program exits. Otherwise, on success, the program continues to the next statement.

# ArgumentParser

* Stdout - an io.Writer to send the output of "--help" to instead of os.Stdout.

* Stderr - an io.Writer to send error messages to instead of os.Stderr.

* Messages - a struct of all the messages that argparse can print to users.
	You can override this to provide translations. The default is the built-in
	English version of these messages. This is still a work in progress.

* HelpSwitches - the switches that arpgarse will interpret as a request for help.
	The default is: []string{"-h", "--help"}

* Root  - a pointer to the root Command object. This is set by argparse for you,
	for convenience.

# Command

The following fields can be set in a Command:

* Name - (optional) the name of the program

* Description - (optional) a short description of the program

* Epilog - (optional) this can be a multi-line string. It is shown after all
    the options in the "--help" output

* Values - the struct that receives the values after parsing

# Values struct

The Values struct needs to have field names that match either the short
name or the long name of the Arguments added to a Command.  Because the
field name needs to be inspected by "argparse", it must start with an upper case
character (so that Go exports those field names to other modules). Also, any embedded
dashes are removed and the field name is expected to be in CamelCase. For example

* "-s" : the field name is S

* "--input": the field name is Input

* "--no-verify": the field name is NoVerify

The fields for switch or positional arguments can be of the scalar types:

* bool - For a switch, if the switch is present, the value is set to true.

* string

* float64

* int

Or they can be the following sclice types. A slice indicates a switch is accepted
more than one, or a positional argument can be appear more than once.

* []bool

* []string

* []float64

* []int

# Argument

The following fields can be set in Argument:

* Switches: (optional) All the accepted versions for this switch. Each one must start
	with at least one hyphen.

* Name: (optional) For positional arguments (after all switch arguments, which start with dashes), the name of
    the argument. While this name is not used by the user when giving the CLI string,
    it is used in the help statement produced by argparse.

* Dest: The name of the field in the Destination where the value will be stored.
    This is only needed if you wish to override the default.

* Description: A description of the argument. Can be multi-line.

* MetaVar: The text to use as the name of the value in the --help output.

* NumArgs: (optional) For positional and switch arguments, speficies how many
	arguments _must_ follow the option.

* NumArgsGlob: (optional) For positional arguments only, a string that specifies
how many values can or must be provided:

    * "\*" - zero or more

    * "+" - one or more

    * "?" - one or zero

    This is not allowed for switch arguments. If neither NumArgs nor NumArgsGlob is given,
    then NumArgs is set to 1.

* Inherit: If true, then all sub-commands of this Command will automatically inherit a copy
	of this Argument. This also means that the Value struct must have a field whose name
	and type work for this Argument. If that is not true, then the New() which adds the
	sub-command will panic. if you Add() a new Inherited argument after already adding
	a sub-command with New(), then the Add() will panic.

* Choices: (optional) A slice (even when the field value is an int) which lists the only
    possible values for the argument value. If a user gives a value that is not in this list,
    an error will be returned to the user. The slice type must match the Value type for
    this Argument: []bool, []string, []int, or []float64

# Examples

For working examples, see the examples/ directory in the source code.
