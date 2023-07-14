package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

// This file implements the type system for argument values.

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type valueStorageType int

const (
	Scalar valueStorageType = iota
	Slice
)

type valueType interface {

	// Parse the text into the destination value
	parse(m *Messages, text string) error

	// If the switch is seen but has no value after it.
	// This is only legal for bools
	seenWithoutValue() error

	defaultSwitchNumArgs() int

	setValue(reflect.Value)
	getValue() reflect.Value

	setChoices(m *Messages, itemsIntf interface{}) error

	storageType() valueStorageType
}

type valueT struct {
	// A "pointer" to where to store the parsed value
	value reflect.Value
}

func (self *valueT) getValue() reflect.Value {
	return self.value
}

func (self *valueT) setValue(valueP reflect.Value) {
	self.value.Set(valueP)
}

// =========================================================== bool

type boolValueT struct {
	valueT
	choices []bool
}

func newBoolValueT(valueP reflect.Value) *boolValueT {
	return &boolValueT{valueT: valueT{valueP}}
}

func (self *boolValueT) defaultSwitchNumArgs() int {
	return 0
}

func (self *boolValueT) seenWithoutValue() error {
	self.value.SetBool(true)
	return nil
}

func (self *boolValueT) parse(m *Messages, text string) error {
	var val bool
	val, err := strconv.ParseBool(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseBooleanFmt, text)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if val == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	self.value.SetBool(val)
	return nil
}

func (self *boolValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]bool)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "bool")
	}
	self.choices = choices
	return nil
}

func (self *boolValueT) storageType() valueStorageType {
	return Scalar
}

// =========================================================== string

type stringValueT struct {
	valueT
	choices []string
}

func newStringValueT(valueP reflect.Value) *stringValueT {
	return &stringValueT{valueT: valueT{valueP}}
}

func (self *stringValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *stringValueT) seenWithoutValue() error {
	return errors.New("Need a string value")
}

func (self *stringValueT) parse(m *Messages, text string) error {
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if text == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	self.value.SetString(text)
	return nil
}

func (self *stringValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]string)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "string")
	}
	self.choices = choices
	return nil
}

func (self *stringValueT) storageType() valueStorageType {
	return Scalar
}

// =========================================================== int

type intValueT struct {
	valueT
	choices []int
}

func newIntValueT(valueP reflect.Value) *intValueT {
	return &intValueT{valueT: valueT{valueP}}
}

func (self *intValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *intValueT) seenWithoutValue() error {
	return errors.New("Need an int value")
}

func text_to_int64(text string) (int64, error) {
	// hex?
	if len(text) > 2 && text[0:2] == "0x" {
		text = text[2:]
		return strconv.ParseInt(text, 16, 64)
	} else if len(text) > 2 && text[0:2] == "0o" {
		// octal with "0o"?
		text = text[2:]
		return strconv.ParseInt(text, 8, 64)
	} else if len(text) > 1 && text[0:1] == "0" {
		// octal with "0"?
		text = text[1:]
		return strconv.ParseInt(text, 8, 64)
	} else {
		// decimal
		return strconv.ParseInt(text, 10, 64)
	}
}

func (self *intValueT) parse(m *Messages, text string) error {
	i64, err := text_to_int64(text)
	i := int(i64)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to an integer: %w", text, err)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if i == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	self.value.SetInt(int64(i))
	return nil
}

func (self *intValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]int)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "int")
	}
	self.choices = choices
	return nil
}

func (self *intValueT) storageType() valueStorageType {
	return Scalar
}

// =========================================================== int64

type int64ValueT struct {
	valueT
	choices []int64
}

func newInt64ValueT(valueP reflect.Value) *int64ValueT {
	return &int64ValueT{valueT: valueT{valueP}}
}

func (self *int64ValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *int64ValueT) seenWithoutValue() error {
	return errors.New("Need an int64 value")
}

func (self *int64ValueT) parse(m *Messages, text string) error {
	i, err := text_to_int64(text)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to an integer: %w", text, err)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if i == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	self.value.SetInt(i)
	return nil
}

func (self *int64ValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]int64)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "int64")
	}
	self.choices = choices
	return nil
}

func (self *int64ValueT) storageType() valueStorageType {
	return Scalar
}

// =========================================================== float

type floatValueT struct {
	valueT
	choices []float64
}

func newFloatValueT(valueP reflect.Value) *floatValueT {
	return &floatValueT{valueT: valueT{valueP}}
}

func (self *floatValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *floatValueT) seenWithoutValue() error {
	return errors.New("Need a float value")
}

func (self *floatValueT) parse(m *Messages, text string) error {
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to a float", text)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if f == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	self.value.SetFloat(f)
	return nil
}

func (self *floatValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]float64)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "float64")
	}
	self.choices = choices
	return nil
}

func (self *floatValueT) storageType() valueStorageType {
	return Scalar
}

// =========================================================== time.Duration

type durationValueT struct {
	valueT
	choices []time.Duration
}

func newDurationValueT(valueP reflect.Value) *durationValueT {
	return &durationValueT{valueT: valueT{valueP}}
}

func (self *durationValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *durationValueT) seenWithoutValue() error {
	// TODO - needs to support i18n
	return errors.New("Need a time duration string")
}

func (self *durationValueT) parse(m *Messages, text string) error {
	d, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf("Cannot parse \"%s\" as a time duration: %s", text, err)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if d.Nanoseconds() == choice.Nanoseconds() {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	self.value.SetInt(int64(d.Nanoseconds()))
	return nil
}

func (self *durationValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]time.Duration)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "time duration string")
	}
	self.choices = choices
	return nil
}

func (self *durationValueT) storageType() valueStorageType {
	return Scalar
}

// =========================================================== bool slice

type boolSliceValueT struct {
	valueT
	choices []bool
}

func newBoolSliceValueT(valueP reflect.Value) *boolSliceValueT {
	return &boolSliceValueT{valueT: valueT{valueP}}
}

func (self *boolSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *boolSliceValueT) seenWithoutValue() error {
	return errors.New("Need a bool value")
}

func (self *boolSliceValueT) parse(m *Messages, text string) error {
	var val bool
	val, err := strconv.ParseBool(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseBooleanFmt, text)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if val == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	itemValue := reflect.ValueOf(val)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}

func (self *boolSliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]bool)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "bool")
	}
	self.choices = choices
	return nil
}

func (self *boolSliceValueT) storageType() valueStorageType {
	return Slice
}

// =========================================================== string slice

type stringSliceValueT struct {
	valueT
	choices []string
}

func newStringSliceValueT(valueP reflect.Value) *stringSliceValueT {
	return &stringSliceValueT{valueT: valueT{valueP}}
}

func (self *stringSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *stringSliceValueT) seenWithoutValue() error {
	return errors.New("Need a string value")
}

func (self *stringSliceValueT) parse(m *Messages, text string) error {
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if text == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	itemValue := reflect.ValueOf(text)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}

func (self *stringSliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]string)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "string")
	}
	self.choices = choices
	return nil
}

func (self *stringSliceValueT) storageType() valueStorageType {
	return Slice
}

// =========================================================== int slice

type intSliceValueT struct {
	valueT
	choices []int
}

func newIntSliceValueT(valueP reflect.Value) *intSliceValueT {
	return &intSliceValueT{valueT: valueT{valueP}}
}

func (self *intSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *intSliceValueT) seenWithoutValue() error {
	return errors.New("Need an int value")
}

func (self *intSliceValueT) parse(m *Messages, text string) error {
	i64, err := text_to_int64(text)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to an integer: %w", text, err)
	}
	i := int(i64)
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if i == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	itemValue := reflect.ValueOf(i)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}

func (self *intSliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]int)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "int")
	}
	self.choices = choices
	return nil
}

func (self *intSliceValueT) storageType() valueStorageType {
	return Slice
}

// =========================================================== int64 slice

type int64SliceValueT struct {
	valueT
	choices []int64
}

func newInt64SliceValueT(valueP reflect.Value) *int64SliceValueT {
	return &int64SliceValueT{valueT: valueT{valueP}}
}

func (self *int64SliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *int64SliceValueT) seenWithoutValue() error {
	return errors.New("Need an int64 value")
}

func (self *int64SliceValueT) parse(m *Messages, text string) error {
	i, err := text_to_int64(text)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to an int64: %w", text, err)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if i == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	itemValue := reflect.ValueOf(i)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}

func (self *int64SliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]int64)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "int64")
	}
	self.choices = choices
	return nil
}

func (self *int64SliceValueT) storageType() valueStorageType {
	return Slice
}

// =========================================================== float slice

type floatSliceValueT struct {
	valueT
	choices []float64
}

func newFloatSliceValueT(valueP reflect.Value) *floatSliceValueT {
	return &floatSliceValueT{valueT: valueT{valueP}}
}

func (self *floatSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *floatSliceValueT) seenWithoutValue() error {
	return errors.New("Need a float value")
}

func (self *floatSliceValueT) parse(m *Messages, text string) error {
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf("Cannot convert \"%s\" to a float", text)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if f == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	itemValue := reflect.ValueOf(f)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}

func (self *floatSliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]float64)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "float64")
	}
	self.choices = choices
	return nil
}

func (self *floatSliceValueT) storageType() valueStorageType {
	return Slice
}

// =========================================================== time.Duration slice

type durationSliceValueT struct {
	valueT
	choices []time.Duration
}

func newDurationSliceValueT(valueP reflect.Value) *durationSliceValueT {
	return &durationSliceValueT{valueT: valueT{valueP}}
}

func (self *durationSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *durationSliceValueT) seenWithoutValue() error {
	return errors.New("Need a time duration value")
}

func (self *durationSliceValueT) parse(m *Messages, text string) error {
	d, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf("Cannot parse \"%s\" as a time duration: %s", text, err)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if d.Nanoseconds() == choice.Nanoseconds() {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	//	self.value.SetInt(int64(d.Nanoseconds()))

	itemValue := reflect.ValueOf(d)
	self.value.Set(reflect.Append(self.value, itemValue))
	return nil
}

func (self *durationSliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]time.Duration)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "time duration string")
	}
	self.choices = choices
	return nil
}

func (self *durationSliceValueT) storageType() valueStorageType {
	return Slice
}
