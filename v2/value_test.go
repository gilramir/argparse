// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package argparse

import (
	"reflect"

	. "gopkg.in/check.v1"
)

type TestValues struct {
	Bool   bool
	String string
	Int    int
	Float  float64
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
	parserVal := NewBoolValueT(valueP)
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
	// Find the pointer to Bool
	field, found := structType.FieldByName("String")
	c.Assert(found, Equals, true)
	valueP := structValue.FieldByIndex(field.Index)

	// Create our valueT
	parserVal := NewStringValueT(valueP)
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
	// Find the pointer to Bool
	field, found := structType.FieldByName("Int")
	c.Assert(found, Equals, true)
	valueP := structValue.FieldByIndex(field.Index)

	// Create our valueT
	parserVal := NewIntValueT(valueP)
	c.Assert(v.Int, Equals, 0)

	// seenWithoutValue does not work
	err := parserVal.seenWithoutValue()
	c.Assert(err, NotNil)

	// parse works
	v.Int = 0
	err = parserVal.parse(&DefaultMessages_en, "100")
	c.Assert(err, IsNil)
	c.Check(v.Int, Equals, 100)

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

func (s *MySuite) TestValueFloat(c *C) {
	v := &TestValues{}

	ptrValue := reflect.ValueOf(v)
	structValue := reflect.Indirect(ptrValue)
	structType := structValue.Type()
	// Find the pointer to Bool
	field, found := structType.FieldByName("Float")
	c.Assert(found, Equals, true)
	valueP := structValue.FieldByIndex(field.Index)

	// Create our valueT
	parserVal := NewFloatValueT(valueP)
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
