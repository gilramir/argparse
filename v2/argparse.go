/*
Argparse is a Go module that makes it easy to write the command-line parsing
part of your program.
It loosely follows the conceptual model of the Python argparse module.

Highlights:

  - You can have nested subcommands.
  - The values for the command-line options are stored in a struct of your
    creation.
  - Argparse can deduce the name of the value field in the struct by looking
    at the name of the option. Or, you can tell it exactly which field to
    use.
  - Argparse will tell you if a particular option was present on the command-line
    or not present, in case you need that information.
  - Options can be inherited by sub-comands, so you need only define them
    once.
  - The built-in help strings are translatable.
*/
package argparse

// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ErrHelp is returned by ParseArgs when the user requested help (for example
// with -h or --help). By the time it is returned, the help text has already
// been written to the parser's Stdout.
var ErrHelp = errors.New("argparse: help requested")

// ErrVersion is returned by ParseArgs when the user requested the version (for
// example with --version). By the time it is returned, the version string has
// already been written to the parser's Stdout.
var ErrVersion = errors.New("argparse: version requested")

type ArgumentParser struct {
	// If this is set, instead of printing the help statement,
	// when --help is requested, to os.Stdout, the output goes here.
	Stdout io.Writer

	// If this is set, instead of printing the usage statement,
	// when a ParseErr is encountered, to os.Stderr, the output goes here.
	Stderr io.Writer

	// Allow the user to modify strings produced by argparse.
	// This is essential for i18n
	Messages Messages

	// The switch strings that can invoke help
	HelpSwitches []string

	// If Version is non-empty, the VersionSwitches (default "--version") will
	// print this string and stop parsing, the same way HelpSwitches print the
	// help. If Version is empty, the version switches are not special.
	Version string

	// The switch strings that can request the version. Only consulted when
	// Version is non-empty.
	VersionSwitches []string

	// The root Command object.
	Root *Command

	// The first time a parse is run, a finalization step need to be
	// performed to fill out inherited Arguments. This flag ensures
	// we do that only once.
	finalized bool
}

// Create a new ArgumentParser, with the Command as its root Command
func New(cmd *Command) *ArgumentParser {
	ap := &ArgumentParser{
		Stdout:          os.Stdout,
		Stderr:          os.Stderr,
		Messages:        DefaultMessages_en,
		HelpSwitches:    []string{"-h", "--help"},
		VersionSwitches: []string{"--version"},
		Root:            cmd,
	}
	cmd.init(nil, ap)
	if cmd.Name == "" {
		cmd.Name = os.Args[0]
	}
	return ap
}

// Add an argument to the root command
func (self *ArgumentParser) Add(arg *Argument) {
	self.Root.Add(arg)
}

// Add a command to the root command
func (self *ArgumentParser) New(c *Command) *Command {
	return self.Root.New(c)
}

// Add a mutually-exclusive group of switch arguments to the root command.
func (self *ArgumentParser) AddMutuallyExclusive(required bool, args ...*Argument) {
	self.Root.AddMutuallyExclusive(required, args...)
}

// Parse the os.Argv arguments and return, having filled out Values.
// On a request for help (-h), print the help and exit with os.Exit(0).
// On a user input error, print the error message and exit with os.Exit(1).
func (self *ArgumentParser) Parse() {
	self.parseRunFunction(true)
}

// ParseArgs parses the given arguments (not os.Args) and returns the Command
// that was triggered, having filled out its Values. Unlike Parse and
// ParseAndExit, it never calls os.Exit and does not run any Function callback;
// it reports problems by returning an error. This makes it suitable for
// embedding argparse in another program (such as a REPL or sub-shell) and for
// testing your own command-line handling.
//
// If the user requested help, the help text is written to the parser's Stdout
// and ErrHelp is returned. On a user input error, the error is returned (it is
// not printed). On success, the error is nil.
func (self *ArgumentParser) ParseArgs(argv []string) (*Command, error) {
	results := self.parseArgv(argv)
	cmd := results.triggeredCommand

	if results.helpRequested {
		helpString := self.helpString(cmd, results.ancestorCommands)
		fmt.Fprintln(self.Stdout, helpString)
		return cmd, ErrHelp
	}
	if results.versionRequested {
		fmt.Fprintln(self.Stdout, self.Version)
		return cmd, ErrVersion
	}
	if results.parseError != nil {
		return cmd, results.parseError
	}
	return cmd, nil
}

// Parse and run the function.
func (self *ArgumentParser) parseRunFunction(shouldReturn bool) {
	cmd, err := self.ParseArgs(os.Args[1:])

	if err == ErrHelp || err == ErrVersion {
		// The help or version text was already written to Stdout by ParseArgs.
		os.Exit(0)
	} else if err != nil {
		fmt.Fprintf(self.Stderr, "%s: %s\n", filepath.Base(self.Root.Name), err.Error())
		os.Exit(1)
	}

	if cmd.Function != nil {
		err := cmd.Function(cmd, cmd.Values)
		if err != nil {
			fmt.Fprintln(self.Stderr, err.Error())
			os.Exit(1)
		}
		// Success
		if shouldReturn {
			return
		} else {
			os.Exit(0)
		}
	}
	// The chosen command had no function to run!
	if !shouldReturn {
		// Print the usage to stderr, and exit with non-0.
		helpString := self.helpString(cmd, cmd.ancestors())
		fmt.Fprintln(self.Stderr, helpString)
		os.Exit(1)
	}
}

// Parse the os.Argv arguments, call the Function for the triggered
// Command, and then exit. An error returned from the Function causes us
// to exit with 1, otherwise, exit with 0.
// On a request for help (-h), print the help and exit with os.Exit(0).
// On a user input error, print the error message and exit with os.Exit(1).
func (self *ArgumentParser) ParseAndExit() {
	self.parseRunFunction(false)
}

func (self *ArgumentParser) parseArgv(argv []string) *parseResults {
	parser := parserState{}
	results := parser.runParser(self, argv)
	return results
}
