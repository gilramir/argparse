// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type TestParseValues struct {
	String string
}

func (self *TestParseValues) Run(values []Destination) (error) {
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
		Long:        "--string",
//		Type:        "",
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
		Name:        "string",
		//Type:        "",
		NumArgs:     '1',
		Help: "Required string value",
	})

	// No string argument passed after subcommand
	argv := []string{"subcommand"}
	err := p0.ParseArgv(argv)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "Expected a required string argument")
}
