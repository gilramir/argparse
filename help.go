// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package argparse

import (
	"os"
	"github.com/gilramir/consolesize"
)



// This should honor width too
func (self *ArgumentParser) usageString(ap *ArgumentParser) string {
	var usage string

	// The name of the program
	if ap.Root.Name == "" {
		usage += os.Args[0]
	} else {
		usage += ap.Root.Name
	}

	usage += "\n"

	if ap.Root.Description != "" {
		usage += "\n" + ap.Root.Description + "\n"
	}

/*
	// Are we a subcommand?
	if self != ap {
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
*/

	return usage
}

func (self *ArgumentParser) helpString(ap *ArgumentParser, width int) string {
	var text string

	if width == 0 {
		wh, err := consolesize.GetConsoleWidthHeight()
		if err != nil {
			width = 80
		} else {
			width = wh.Width
		}
	}

	text = self.usageString(ap) + "\n"
	formatter := &HelpFormatter{}

/*
	if len(self.subParsers) > 0 {
		text += self.StringSubCommands + ":\n\n"

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
				subParser.Description)
		}
	}
*/

	for _, arg := range self.Root.switchArguments {
		//argumentStrings := self.HelpSwitches()
		argumentStrings := arg.Switches
		formatter.addOption(argumentStrings, arg.Help)
	}

	// The argument should know it's a positional argument, and thus how to format itels
	for _, arg := range self.Root.positionalArguments {
		var argName string
		if arg.MetaVar != "" {
			argName = arg.MetaVar
		} else {
			argName = arg.Name
		}

		if arg.NumArgsGlob == "?" {
			argName = "[" + argName + "]"
		} else if arg.NumArgsGlob == "+" {
			argName = argName + "[ ... ]"
		} else if arg.NumArgsGlob == "*" {
			argName = "[" + argName + "[ ... ] ]"
		}
		formatter.addOption([]string{argName}, arg.Help)
	}

	text += formatter.produceString(width)

	/*
	if len(self.switchArguments) > 0 {
		text += "Options:\n"
		argumentsLabelPrinted = true
		for _, argument := range self.switchArguments {
			lhs := argument.helpString()
			indent := startRHS - len(lhs)
			text += "\t" + lhs + strings.Repeat(" ", indent) + argument.Help + "\n"
			// Show Choices if available
			if len(argument.Choices) > 0 {
				text += "\t" + strings.Repeat(" ", len(lhs)+indent) +
					"Possible choices: " + argument.getChoicesString() + "\n"
			}
		}
	}
	if len(self.positionalArguments) > 0 {
		if argumentsLabelPrinted {
			text += "\n"
			argumentsLabelPrinted = true
		}
		text += "Positional Arguments:\n"
		for _, argument := range self.positionalArguments {
			lhs := argument.getMetaVar()
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
*/
	return text
}
