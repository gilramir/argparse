// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type TestParseValues struct {
	String      string
	Strings		[]string
	PassThrough []string
}

func (self *TestParseValues) Run(values []Destination) error {
	return nil
}

func (s *MySuite) TestParseHelpNoOptions(c *C) {
	var buffer bytes.Buffer

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      &TestParseValues{},
		Stdout:           &buffer,
	}

	argv := []string{"--help"}
	err := p.ParseArgv(argv)
	c.Assert(err, IsNil)

	c.Check(buffer.String(), Equals, `progname

No options
`)

}

func (s *MySuite) TestParseHelpOptions(c *C) {
	var buffer bytes.Buffer

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      &TestParseValues{},
		Stdout:           &buffer,
	}
	p.AddArgument(&Argument{
		Long: "--string",
		Help: "Pass a string value",
	})

	argv := []string{"--help"}
	err := p.ParseArgv(argv)
	c.Assert(err, IsNil)

	c.Check(buffer.String(), Equals, `progname

Options:
	--string=STRING     Pass a string value
`)

}

func (s *MySuite) TestParseRequiredPositionalArgument(c *C) {
	p0 := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      &TestParseValues{},
	}
	p1 := p0.AddParser(&ArgumentParser{
		Name:        "subcommand",
		Destination: &TestParseValues{},
	})
	p1.AddArgument(&Argument{
		Name:    "string",
		Help:    "Required string value",
	})

	// No string argument passed after subcommand
	argv := []string{"subcommand"}
	err := p0.ParseArgv(argv)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "Expected a required string argument")
}

func (s *MySuite) TestParseOneString(c *C) {
	values := &TestParseValues{}

	p0 := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p0.AddArgument(&Argument{
		Name:    "strings",
		Help:    "Required string value",
	})

	// No string argument passed after subcommand
	argv := []string{"foo", "bar", "baz"}
	err := p0.ParseArgv(argv)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "Unexpected positional argument: bar")
}

func (s *MySuite) TestParseOneOrMoreString(c *C) {
	values := &TestParseValues{}

	p0 := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p0.AddArgument(&Argument{
		Name:    "strings",
		Help:    "Required string value",
		NumArgs: '+',
	})

	// No string argument passed after subcommand
	argv := []string{"foo", "bar", "baz"}
	err := p0.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.Strings, DeepEquals, []string{"foo", "bar", "baz"})
}

func (s *MySuite) TestParsePassThroughAfterPositional(c *C) {
	values := &TestParseValues{}

	p0 := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p0.AddArgument(&Argument{
		ParseCommand: PassThrough,
		String:        "--",
		Dest:         "PassThrough",
	})
	p0.AddArgument(&Argument{
		Name:    "string",
//		NumArgs: '1',
		Help:    "Required string value",
	})

	// No string argument passed after subcommand
	argv := []string{"foo", "--", "a", "b", "c"}
	err := p0.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.String, Equals, "foo")
	c.Assert(len(values.PassThrough), Equals, 3)
	c.Assert(values.PassThrough, DeepEquals, []string{"a", "b", "c"})
}

func (s *MySuite) TestParsePassThroughAfterSwitch(c *C) {
	values := &TestParseValues{}

	p0 := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p0.AddArgument(&Argument{
		ParseCommand: PassThrough,
		String:        "--",
		Dest:         "PassThrough",
	})
	p0.AddArgument(&Argument{
		Long:    "--string",
		Help:    "Required string value",
	})

	// No string argument passed after subcommand
	argv := []string{"--string", "xxx", "--", "a", "b", "c"}
	err := p0.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.String, Equals, "xxx")
	c.Assert(len(values.PassThrough), Equals, 3)
	c.Assert(values.PassThrough, DeepEquals, []string{"a", "b", "c"})
}
