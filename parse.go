// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	"fmt"

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

	nextPositionalArgument int
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
	//	var lastSubParser *ArgumentParser

	for argToken := range parser.tokenChan {
		switch argToken.typ {
		case tokArgument:
			lastArgument = argToken.argument
			lastArgLabel = argToken.argumentLabel
			if lastArgument.NumArgs == numArgs0 {
				lastArgument.Seen()
			}

		case tokValue:
			if lastArgument == nil {
				panic("Found value without a preceding argument")
			}
			err := lastArgument.Parse(argToken.value)
			if err != nil {
				results.parseError = errors.Wrapf(err,
					"While parsing value for %s", lastArgLabel)
				return results
			}
		case tokValueNotPresent:
			if lastArgument == nil {
				panic("Found ValueNotPresent without a preceding argument")
			}
			err := lastArgument.SeenWithoutValue()
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
	if len(self.positionalArguments) > 0 && parser.nextPositionalArgument < len(self.positionalArguments) {
		arg := self.positionalArguments[parser.nextPositionalArgument]
		if arg.NumArgs == numArgs1 || arg.NumArgs == numArgsStar {
			results.parseError = errors.Errorf("Expected %s argument (required)", arg.prettyName())
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

	if len(arg) > 1 && arg[0] == '-' {
		if len(arg) > 2 {
			if arg[0:2] == "--" {
				return self.stateLongOption
			}

			// Grouped short arguments
			panic("Cannot handle grouped short arguments yet")
		} else {
			// Short arguments, including "--" (two hyphens)
			return self.stateShortOption
		}

		parser.emitWithValue(tokError, fmt.Sprintf("Unknown argument: %s", arg))
		return nil
	}

	// Positional argument?
	if parser.nextPositionalArgument == 0 && len(self.positionalArguments) > 0 {
		return self.statePositionalArgument
	}

	parser.emitWithValue(tokError, fmt.Sprintf("Unexpected argument: %s", arg))
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

func (self *ArgumentParser) stateLongOption(parser *parserState) stateFunc {
	arg := parser.args[parser.pos]
	if arg == "" {
		parser.emitWithValue(tokError, "<empty string>")
		return nil
	}
	longName := arg[2:]
	if longName == "help" {
		parser.emitToken(tokHelp)
		return nil
	}
	for _, arg := range self.switchArguments {
		if arg.Long == "--"+longName {
			parser.emitWithArgument(tokArgument, arg, arg.Long)
			parser.pos += 1
			parser.lastSwitch = "--" + longName
			switch arg.NumArgs {
			case numArgs0:
				return self.stateArgument
			case numArgs1:
				return self.stateOneValue
			case numArgsMaybe:
				return self.stateMaybeOneValue
			case numArgsStar:
				return self.stateMultipleValues
			default:
				panic(fmt.Sprintf("Unexpected num args: %v", arg.NumArgs))
			}
		}
	}
	// Didn't find a long arg with that name
	parser.emitWithValue(tokError,
		fmt.Sprintf("No such switch: --%s", longName))
	return nil
}

func (self *ArgumentParser) stateShortOption(parser *parserState) stateFunc {
	arg := parser.args[parser.pos]
	if arg == "" {
		parser.emitWithValue(tokError, "<empty string>")
		return nil
	}
	shortName := arg[1:2]
	for _, arg := range self.switchArguments {
		if arg.Short == "-"+shortName {
			parser.emitWithArgument(tokArgument, arg, arg.Short)
			parser.pos += 1
			parser.lastSwitch = "-" + shortName
			switch arg.NumArgs {
			case numArgs0:
				return self.stateArgument
			case numArgs1:
				return self.stateOneValue
			case numArgsStar:
				return self.stateMultipleValues
			default:
				panic(fmt.Sprintf("Unexpected num args: %v", arg.NumArgs))
			}
		}
	}

	parser.emitWithValue(tokError, fmt.Sprintf("Unknown argument: %s", arg))
	return nil
}

func (self *ArgumentParser) statePositionalArgument(parser *parserState) stateFunc {
	if parser.pos == len(parser.args) {
		// End of the list
		return nil
	}
	arg := parser.args[parser.pos]
	if len(self.positionalArguments) > parser.nextPositionalArgument {
		posArg := self.positionalArguments[parser.nextPositionalArgument]
		parser.emitWithArgument(tokArgument, posArg, posArg.Name)
		parser.nextPositionalArgument++
		parser.emitWithValue(tokValue, arg)
		parser.pos += 1
		return self.statePositionalArgument
	} else {
		parser.emitWithValue(tokError, fmt.Sprintf("Unexpected argument %s", arg))
		return nil
	}
}
