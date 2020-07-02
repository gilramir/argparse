// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package argparse

import (
	"github.com/gilramir/consolesize"
	"os"
)

// TODO - the help should show Choices, if available

// This should honor width too
func (self *ArgumentParser) usageString(cmd *Command,
	ancestorCommands []*Command) string {
	var usage string

	// Show the names of the command(s)
	commands := make([]*Command, len(ancestorCommands))
	copy(commands, ancestorCommands)
	commands = append(commands, cmd)

	for i, iCmd := range commands {
		if i == 0 {
			if iCmd.Name == "" {
				usage += os.Args[0]
			} else {
				usage = iCmd.Name
			}
		} else {
			usage += " " + iCmd.Name
		}
	}

	// Do we have a description to report
	if cmd.Description != "" {
		usage += "\n\n" + cmd.Description + "\n"
	}

	return usage
}

func (self *ArgumentParser) helpString(cmd *Command,
	ancestorCommands []*Command) string {
	var text string

	width := 80
	wh, err := consolesize.GetConsoleWidthHeight()
	if err != nil {
		width = wh.Width
	}

	text = self.usageString(cmd, ancestorCommands) + "\n"
	formatter := &helpFormatter{}

	// Switch arguments

	for _, arg := range cmd.switchArguments {
		argumentStrings := arg.Switches
		formatter.addOption(argumentStrings, arg.Help)
	}
	formatter.addOption(self.HelpSwitches, self.Messages.HelpDescription)

	// Positional arguments

	for _, arg := range cmd.positionalArguments {
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

	// Sub-commands

	if len(cmd.subCommands) > 0 {
		text += "\n" + self.Messages.SubCommands + ":\n\n"

		subFormatter := &helpFormatter{}
		for _, subCommand := range cmd.subCommands {
			subFormatter.addOption([]string{subCommand.Name}, subCommand.Description)
		}
		text += subFormatter.produceString(width)
	}

	return text
}
