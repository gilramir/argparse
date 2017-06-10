package argparse

import (
	"fmt"
	"io"
	//        "log"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	kNilRune = '\x00'

	PassThrough = 1
)

type Argument struct {
	Short        string
	Long         string
	Name         string
	Help         string
	Metavar      string
	Dest         string

	ParseCommand    int
        String          string

	// Number of arguments that can or should appear, only for positional arguments
	NumArgs numArgsType

	// The golang field type (Kind) where the parsed value will be stored
	typeKind reflect.Kind
	// If typeKind is a Slice, then it's a slice of what?
	sliceKind reflect.Kind

	// A "pointer" to where to store the parsed value
	value reflect.Value
}

func (self *Argument) sanityCheck(dest Destination) {
	var err error
	// Ensure that there is some name field set
	err = self._sanityCheckName()
	if err != nil {
		panic(err.Error())
	}
	// Check the type of value in the destination struct
	// This is the side-effect of setting self.typeKind and self.value
	err = self._sanityCheckDestination(dest)
	if err != nil {
		panic(err.Error())
	}

	err = self._sanityCheckNumArgs()
	if err != nil {
		panic(err.Error())
	}
}

func (self *Argument) _sanityCheckName() error {

        if self.ParseCommand != 0 {
            if self.Short != "" || self.Long != "" || self.Name != "" {
                return errors.New("A ParseCommand cannot have a Short, Long, or Name field")
            }
            if self.String == "" {
                return errors.New("A ParseCommand must have a String field")
            }
            return nil
        }

        if self.String != "" {
                return errors.New("String cannot be set if ParseCommand is not set")
        }

	if self.Short != "" && self.Short[0] != '-' {
		return errors.New("The Short version of the argument must begin with '-'")
	}
	if self.Long != "" && len(self.Long) < 2 {
		return errors.New("The Long version of the argument must begin with '--'")
	}
	if self.Long != "" && self.Long[0:2] != "--" {
		return errors.New("The Long version of the argument must begin with '--'")
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

// Check that there is a field in the destination struct that correponds
// to this argument.
func (self *Argument) _sanityCheckDestination(dest Destination) error {
	// TODO - some sanity checks here would be great
        // TODO - if PassThrough, check that it's a slice

	ptrValue := reflect.ValueOf(dest)
	structValue := reflect.Indirect(ptrValue)
	structType := structValue.Type()

	var field reflect.StructField
	var found bool
	var needles []string

	if self.Dest != "" {
		field, found = structType.FieldByName(self.Dest)
		if !found {
			return errors.New(fmt.Sprintf("Could not find destination field for argument %s, given as %s",
				self.prettyName(), self.Dest))
		}
	} else {
		if self.Short != "" {
			shortStructName := argumentVariableName(self.Short[1:len(self.Short)])
			needles = append(needles, shortStructName)
			field, found = structType.FieldByName(shortStructName)
		}
		if !found && self.Long != "" {
			longStructName := argumentVariableName(self.Long[2:len(self.Long)])
			needles = append(needles, longStructName)
			field, found = structType.FieldByName(longStructName)
		}
		if !found && self.Name != "" {
			structName := argumentVariableName(self.Name)
			needles = append(needles, structName)
			field, found = structType.FieldByName(structName)
		}
		if !found {
			return errors.New(fmt.Sprintf("Could not find destination field for argument %s; checked %s",
				self.prettyName(), strings.Join(needles, ",")))
		}
	}

	// By using the index of the field within the struct type,
	// we can get the corresponding struct value
	self.value = structValue.FieldByIndex(field.Index)
	self.typeKind = field.Type.Kind()
	if self.typeKind == reflect.Slice {
		// Get the type of slice
		self.sliceKind = self.value.Type().Elem().Kind()
	}
	//        log.Printf("field=%v index=%d typeKind=%v value=%v sliceKind=%v",
	//            field, field.Index, self.typeKind, self.value, self.sliceKind)
	return nil
}

func (self *Argument) _sanityCheckNumArgs() error {
	// Was NumArgs not given?
	if self.NumArgs == kNilRune {
		//switch reflect.TypeOf(self.Type).Kind() {
		switch self.typeKind {
		case reflect.Bool:
			self.NumArgs = '0'
		case reflect.Int:
			self.NumArgs = '1'
		case reflect.String:
			self.NumArgs = '1'
		case reflect.Slice:
			// What kind of slice is it?
			switch self.sliceKind {
			case reflect.String:
				self.NumArgs = '1'
			default:
				return errors.New(fmt.Sprintf("Argument %s cannot be of type []%s",
					self.prettyName(), self.sliceKind.String()))
			}
		default:
			return errors.New(fmt.Sprintf("Argument %s cannot be of type %s",
				self.prettyName(), self.typeKind.String()))
		}
	} else {
		// TODO - only allow setting NumArgs for positional arguments
		/// check short/long/named
	}
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
	return !self.isCommand() && (self.Short != "" || self.Long != "")
}

func (self *Argument) isPositional() bool {
	return !self.isCommand() && !(self.isSwitch())
}

func (self *Argument) isCommand() bool {
	return self.ParseCommand != 0
}

func (self *Argument) parse(text string) error {
	var err error

	//switch reflect.TypeOf(self.Type).Kind() {
	switch self.typeKind {
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
		//switch reflect.TypeOf(self.Type).Elem().Kind() {
		switch self.sliceKind {
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

func (self *Argument) seen() {
	//switch reflect.TypeOf(self.Type).Kind() {
	switch self.typeKind {
	case reflect.Bool:
		self.value.SetBool(true)
	default:
		panic("Should not reach")
	}
}

func (self *Argument) seenWithoutValue() error {
	panic("not yet implemented")
	return nil
}

func (self *Argument) getMetavar() string {
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

func (self *Argument) helpString() string {
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
		text += "=" + self.getMetavar()
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
	//	fmt.Printf("%sType: %q\n", spaces, self.Type)
	fmt.Printf("%sHelp: %s\n", spaces, self.Help)
	if self.Metavar != "" {
		fmt.Printf("%sMetavar: %s\n", spaces, self.Metavar)
	}
	if self.NumArgs != '0' && self.NumArgs != '1' {
		fmt.Printf("%sNumArgs: %s\n", spaces, string(self.NumArgs))
	}
	fmt.Printf("\n")
}
