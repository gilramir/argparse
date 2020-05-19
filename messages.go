package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

// Strings that can be printed out to the user. They can be
// overridden for i18n
type Messages struct {
	// "Sub-Commands"
	SubCommands string

	// "Options'
	Options string

	// The description for the help options (-h / --help):
	// "See this list of options"
	HelpDescription string

	// Error when parsing a boolean
	// "Cannot convert \"%s\" to a boolean"
	CannotParseBooleanFmt string
}

var DefaultMessages_en = Messages{
	SubCommands: "Sub-Commands",
	Options: "Options",
	HelpDescription: "See this list of options",

	CannotParseBooleanFmt: "Cannot convert \"%s\" to a boolean",
}


