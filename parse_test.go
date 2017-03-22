// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type TestParseValues struct {
	String string
}

func (self *TestParseValues) Run(values []Destination) (cliOK bool, err error) {
	return true, nil
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
		Type:        "",
		Description: "Pass a string value",
	})

	argv := []string{"--help"}
	err := p.ParseArgv(argv)
	c.Assert(err, IsNil)

	c.Check(buffer.String(), Equals, `progname

Options:
	--string=STRING     Pass a string value
`)

}
