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

	// The Choices slice is of the wrong type.
	// "Choices should be []%s"
	ChoicesOfWrongTypeFmt string

	// The given value is not a valid choice
	// "Not a valid choice. Should be one of: %v"
	// TODO This should be changed to have %s and %v, to show the incorrect value
	ShouldBeAValidChoiceFmt string
}

var DefaultMessages_en = Messages{
	SubCommands:     "Sub-Commands",
	Options:         "Options",
	HelpDescription: "See this list of options",

	CannotParseBooleanFmt:   "Cannot convert \"%s\" to a boolean",
	ChoicesOfWrongTypeFmt:   "Choices should be []%s",
	ShouldBeAValidChoiceFmt: "Not a valid choice. Should be one of: %v",
}
