package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/gilramir/consolesize"
	"github.com/gilramir/unicodemonowidth"
)

// argumentHelp returns the help text for an argument, with its allowed values
// (Choices) and default value appended, if any.
func (self *ArgumentParser) argumentHelp(arg *Argument) string {
	parts := []string{}
	if arg.Help != "" {
		parts = append(parts, arg.Help)
	}
	if arg.Choices != nil {
		parts = append(parts,
			fmt.Sprintf(self.Messages.HelpChoicesFmt, formatChoices(arg.Choices)))
	}
	if arg.hasDefault {
		parts = append(parts,
			fmt.Sprintf(self.Messages.HelpDefaultFmt, arg.defaultValue))
	}
	return strings.Join(parts, " ")
}

// formatChoices renders a Choices slice (held as an interface{}) as a
// comma-separated list.
func formatChoices(choices interface{}) string {
	v := reflect.ValueOf(choices)
	if v.Kind() != reflect.Slice {
		return fmt.Sprintf("%v", choices)
	}
	parts := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		parts[i] = fmt.Sprintf("%v", v.Index(i).Interface())
	}
	return strings.Join(parts, ", ")
}

// This should honor width too
func (self *ArgumentParser) usageString(cmd *Command, width int, ancestorCommands []*Command) string {
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

	usage += "\n\n"
	// Do we have a description to report
	if cmd.Description != "" {
		descWords := unicodemonowidth.WhitespaceSplit(cmd.Description)
		descRows := unicodemonowidth.WrapPrintedWords(descWords, width)
		for _, row := range descRows {
			usage += row + "\n"
		}
	}

	return usage
}

func (self *ArgumentParser) helpString(cmd *Command,
	ancestorCommands []*Command) string {
	var text string

	width := 80
	wh, err := consolesize.GetConsoleWidthHeight()
	if err == nil {
		width = wh.Width
	}

	text = self.usageString(cmd, width, ancestorCommands) + "\n"
	formatter := &helpFormatter{}

	// Switch arguments

	for _, arg := range cmd.switchArguments {
		argumentStrings := arg.Switches
		if !arg.isPositional() && arg.NumArgs > 0 {
			// set a default metavar?
			var metavar string
			if arg.MetaVar == "" {
				// Use the upper-case version of the first switch, with
				// no dashes at the front.
				metavar = strings.TrimLeft(strings.ToUpper(arg.Switches[0]), "-")
			} else {
				metavar = arg.MetaVar
			}
			// Add the metavar to the last one
			idx := len(argumentStrings) - 1
			argumentStrings[idx] = argumentStrings[idx] + "=" + metavar
		}
		formatter.addOption(argumentStrings, self.argumentHelp(arg))
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
		formatter.addOption([]string{argName}, self.argumentHelp(arg))
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

	// Do we have an Epilog to report?
	if cmd.Epilog != "" {
		indent := ""
		if width > 40 {
			// Indent 4, with 4 spaces on the RHS too.
			indent = "    "
			width -= 8
		}
                // Don't split the epilog; treat it as raw.
                /*
		epilogWords := unicodemonowidth.WhitespaceSplit(cmd.Epilog)
		epilogRows := unicodemonowidth.WrapPrintedWords(epilogWords, width)
		text += "\n\n"
                */
                text += "\n"
                epilogRows := strings.Split(cmd.Epilog, "\n")
		for _, row := range epilogRows {
			text += indent + row + "\n"
		}
	}

	return text
}
