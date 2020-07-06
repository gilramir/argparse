package argparse

// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"fmt"
	"io"
	"os"
)

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
		Stdout:       os.Stdout,
		Stderr:       os.Stderr,
		Messages:     DefaultMessages_en,
		HelpSwitches: []string{"-h", "--help"},
		Root:         cmd,
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

// Parse the os.Argv arguments and return, having filled out Values.
// On a request for help (-h), print the help and exit with os.Exit(0).
// On a user input error, print the error message and exit with os.Exit(1).
func (self *ArgumentParser) Parse() {
	results := self.ParseArgv(os.Args[1:])

	cmd := results.triggeredCommand

	if results.helpRequested {
		helpString := self.helpString(cmd, results.ancestorCommands)
		fmt.Fprintln(self.Stdout, helpString)
		os.Exit(0)
	} else if results.parseError != nil {
		fmt.Fprintln(self.Stderr, results.parseError.Error())
		os.Exit(1)
	}

	if cmd.Function != nil {
		err := cmd.Function(cmd, cmd.Values)
		if err != nil {
			fmt.Fprintln(self.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

/*
func (self *ArgumentParser) ParseOsArgs() (Values) {
	// Sanity check
*/
/*
	if self.root.Function != nil {
		panic fmt.Sprintf("ParseOsArgv can't be used with a root ArgumentParser "\
			"that has a callback Function. Use ExecuteOsArgv instead.")
	}
*/
/*
	return self.ParseArgv( os.Args[1:] )
}
*/

/*
	// Is there a function?
	if results.triggeredParser.Function == nil {
		return nil
	}

	err := results.triggeredParser.Function(nil)
	switch err.(type) {
	case ParseError:
		fmt.Fprintf(pt.Stderr, "\n%s\n",
			results.triggeredParser.helpString(pt.rootParser, 0))
	}
	return err
}
*/

func (self *ArgumentParser) ParseArgv(argv []string) *parseResults {
	parser := parserState{}
	results := parser.runParser(self, argv)
	return results
}

/*
	valueStack := make([]Values, 0, 1)
	err := self.root.startParseArgv(self, argv, &valueStack)

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	return valueStack[len(valueStack)-1]
}*/
/*
func (self *ArgumentParser) CommandParseOsArgs() (*Command, Values, error) {
	// Sanity check
	/*
	if self.root.Function != nil {
		panic fmt.Sprintf("ParseOsArgv can't be used with a root ArgumentParser "\
			"that has a callback Function. Use ExecuteOsArgv instead.")
	}
	return self.ParseArgv( os.Args[1:] )
}

func (self *ArgumentParser) CommandParseArgv(argv []string) (*Command, Values, error) {
	// Sanity check
	/*
	if self.root.Function != nil {
		panic fmt.Sprintf("ParseArgv can't be used with a root ArgumentParser "\
			"that has a callback Function. Use ExecuteArgv instead.")
	}

	valueStack := make([]Values, 0, 1)
	err := self.root.startParseArgv(self, argv, &valueStack)
	return valueStack[len(valueStack)-1], err
}
*/

/*
func (self *ArgumentParser) ExecuteOsArgs() (error) {
	// Sanity check
	if self.root.Function == nil {
		panic fmt.Sprintf("ExecuteOsArgs can't be used with a root ArgumentParser "\
			"that has no callback Function. Use ParseOsArgs instead.")
	}
	return self.ExecuteArgv( os.Args[1:] )
}

func (self *ArgumentParser) ExeuteArgv(argv []string) (error) {
	// Sanity check
	if self.root.Function == nil {
		panic fmt.Sprintf("ExecuteArgv can't be used with a root ArgumentParser "\
			"that has no callback Function. Use ParseArgv instead.")
	}

	valueStack := make([]OptionValues, 0, 1)
	return self.root.startParseArgv(self, argv, &valueStack)
}
*/
