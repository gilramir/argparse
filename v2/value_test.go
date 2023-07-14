package argparse

// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"reflect"
	"time"

	. "gopkg.in/check.v1"
)

type TestValues struct {
	Bool     bool
	String   string
	Int      int
	Int64    int64
	Float    float64
	Duration time.Duration
}

func (s *MySuite) TestValueBool(c *C) {
	v := &TestValues{}

	ptrValue := reflect.ValueOf(v)
	structValue := reflect.Indirect(ptrValue)
	structType := structValue.Type()
	// Find the pointer to Bool
	field, found := structType.FieldByName("Bool")
	c.Assert(found, Equals, true)
	valueP := structValue.FieldByIndex(field.Index)

	// Create our valueT
	parserVal := newBoolValueT(valueP)
	c.Assert(v.Bool, Equals, false)

	// seenWithoutValue works
	err := parserVal.seenWithoutValue()
	c.Assert(err, IsNil)
	c.Check(v.Bool, Equals, true)

	// parse works
	v.Bool = false
	err = parserVal.parse(&DefaultMessages_en, "true")
	c.Assert(err, IsNil)
	c.Check(v.Bool, Equals, true)

	err = parserVal.parse(&DefaultMessages_en, "false")
	c.Assert(err, IsNil)
	c.Check(v.Bool, Equals, false)

	err = parserVal.parse(&DefaultMessages_en, "incorrect")
	c.Assert(err, NotNil)

	// Set Choices
	err = parserVal.setChoices(&DefaultMessages_en, []bool{true})
	c.Assert(err, IsNil)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "true")
	c.Assert(err, IsNil)
	c.Check(v.Bool, Equals, true)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "false")
	c.Check(err, NotNil)
}

func (s *MySuite) TestValueString(c *C) {
	v := &TestValues{}

	ptrValue := reflect.ValueOf(v)
	structValue := reflect.Indirect(ptrValue)
	structType := structValue.Type()
	// Find the pointer to String
	field, found := structType.FieldByName("String")
	c.Assert(found, Equals, true)
	valueP := structValue.FieldByIndex(field.Index)

	// Create our valueT
	parserVal := newStringValueT(valueP)
	c.Assert(v.String, Equals, "")

	// seenWithoutValue does not work
	err := parserVal.seenWithoutValue()
	c.Assert(err, NotNil)

	// parse works
	v.String = ""
	err = parserVal.parse(&DefaultMessages_en, "dog")
	c.Assert(err, IsNil)
	c.Check(v.String, Equals, "dog")

	err = parserVal.parse(&DefaultMessages_en, "")
	c.Assert(err, IsNil)
	c.Check(v.String, Equals, "")

	err = parserVal.parse(&DefaultMessages_en, "a b c")
	c.Assert(err, IsNil)
	c.Check(v.String, Equals, "a b c")

	// Set Choices
	err = parserVal.setChoices(&DefaultMessages_en, []string{"x", "y"})
	c.Assert(err, IsNil)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "x")
	c.Assert(err, IsNil)
	c.Check(v.String, Equals, "x")

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "y")
	c.Assert(err, IsNil)
	c.Check(v.String, Equals, "y")

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "z")
	c.Check(err, NotNil)
}

func (s *MySuite) TestValueInt(c *C) {
	v := &TestValues{}

	ptrValue := reflect.ValueOf(v)
	structValue := reflect.Indirect(ptrValue)
	structType := structValue.Type()
	// Find the pointer to Int
	field, found := structType.FieldByName("Int")
	c.Assert(found, Equals, true)
	valueP := structValue.FieldByIndex(field.Index)

	// Create our valueT
	parserVal := newIntValueT(valueP)
	c.Assert(v.Int, Equals, 0)

	// seenWithoutValue does not work
	err := parserVal.seenWithoutValue()
	c.Assert(err, NotNil)

	// parse works
	v.Int = 0
	err = parserVal.parse(&DefaultMessages_en, "100")
	c.Assert(err, IsNil)
	c.Check(v.Int, Equals, 100)

	// hex
	err = parserVal.parse(&DefaultMessages_en, "0xff")
	c.Assert(err, IsNil)
	c.Check(v.Int, Equals, 255)

	// octal with "0o"
	err = parserVal.parse(&DefaultMessages_en, "0o776")
	c.Assert(err, IsNil)
	c.Check(v.Int, Equals, 510)

	// octal with "0"
	err = parserVal.parse(&DefaultMessages_en, "0775")
	c.Assert(err, IsNil)
	c.Check(v.Int, Equals, 509)

	err = parserVal.parse(&DefaultMessages_en, "-55")
	c.Assert(err, IsNil)
	c.Check(v.Int, Equals, -55)

	err = parserVal.parse(&DefaultMessages_en, "abc")
	c.Assert(err, NotNil)

	err = parserVal.parse(&DefaultMessages_en, "5.7")
	c.Assert(err, NotNil)

	// Set Choices
	err = parserVal.setChoices(&DefaultMessages_en, []int{3, 5})
	c.Assert(err, IsNil)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "3")
	c.Assert(err, IsNil)
	c.Check(v.Int, Equals, 3)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "5")
	c.Assert(err, IsNil)
	c.Check(v.Int, Equals, 5)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "17")
	c.Check(err, NotNil)
}

func (s *MySuite) TestValueInt64(c *C) {
	v := &TestValues{}

	ptrValue := reflect.ValueOf(v)
	structValue := reflect.Indirect(ptrValue)
	structType := structValue.Type()
	// Find the pointer to Int64
	field, found := structType.FieldByName("Int64")
	c.Assert(found, Equals, true)
	valueP := structValue.FieldByIndex(field.Index)

	// Create our valueT
	parserVal := newInt64ValueT(valueP)
	c.Assert(v.Int64, Equals, int64(0))

	// seenWithoutValue does not work
	err := parserVal.seenWithoutValue()
	c.Assert(err, NotNil)

	// parse works
	v.Int64 = 64
	err = parserVal.parse(&DefaultMessages_en, "100")
	c.Assert(err, IsNil)
	c.Check(v.Int64, Equals, int64(100))

	// hex
	err = parserVal.parse(&DefaultMessages_en, "0xffff")
	c.Assert(err, IsNil)
	c.Check(v.Int64, Equals, int64(65535))

	// octal with "0o"
	err = parserVal.parse(&DefaultMessages_en, "0o776")
	c.Assert(err, IsNil)
	c.Check(v.Int64, Equals, int64(510))

	// octal with "0"
	err = parserVal.parse(&DefaultMessages_en, "0775")
	c.Assert(err, IsNil)
	c.Check(v.Int64, Equals, int64(509))

	err = parserVal.parse(&DefaultMessages_en, "-55")
	c.Assert(err, IsNil)
	c.Check(v.Int64, Equals, int64(-55))

	err = parserVal.parse(&DefaultMessages_en, "abc")
	c.Assert(err, NotNil)

	err = parserVal.parse(&DefaultMessages_en, "5.7")
	c.Assert(err, NotNil)

	// Set Choices
	err = parserVal.setChoices(&DefaultMessages_en, []int64{3, 5})
	c.Assert(err, IsNil)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "3")
	c.Assert(err, IsNil)
	c.Check(v.Int64, Equals, int64(3))

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "5")
	c.Assert(err, IsNil)
	c.Check(v.Int64, Equals, int64(5))

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "17")
	c.Check(err, NotNil)
}

func (s *MySuite) TestValueFloat(c *C) {
	v := &TestValues{}

	ptrValue := reflect.ValueOf(v)
	structValue := reflect.Indirect(ptrValue)
	structType := structValue.Type()
	// Find the pointer to Float
	field, found := structType.FieldByName("Float")
	c.Assert(found, Equals, true)
	valueP := structValue.FieldByIndex(field.Index)

	// Create our valueT
	parserVal := newFloatValueT(valueP)
	c.Assert(v.Float, Equals, 0.0)

	// seenWithoutValue does not work
	err := parserVal.seenWithoutValue()
	c.Assert(err, NotNil)

	// parse works
	v.Float = 0
	err = parserVal.parse(&DefaultMessages_en, "100")
	c.Assert(err, IsNil)
	c.Check(v.Float, Equals, 100.0)

	// Do we need to worry about i18n? (, instead of .)?
	err = parserVal.parse(&DefaultMessages_en, "-55.2")
	c.Assert(err, IsNil)
	c.Check(v.Float, Equals, -55.2)

	err = parserVal.parse(&DefaultMessages_en, "abc")
	c.Assert(err, NotNil)

	// Set Choices
	err = parserVal.setChoices(&DefaultMessages_en, []float64{33.3, 55.0})
	c.Assert(err, IsNil)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "33.3")
	c.Assert(err, IsNil)
	c.Check(v.Float, Equals, 33.3)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "33.300")
	c.Assert(err, IsNil)
	c.Check(v.Float, Equals, 33.3)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "55")
	c.Assert(err, IsNil)
	c.Check(v.Float, Equals, 55.0)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "17.22")
	c.Check(err, NotNil)
}

func (s *MySuite) TestValueDuration(c *C) {
	v := &TestValues{}

	ptrValue := reflect.ValueOf(v)
	structValue := reflect.Indirect(ptrValue)
	structType := structValue.Type()
	// Find the pointer to Duration
	field, found := structType.FieldByName("Duration")
	c.Assert(found, Equals, true)
	valueP := structValue.FieldByIndex(field.Index)

	// Create our valueT
	parserVal := newDurationValueT(valueP)
	c.Assert(v.Duration.Seconds(), Equals, 0.0)

	// seenWithoutValue does not work
	err := parserVal.seenWithoutValue()
	c.Assert(err, NotNil)

	// set v.Int to some value, and then check that
	// parsing sets it do a different value
	v.Duration, err = time.ParseDuration("1s")
	c.Assert(err, IsNil)
	c.Assert(v.Duration.Seconds(), Equals, 1.0)

	err = parserVal.parse(&DefaultMessages_en, "30s")
	c.Assert(err, IsNil)
	c.Check(v.Duration.Seconds(), Equals, 30.0)

	err = parserVal.parse(&DefaultMessages_en, "5m")
	c.Assert(err, IsNil)
	c.Check(v.Duration.Seconds(), Equals, 300.0)

	err = parserVal.parse(&DefaultMessages_en, "not-a-duration")
	c.Assert(err, NotNil)

	err = parserVal.parse(&DefaultMessages_en, "5")
	c.Assert(err, NotNil)

	d1min, _ := time.ParseDuration("1m")
	d2min, _ := time.ParseDuration("2m")

	// Set Choices
	err = parserVal.setChoices(&DefaultMessages_en, []time.Duration{d1min, d2min})
	c.Assert(err, IsNil)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "60s")
	c.Assert(err, IsNil)
	c.Check(v.Duration.Seconds(), Equals, 60.0)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "120s")
	c.Assert(err, IsNil)
	c.Check(v.Duration.Seconds(), Equals, 120.0)

	// Test choice
	err = parserVal.parse(&DefaultMessages_en, "2s")
	c.Check(err, NotNil)
}
