// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	. "gopkg.in/check.v1"
)

type TestArgumentValues struct {
	Bool    bool
	String  string
	Strings []string
	Integer int
}

func (self *TestArgumentValues) Run(values []Destination) (error) {
	return nil
}

func (s *MySuite) TestArgumentIsSwitch(c *C) {
	switchArgShort := &Argument{
		Short: "-s",
	}
	c.Check(switchArgShort.isSwitch(), Equals, true)
	c.Check(switchArgShort.isPositional(), Equals, false)
	c.Check(switchArgShort.isCommand(), Equals, false)

	switchArgLong := &Argument{
		Long: "--switch",
	}
	c.Check(switchArgLong.isSwitch(), Equals, true)
	c.Check(switchArgLong.isPositional(), Equals, false)
	c.Check(switchArgLong.isCommand(), Equals, false)
}

func (s *MySuite) TestArgumentIsPositional(c *C) {
	positionalArg := &Argument{
		Name: "s",
	}
	c.Check(positionalArg.isSwitch(), Equals, false)
	c.Check(positionalArg.isPositional(), Equals, true)
	c.Check(positionalArg.isCommand(), Equals, false)
}

func (s *MySuite) TestArgumentIsCommand(c *C) {
	commandArg := &Argument{
                ParseCommand: PassThrough,
		Short: "--",
	}
	c.Check(commandArg.isSwitch(), Equals, false)
	c.Check(commandArg.isPositional(), Equals, false)
	c.Check(commandArg.isCommand(), Equals, true)
}

func (s *MySuite) TestArgumentPrettyName(c *C) {
	arg := &Argument{
		Short: "-s",
	}
	c.Check(arg.prettyName(), Equals, "-s")
}

func (s *MySuite) TestArgumentToSafeCamelCase(c *C) {
	var newString string
	var err error

	newString, err = toSafeCamelCase("No-fun")
	c.Assert(err, IsNil)
	c.Check(newString, Equals, "NoFun")

	newString, err = toSafeCamelCase("No__fun-really")
	c.Assert(err, IsNil)
	c.Check(newString, Equals, "NoFunReally")
}

func (s *MySuite) TestArgumentParseString(c *C) {
	v := &TestArgumentValues{}

	arg := &Argument{
		Name: "string",
		//Type: "",
	}
	arg.sanityCheck(v)
	arg.Parse("foo")
	c.Check(v.String, Equals, "foo")
}

func (s *MySuite) TestArgumentParseStringSlice(c *C) {
	v := &TestArgumentValues{}

	arg := &Argument{
		Name: "strings",
//		Type: []string{},
	}
	arg.sanityCheck(v)
	arg.Parse("foo")
	arg.Parse("bar")
	c.Check(len(v.Strings), Equals, 2)
	c.Check(v.Strings[0], Equals, "foo")
	c.Check(v.Strings[1], Equals, "bar")
}

func (s *MySuite) TestArgumentParseBool(c *C) {
	v := &TestArgumentValues{}

	arg := &Argument{
		Name: "bool",
//		Type: false,
	}
	arg.sanityCheck(v)
	arg.Parse("true")
	c.Check(v.Bool, Equals, true)
}

func (s *MySuite) TestArgumentParseInt(c *C) {
	v := &TestArgumentValues{}

	arg := &Argument{
		Name: "integer",
//		Type: 0,
	}
	arg.sanityCheck(v)
	arg.Parse("42")
	c.Check(v.Integer, Equals, 42)
}
