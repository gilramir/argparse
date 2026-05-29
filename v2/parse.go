package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"errors"
	"fmt"

	//	"log"
	"strings"
)

// This is returned by the parser
type parseResults struct {
	parseError       error
	helpRequested    bool
	versionRequested bool
	triggeredCommand *Command
	//ancestorValues		[]Values
	ancestorCommands []*Command
}

type tokenType int

const (
	tokError tokenType = iota
	tokArgument
	tokValue
	tokValueNotPresent
	tokSubParser
	tokHelp
	tokVersion
)

type argToken struct {
	typ      tokenType
	pos      int
	value    string
	argument *Argument

	// The name that was given to the argument by the user,
	// since an arg might have a short and a long value
	// in its definition.
	argumentLabel string

	command *Command
}

// This is the parser
type parserState struct {
	ap         *ArgumentParser
	pos        int
	args       []string
	tokenChan  chan argToken
	tokens     []argToken
	lastSwitch string

	cmd *Command
	// If there are sub commands that could be present,
	// this starts as true. Once an arg is parsed, no
	// subparsers can be accepted, so it's changed to false.
	// Switching to a new command changes that of course.
	subCommandAllowed bool

	nextPositionalArgument          int
	numEvaluatedPositionalArguments int
	// How many values have been consumed for the positional argument
	// currently at nextPositionalArgument. Used to know when a fixed-count
	// (NumArgs > 1) positional is full and we should advance to the next one.
	numValuesForCurrentPositional int

	needNValues int
	// when we need to keep track of an *Argument across state transitions
	//	stickyArg *Argument
}

// Each parser state is a function
type stateFunc func() stateFunc

func (self *parserState) emitWithArgument(typ tokenType, argument *Argument, label string) {
	self.tokenChan <- argToken{
		typ:           typ,
		pos:           self.pos,
		argument:      argument,
		argumentLabel: label,
	}
}
func (self *parserState) emitWithValue(typ tokenType, value string) {
	self.tokenChan <- argToken{
		typ:   typ,
		pos:   self.pos,
		value: value,
	}
}
func (self *parserState) emitParser(cmd *Command) {
	self.tokenChan <- argToken{
		typ:     tokSubParser,
		pos:     self.pos,
		command: cmd,
	}
}
func (self *parserState) emitToken(typ tokenType) {
	self.tokenChan <- argToken{
		typ: typ,
		pos: self.pos,
	}
}

// The entrance to the parser
func (self *parserState) runParser(ap *ArgumentParser, argv []string) *parseResults {
	// Initialize the results
	results := &parseResults{
		triggeredCommand: ap.Root,
	}

	// Initialize our state
	self.ap = ap
	self.args = argv
	self.tokenChan = make(chan argToken)

	self.subCommandAllowed = len(ap.Root.subCommands) > 0
	self.cmd = ap.Root

	// The parsing happens in a goroutine
	go self._parse()

	var lastArgLabel string
	var lastArgument *Argument

	for argToken := range self.tokenChan {
		switch argToken.typ {
		case tokArgument:
			self.cmd.Seen[argToken.argument.Dest] = true
			lastArgument = argToken.argument
			lastArgLabel = argToken.argumentLabel
			// If the argument is a boolean argument (no value), then
			// we mark it as seen and move on.
			if lastArgument.NumArgs == 0 {
				err := lastArgument.value.seenWithoutValue(&ap.Messages)
				if err != nil {
					panic(fmt.Sprintf("not reached for arg %s: %s",
						lastArgLabel, err))
				}
			}

		case tokValue:
			if lastArgument == nil {
				panic("Found value without a preceding argument")
			}

			// Parse the text and validate against the Choices, if there
			// are any set for this Argument
			err := lastArgument.value.parse(&ap.Messages, argToken.value)
			if err != nil {
				results.parseError = fmt.Errorf(
					ap.Messages.WhileParsingValueForFmt, lastArgLabel, err)
				return results
			}
		case tokValueNotPresent:
			if lastArgument == nil {
				panic("Found ValueNotPresent without a preceding argument")
			}
			// only bools can have no value
			err := lastArgument.value.seenWithoutValue(&ap.Messages)
			if err != nil {
				results.parseError = fmt.Errorf(
					ap.Messages.ArgumentErrorFmt, lastArgLabel, err)
				return results
			}
		case tokSubParser:
			results.ancestorCommands = append(results.ancestorCommands,
				results.triggeredCommand)
			results.triggeredCommand = argToken.command
		case tokHelp:
			results.helpRequested = true
			return results
		case tokVersion:
			results.versionRequested = true
			return results
		case tokError:
			results.parseError = errors.New(argToken.value)
			return results
		default:
			panic("Unhandled argToken type")
		}
		// XXX - maybe don't need self.tokens ?
		self.tokens = append(self.tokens, argToken)
	}

	// No need to wait for the goroutine to finish. The closing of the
	// channel means the goroutine finished. umm.. what about the early returns
	// for help and errors

	// Did we find all required parameters?
	// TODO - switchArgumants

	// If fewer values were given than the positional arguments require, that's
	// an error regardless of how the shortfall is distributed.
	cmd := results.triggeredCommand
	if self.numEvaluatedPositionalArguments < cmd.numRequiredPositionalArguments {
		// nextPositionalArgument points at the first positional that is still
		// short; clamp in case every defined positional was already visited.
		idx := self.nextPositionalArgument
		if idx >= len(cmd.positionalArguments) {
			idx = len(cmd.positionalArguments) - 1
		}
		arg := cmd.positionalArguments[idx]
		results.parseError = fmt.Errorf(ap.Messages.ExpectedRequiredArgumentFmt, arg.PrettyName())
		return results
	}

	// Propagate inherited argument values
	if len(results.ancestorCommands) > 0 {
		cmdStack := make([]*Command, len(results.ancestorCommands)+1)
		copy(cmdStack, results.ancestorCommands)
		cmdStack[len(cmdStack)-1] = cmd
		cmdStack[0].propagateInherited(cmdStack, 0)
	}

	// Were all the required switches given? Check every command along the
	// path that was traversed (the ancestors and the triggered command);
	// inherited values have already been propagated above.
	pathCommands := append(append([]*Command{}, results.ancestorCommands...), cmd)
	for _, pathCmd := range pathCommands {
		for _, arg := range pathCmd.switchArguments {
			if arg.Required && !pathCmd.Seen[arg.Dest] {
				results.parseError = fmt.Errorf(
					ap.Messages.MissingRequiredSwitchFmt, arg.PrettyName())
				return results
			}
		}
	}

	// Enforce mutually-exclusive groups along the path.
	for _, pathCmd := range pathCommands {
		for _, group := range pathCmd.mutexGroups {
			seen := 0
			for _, arg := range group.args {
				if pathCmd.Seen[arg.Dest] {
					seen++
				}
			}
			if seen > 1 {
				results.parseError = fmt.Errorf(
					ap.Messages.MutuallyExclusiveFmt, group.prettyNames())
				return results
			}
			if group.required && seen == 0 {
				results.parseError = fmt.Errorf(
					ap.Messages.MutuallyExclusiveRequiredFmt, group.prettyNames())
				return results
			}
		}
	}

	return results
}

// This is the engine of the state machine
func (self *parserState) _parse() {
	defer close(self.tokenChan)

	// Start at the initial state, and get the next state,
	// ove and over again, entil we reach the final state (nil)
	var state stateFunc
	for state = self.stateArgument; state != nil; {
		state = state()
	}
}

func (self *parserState) stateArgument() stateFunc {
	if self.pos == len(self.args) {
		// End of the list
		return nil
	}

	arg := self.args[self.pos]
	if arg == "" {
		self.emitWithValue(tokError, self.ap.Messages.EmptyArgument)
		return nil
	}

	// Is it a sub-command?
	if self.subCommandAllowed {
		for _, subCommand := range self.cmd.subCommands {
			if arg == subCommand.Name {
				self.cmd.CommandSeen[arg] = true
				self.emitParser(subCommand)
				self.pos += 1
				// The subparser can have its own subparsers
				self.subCommandAllowed = len(subCommand.subCommands) > 0
				// Start parsing in the subCommand!
				self.cmd = subCommand
				return self.stateArgument
			}
		}
	}

	// Is it a switch argument?
	if len(arg) > 1 && arg[0] == '-' {
		return self.stateSwitchArgument
		/*		self.emitWithValue(tokError, fmt.Sprintf("Unknown argument: %s", arg))
				return nil*/
	}

	// Positional argument?
	if self.nextPositionalArgument == 0 && len(self.cmd.positionalArguments) > 0 {
		return self.statePositionalArgument
	}

	self.emitWithValue(tokError, fmt.Sprintf(self.ap.Messages.UnexpectedArgumentFmt, arg))
	return nil
}

func (self *parserState) stateMaybeOneValue() stateFunc {
	if self.pos == len(self.args) {
		// Fine, we're finished.
		self.emitToken(tokValueNotPresent)
		return nil
	}

	// Does the next token start with a hyphen?
	nextArg := self.args[self.pos]
	if len(nextArg) > 1 && nextArg[0] == '-' {
		// Okay, we have no value
		self.emitToken(tokValueNotPresent)
		return self.stateArgument
	}

	// We do have a value
	self.emitWithValue(tokValue, nextArg)
	self.pos += 1
	return self.stateArgument
}

func (self *parserState) stateOneValue() stateFunc {
	if self.pos == len(self.args) {
		self.emitWithValue(tokError, fmt.Sprintf(self.ap.Messages.ExpectedValueAfterFmt, self.lastSwitch))
		return nil
	}

	self.emitWithValue(tokValue, self.args[self.pos])
	self.pos += 1
	return self.stateArgument
}

func (self *parserState) stateMultipleValues() stateFunc {
	if self.pos == len(self.args) {
		// The command-line ended before all the required values were given.
		if self.needNValues > 0 {
			self.emitWithValue(tokError,
				fmt.Sprintf(self.ap.Messages.ExpectedValueAfterFmt, self.lastSwitch))
		}
		return nil
	}
	if self.needNValues < 1 {
		panic("Should not reach")
	}

	self.emitWithValue(tokValue, self.args[self.pos])
	self.pos += 1
	self.needNValues--
	if self.needNValues > 0 {
		return self.stateMultipleValues
	} else {
		return self.stateArgument
	}
}

func (self *parserState) stateSwitchArgument() stateFunc {
	text := self.args[self.pos]
	if text == "" {
		self.emitWithValue(tokError, self.ap.Messages.EmptyArgument)
		return nil
	}
	// "--" is special... it means the rest of the line is a positional argument
	if text == "--" {
		// If there's a pass-through argument, "--" sends the entire remainder
		// to it verbatim (this works even with no other positional arguments).
		if self.cmd.passThroughArg != nil {
			self.pos += 1
			return self.statePassThrough
		}
		// Positional argument?
		if self.nextPositionalArgument == 0 && len(self.cmd.positionalArguments) > 0 {
			self.pos += 1
			return self.statePositionalArgument
		} else {
			self.emitWithValue(tokError,
				self.ap.Messages.DashDashWithoutPositional)
			return nil
		}
	}

	// Check for '=', as in --value=foo
	// XXX - add check to sanity check, ensureing '=' is not in switch name
	equalsIndex := strings.Index(text, "=")
	var rhs string
	if equalsIndex == 0 {
		self.emitWithValue(tokError, self.ap.Messages.SwitchNameCannotBeginWithEquals)
		return nil
	} else if equalsIndex > 0 {
		rhs = text[equalsIndex+1:]
		text = text[:equalsIndex]
	}

	// Check the help switches
	for _, hw := range self.ap.HelpSwitches {
		if text == hw {
			if rhs == "" {
				self.emitToken(tokHelp)
				return nil
			} else {
				self.emitWithValue(tokError, fmt.Sprintf(self.ap.Messages.DoesNotAcceptValueFmt, hw))
				return nil
			}
		}
	}

	// Check the version switches (only when a Version is configured)
	if self.ap.Version != "" {
		for _, vw := range self.ap.VersionSwitches {
			if text == vw {
				if rhs == "" {
					self.emitToken(tokVersion)
					return nil
				} else {
					self.emitWithValue(tokError, fmt.Sprintf(self.ap.Messages.DoesNotAcceptValueFmt, vw))
					return nil
				}
			}
		}
	}
	match := false
	var arg *Argument
	for _, arg = range self.cmd.switchArguments {
		for _, possibility := range arg.Switches {
			// Does it directly match a switch?
			if text == possibility {
				match = true
				break
			}
			/*
				// We could still have -j4, which is a short option
				// with an adjoining number; this is only valid for short options
				// with  numeric arguments
				if arg.typeKind == reflect.Int &&		// dest is an Int
					text[1] != '-' &&			// short option
					rhs == "" &&				// There wasn't an =
					len(possibility) < len(text) &&
					text[:len(possibility)] == possibility {

					rhs = text[len(possibility):]
					text = text[:len(possibility)]
					match = true
					break
				}
				// TODO - this might be too early to do this
				// Or we could have a group of short booleans IFF the option name is
				// onlyone character long; if -x is a boolean
				// and -y is a boolan, than -xy (and -yx) are valid
				if arg.typeKind == reflect.Bool &&		// dest is an Boolean
					len(possibility) == 2 &&		// switch is 2 chars long
					text[1] != '-' &&			// short option given
					rhs == "" &&				// There wasn't an =
					text[:2] == possibility {

					// Emit this one
					self.emitWithArgument(tokArgument, arg, text[:2])
					self.lastSwitch = text[:2]

					// All other characters in the given switch must also be one-character
					// short-option booleans
					// TODO- think about utf-8 here
					all_others_good := true
					for _, r := range text[2:] {
						found := false
						for _, iArg := range self.cmd.switchArguments {
							if iArg.NumArgs == numArgs0 {
								for _, iSwitch := range iArg.Switches {
									if len(iSwitch) == 2 && rune(iSwitch[1]) == r {
										found = true
										// Emit this one
										self.emitWithArgument(tokArgument, iArg, iSwitch)
										self.lastSwitch = iSwitch
										break
									}
								}
							}
							if found {
								break
							}
						}
						if ! found {
							all_others_good = false
							break
						}
					}
					if !all_others_good {
						self.emitWithValue(tokError,
							fmt.Sprintf("The %s switch is valid but not as '%s'",
								text[:2], text))
						return nil
					}

					// We finished the parse and need to return successfully now
					self.pos += 1
					return self.stateArgument
				}
			*/
		}

		if match {
			break
		}
	}
	// Didn't match ?
	if !match {
		// Didn't find a switch with that name
		self.emitWithValue(tokError, fmt.Sprintf(self.ap.Messages.NoSuchSwitchFmt, text))
		return nil
	}

	self.emitWithArgument(tokArgument, arg, text)
	self.lastSwitch = text
	if rhs == "" {
		self.pos += 1
		if arg.NumArgs == 0 {
			return self.stateArgument
		} else if arg.NumArgs == 1 {
			return self.stateOneValue
		} else if arg.NumArgs > 1 {
			self.needNValues = arg.NumArgs
			return self.stateMultipleValues
		} else if arg.NumArgs == -1 {
			panic("not reached")
		} else {
			// ???
			panic(fmt.Sprintf("Unexpected num args: %v", arg.NumArgs))
		}
	} else {
		if arg.NumArgs == 0 {
			self.emitWithValue(tokError,
				fmt.Sprintf(self.ap.Messages.SwitchDoesNotTakeValueFmt, text))
			return nil
		} else if arg.NumArgs == 1 {
			self.emitWithValue(tokValue, rhs)
			self.pos += 1
			return self.stateArgument
		} else if arg.NumArgs > 1 {
			self.emitWithValue(tokValue, rhs)
			self.pos += 1
			self.needNValues = arg.NumArgs - 1
			return self.stateMultipleValues
		} else if arg.NumArgs == -1 {
			panic("not reached")
		} else {
			// ???
			panic(fmt.Sprintf("Unexpected num args: %v", arg.NumArgs))
		}
	}
	panic("not reached")
}

func (self *parserState) statePositionalArgument() stateFunc {
	if self.pos == len(self.args) {
		// End of the list
		return nil
	}

	// If the positional we're about to fill is the pass-through argument,
	// switch to capturing the rest of the command-line verbatim.
	if self.cmd.passThroughArg != nil &&
		self.nextPositionalArgument < len(self.cmd.positionalArguments) &&
		self.cmd.positionalArguments[self.nextPositionalArgument] == self.cmd.passThroughArg {
		return self.statePassThrough
	}

	arg := self.args[self.pos]

	// Consume this token as a value for the current positional argument if
	// there is still room: either the final positional is unbounded ("*"/"+",
	// numMax == -1) or we have not yet reached the maximum number of values.
	if self.cmd.numMaxPositionalArguments == -1 ||
		self.numEvaluatedPositionalArguments < self.cmd.numMaxPositionalArguments {
		posArg := self.cmd.positionalArguments[self.nextPositionalArgument]
		self.emitWithArgument(tokArgument, posArg, posArg.Name)
		self.emitWithValue(tokValue, arg)
		self.pos += 1
		self.consumedPositionalValue(posArg)
		return self.statePositionalArgument
	}

	// There is no more room for positional values. A switch argument may
	// follow a fixed number of positional arguments.
	if len(arg) > 1 && arg[0] == '-' {
		return self.stateSwitchArgument
	}
	self.emitWithValue(tokError, fmt.Sprintf(self.ap.Messages.UnexpectedPositionalArgumentFmt, arg))
	return nil
}

// consumedPositionalValue updates the bookkeeping after a value has been
// emitted for the positional argument at nextPositionalArgument. If that
// positional has now received all the values it accepts, advance to the next
// positional argument.
func (self *parserState) consumedPositionalValue(posArg *Argument) {
	self.numEvaluatedPositionalArguments++
	self.numValuesForCurrentPositional++

	// "*" and "+" accept an unlimited number of values and are always the last
	// positional argument, so they never become "full". "?" accepts at most
	// one. A fixed NumArgs is full once it has collected NumArgs values.
	var full bool
	switch posArg.NumArgsGlob {
	case "*", "+":
		full = false
	case "?":
		full = true
	default:
		full = self.numValuesForCurrentPositional >= posArg.NumArgs
	}
	if full {
		self.nextPositionalArgument++
		self.numValuesForCurrentPositional = 0
	}
}

// statePassThrough consumes every remaining token verbatim into the command's
// pass-through argument, with no switch or "--" interpretation.
func (self *parserState) statePassThrough() stateFunc {
	passArg := self.cmd.passThroughArg
	for self.pos < len(self.args) {
		self.emitWithArgument(tokArgument, passArg, passArg.Name)
		self.emitWithValue(tokValue, self.args[self.pos])
		self.pos++
	}
	return nil
}
