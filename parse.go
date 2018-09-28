package argparse

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type numArgsType rune

const (
	// These probably can't be used by user
	numArgs0 numArgsType = '0' // short or long: no args
	numArgs1 numArgsType = '1' // short or long: 1 arg

	// These can be used by user
	numArgsStar  numArgsType = '*' // named: 0 or more present
	numArgsPlus  numArgsType = '+' // named: 1 or more present
	numArgsMaybe numArgsType = '?' // named: 0 or 1 present
)

type tokenType int

type argToken struct {
	typ      tokenType
	pos      int
	value    string
	argument *Argument

	// The name that was given to the argument by the user,
	// since an arg might have a short and a long value
	// in its definition.
	argumentLabel string

	parser *ArgumentParser
}

const (
	tokError tokenType = iota
	tokArgument
	tokValue
	tokValueNotPresent
	tokSubParser
	tokHelp
)

type parserState struct {
	pos        int
	args       []string
	tokenChan  chan argToken
	tokens     []argToken
	lastSwitch string

	// If there are sub parsers that could be present,
	// this starts as true. Once an arg is parsed, no
	// subparsers can be accepted, so it's changed to false.
	subParserAllowed bool

	nextPositionalArgument          int
	numEvaluatedPositionalArguments int

	// when we need to keep track of an *Argument across state transitions
	stickyArg *Argument
}

type stateFunc func(*parserState) stateFunc

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
func (self *parserState) emitParser(p *ArgumentParser) {
	self.tokenChan <- argToken{
		typ:    tokSubParser,
		pos:    self.pos,
		parser: p,
	}
}
func (self *parserState) emitToken(typ tokenType) {
	self.tokenChan <- argToken{
		typ: typ,
		pos: self.pos,
	}
}

type parseResults struct {
	parseError      error
	helpRequested   bool
	triggeredParser *ArgumentParser
	ancestors       []Destination
}

func (self *ArgumentParser) parseArgv(argv []string) *parseResults {
	results := &parseResults{
		triggeredParser: self,
	}

	var parser parserState
	parser.args = argv
	parser.tokenChan = make(chan argToken)

	parser.subParserAllowed = len(self.subParsers) > 0

	go self._parse(&parser)

	var lastArgLabel string
	var lastArgument *Argument

	for argToken := range parser.tokenChan {
		switch argToken.typ {
		case tokArgument:
			lastArgument = argToken.argument
			lastArgLabel = argToken.argumentLabel
			// If the argument is a boolean argument (no value), then
			// we mark it as seen and move on.
			if lastArgument.NumArgs == numArgs0 {
				lastArgument.seen()
			}

		case tokValue:
			if lastArgument == nil {
				panic("Found value without a preceding argument")
			}
			// Does the lastArgument have a Choices slice which limits
			// the valid values?
			if len(lastArgument.Choices) > 0 {
				good := false
				for _, choice := range lastArgument.Choices {
					if argToken.value == choice {
						good = true
						break
					}
				}
				if !good {
					results.parseError = errors.Errorf(
						"The possible values for %s are %s", lastArgLabel,
						lastArgument.getChoicesString())
					return results
				}
			}

			err := lastArgument.parse(argToken.value)
			if err != nil {
				results.parseError = errors.Wrapf(err,
					"While parsing value for %s", lastArgLabel)
				return results
			}
		case tokValueNotPresent:
			if lastArgument == nil {
				panic("Found ValueNotPresent without a preceding argument")
			}
			// TODO - why?
			err := lastArgument.seenWithoutValue()
			if err != nil {
				results.parseError = errors.Wrapf(err,
					"%s argument", lastArgLabel)
				return results
			}
		case tokSubParser:
			results.ancestors = append(results.ancestors,
				results.triggeredParser.Destination)
			results.triggeredParser = argToken.parser
		case tokHelp:
			results.helpRequested = true
			return results
		case tokError:
			results.parseError = errors.New(argToken.value)
			return results
		default:
			panic("Unhandled argToken type")
		}
		// XXX - maybe don't need parser.tokens ?
		parser.tokens = append(parser.tokens, argToken)
	}

	// Did we find all required parameters?
	// TODO - switchArgumants

	// If there aren't enough positional arguments, check the next known argument to see if it is required
	if len(results.triggeredParser.positionalArguments) > 0 && parser.numEvaluatedPositionalArguments < results.triggeredParser.numRequiredPositionalArguments {
		arg := results.triggeredParser.positionalArguments[parser.nextPositionalArgument]
		if arg.NumArgs == '1' || arg.NumArgs == '+' {
			results.parseError = errors.Errorf("Expected a required '%s' argument", arg.prettyName())
			return results
		}
	}

	return results
}

func (self *ArgumentParser) _parse(parser *parserState) {
	defer close(parser.tokenChan)

	var state stateFunc
	for state = self.stateArgument; state != nil; {
		state = state(parser)
	}
}

func (self *ArgumentParser) stateArgument(parser *parserState) stateFunc {
	if parser.pos == len(parser.args) {
		// End of the list
		return nil
	}

	arg := parser.args[parser.pos]
	if arg == "" {
		parser.emitWithValue(tokError, "<empty string>")
		return nil
	}

	// Is it a subparser?
	if parser.subParserAllowed {
		for _, subParser := range self.subParsers {
			if arg == subParser.Name {
				parser.emitParser(subParser)
				parser.pos += 1
				// The subparser can have its own subparsers
				parser.subParserAllowed = len(subParser.subParsers) > 0
				// Start parsing in the subParser!
				return subParser.stateArgument
			}
		}
	}

	// Is it a parse command?
	for _, commandArg := range self.commandArguments {
		if arg == commandArg.String {
			// A little ugly, since stateCommandArgument is going to check self.commandArguments again
			return self.stateCommandArgument
		}
	}

	// Is it a switch argument?
	if len(arg) > 1 && arg[0] == '-' {
		return self.stateOption
		parser.emitWithValue(tokError, fmt.Sprintf("Unknown argument: %s", arg))
		return nil
	}

	// Positional argument?
	if parser.nextPositionalArgument == 0 && len(self.positionalArguments) > 0 {
		return self.statePositionalArgument
	}

	parser.emitWithValue(tokError, fmt.Sprintf("Unexpected switch argument: %s", arg))
	return nil
}

func (self *ArgumentParser) stateMaybeOneValue(parser *parserState) stateFunc {
	if parser.pos == len(parser.args) {
		// Fine, we're finished.
		parser.emitToken(tokValueNotPresent)
		return nil
	}

	// Does the next token start with a hyphen?
	nextArg := parser.args[parser.pos]
	if len(nextArg) > 1 && nextArg[0] == '-' {
		// Okay, we have no value
		parser.emitToken(tokValueNotPresent)
		return self.stateArgument
	}

	// We do have a value
	parser.emitWithValue(tokValue, nextArg)
	parser.pos += 1
	return self.stateArgument
}

func (self *ArgumentParser) stateOneValue(parser *parserState) stateFunc {
	if parser.pos == len(parser.args) {
		parser.emitWithValue(tokError, fmt.Sprintf("Expected a value after %s", parser.lastSwitch))
		return nil
	}

	parser.emitWithValue(tokValue, parser.args[parser.pos])
	parser.pos += 1
	return self.stateArgument
}

func (self *ArgumentParser) stateMultipleValues(parser *parserState) stateFunc {
	if parser.pos == len(parser.args) {
		return nil
	}

	parser.emitWithValue(tokValue, parser.args[parser.pos])
	parser.pos += 1
	return self.stateMultipleValues
}

func (self *ArgumentParser) stateOption(parser *parserState) stateFunc {
	text := parser.args[parser.pos]
	if text == "" {
		parser.emitWithValue(tokError, "<empty string>")
		return nil
	}

	// Check for '=', as in --value=foo
	// XXX - add check to sanity check, ensureing '=' is not in switch name
	equalsIndex := strings.Index(text, "=")
	var rhs string
	if equalsIndex == 0 {
		parser.emitWithValue(tokError, "A switch name cannot begin with '='")
		return nil
	} else if equalsIndex > 0 {
		rhs = text[equalsIndex+1:]
		text = text[:equalsIndex]
	}

	if text == "--help" || text == "-h" {
		if rhs == "" {
			parser.emitToken(tokHelp)
			return nil
		} else {
			parser.emitWithValue(tokError, "-h/--help does not accept a value")
			return nil
		}
	}
	match := false
	var arg *Argument
	for _, arg = range self.switchArguments {
		for _, possibility := range arg.Switches {
			// Does it directly match a switch?
			if text == possibility {
				match = true
				break
			}
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
				parser.emitWithArgument(tokArgument, arg, text[:2])
				parser.lastSwitch = text[:2]

				// All other characters in the given switch must also be one-character
				// short-option booleans
				// TODO- think about utf-8 here
				all_others_good := true
				for _, r := range text[2:] {
					found := false
					for _, iArg := range self.switchArguments {
						if iArg.NumArgs == numArgs0 {
							for _, iSwitch := range iArg.Switches {
								if len(iSwitch) == 2 && rune(iSwitch[1]) == r {
									found = true
									// Emit this one
									parser.emitWithArgument(tokArgument, iArg, iSwitch)
									parser.lastSwitch = iSwitch
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
					parser.emitWithValue(tokError,
						fmt.Sprintf("The %s switch is valid but not as '%s'",
							text[:2], text))
					return nil
				}

				// We finished the parse and need to return successfully now
				parser.pos += 1
				return self.stateArgument
			}
			
		}

		if match {
			break
		}
	}	
	// Didn't match ?
	if ! match {
		// Didn't find a switch with that name
		parser.emitWithValue(tokError, fmt.Sprintf("No such switch: %s", text))
		return nil
	}

	parser.emitWithArgument(tokArgument, arg, text)
	parser.lastSwitch = text
	if rhs == "" {
		parser.pos += 1
		switch arg.NumArgs {
		case numArgs0:
			return self.stateArgument
		case numArgs1:
			return self.stateOneValue
		case numArgsMaybe:
			panic("not reached")
			//				return self.stateMaybeOneValue
		case numArgsStar:
			panic("not reached")
			//				return self.stateMultipleValues
		default:
			panic(fmt.Sprintf("Unexpected num args: %v", arg.NumArgs))
		}
	} else {
		switch arg.NumArgs {
		case numArgs0:
			parser.emitWithValue(tokError,
				fmt.Sprintf("The %s switch does not take a value", text))
			return nil
		case numArgs1:
			parser.emitWithValue(tokValue, rhs)
			parser.pos += 1
			return self.stateArgument
		case numArgsMaybe:
			panic("not reached")
		case numArgsStar:
			panic("not reached")
		default:
			panic(fmt.Sprintf("Unexpected num args: %v", arg.NumArgs))
		}
	}
	panic("not reached")
	//return nil
}


func (self *ArgumentParser) statePositionalArgument(parser *parserState) stateFunc {
	if parser.pos == len(parser.args) {
		// End of the list
		return nil
	}

	//        log.Printf("nextPositional=%d numEvaluated=%d numRequired=%d numMax=%d",
	//            parser.nextPositionalArgument, parser.numEvaluatedPositionalArguments, self.numRequiredPositionalArguments, self.numMaxPositionalArguments)

	// Is there more than enough required positional arguments, but there could be more?
	//if parser.numEvaluatedPositionalArguments > self.numRequiredPositionalArguments && self.numMaxPositionalArguments == -1 {
	if self.numMaxPositionalArguments == -1 {
		arg := parser.args[parser.pos]
		// Check for a command argument; it has precedence over an optional positional argument
		for _, commandArg := range self.commandArguments {
			if arg == commandArg.String {
				// A little ugly, since stateCommandArgument is going to check self.commandArguments again
				return self.stateCommandArgument
			}
		}
		// It was not a command argument; it was a positional argument
		posArg := self.positionalArguments[parser.nextPositionalArgument]
		parser.emitWithArgument(tokArgument, posArg, posArg.Name)
		parser.emitWithValue(tokValue, arg)
		parser.pos += 1
		if posArg.NumArgs == '1' || posArg.NumArgs == '?' {
			parser.nextPositionalArgument++
		}
		parser.numEvaluatedPositionalArguments++
		return self.statePositionalArgument
	}

	// We still have required positional arguments to check
	if parser.numEvaluatedPositionalArguments < self.numRequiredPositionalArguments {
		arg := parser.args[parser.pos]
		posArg := self.positionalArguments[parser.nextPositionalArgument]
		parser.emitWithArgument(tokArgument, posArg, posArg.Name)
		// If only one arg is allowed, then go to the next positional argument
		if posArg.NumArgs == '1' {
			parser.nextPositionalArgument++
		}
		parser.emitWithValue(tokValue, arg)
		parser.pos += 1
		parser.numEvaluatedPositionalArguments++
		return self.statePositionalArgument
	} else if parser.numEvaluatedPositionalArguments < self.numMaxPositionalArguments {
		arg := parser.args[parser.pos]
		posArg := self.positionalArguments[parser.nextPositionalArgument]
		parser.emitWithArgument(tokArgument, posArg, posArg.Name)
		// If only one arg is allowed, then go to the next positional argument
		if posArg.NumArgs == '?' {
			parser.nextPositionalArgument++
		}
		parser.emitWithValue(tokValue, arg)
		parser.pos += 1
		parser.numEvaluatedPositionalArguments++
		return self.statePositionalArgument
	} else if len(self.commandArguments) > 0 {
		// Is it a command argument?
		return self.stateCommandArgument
	} else {
		arg := parser.args[parser.pos]
		parser.emitWithValue(tokError, fmt.Sprintf("Unexpected positional argument: %s", arg))
		return nil
	}
}

func (self *ArgumentParser) stateCommandArgument(parser *parserState) stateFunc {
	arg := parser.args[parser.pos]

	for _, commandArg := range self.commandArguments {
		if arg == commandArg.String {
			switch commandArg.ParseCommand {
			case PassThrough:
				parser.stickyArg = commandArg
				parser.pos += 1
				return self.statePassThrough
			default:
				parser.emitWithValue(tokError, fmt.Sprintf("Unexpected ParseCommand value %s = %d", arg,
					commandArg.ParseCommand))
				return nil
			}
		}
	}
	// Didn't match
	parser.emitWithValue(tokError, fmt.Sprintf("Unexpected argument: %s", arg))
	return nil
}

// Consume the rest of the args
func (self *ArgumentParser) statePassThrough(parser *parserState) stateFunc {
	for ; parser.pos < len(parser.args); parser.pos++ {
		arg := parser.args[parser.pos]
		parser.emitWithArgument(tokArgument, parser.stickyArg, parser.stickyArg.String)
		parser.emitWithValue(tokValue, arg)
	}
	return nil
}
