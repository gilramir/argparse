// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	. "gopkg.in/check.v1"
)

type APTestOptions struct {
	Bool1   bool
	String1 string
}

func (s *MySuite) TestChooseHelpLongAtRoot(c *C) {

	opts := &APTestOptions{}
	ap := New(&Command{
		Description: "This is a test program",
		Values:      opts,
	})

	argv := []string{"--help"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(results.helpRequested, Equals, true)
	c.Check(results.triggeredCommand, Equals, ap.Root)
}

func (s *MySuite) TestChooseHelpShortAtRoot(c *C) {

	opts := &APTestOptions{}
	ap := New(&Command{
		Description: "This is a test program",
		Values:      opts,
	})

	argv := []string{"-h"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(results.helpRequested, Equals, true)
	c.Check(results.triggeredCommand, Equals, ap.Root)
}

func (s *MySuite) TestChooseHelpCustomAtRoot(c *C) {

	opts := &APTestOptions{}
	ap := New(&Command{
		Description: "This is a test program",
		Values:      opts,
	})
	ap.HelpSwitches = []string{"--custom-help"}

	argv := []string{"--custom-help"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(results.helpRequested, Equals, true)
	c.Check(results.triggeredCommand, Equals, ap.Root)
}

/*
type TestArgParseValues struct {
	String string
	ran    bool
}

func (self *TestArgParseValues) Run(values []Destination) (error) {
	self.ran = true
	return nil
}

//

func (s *MySuite) TestAddParser(c *C) {

	v0 := &TestArgParseValues{}
	v1 := &TestArgParseValues{}

	p := &ArgumentParser{
		Name:             "program-name",
		ShortDescription: "This is a simple program",
		Destination:      v0,
	}

	p1 := p.AddParser(&ArgumentParser{
		Name:             "subcommand",
		ShortDescription: "This is a subcommand",
		Destination:      v1,
	})

	c.Check(p, NotNil)
	c.Check(p1, NotNil)
	c.Assert(len(p.subParsers), Equals, 1)
	c.Check(p.subParsers[0], Equals, p1)
	c.Check(len(p1.subParsers), Equals, 0)
}

func (s *MySuite) TestArgRoot(c *C) {

	v0 := &TestArgParseValues{}
	v1 := &TestArgParseValues{}

	p := &ArgumentParser{
		Name:             "program-name",
		ShortDescription: "This is a simple program",
		Destination:      v0,
	}

	p.AddParser(&ArgumentParser{
		Name:             "subcommand",
		ShortDescription: "This is a subcommand",
		Destination:      v1,
	})

	argv := []string{}
	err := p.ParseArgv(argv)

	c.Assert(err, IsNil)
	c.Check(v0.ran, Equals, true)
	c.Check(v1.ran, Equals, false)
}

func (s *MySuite) TestArgChild1(c *C) {

	v0 := &TestArgParseValues{}
	v1 := &TestArgParseValues{}

	p := &ArgumentParser{
		Name:             "program-name",
		ShortDescription: "This is a simple program",
		Destination:      v0,
	}

	p.AddParser(&ArgumentParser{
		Name:             "subcommand",
		ShortDescription: "This is a subcommand",
		Destination:      v1,
	})

	argv := []string{"subcommand"}
	err := p.ParseArgv(argv)

	c.Assert(err, IsNil)
	c.Check(v0.ran, Equals, false)
	c.Check(v1.ran, Equals, true)
}

func (s *MySuite) TestArgString(c *C) {

	v := &TestArgParseValues{}

	p := &ArgumentParser{
		Name:             "program-name",
		ShortDescription: "This is a simple program",
		Destination:      v,
	}
	p.AddArgument(&Argument{
		Switches: []string{"--string"},
		//Type: "",
	})

	argv := []string{"--string", "foo"}
	results := p.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(results.triggeredParser, Equals, p)
	c.Check(v.String, Equals, "foo")
}

type TestArgParseSubcommandValues struct {
	String       string
	numAncestors int
	parentString string
}

func (self *TestArgParseSubcommandValues) Run(values []Destination) (err error) {
	self.numAncestors = len(values)
	if len(values) == 1 {
		self.parentString = values[0].(*TestArgParseValues).String
	}
	return nil
}

func (s *MySuite) TestArgDestinations(c *C) {

	v0 := &TestArgParseValues{}
	v1 := &TestArgParseSubcommandValues{}

	p0 := &ArgumentParser{
		Name:             "program-name",
		ShortDescription: "This is a simple program",
		Destination:      v0,
	}
	p0.AddArgument(&Argument{
		Switches: []string{"--string"},
		//Type: "",
	})

	p1 := p0.AddParser(&ArgumentParser{
		Name:             "subcommand",
		ShortDescription: "This is a subcommand",
		Destination:      v1,
	})
	p1.AddArgument(&Argument{
		Switches: []string{"--string"},
		//Type: "",
	})

	argv := []string{"--string", "foo", "subcommand", "--string", "bar"}
	err := p0.ParseArgv(argv)

	c.Assert(err, IsNil)
	c.Check(v1.String, Equals, "bar")
	c.Check(v1.numAncestors, Equals, 1)
	c.Check(v1.parentString, Equals, "foo")
}

func (s *MySuite) TestArgParseSubcommandHelp(c *C) {

	p0 := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      &TestParseValues{},
	}
	p1 := p0.AddParser(&ArgumentParser{
		Name:             "sub1",
		ShortDescription: "This is subcommand #1",
		Destination:      &TestParseValues{},
	})
	p2 := p1.AddParser(&ArgumentParser{
		Name:             "sub2",
		ShortDescription: "This is subcommand #2",
		Destination:      &TestParseValues{},
	})
	p3 := p2.AddParser(&ArgumentParser{
		Name:             "sub3",
		ShortDescription: "This is subcommand #3",
		Destination:      &TestParseValues{},
	})

	c.Check(p3.usageString(), Equals, `progname sub1 sub2 sub3

`)

}
*/
