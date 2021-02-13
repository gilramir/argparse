package argparse

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	//	"strconv"
	"strings"
)

type Argument struct {
	// Any number of switch patterns, each starting with at lease one hypen
	Switches []string

	// The name of the positional argument. No starting hyphens.
	Name string

	// The help string to display to the user
	Help string

	// The name of the value field to show in the usage statement,
	// for non-boolean switches
	MetaVar string

	// The name of the destination field for the value of the switch
	// or positional argument, if it is named differently from any of the
	// Switches, or Name.
	Dest string

	// Number of arguments that can or should appear
	// If NumArgs is 0 (never initialized), and NumArgsGlob is "",
	// the value of NumArgs is set to 1, unless this is a Bool, in which case
	// it's set to 0.
	// If NumArgs is not 0 or 1, then NumArgsGlob must be "", in which case
	// the number of args is exactly NumArgs.
	// If NumArgs is 0, and NumArgsGlob is not "", then it must be one
	// of "+" ("one or more"), "?" ("zero or one"), or "*" ("zero or more"),
	// and then NumArgs is set to -1
	NumArgs     int
	NumArgsGlob string

	// Is the user required to provide this argument? This is only
	// checked for Switch arguments; positional arguments make use of
	// NumArgs or NumArgsGlob
	//	Required bool

	// Will a sub-command inherit this argument definition if one is not
	// defined for that sub-command, *and* if the Value struct for that
	// Command has a suitable field?
	Inherit bool

	// For non-boolean options, the valid values that the user can provide.
	// If Choices is given, and the user provides a value not in this list,
	// the user will be presented with an error.
	Choices interface{}

	// The methods for the specific storage type of this value of the Argument
	// (bool, int, string, float64, etc.)
	value valueType
}

func (self *Argument) deepCopy() *Argument {
	arg := &Argument{
		Switches:    make([]string, len(self.Switches)),
		Name:        self.Name,
		Help:        self.Help,
		MetaVar:     self.MetaVar,
		Dest:        self.Dest,
		NumArgs:     self.NumArgs,
		NumArgsGlob: self.NumArgsGlob,
		//		Required: self.Required,
		Inherit: self.Inherit,
		Choices: self.Choices,
	}
	copy(arg.Switches, self.Switches)
	return arg
}

func (self *Argument) init(dest Values, messages *Messages) {
	var err error
	// Ensure that there is some name field set
	err = self.sanityCheckNameAndSwitches()
	if err != nil {
		panic(err.Error())
	}
	// Check the type of value in the destination struct
	// This is the side-effect of setting self.typeKind and self.value
	err = self.sanityCheckValueType(dest)
	if err != nil {
		panic(err.Error())
	}

	// Any Choices?
	if self.Choices != nil {
		err = self.value.setChoices(messages, self.Choices)
		if err != nil {
			panic(fmt.Sprintf("Argument %s: %s", self.PrettyName(),
				err.Error()))
		}
	}
}

func (self *Argument) sanityCheckNameAndSwitches() error {
	for i, switchName := range self.Switches {
		if len(switchName) == 0 {
			return fmt.Errorf("Switch #%d is an empty string", i+1)
		}
		if switchName[0] != '-' {
			return fmt.Errorf("Switch #%d '%s' should begin with '-'", i+1, switchName)
		}
	}

	if self.Name != "" && self.Name[0] == '-' {
		return errors.New("The Name of a positional argument cannot begin with '-'")
	}

	if len(self.Switches) == 0 && self.Name == "" {
		return errors.New("No Switches or Name given for Argument")
	}
	if len(self.Switches) > 0 && self.Name != "" {
		return errors.New("Name cannot be given if Switches is given")
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

// Check that there is a field in the Values struct that corresponds
// to this argument.
func (self *Argument) sanityCheckValueType(dest Values) error {

        // dest is Values, which is an interface{} type.
        // TypeOf(dest) gives us the dynamic type, the pointer to a
        // user-defined struct that was given to argparse.
	userStructPtrValue := reflect.ValueOf(dest)

        // Indirecting points us to the user's struct
        userStructValue := reflect.Indirect(userStructPtrValue)
        userStructType := userStructValue.Type()

	var field reflect.StructField
	var found bool
	var needles []string

	if self.Dest != "" {
		field, found = userStructType.FieldByName(self.Dest)
		if !found {
			return errors.New(fmt.Sprintf("Could not find destination field for argument %s, given as %s",
				self.PrettyName(), self.Dest))
		}
	} else {
		for _, switchName := range self.Switches {
			structName := argumentVariableName(switchName[1:])
			needles = append(needles, structName)
			field, found = userStructType.FieldByName(structName)
			if found {
				self.Dest = field.Name
				break
			}
		}
		if !found && self.Name != "" {
			structName := argumentVariableName(self.Name)
			needles = append(needles, structName)
			field, found = userStructType.FieldByName(structName)
			if found {
				self.Dest = field.Name
			}
		}
		if !found {
			return errors.New(fmt.Sprintf("Could not find destination field for argument %s; checked %s",
				self.PrettyName(), strings.Join(needles, ",")))
		}
	}

	// By using the index of the field within the struct type,
	// we can get the corresponding struct value
	fieldValue := userStructValue.FieldByIndex(field.Index)
        fieldType := fieldValue.Type()
        fieldTypeKind := fieldType.Kind()
//	typeKind := field.Type.Kind()
//        fmt.Printf("field %s has type %s but kind type %s\n", self.Dest, field.Type.String(),
//        typeKind.String())
	switch fieldType.String() {
        case "time.Duration":
            self.value = newDurationValueT(fieldValue)
            return nil
        }

	switch fieldTypeKind {
	case reflect.Bool:
		self.value = newBoolValueT(fieldValue)
	case reflect.String:
		self.value = newStringValueT(fieldValue)
	case reflect.Int:
		self.value = newIntValueT(fieldValue)
	case reflect.Float64:
		self.value = newFloatValueT(fieldValue)
	case reflect.Slice:
		sliceKind := fieldValue.Type().Elem().Kind()
		switch sliceKind {
		case reflect.Bool:
			self.value = newBoolSliceValueT(fieldValue)
		case reflect.Int:
			self.value = newIntSliceValueT(fieldValue)
		case reflect.String:
			self.value = newStringSliceValueT(fieldValue)
		case reflect.Float64:
			self.value = newFloatSliceValueT(fieldValue)
		default:
			return fmt.Errorf("Argument %s cannot be of type []%s",
				self.PrettyName(), sliceKind.String())
		}
	default:
		return errors.New(fmt.Sprintf("Argument %s cannot be of type %s",
			self.PrettyName(), fieldType.String()))
	}
	return nil
}

func (self *Argument) PrettyName() string {
	if len(self.Switches) > 0 {
		return strings.Join(self.Switches, "/")
	} else if self.Name != "" {
		return self.Name
	}
	panic("Argument has neither Switches or Name.")
}

func (self *Argument) isSwitch() bool {
	if len(self.Switches) > 0 {
		if self.Name == "" {
			return true
		} else {
			panic(fmt.Sprintf("Argument %s has Switches and Name",
				self.PrettyName()))
		}
	} else {
		if self.Name == "" {
			panic(fmt.Sprintf("Argument %s has neither Switches nor Name",
				self.PrettyName()))
		} else {
			return false
		}
	}
}

func (self *Argument) isPositional() bool {
	return !self.isSwitch()
}

/*

func (self *Argument) getChoicesString() string {
	switch self.typeKind {
	case reflect.Bool:
		panic("not reached")

	case reflect.Int:
		return nonQuotedListString(self.Choices)

	case reflect.String:
		return quotedListString(self.Choices)

	case reflect.Slice:
		switch self.sliceKind {
		case reflect.String:
			return quotedListString(self.Choices)
		default:
			panic("Should not reach")
		}
	default:
		panic("Should not reach")
	}
	return ""
}

func nonQuotedListString(choices []string) string {
	var text string

	// Special handling for the 2-choice case
	if len(choices) == 2 {
		return choices[0] + " and " + choices[1]
	}

	for i := 0; i < len(choices); i++ {
		if i == 0 {
			text += choices[i]
		} else if i == len(choices)-1 {
			text += ", and " + choices[i]
		} else {
			text += ", " + choices[i]
		}
	}
	return text
}

func quotedListString(choices []string) string {
	var text string

	// Special handling for the 2-choice case
	if len(choices) == 2 {
		return "'" + choices[0] + "' and '" + choices[1] + "'"
	}

	for i := 0; i < len(choices); i++ {
		if i == 0 {
			text += "'" + choices[i] + "'"
		} else if i == len(choices)-1 {
			text += ", and '" + choices[i] + "'"
		} else {
			text += ", '" + choices[i] + "'"
		}
	}
	return text
}

func (self *Argument) getMetavar() string {
	if self.Metavar != "" {
		return self.Metavar
	} else if self.Name != "" {
		return strings.ToUpper(self.Name)
	} else if len(self.Switches) > 0 {
		firstSwitch := strings.ToUpper(self.Switches[0])
		if firstSwitch[1] == '-' {
			return strings.ToUpper(firstSwitch[2:])
		} else {
			return strings.ToUpper(firstSwitch[1:])
		}
	} else {
		panic("Should not reach")
	}
	return ""
}

func (self *Argument) helpString() string {
	var text string

	if len(self.Switches) > 0 {
		text += strings.Join(self.Switches, ",")
	}
	if self.NumArgs != numArgs0 {
		text += "=" + self.getMetavar()
	}
	return text
}

func (self *Argument) helpStrings() []string {
	strings := make([]string, 0, len(self.Switches))

	for i, switchText := range self.Switches {
		text := switchText
		if self.NumArgs != numArgs0 {
			text += "=" + self.getMetavar()
		}
		if i < len(self.Switches) - 1 {
			text = text + ","
		}
		strings = append(strings, text)
	}
	return strings
}
*/
