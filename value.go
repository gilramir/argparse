package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	//"strings"
)

type valueType interface {

	// Parse the text into the destination value
	parse(m *Messages, text string) error

	// If the switch is seen but has no value after it.
	// This is only legal for bools
	seenWithoutValue() error

	defaultSwitchNumArgs() int
}

type valueT struct {
	// The golang field type (Kind) where the parsed value will be stored
//	typeKind reflect.Kind

	// If typeKind is a Slice, then it's a slice of what?
//	sliceKind reflect.Kind

	// A "pointer" to where to store the parsed value
	value reflect.Value
}

// =========================================================== bool

type boolValueT struct {
	valueT
}

func NewBoolValueT( valueP reflect.Value ) boolValueT {
	return boolValueT{ valueT: valueT { valueP } }
}

func (self boolValueT) defaultSwitchNumArgs() int {
	return 0
}

func (self boolValueT) seenWithoutValue() (error) {
	self.value.SetBool(true)
	return nil
}

func (self boolValueT) parse(m *Messages, text string) error {
	var val bool
	val, err := strconv.ParseBool(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseBooleanFmt, text)
	}
	self.value.SetBool(val)
	return nil
}


// =========================================================== string

type stringValueT struct {
	valueT
}

func NewStringValueT( valueP reflect.Value ) stringValueT {
	return stringValueT{ valueT: valueT { valueP } }
}

func (self stringValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self stringValueT) seenWithoutValue() (error) {
	return errors.New("Need a string value")
}

func (self stringValueT) parse(m *Messages, text string) error {
	self.value.SetString(text)
	return nil
}

// =========================================================== int

type intValueT struct {
	valueT
}

func NewIntValueT( valueP reflect.Value ) intValueT {
	return intValueT{ valueT: valueT { valueP } }
}

func (self intValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self intValueT) seenWithoutValue() (error) {
	return errors.New("Need an int value")
}

func (self intValueT) parse(m *Messages, text string) error {
	i, err := strconv.Atoi(text)
	//i, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to an integer", text)
	}
	self.value.SetInt(int64(i))
	return nil
}

// =========================================================== float

type floatValueT struct {
	valueT
}

func NewFloatValueT( valueP reflect.Value ) floatValueT {
	return floatValueT{ valueT: valueT { valueP } }
}

func (self floatValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self floatValueT) seenWithoutValue() (error) {
	return errors.New("Need an float value")
}

func (self floatValueT) parse(m *Messages, text string) error {
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to an float", text)
	}
	self.value.SetFloat(f)
	return nil
}

// =========================================================== bool slice

type boolSliceValueT struct {
	valueT
}

func NewBoolSliceValueT( valueP reflect.Value ) boolSliceValueT {
	return boolSliceValueT{ valueT: valueT { valueP } }
}

func (self boolSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self boolSliceValueT) seenWithoutValue() (error) {
	return errors.New("Need a bool value")
}

func (self boolSliceValueT) parse(m *Messages, text string) error {
	var val bool
	val, err := strconv.ParseBool(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseBooleanFmt, text)
	}
	itemValue := reflect.ValueOf(val)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}


// =========================================================== string slice

type stringSliceValueT struct {
	valueT
}

func NewStringSliceValueT( valueP reflect.Value ) stringSliceValueT {
	return stringSliceValueT{ valueT: valueT { valueP } }
}

func (self stringSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self stringSliceValueT) seenWithoutValue() (error) {
	return errors.New("Need a string value")
}

func (self stringSliceValueT) parse(m *Messages, text string) error {
	itemValue := reflect.ValueOf(text)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}

// =========================================================== int slice

type intSliceValueT struct {
	valueT
}

func NewIntSliceValueT( valueP reflect.Value ) intSliceValueT {
	return intSliceValueT{ valueT: valueT { valueP } }
}

func (self intSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self intSliceValueT) seenWithoutValue() (error) {
	return errors.New("Need an int value")
}

func (self intSliceValueT) parse(m *Messages, text string) error {
	i, err := strconv.Atoi(text)
	//i, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to an integer", text)
	}
//	self.value.SetInt(int64(i))
	itemValue := reflect.ValueOf(i)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}

// =========================================================== float slice

type floatSliceValueT struct {
	valueT
}

func NewFloatSliceValueT( valueP reflect.Value ) floatSliceValueT {
	return floatSliceValueT{ valueT: valueT { valueP } }
}

func (self floatSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self floatSliceValueT) seenWithoutValue() (error) {
	return errors.New("Need an float value")
}

func (self floatSliceValueT) parse(m *Messages, text string) error {
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to a float", text)
	}
	itemValue := reflect.ValueOf(f)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}
