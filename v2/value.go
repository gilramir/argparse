package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

// This file implements the type system for argument values.

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// textUnmarshalerType is the reflect.Type of encoding.TextUnmarshaler, used to
// detect fields (and slice elements) of custom types that know how to parse
// themselves from a string.
var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

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
	seenWithoutValue(m *Messages) error

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

func (self *boolValueT) seenWithoutValue(m *Messages) error {
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

func (self *stringValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedStringValue)
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

func (self *intValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedIntValue)
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
		return fmt.Errorf(m.CannotParseIntegerFmt, text, err)
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

func (self *int64ValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedInt64Value)
}

func (self *int64ValueT) parse(m *Messages, text string) error {
	i, err := text_to_int64(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseIntegerFmt, text, err)
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

func (self *floatValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedFloatValue)
}

func (self *floatValueT) parse(m *Messages, text string) error {
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf(m.CannotParseFloatFmt, text)
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

func (self *durationValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedDurationValue)
}

func (self *durationValueT) parse(m *Messages, text string) error {
	d, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseDurationFmt, text, err)
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

func (self *boolSliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedBoolValue)
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

func (self *stringSliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedStringValue)
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

func (self *intSliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedIntValue)
}

func (self *intSliceValueT) parse(m *Messages, text string) error {
	i64, err := text_to_int64(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseIntegerFmt, text, err)
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

func (self *int64SliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedInt64Value)
}

func (self *int64SliceValueT) parse(m *Messages, text string) error {
	i, err := text_to_int64(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseIntegerFmt, text, err)
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

func (self *floatSliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedFloatValue)
}

func (self *floatSliceValueT) parse(m *Messages, text string) error {
	f, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf(m.CannotParseFloatFmt, text)
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

func (self *durationSliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedDurationValue)
}

func (self *durationSliceValueT) parse(m *Messages, text string) error {
	d, err := time.ParseDuration(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseDurationFmt, text, err)
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

// =========================================================== uint (all unsigned kinds)

func text_to_uint64(text string) (uint64, error) {
	if len(text) > 2 && text[0:2] == "0x" {
		return strconv.ParseUint(text[2:], 16, 64)
	} else if len(text) > 2 && text[0:2] == "0o" {
		return strconv.ParseUint(text[2:], 8, 64)
	} else if len(text) > 1 && text[0:1] == "0" {
		return strconv.ParseUint(text[1:], 8, 64)
	}
	return strconv.ParseUint(text, 10, 64)
}

// uintValueT handles a scalar field of any unsigned integer kind (uint,
// uint8, uint16, uint32, uint64).
type uintValueT struct {
	valueT
	choices []uint64
}

func newUintValueT(valueP reflect.Value) *uintValueT {
	return &uintValueT{valueT: valueT{valueP}}
}

func (self *uintValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *uintValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedUintValue)
}

func (self *uintValueT) parse(m *Messages, text string) error {
	u, err := text_to_uint64(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseUintFmt, text, err)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if u == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	self.value.SetUint(u)
	return nil
}

func uintChoices(m *Messages, choicesIntf interface{}) ([]uint64, error) {
	choices, ok := choicesIntf.([]uint)
	if !ok {
		return nil, fmt.Errorf(m.ChoicesOfWrongTypeFmt, "uint")
	}
	c64 := make([]uint64, len(choices))
	for i, v := range choices {
		c64[i] = uint64(v)
	}
	return c64, nil
}

func (self *uintValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	c, err := uintChoices(m, choicesIntf)
	if err != nil {
		return err
	}
	self.choices = c
	return nil
}

func (self *uintValueT) storageType() valueStorageType {
	return Scalar
}

// =========================================================== uint slice

type uintSliceValueT struct {
	valueT
	choices []uint64
}

func newUintSliceValueT(valueP reflect.Value) *uintSliceValueT {
	return &uintSliceValueT{valueT: valueT{valueP}}
}

func (self *uintSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *uintSliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedUintValue)
}

func (self *uintSliceValueT) parse(m *Messages, text string) error {
	u, err := text_to_uint64(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseUintFmt, text, err)
	}
	if len(self.choices) > 0 {
		ok := false
		for _, choice := range self.choices {
			if u == choice {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf(m.ShouldBeAValidChoiceFmt, self.choices)
		}
	}
	elem := reflect.New(self.value.Type().Elem()).Elem()
	elem.SetUint(u)
	self.value.Set(reflect.Append(self.value, elem))
	return nil
}

func (self *uintSliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	c, err := uintChoices(m, choicesIntf)
	if err != nil {
		return err
	}
	self.choices = c
	return nil
}

func (self *uintSliceValueT) storageType() valueStorageType {
	return Slice
}

// =========================================================== signed int slice (int8/16/32)

// signedIntSliceValueT handles a slice field whose element is a sub-width
// signed integer (int8, int16, int32). []int and []int64 keep their own types.
type signedIntSliceValueT struct {
	valueT
	choices []int64
}

func newSignedIntSliceValueT(valueP reflect.Value) *signedIntSliceValueT {
	return &signedIntSliceValueT{valueT: valueT{valueP}}
}

func (self *signedIntSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *signedIntSliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedIntValue)
}

func (self *signedIntSliceValueT) parse(m *Messages, text string) error {
	i, err := text_to_int64(text)
	if err != nil {
		return fmt.Errorf(m.CannotParseIntegerFmt, text, err)
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
	elem := reflect.New(self.value.Type().Elem()).Elem()
	elem.SetInt(i)
	self.value.Set(reflect.Append(self.value, elem))
	return nil
}

func (self *signedIntSliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]int)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "int")
	}
	c64 := make([]int64, len(choices))
	for i, v := range choices {
		c64[i] = int64(v)
	}
	self.choices = c64
	return nil
}

func (self *signedIntSliceValueT) storageType() valueStorageType {
	return Slice
}

// =========================================================== float32 slice

type float32SliceValueT struct {
	valueT
	choices []float64
}

func newFloat32SliceValueT(valueP reflect.Value) *float32SliceValueT {
	return &float32SliceValueT{valueT: valueT{valueP}}
}

func (self *float32SliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *float32SliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedFloatValue)
}

func (self *float32SliceValueT) parse(m *Messages, text string) error {
	f, err := strconv.ParseFloat(text, 32)
	if err != nil {
		return fmt.Errorf(m.CannotParseFloatFmt, text)
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
	elem := reflect.New(self.value.Type().Elem()).Elem()
	elem.SetFloat(f)
	self.value.Set(reflect.Append(self.value, elem))
	return nil
}

func (self *float32SliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	choices, ok := choicesIntf.([]float32)
	if !ok {
		return fmt.Errorf(m.ChoicesOfWrongTypeFmt, "float32")
	}
	c64 := make([]float64, len(choices))
	for i, v := range choices {
		c64[i] = float64(v)
	}
	self.choices = c64
	return nil
}

func (self *float32SliceValueT) storageType() valueStorageType {
	return Slice
}

// =========================================================== custom (encoding.TextUnmarshaler)

// textUnmarshalerValueT handles a scalar field whose type implements
// encoding.TextUnmarshaler.
type textUnmarshalerValueT struct {
	valueT
}

func newTextUnmarshalerValueT(valueP reflect.Value) *textUnmarshalerValueT {
	return &textUnmarshalerValueT{valueT: valueT{valueP}}
}

func (self *textUnmarshalerValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *textUnmarshalerValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedValue)
}

func (self *textUnmarshalerValueT) parse(m *Messages, text string) error {
	u := self.value.Addr().Interface().(encoding.TextUnmarshaler)
	if err := u.UnmarshalText([]byte(text)); err != nil {
		return fmt.Errorf(m.CannotParseValueFmt, text, err)
	}
	return nil
}

func (self *textUnmarshalerValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	return errors.New(m.ChoicesNotSupportedForCustomType)
}

func (self *textUnmarshalerValueT) storageType() valueStorageType {
	return Scalar
}

// =========================================================== custom slice

// textUnmarshalerSliceValueT handles a slice field whose element type
// implements encoding.TextUnmarshaler.
type textUnmarshalerSliceValueT struct {
	valueT
}

func newTextUnmarshalerSliceValueT(valueP reflect.Value) *textUnmarshalerSliceValueT {
	return &textUnmarshalerSliceValueT{valueT: valueT{valueP}}
}

func (self *textUnmarshalerSliceValueT) defaultSwitchNumArgs() int {
	return 1
}

func (self *textUnmarshalerSliceValueT) seenWithoutValue(m *Messages) error {
	return errors.New(m.NeedValue)
}

func (self *textUnmarshalerSliceValueT) parse(m *Messages, text string) error {
	// Make a new element, unmarshal into it, then append it to the slice.
	elemPtr := reflect.New(self.value.Type().Elem())
	u := elemPtr.Interface().(encoding.TextUnmarshaler)
	if err := u.UnmarshalText([]byte(text)); err != nil {
		return fmt.Errorf(m.CannotParseValueFmt, text, err)
	}
	self.value.Set(reflect.Append(self.value, elemPtr.Elem()))
	return nil
}

func (self *textUnmarshalerSliceValueT) setChoices(m *Messages, choicesIntf interface{}) error {
	return errors.New(m.ChoicesNotSupportedForCustomType)
}

func (self *textUnmarshalerSliceValueT) storageType() valueStorageType {
	return Slice
}
