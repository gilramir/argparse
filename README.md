# argparse
Argparse is a Go module that helps you parse command-line arguments.
It loosely follows the conceptual model of the Python argparse module.

Highlights:

* You can have nested subcommands.
* The values for the command-line options are stored in a struct of your
  creation.
* Argparse can deduce the name of the value field in the struct by looking
        at the name of the option. Or, you can tell it exactly which field to
        use.
* Argparse will tell you if a particular option was present on the command-line
        or not present, in case you need that information.
* Options can be inherited by sub-comands, and you need only define them
        once.

See [the GoDoc documentation for argparse](https://godoc.org/github.com/gilramir/argparse/v2)

# Usage

1. Define the struct that will hold the values from the parse of the command-line.

        type MyOptions struct {
                Count       int
                Expiration  time.Duration
                Verbose     bool
                Names       []string
        }

2. Instantiate an Argparse object with a root Command object. This lets
you describe your program, and points argparse to the value struct object.

        opts := &MyOptions{}
        ap := argparse.New(&argparse.Command{
                Description:    "This is an example program",
                Values:         opts,
        })

3. Add options to root Command object (via the Argparse object):

        // These are switch arguments
        ap.Add(&argparse.Argument{
                Switches:       []string{"--count"},
                MetaVar:        "N",
                Help:           "How many items",
        })

        ap.Add(&argparse.Argument{
                Switches:       []string{"--expiration", "-x"},
                Help:           "How long: #(h|m|s|ms|us|ns)",
        })

        ap.Add(&argparse.Argument{
                Switches:       []string{"-v", "--verbose"},
                Help:           "Set verbose mode",
        })

        // This is a positional argument
        ap.Add(&argparse.Argument{
                Name:           "names",
                Help:           "Some names passed into the program",
                // We require one or more names
                NumArgsGlob:    '+',
        })

When each Argument is added, the argparse logic looks in the Command's
Value struct for a field name that matches either a "Switches" or "Name" value
of the argument.  If it fails to find a matching field, the code will panic().

4. Perform the parse.

        ap.Parse()

If the user requests help, the help text will be given, and the program exits.
If the users gives an illegal command-line, the error message is shown, and the
program exits. Otherwise, on success, the program continues to the next statement.

5. Use the values.

        if opts.Verbose {
                // do something


## Default values

Because you supply the struct that will be used to hold the values seen on the
command-line, you can set the initial values to anything you want. Those are
thus the default values.

You can know if the user actually provided an option on the command-line by
checking the argparse.Command.Seen map, which is filled in after the parsing
happens. The argparse object's command object, the root of the command tree,
can be accessed by the **Root** field.

The **Seen** map uses the name of the struct field as keys.
For example, this tells you if "--count" was ggiven:

        if ap.Root.Seen["Count"]

## Sub-commands

Sub-commands are also supported. You add a Command as a child of its parent
Command, and then can add arguments to that new Command

        // Add "open" as a sub-command ot the root parser "ap".
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

With sub-commands, the Function argument is a callback to your code, to
run when the sub-command is chosen. The callback accepts two arguments:
a pointer to leaf argparse.Command object that was triggered by the
command-line, and the Values associated with that Command.

Since sub-commands have callback functions, it's usually better to
perform the parse with ParseAndExit. In that way, the program exits
after completing the sub-commands callback function.

        ap.ParseAndExit()

To use the Values, you will need to coerce them from the argparse.Values
interface to the actual struct-pointer that they are:

    func DoOpen(cmd *argparse.Command, values argparse.Values) error {
        opts := values.(*OpenOptions)
    }

This of course gives you the values of the arguments, be they default values
or values given by the user. The Command object has a Seen map which
tells you if a command was actually given by the user (if it was "seen" by the
parser). You give it the string name of the field in the Values struct:

        cmd.Seen["Verbose"]


# Translation

Once you create your ArgumentParser object with the argparse.New() function:

        ap := argparse.New(&argparse.Command{
                Description:    "This is an example program",
                Values:         opts,
        })

The following fields can bet changed:

* **Stdout** - an io.Writer to send the output of "--help" to instead of os.Stdout.

* **Stderr** - an io.Writer to send error messages to instead of os.Stderr.

* **HelpSwitches** - the switches that arpgarse will interpret as a request for help.
        The default is: []string{"-h", "--help"}

* **Messages** - a struct of all the messages that argparse can print to users.
        You can override this to provide translations. The default is the built-in
        English version of these messages. Not all strings are supported
        via this mechnism; it's still a work in progress.


# Values struct and field names

The Values struct needs to have field names that match either the short
name or the long name of the Arguments added to a Command.  Because the
field name needs to be inspected by "argparse", it must start with an upper case
character (so that Go exports those field names to other modules). Also, any embedded
dashes are removed and the field name is expected to be in CamelCase. For example

* "-s" : the field name is S

* "--input": the field name is Input

* "--no-verify": the field name is NoVerify

The fields for switch or positional arguments can be of the scalar types:

* **bool** - For a switch, if the switch is present, the value is set to true.

* **string**

* **float64**

* **int**, **int64**

* **time.Duration** - parsed by time.ParseDuration()

Or they can be the following slice types. A slice value for a switch argument
indicates is accepted more than once on the command-line. A slice value
for a positional argument means the positional argument can be appear more than once.
In that case, **NumArgs** or **NumArgsGlob**, should be used to define how
many times it must appear, otherwise argparse still accepts it only once.

* **[]bool**

* **[]string**

* **[]float64**

* **[]int**, **[]int64**

* **[]time.Duration** - each time.Duration is parsed by time.ParseDuration()

# Argument

The following fields can be set in Argument:

* **Switches**: (optional) All the accepted versions for this switch. Each one must start
        with at least one hyphen.

* **Name**: (optional) For positional arguments (after all switch arguments, which start with dashes), the name of
    the argument. While this name is not used by the user when giving the CLI string,
    it is used in the help statement produced by argparse.

* **Dest**: The name of the field in the Destination where the value will be stored.
    This is only needed if you wish to override the default.

* **Description**: A description of the argument. Can be multi-line.

* **MetaVar**: The text to use as the name of the value in the --help output.

* **NumArgs**: (optional) For positional and switch arguments, specifies how many
        arguments _must_ follow the option.

* **NumArgsGlob**: (optional) For positional arguments only, a string that specifies
how many values can or must be provided:

    * "\*" - zero or more

    * "+" - one or more

    * "?" - one or zero

    This is not allowed for switch arguments. If neither NumArgs nor NumArgsGlob is given,
    then NumArgs is set to 1.

* **Inherit**: If true, then all sub-commands of this Command will automatically inherit a copy
        of this Argument. This also means that the Value struct must have a field whose name
        and type work for this Argument. If that is not true, then the New() which adds the
        sub-command will panic. if you Add() a new Inherited argument after already adding
        a sub-command with New(), then the Add() will panic.

* **Choices**: (optional) A slice (even when the field value is an int) which lists the only
    possible values for the argument value. If a user gives a value that is not in this list,
    an error will be returned to the user. The slice type must match the Value type for
    this Argument: []bool, []string, []int, or []float64

* **Function**: If this is not nil, then if this is the "triggered" command or sub-command,
        then this function is called. The type is:

        type ParserCallback func (Values) error

# Inheritance by Sub-commands

It's often the case that some arguments can be given at any level in the
command-hierarchy. For example, a "-v"/"--verbose" option could be given for a
root command, or a sub-command. Instead of having to define the same argument
at each level, argparse lets you define it at the root level and let the
sub-command inherit them.

You will probably also want to have your Values structs inherit the values
through Go's composition.

The examples/two_levels_with_defaults example shows this.  The root options
have Verbose and Debug, and the OpenOptions and CloseOptions inherit them
through composition.


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

The root command parser defines the arguments, and also sets their Inherit flag
to true:

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

After that, the parsers for the open and close sub-commands do not need to
define the arguments. Because the parent command (the root) had these arguments
with Inherit = true, and because the Values struct for open and close have
destiation fields for Verbose and Debug, argparse will copy the argument
definitions from the root command to the open and close commands.

# Notes

If the parser sees "--" on the command-line, it denotes the beginning of a positional
argument, and no other switch arguments will be processed.

# Examples

For working examples, see the examples/ directory in the source code.
