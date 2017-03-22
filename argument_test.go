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

func (self *TestArgumentValues) Run(values []Destination) (cliOK bool, err error) {
	return true, nil
}

func (s *MySuite) TestArgumentPrettyname(c *C) {
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
		Type: "",
	}
	arg.sanityCheck(v)
	arg.Parse("foo")
	c.Check(v.String, Equals, "foo")
}

func (s *MySuite) TestArgumentParseStringSlice(c *C) {
	v := &TestArgumentValues{}

	arg := &Argument{
		Name: "strings",
		Type: []string{},
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
		Type: false,
	}
	arg.sanityCheck(v)
	arg.Parse("true")
	c.Check(v.Bool, Equals, true)
}

func (s *MySuite) TestArgumentParseInt(c *C) {
	v := &TestArgumentValues{}

	arg := &Argument{
		Name: "integer",
		Type: 0,
	}
	arg.sanityCheck(v)
	arg.Parse("42")
	c.Check(v.Integer, Equals, 42)
}
