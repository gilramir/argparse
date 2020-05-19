package argparse

import (
	"errors"
	"fmt"
//	"log"
//	"reflect"
	"strings"
)

// This is returned by the parser
type parseResults struct {
	parseError		error
	helpRequested		bool
	triggeredCommand	*Command
	ancestorValues		[]Values
	ancestorCommands        []*Command
}

type tokenType int

const (
	tokError tokenType = iota
	tokArgument
	tokValue
	tokValueNotPresent
	tokSubParser
	tokHelp
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
	ap		*ArgumentParser
	pos		int
	args		[]string
	tokenChan	chan argToken
	tokens		[]argToken
	lastSwitch	string

	cmd		*Command
	// If there are sub commands that could be present,
	// this starts as true. Once an arg is parsed, no
	// subparsers can be accepted, so it's changed to false.
	// Switching to a new command changes that of course.
	subCommandAllowed bool

	nextPositionalArgument          int
	numEvaluatedPositionalArguments int

	// when we need to keep track of an *Argument across state transitions
	stickyArg *Argument
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
		typ:    tokSubParser,
		pos:    self.pos,
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
			self.cmd.Seen[ argToken.argument.Dest ] = true
			lastArgument = argToken.argument
			lastArgLabel = argToken.argumentLabel
			// If the argument is a boolean argument (no value), then
			// we mark it as seen and move on.
			if lastArgument.NumArgs == 0 {
				err := lastArgument.value.seenWithoutValue()
				if err != nil {
					panic( fmt.Sprintf("not reached for arg %s: %s",
						lastArgLabel, err) )
				}
			}

		case tokValue:
			if lastArgument == nil {
				panic("Found value without a preceding argument")
			}
			// Does the lastArgument have a Choices slice which limits
			// the valid values?
//			if len(lastArgument.Choices) > 0 {
//				good := false
//				for _, choice := range lastArgument.Choices {
//					if argToken.value == choice {
//						good = true
//						break
//					}
//				}
//				if !good {
//					results.parseError = fmt.Errorf(
//						"The possible values for %s are %s", lastArgLabel,
//						lastArgument.getChoicesString())
//					return results
//				}
//			}

			err := lastArgument.value.parse(&ap.Messages, argToken.value)
			if err != nil {
				results.parseError = fmt.Errorf(
					"While parsing value for %s: %w", lastArgLabel, err)
				return results
			}
		case tokValueNotPresent:
			if lastArgument == nil {
				panic("Found ValueNotPresent without a preceding argument")
			}
			// only bools can have no value
			err := lastArgument.value.seenWithoutValue()
			if err != nil {
				results.parseError = fmt.Errorf(
					"%s argument: %w", lastArgLabel, err)
				return results
			}
		case tokSubParser:
// FIXME 
//			results.ancestors = append(results.ancestors,
//				results.triggeredParser.Destination) 
			results.triggeredCommand = argToken.command
		case tokHelp:
			results.helpRequested = true
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
	// channel means the goroutine finished.

	// Did we find all required parameters?
	// TODO - switchArgumants

	// If there aren't enough positional arguments, check the next known argument to see if it is required
	cmd := results.triggeredCommand
	if len(cmd.positionalArguments) > 0 && self.numEvaluatedPositionalArguments < cmd.numRequiredPositionalArguments {
		arg := results.triggeredCommand.positionalArguments[self.nextPositionalArgument]
		if arg.NumArgs == 1 || arg.NumArgsGlob == "+" {
			results.parseError = fmt.Errorf("Expected a required '%s' argument", arg.PrettyName())
			return results
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
		self.emitWithValue(tokError, "<empty string>")
		return nil
	}

	// Is it a subparser?
	if self.subCommandAllowed {
		for _, subCommand := range self.cmd.subCommands {
			if arg == subCommand.Name {
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

	/*
	// Is it a parse command?
	for _, commandArg := range self.commandArguments {
		if arg == commandArg.String {
			// A little ugly, since stateCommandArgument is going to check self.commandArguments again
			return self.stateCommandArgument
		}
	}
	*/

	// Is it a switch argument?
	if len(arg) > 1 && arg[0] == '-' {
		return self.stateOption
/*		self.emitWithValue(tokError, fmt.Sprintf("Unknown argument: %s", arg))
		return nil*/
	}

	// Positional argument?
	if self.nextPositionalArgument == 0 && len(self.cmd.positionalArguments) > 0 {
		return self.statePositionalArgument
	}

	self.emitWithValue(tokError, fmt.Sprintf("Unexpected argument: %s", arg))
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
		self.emitWithValue(tokError, fmt.Sprintf("Expected a value after %s", self.lastSwitch))
		return nil
	}

	self.emitWithValue(tokValue, self.args[self.pos])
	self.pos += 1
	return self.stateArgument
}

func (self *parserState) stateMultipleValues() stateFunc {
	if self.pos == len(self.args) {
		return nil
	}

	self.emitWithValue(tokValue, self.args[self.pos])
	self.pos += 1
	return self.stateMultipleValues
}

func (self *parserState) stateOption() stateFunc {
	text := self.args[self.pos]
	if text == "" {
		self.emitWithValue(tokError, "<empty string>")
		return nil
	}

	// Check for '=', as in --value=foo
	// XXX - add check to sanity check, ensureing '=' is not in switch name
	equalsIndex := strings.Index(text, "=")
	var rhs string
	if equalsIndex == 0 {
		self.emitWithValue(tokError, "A switch name cannot begin with '='")
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
				self.emitWithValue(tokError, hw + " does not accept a value")
				return nil
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
	if ! match {
		// Didn't find a switch with that name
		self.emitWithValue(tokError, fmt.Sprintf("No such switch: %s", text))
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
			panic("not implemented yet")
		} else if arg.NumArgs == -1 {
			panic("not reached")
		} else {
			// ???
			panic(fmt.Sprintf("Unexpected num args: %v", arg.NumArgs))
		}
	} else {
		if arg.NumArgs == 0 {
			self.emitWithValue(tokError,
				fmt.Sprintf("The %s switch does not take a value", text))
			return nil
		} else if arg.NumArgs == 1 {
			self.emitWithValue(tokValue, rhs)
			self.pos += 1
			return self.stateArgument
		} else if arg.NumArgs > 1 {
			panic("not implemented yet")
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

	//        log.Printf("nextPositional=%d numEvaluated=%d numRequired=%d numMax=%d",
	//            self.nextPositionalArgument, self.numEvaluatedPositionalArguments, self.numRequiredPositionalArguments, self.numMaxPositionalArguments)

	// Is there more than enough required positional arguments, but there could be more?
	//if self.numEvaluatedPositionalArguments > self.numRequiredPositionalArguments && self.numMaxPositionalArguments == -1 {
	if self.cmd.numMaxPositionalArguments == -1 {
		arg := self.args[self.pos]
		posArg := self.cmd.positionalArguments[self.nextPositionalArgument]
		self.emitWithArgument(tokArgument, posArg, posArg.Name)
		self.emitWithValue(tokValue, arg)
		self.pos += 1
		if posArg.NumArgs == 1 || posArg.NumArgsGlob == "?" {
			self.nextPositionalArgument++
		}
		self.numEvaluatedPositionalArguments++
		return self.statePositionalArgument
	}

	// We still have required positional arguments to check
	if self.numEvaluatedPositionalArguments < self.cmd.numRequiredPositionalArguments {
		arg := self.args[self.pos]
		posArg := self.cmd.positionalArguments[self.nextPositionalArgument]
		self.emitWithArgument(tokArgument, posArg, posArg.Name)
		// If only one arg is allowed, then go to the next positional argument
		if posArg.NumArgs == 1 {
			self.nextPositionalArgument++
		}
		self.emitWithValue(tokValue, arg)
		self.pos += 1
		self.numEvaluatedPositionalArguments++
		return self.statePositionalArgument
	} else if self.numEvaluatedPositionalArguments < self.cmd.numMaxPositionalArguments {
		arg := self.args[self.pos]
		posArg := self.cmd.positionalArguments[self.nextPositionalArgument]
		self.emitWithArgument(tokArgument, posArg, posArg.Name)
		// If only one arg is allowed, then go to the next positional argument
		if posArg.NumArgsGlob == "?" {
			self.nextPositionalArgument++
		}
		self.emitWithValue(tokValue, arg)
		self.pos += 1
		self.numEvaluatedPositionalArguments++
		return self.statePositionalArgument
	} else {
		arg := self.args[self.pos]
		self.emitWithValue(tokError, fmt.Sprintf("Unexpected positional argument: %s", arg))
		return nil
	}
}

// Consume the rest of the args
/*
func (self *parserState) statePassThrough() stateFunc {
	for ; self.pos < len(self.args); self.pos++ {
		arg := self.args[self.pos]
		self.emitWithArgument(tokArgument, self.stickyArg, self.stickyArg.String)
		self.emitWithValue(tokValue, arg)
	}
	return nil
}
*/
