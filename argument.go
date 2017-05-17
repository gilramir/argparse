// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Argument struct {
	Type        interface{}
	Short       string
	Long        string
	Name        string
	Description string
	Metavar     string
	NumArgs     numArgsType
	value       reflect.Value
}

func (self *Argument) sanityCheck(dest Destination) {
	var err error
	err = self._sanityCheckName()
	if err != nil {
		panic(err.Error())
	}
	err = self._sanityCheckType()
	if err != nil {
		panic(err.Error())
	}
	err = self._sanityCheckNumArgs()
	if err != nil {
		panic(err.Error())
	}
	err = self._sanityCheckDestination(dest)
	if err != nil {
		panic(err.Error())
	}
}

func (self *Argument) _sanityCheckName() error {
	if self.Short != "" && self.Short[0] != '-' {
		return errors.New("The Short version of the argument must begin with '-'")
	}
	if self.Long != "" && len(self.Long) < 2 {
		return errors.New("The Long version of the argument must begin with '--'")
	}
	if self.Long != "" && self.Long[0:2] != "--" {
		return errors.New("The Short version of the argument must begin with '--'")
	}
	if self.Name != "" && self.Name[0] == '-' {
		return errors.New("The Name of a positional argument cannot begin with '-'")
	}

	if self.Short == "" && self.Long == "" {
		if self.Name == "" {
			return errors.New("No name/short/long given for Argument")
		}
	}
	if self.Short != "" || self.Long != "" {
		if self.Name != "" {
			return errors.New("Name cannot be given if short/long is given")
		}
	}
	return nil
}

func (self *Argument) _sanityCheckType() error {
	if self.Type == nil {
		return errors.New(fmt.Sprintf("Argument %s needs a Type field", self.prettyName()))
	}
	switch reflect.TypeOf(self.Type).Kind() {
	case reflect.Bool:
		// no-op
	case reflect.Int:
		// no-op
	case reflect.String:
		// no-op
	case reflect.Slice:
		// What kind of slice is it?
		switch reflect.TypeOf(self.Type).Elem().Kind() {
		case reflect.String:
			// no -op
		default:
			return errors.New(fmt.Sprintf("Argument %s cannot be of type %s",
				self.prettyName(), reflect.TypeOf(self.Type)))
		}
	default:
		return errors.New(fmt.Sprintf("Argument %s cannot be of type %s",
			self.prettyName(), reflect.TypeOf(self.Type)))
	}
	return nil
}

func (self *Argument) _sanityCheckNumArgs() error {
	// Was NumArgs not given?
	if self.NumArgs == '\x00' {
		switch reflect.TypeOf(self.Type).Kind() {
		case reflect.Bool:
			self.NumArgs = '0'
		case reflect.Int:
			self.NumArgs = '1'
		case reflect.String:
			self.NumArgs = '1'
		case reflect.Slice:
			// What kind of slice is it?
			switch reflect.TypeOf(self.Type).Elem().Kind() {
			case reflect.String:
				self.NumArgs = '1'
			default:
				return errors.New(fmt.Sprintf("Argument %s cannot be of type %s",
					self.prettyName(), reflect.TypeOf(self.Type)))
			}
		default:
			return errors.New(fmt.Sprintf("Argument %s cannot be of type %s",
				self.prettyName(), reflect.TypeOf(self.Type)))
		}
	} else {
		// XXX
		/// check short/long/named
	}
	return nil
}

// Capitalize the first rune of a string; utf-8 compatible
func firstRuneUpper(orig string) (string, error) {
	if len(orig) == 0 {
		return "", nil
	}
	reader := strings.NewReader(orig)
	ch, size, err := reader.ReadRune()
	if err != nil {
		return "", err
	}
	newString := strings.ToUpper(string(ch)) + orig[size:len(orig)]
	return newString, nil
}

// Remove punctuation and convert the next character to CamelCase.
// e.g., No-checkout => NoCheckout
func toSafeCamelCase(orig string) (string, error) {
	var newString string
	var capitalizeNext bool
	var chString string

	if len(orig) == 0 {
		return "", nil
	}

	reader := strings.NewReader(orig)
	ch, size, err := reader.ReadRune()
	if err != nil {
		return "", err
	}

	for size != 0 {
		if capitalizeNext {
			chString = strings.ToUpper(string(ch))
			capitalizeNext = false
		} else {
			chString = string(ch)
		}

		if ch == '-' || ch == '.' || ch == '_' {
			capitalizeNext = true
		} else {
			newString += chString
		}
		ch, size, err = reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
	}

	return newString, nil
}

func argumentVariableName(orig string) string {
	newString, err := firstRuneUpper(orig)
	if err != nil {
		panic(fmt.Sprintf("Error converting(1) '%s': %s", orig, err.Error()))
	}
	newString, err = toSafeCamelCase(newString)
	if err != nil {
		panic(fmt.Sprintf("Error converting(2) '%s': %s", newString, err.Error()))
	}
	return newString
}

// Check that we have we
func (self *Argument) _sanityCheckDestination(dest Destination) error {
	ptrValue := reflect.ValueOf(dest)
	structValue := reflect.Indirect(ptrValue)
	type_ := structValue.Type()

	var field reflect.StructField
	var found bool

	if self.Short != "" {
		shortStructName := argumentVariableName(self.Short[1:len(self.Short)])
		field, found = type_.FieldByName(shortStructName)
	}
	if !found && self.Long != "" {
		longStructName := argumentVariableName(self.Long[2:len(self.Long)])
		field, found = type_.FieldByName(longStructName)
	}
	if !found && self.Name != "" {
		structName := argumentVariableName(self.Name)
		field, found = type_.FieldByName(structName)
	}
	if !found {
		return errors.New(fmt.Sprintf("Could not find destination variable for argument %s",
			self.prettyName()))
	}

	fieldIndex := field.Index
	self.value = structValue.FieldByIndex(fieldIndex)
	return nil
}

func (self *Argument) prettyName() string {
	if self.Long == "" {
		if self.Short != "" {
			return self.Short
		} else if self.Name != "" {
			return self.Name
		} else {
			panic("Unexpected")
		}
	} else if self.Short == "" {
		if self.Long != "" {
			return self.Long
		} else if self.Name != "" {
			return self.Name
		} else {
			panic("Unexpected")
		}
	} else {
		if self.Name != "" {
			panic("Unexpected")
		}
		return self.Short + "/" + self.Long
	}
	panic("Not reached")
	return ""
}

func (self *Argument) isSwitch() bool {
	return self.Short != "" || self.Long != ""
}

func (self *Argument) isPositional() bool {
	return !(self.isSwitch())
}

func (self *Argument) Parse(text string) error {
	var err error

	switch reflect.TypeOf(self.Type).Kind() {
	case reflect.Bool:
		var boolValue bool
		boolValue, err = strconv.ParseBool(text)
		if err != nil {
			return errors.Errorf("Cannot convert \"%s\" to a boolean", text)
		}
		self.value.SetBool(boolValue)
	case reflect.Int:
		var i int
		i, err = strconv.Atoi(text)
		if err != nil {
			return errors.Errorf("Cannot convert \"%s\" to an integer", text)
		}
		self.value.SetInt(int64(i))
	case reflect.String:
		self.value.SetString(text)
	case reflect.Slice:
		switch reflect.TypeOf(self.Type).Elem().Kind() {
		case reflect.String:
			newValue := reflect.ValueOf(text)
			self.value.Set(reflect.Append(self.value, newValue))
		default:
			panic("Should not reach")
		}
	default:
		panic("Should not reach")
	}
	return nil
}

func (self *Argument) Seen() {
	switch reflect.TypeOf(self.Type).Kind() {
	case reflect.Bool:
		self.value.SetBool(true)
	default:
		panic("Should not reach")
	}
}

func (self *Argument) SeenWithoutValue() error {
	panic("not yet implemented")
	return nil
}

func (self *Argument) GetMetavar() string {
	if self.Metavar != "" {
		return self.Metavar
	} else if self.Name != "" {
		return strings.ToUpper(self.Name)
	} else if self.Long != "" {
		return strings.ToUpper(self.Long[2:])
	} else if self.Short != "" {
		return strings.ToUpper(self.Short[1:])
	} else {
		panic("Should not reach")
	}
	return ""
}

func (self *Argument) HelpString() string {
	var text string

	if self.Short != "" {
		text = self.Short
		if self.Long != "" {
			text += ","
		}
	}
	if self.Long != "" {
		text += self.Long
	}
	if self.NumArgs != numArgs0 {
		text += "=" + self.GetMetavar()
	}
	return text
}

func (self *Argument) dump(spaces string) {
	if self.Short != "" {
		fmt.Printf("%sShort: %s\n", spaces, self.Short)
	}
	if self.Long != "" {
		fmt.Printf("%sLong: %s\n", spaces, self.Long)
	}
	if self.Name != "" {
		fmt.Printf("%sName: %s\n", spaces, self.Name)
	}
	fmt.Printf("%sType: %q\n", spaces, self.Type)
	fmt.Printf("%sDescription: %s\n", spaces, self.Description)
	if self.Metavar != "" {
		fmt.Printf("%sMetavar: %s\n", spaces, self.Metavar)
	}
	if self.NumArgs != '0' && self.NumArgs != '1' {
		fmt.Printf("%sNumArgs: %s\n", spaces, string(self.NumArgs))
	}
	fmt.Printf("\n")
}
