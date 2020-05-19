package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"fmt"
//	"log"
)

type Values interface {}

type ParserCallback func (Values) error

type Command struct {
	// The name of the program or subcommand
	Name string

	Description string
	/*
	// One-line description of the program
	ShortDescription string

	// This can be a multi-line, longer explanation of
	// the program.
	LongDescription string
*/
	// This can be a multi-line string that is shown
	// after all the options in the --help output
	Epilog string

	// The struct that will receive the values after parsing
	Values Values

	// The function to call when this parser is selected
	Function ParserCallback

	// Was an option seen during the parse? The key is the name
	// of the destination variable.
	Seen map[string]bool

	// Internal fields
	subCommands         []*Command
	switchArguments     []*Argument
	positionalArguments []*Argument

	numRequiredPositionalArguments int
	// -1 if there is no max (i.e., if the final NumArgs is * or +)
	numMaxPositionalArguments int
}

func (self *Command) init() {
	self.Seen = make(map[string]bool)
}

func (self *Command) New(cmd *Command) *Command {
	cmd.init()
	self.subCommands = append(self.subCommands, cmd)
	return cmd
}

// TODO - check that it's not a HelpSwitch; Command will need to know HelpSwitches
func (self *Command) Add(arg *Argument) {
	if self.Values == nil {
		panic(fmt.Sprintf("There is no Values field set for Command %s", self.Name))
	}
	arg.init(self.Values)

	if arg.isPositional() {
		if len(self.positionalArguments) > 0 {
			prevArg := self.positionalArguments[len(self.positionalArguments)-1]
			if prevArg.NumArgsGlob == "*" || prevArg.NumArgsGlob == "+" ||
					prevArg.NumArgsGlob == "?" {
				panic(fmt.Sprintf(
					"Cannot add a positional argument after argument " +
					"'%s' which has a variable number of values.",
					prevArg.PrettyName()))
			}
		}

		self.positionalArguments = append(self.positionalArguments, arg)
		// If the user didn't set it, it's 1.
		if arg.NumArgs == 0 && arg.NumArgsGlob == "" {
			arg.NumArgs = 1
		}
		if arg.NumArgs > 0 {
			self.numRequiredPositionalArguments += arg.NumArgs
			self.numMaxPositionalArguments += arg.NumArgs
		} else if arg.NumArgsGlob == "+" {
			arg.NumArgs = -1
			self.numRequiredPositionalArguments++
			self.numMaxPositionalArguments += 2
		} else if arg.NumArgsGlob == "?" {
			arg.NumArgs = -1
			self.numMaxPositionalArguments++
		} else if arg.NumArgsGlob == "*" {
			arg.NumArgs = -1
			self.numMaxPositionalArguments = -1
		} else {
			panic("Not reached")
		}
	} else if arg.isSwitch() {
		if arg.NumArgsGlob != "" {
			panic(fmt.Sprintf(
				"Cannot add a switch argument (%s) with a NumArgGlobs " +
				"pattern (%s)", arg.PrettyName(), arg.NumArgsGlob))
		}
		if arg.NumArgs == 0 {
			arg.NumArgs = arg.value.defaultSwitchNumArgs()
		}
		self.switchArguments = append(self.switchArguments, arg)
	} else {
		panic(fmt.Sprintf("Cannot determine argument type for %v", arg))
	}
}

