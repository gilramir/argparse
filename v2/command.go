package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"fmt"
	//	"log"
)

type Values interface{}

type ParserCallback func(*Command, Values) error

type Command struct {
	// The name of the program or subcommand
	Name string

	// This can be a multi-line string that is shown
	// after after the command name, and before the options.
	Description string

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

	// Was a sub-command seen during the parse?
	CommandSeen map[string]bool

	// Internal fields
	subCommands         []*Command
	switchArguments     []*Argument
	positionalArguments []*Argument

	numRequiredPositionalArguments int
	// -1 if there is no max (i.e., if the final NumArgsGlob is "*" or "+")
	numMaxPositionalArguments int

	// Pointer to the ArgumentParser
	ap *ArgumentParser
}

func (self *Command) init(parent *Command, ap *ArgumentParser) {
	self.Seen = make(map[string]bool)
	self.CommandSeen = make(map[string]bool)
	self.ap = ap

	// Nothing futher for the root Command
	if parent == nil {
		return
	}

	// Does the parent have any arguments to inherit?
	for _, arg := range parent.switchArguments {
		if arg.Inherit {
			newArg := arg.deepCopy()
			self.Add(newArg)
		}
	}
}

func (self *Command) propagateInherited(cmds []*Command, myIndex int) {
	if self != cmds[myIndex] {
		panic(fmt.Sprintf("Expected %v at %d but got %v", self, myIndex,
			cmds[myIndex]))
	}
	if len(self.switchArguments) == 0 {
		return
	}

	nextIndex := myIndex + 1
	if len(cmds) <= nextIndex {
		return
	}
	nextCmd := cmds[nextIndex]

	if len(nextCmd.switchArguments) == 0 {
		return
	}

	for _, arg := range self.switchArguments {
		if arg.Inherit && self.Seen[arg.Dest] && !nextCmd.Seen[arg.Dest] {
			// Propagate
			found := false
			for _, nextCmdArg := range nextCmd.switchArguments {
				if nextCmdArg.Dest == arg.Dest {
					found = true
					nextCmdArg.value.setValue(arg.value.getValue())
					nextCmd.Seen[arg.Dest] = true
					break
				}
			}
			if !found {
				panic(fmt.Sprintf("Arg %s inherited from %s to %s can't be found",
					arg.Dest, self.Name, nextCmd.Name))
			}
		}
	}

	if len(cmds) > nextIndex {
		nextCmd.propagateInherited(cmds, nextIndex)
	}
}

func (self *Command) New(cmd *Command) *Command {

	// Check for duplicates
	for _, other := range self.subCommands {
		if other.Name == cmd.Name {
			panic(fmt.Sprintf("Sub-command %s already exists in %s",
				other.Name, self.Name))
		}
	}

	cmd.init(self, self.ap)

	self.subCommands = append(self.subCommands, cmd)
	return cmd
}

// TODO - check that it's not a HelpSwitch; Command will need to know HelpSwitches
func (self *Command) Add(arg *Argument) {

	// Arguments with Inherit == true cannot be added after a sub-command is already
	// added
	if arg.Inherit && len(self.subCommands) > 0 {
		panic(fmt.Sprintf("Cannot add argument %s because it's Inherit flag is true "+
			"and the Command %s already has sub-commands", arg.PrettyName(),
			self.Name))
	}

	// Check for a duplicate
	if arg.isPositional() {
		for _, other := range self.positionalArguments {
			if other.Name == arg.Name {
				panic(fmt.Sprintf("%s is already used by a "+
					"positional argument in this Command.", arg.Name))
			}
		}
	} else {
		for _, other := range self.switchArguments {
			for _, otherSwitch := range other.Switches {
				for _, thisSwitch := range arg.Switches {
					if otherSwitch == thisSwitch {
						panic(fmt.Sprintf("%s is already used by a "+
							"switch argument in this Command.", thisSwitch))
					}
				}
			}
		}
	}

	// Sanity check
	if self.Values == nil {
		panic(fmt.Sprintf("There is no Values field set for Command %s", self.Name))
	}

	// set arg.value
	arg.init(self.Values, &self.ap.Messages)

	if arg.isPositional() {
		if len(self.positionalArguments) > 0 {
			prevArg := self.positionalArguments[len(self.positionalArguments)-1]
			// Can't add a positional argument after * or +, as
			// we don't know when those argument lists will end.
			// But we can add after a ?; it either exists or not.
			if prevArg.NumArgsGlob == "*" || prevArg.NumArgsGlob == "+" {
				panic(fmt.Sprintf(
					"Cannot add a positional argument after argument "+
						"'%s' which has an unlimited number of values.",
					prevArg.PrettyName()))
			}
		}

		self.positionalArguments = append(self.positionalArguments, arg)
		// If the user didn't set it, it's 1.
		if arg.NumArgs == 0 && arg.NumArgsGlob == "" {
			arg.NumArgs = 1
		}
		var canBeMoreThanOne bool
		if arg.NumArgs > 0 {
			self.numRequiredPositionalArguments += arg.NumArgs
			self.numMaxPositionalArguments += arg.NumArgs
			canBeMoreThanOne = true
		} else if arg.NumArgsGlob == "+" {
			// + : one or more
			arg.NumArgs = -1
			self.numRequiredPositionalArguments++
			self.numMaxPositionalArguments = -1
			canBeMoreThanOne = true
		} else if arg.NumArgsGlob == "?" {
			// ? : zero or one
			arg.NumArgs = -1
			self.numMaxPositionalArguments++
			canBeMoreThanOne = false
		} else if arg.NumArgsGlob == "*" {
			// * : zero or more
			arg.NumArgs = -1
			self.numMaxPositionalArguments = -1
			canBeMoreThanOne = true
		} else {
			panic("Not reached")
		}

		// If the positional argument accepts more than one value,
		// the destination must be a slice
		if arg.NumArgs == -1 && canBeMoreThanOne && arg.value.storageType() != Slice {
			panic(fmt.Sprintf(
				"Cannot use positional argument %s with a non-slice destination variable because NumArgsGlob is %s", arg.PrettyName(), arg.NumArgsGlob))
		} else if arg.NumArgs > 1 && arg.value.storageType() != Slice {
			panic(fmt.Sprintf(
				"Cannot use positional argument %s with a non-slice destination variable because NumArgs is %d", arg.PrettyName(), arg.NumArgs))
		}

	} else if arg.isSwitch() {
		if arg.NumArgsGlob != "" {
			panic(fmt.Sprintf(
				"Cannot add a switch argument (%s) with a NumArgGlobs "+
					"pattern (%s)", arg.PrettyName(), arg.NumArgsGlob))
		}
		if arg.NumArgs == 0 {
			arg.NumArgs = arg.value.defaultSwitchNumArgs()
		}
		// If the switch argument accepts more than one value,
		// the destination must be a slice
		if arg.NumArgs > 1 && arg.value.storageType() != Slice {
			panic(fmt.Sprintf(
				"Cannot use switch argument %s with a non-slice destination variable because NumArgs is %d", arg.PrettyName(), arg.NumArgs))
		}
		self.switchArguments = append(self.switchArguments, arg)
	} else {
		panic(fmt.Sprintf("Cannot determine argument type for %v", arg))
	}
}
