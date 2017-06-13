// The argparse module is a simple way to add a command-line
// parser to your CLI program. It is modeled somewhat after the
// Python module of the same name. Noticeably, it keeps you organized by
// requiring that the fields whose values are set from the command-line are
// members of a struct, instead of global or local-scope variables.
package argparse

// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Destination interface {
	Run([]Destination) (err error)
}

// The ArgumentParser struct is the top-level, root node of
// the command-line option parsing.
type ArgumentParser struct {
	// The name of the program or subcommand
	Name string

	// One-line description of the program
	ShortDescription string

	// This can be a multi-line, longer explanation of
	// the program.
	LongDescription string

	// This can be a multi-line string that is shown
	// after all the options in the --help output
	Epilog string

	// The struct that will receive the values after parsing
	Destination Destination

	// If this is set, instead of printing the help statement,
	// when --help is requested, to os.Stdout, the output goes here.
	Stdout io.Writer

	// If this is set, instead of printing the usage statement,
	// when a ParseErr is encountered, to os.Stderr, the output goes here.
	Stderr io.Writer

	// Internal fields
	subParsers          []*ArgumentParser
	switchArguments     []*Argument
	positionalArguments []*Argument
        commandArguments    []*Argument

        numRequiredPositionalArguments  int
        // -1 if there is no max (i.e., if the final NumArgs is * or +)
        numMaxPositionalArguments int

	// This causes a circular reference, keeping
	// the tree of ArgumentParsers to NOT be garbage collected
	parentParser *ArgumentParser
}

func (self *ArgumentParser) AddParser(p *ArgumentParser) *ArgumentParser {
	if p.Stdout != nil {
		panic("Only the root parser can set Stdout")
	}

	self.subParsers = append(self.subParsers, p)
	p.parentParser = self
	return p
}

func (self *ArgumentParser) AddArgument(arg *Argument) {
	if self.Destination == nil {
		panic(fmt.Sprintf("There is no Destination set for ArgumentParser %s", self.Name))
	}
	arg.sanityCheck(self.Destination)
	if arg.isCommand() {
		self.commandArguments = append(self.commandArguments, arg)
        } else if arg.isPositional() {
                if len(self.positionalArguments) > 0 {
                    prevNumArgs := self.positionalArguments[len(self.positionalArguments)-1].NumArgs
                    if prevNumArgs == '*' || prevNumArgs == '+' || prevNumArgs == '?' {
                        panic(fmt.Sprintf("Cannot add an Argument after a NumArgs=%s Argument", prevNumArgs))
                    }
                }

		self.positionalArguments = append(self.positionalArguments, arg)
                if arg.NumArgs == '1' || arg.NumArgs == '+' {
                    self.numRequiredPositionalArguments++
                }

                if arg.NumArgs == '1' || arg.NumArgs == '?' {
                    self.numMaxPositionalArguments++
                } else {
                    self.numMaxPositionalArguments = -1
                }
        } else if arg.isSwitch() {
		self.switchArguments = append(self.switchArguments, arg)
	} else {
            panic(fmt.Sprintf("Cannot determine type for %v", arg))
        }
}

func (self *ArgumentParser) ParseArgs() error {
	return self.ParseArgv(os.Args[1:])
}

func (self *ArgumentParser) ParseArgv(argv []string) error {

	results := self.parseArgv(argv)
	if results.parseError != nil {
		// XXX - usage statement
		return results.parseError
	}
	if results.helpRequested {
		var output io.Writer

		if self.Stdout == nil {
			output = os.Stdout
		} else {
			output = self.Stdout
		}
		fmt.Fprintf(output, results.triggeredParser.helpString())
		return nil
	}

	// The parser doesn't have a destination?
	if results.triggeredParser.Destination == nil {
		// show usage statement
                if self.parentParser == nil {
                    // The root command is allowed to have no Destination if it also
                    // doesn't have any subcommands
                    if len(self.subParsers) == 0 {
                        return nil
                    } else {
                        return errors.New("The ArgumentParser needs a Destination")
                    }
                } else {
                        return errors.Errorf("The ArgumentParser '%s' needs a Destination", self.Name)
                }
	}

	err := results.triggeredParser.Destination.Run(results.ancestors)
	switch err.(type) {
	case ParseErr:
		var output io.Writer

		if self.Stderr == nil {
			output = os.Stderr
		} else {
			output = self.Stderr
		}
		fmt.Fprintf(output, "\n%s\n", results.triggeredParser.helpString())
	}
	return err
}

func (self *ArgumentParser) usageString() string {
	var usage string

	var rootParser *ArgumentParser
	for rootParser = self; rootParser.parentParser != nil; {
		rootParser = rootParser.parentParser
	}

	// The name of the program
	if rootParser.Name == "" {
		usage += os.Args[0]
	} else {
		usage += rootParser.Name
	}

	// Are we a subcommand?
	if self.parentParser != nil {
		var subcommandNames string
		for parser := self; parser.parentParser != nil; {
			if subcommandNames == "" {
				subcommandNames = parser.Name
			} else {
				subcommandNames = parser.Name + " " + subcommandNames
			}
			parser = parser.parentParser
		}
		usage += " " + subcommandNames
	}
	usage += "\n\n"

	return usage
}

func (self *ArgumentParser) helpString() string {
	var text string

	text = self.usageString()

	if len(self.subParsers) > 0 {
		text += "Sub-Commands:\n\n"

		longestSubcommandLen := 0

		// Find the longest length of a subcommand name
		for _, subParser := range self.subParsers {
			if len(subParser.Name) > longestSubcommandLen {
				longestSubcommandLen = len(subParser.Name)
			}
		}

		indentation := longestSubcommandLen + 4

		for _, subParser := range self.subParsers {
			padding := strings.Repeat(" ", indentation-len(subParser.Name))
			text += fmt.Sprintf("    %s%s%s\n", subParser.Name, padding,
				subParser.ShortDescription)
		}
	}

	argumentsLabelPrinted := false

	maxLHS := 0
	for _, argument := range self.switchArguments {
		length := len(argument.helpString())
		if length > maxLHS {
			maxLHS = length
		}
	}
	for _, argument := range self.positionalArguments {
		length := len(argument.Metavar)
		if argument.NumArgs == numArgsMaybe {
			length += 2
		}
		if length > maxLHS {
			maxLHS = length
		}
	}
	var startRHS int
	if maxLHS < 20 {
		startRHS = 20
	} else if maxLHS < 30 {
		startRHS = 30
	} else if maxLHS < 40 {
		startRHS = 40
	} else {
		startRHS = maxLHS + 2
	}

	if len(self.switchArguments) > 0 {
		text += "Options:\n"
		argumentsLabelPrinted = true
		for _, argument := range self.switchArguments {
			lhs := argument.helpString()
			indent := startRHS - len(lhs)
			text += "\t" + lhs + strings.Repeat(" ", indent) + argument.Help + "\n"
		}
	}
	if len(self.positionalArguments) > 0 {
		if argumentsLabelPrinted {
			text += "\n"
			argumentsLabelPrinted = true
		}
		text += "Positional Arguments:\n"
		for _, argument := range self.positionalArguments {
			lhs := argument.getMetavar()
			if argument.NumArgs == numArgsMaybe {
				lhs = "[" + lhs + "]"
			}
			indent := startRHS - len(lhs)
			text += "\t" + lhs + strings.Repeat(" ", indent) + argument.Help + "\n"
		}
	}

	if len(self.subParsers)+len(self.switchArguments)+len(self.positionalArguments) == 0 {
		text += "No options\n"
	}

	return text
}

// Print a textual representation of the parser tree to stdout.
// This can be useful for developers if they have issues with their parser.
func (self *ArgumentParser) Dump() {
	self.dump("")
}

func (self *ArgumentParser) dump(spaces string) {
	fmt.Printf("%sName: %s\n", spaces, self.Name)
	if self.ShortDescription != "" {
		fmt.Printf("%sShortDescription: %s\n", spaces, self.ShortDescription)
	}
	if self.LongDescription != "" {
		fmt.Printf("%sLongDescription: %s\n", spaces, self.LongDescription)
	}
	if self.Epilog != "" {
		fmt.Printf("%sEpilog: %s\n", spaces, self.Epilog)
	}
	fmt.Printf("%sDestination: %v\n", spaces, self.Destination)

	var subSpaces string = spaces + "    "
	for _, arg := range self.switchArguments {
		arg.dump(subSpaces)
	}
	for _, arg := range self.positionalArguments {
		arg.dump(subSpaces)
	}
	fmt.Printf("\n")

	subSpaces += "    "
	for _, subParser := range self.subParsers {
		subParser.dump(subSpaces)
	}
}
