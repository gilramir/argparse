// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse
/*
import (
	. "gopkg.in/check.v1"
)
*/
/*
type TestArgumentValues struct {
	Bool    bool
	String  string
	Strings []string
	Integer int
}

func (self *TestArgumentValues) Run(values []Destination) error {
	return nil
}

func (s *MySuite) TestArgumentIsSwitch(c *C) {
	switchArgShort := &Argument{
		Switches: []string{"-s"},
	}
	c.Check(switchArgShort.isSwitch(), Equals, true)
	c.Check(switchArgShort.isPositional(), Equals, false)
	c.Check(switchArgShort.isCommand(), Equals, false)

	switchArgLong := &Argument{
		Switches: []string{"--switch"},
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
		Switches:        []string{"--"},
	}
	c.Check(commandArg.isSwitch(), Equals, false)
	c.Check(commandArg.isPositional(), Equals, false)
	c.Check(commandArg.isCommand(), Equals, true)
}

func (s *MySuite) TestArgumentPrettyName(c *C) {
	arg := &Argument{
		Switches: []string{"-s"},
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
	arg.parse("foo")
	c.Check(v.String, Equals, "foo")
}

func (s *MySuite) TestArgumentParseStringSlice(c *C) {
	v := &TestArgumentValues{}

	arg := &Argument{
		Name: "strings",
		//		Type: []string{},
	}
	arg.sanityCheck(v)
	arg.parse("foo")
	arg.parse("bar")
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
	arg.parse("true")
	c.Check(v.Bool, Equals, true)
}

func (s *MySuite) TestArgumentParseInt(c *C) {
	v := &TestArgumentValues{}

	arg := &Argument{
		Name: "integer",
		//		Type: 0,
	}
	arg.sanityCheck(v)
	arg.parse("42")
	c.Check(v.Integer, Equals, 42)
}

func (s *MySuite) TestNonQuotedListString(c *C) {
	input := []string{"1", "5", "10"}
	c.Check(nonQuotedListString(input), Equals, "1, 5, and 10")
}

func (s *MySuite) TestQuotedListString(c *C) {
	input := []string{"a", "x", "zz"}
	c.Check(quotedListString(input), Equals, "'a', 'x', and 'zz'")
}
*/
