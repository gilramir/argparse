// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type TestParseValues struct {
	String      string
	Strings     []string
	PassThrough []string
	Integer     int
	J           int
	X           bool
	Y           bool
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
		Switches: []string{"--string"},
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
		Name: "string",
		Help: "Required string value",
	})

	// No string argument passed after subcommand
	argv := []string{"subcommand"}
	err := p0.ParseArgv(argv)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "Expected a required 'string' argument")
}

func (s *MySuite) TestParseOptionalPositionalArgumentPresent(c *C) {
	values := &TestParseValues{}

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p.AddArgument(&Argument{
		Name: "string",
		Help: "Required string value",
	})
	p.AddArgument(&Argument{
		Name:    "integer",
		NumArgs: '?',
		Help:    "Optional integer value",
	})

	// No string argument passed after subcommand
	argv := []string{"string_value", "123"}
	err := p.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.String, Equals, "string_value")
	c.Check(values.Integer, Equals, 123)
}

func (s *MySuite) TestParseOptionalPositionalArgumentAbsent(c *C) {
	values := &TestParseValues{}

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p.AddArgument(&Argument{
		Name: "string",
		Help: "Required string value",
	})
	p.AddArgument(&Argument{
		Name:    "integer",
		NumArgs: '?',
		Help:    "Optional integer value",
	})

	// No string argument passed after subcommand
	argv := []string{"string_value"}
	err := p.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.String, Equals, "string_value")
}

func (s *MySuite) TestParseOneString(c *C) {
	values := &TestParseValues{}

	p0 := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p0.AddArgument(&Argument{
		Name: "strings",
		Help: "Required string value",
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
		String:       "--",
		Dest:         "PassThrough",
	})
	p0.AddArgument(&Argument{
		Name: "string",
		Help: "Required string value",
	})

	// No string argument passed after subcommand
	argv := []string{"foo", "--", "a", "b", "c"}
	err := p0.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.String, Equals, "foo")
	c.Assert(len(values.PassThrough), Equals, 3)
	c.Assert(values.PassThrough, DeepEquals, []string{"a", "b", "c"})
}

func (s *MySuite) TestParsePassThroughAfterPositionalMultiValue(c *C) {
	values := &TestParseValues{}

	p0 := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p0.AddArgument(&Argument{
		ParseCommand: PassThrough,
		String:       "--",
		Dest:         "PassThrough",
	})
	p0.AddArgument(&Argument{
		Name:    "strings",
		NumArgs: '+',
		Help:    "Required string value",
	})

	// No string argument passed after subcommand
	argv := []string{"foo", "x", "y", "--", "a", "b", "c"}
	err := p0.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.Strings, DeepEquals, []string{"foo", "x", "y"})
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
		String:       "--",
		Dest:         "PassThrough",
	})
	p0.AddArgument(&Argument{
		Switches: []string{"--string"},
		Help: "Required string value",
	})

	// No string argument passed after subcommand
	argv := []string{"--string", "xxx", "--", "a", "b", "c"}
	err := p0.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.String, Equals, "xxx")
	c.Assert(len(values.PassThrough), Equals, 3)
	c.Assert(values.PassThrough, DeepEquals, []string{"a", "b", "c"})
}

func (s *MySuite) TestParseShortWithNumber(c *C) {
	values := &TestParseValues{}

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p.AddArgument(&Argument{
		Switches: []string{"-j"},
	})

	// Pass a number adjoined with "-j"
	argv := []string{"-j4"}
	err := p.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.J, Equals, 4)
}

func (s *MySuite) TestParseShortGroupedBooleans(c *C) {
	values := &TestParseValues{}

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p.AddArgument(&Argument{
		Switches: []string{"-x"},
	})
	p.AddArgument(&Argument{
		Switches: []string{"-y"},
	})

	// Pass a number adjoined with "-j"
	argv := []string{"-yx"}
	err := p.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.X, Equals, true)
	c.Check(values.Y, Equals, true)
}

func (s *MySuite) TestParseShortGroupedErrors(c *C) {
	values := &TestParseValues{}

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p.AddArgument(&Argument{
		Switches: []string{"-j"},
	})
	p.AddArgument(&Argument{
		Switches: []string{"-x"},
	})
	p.AddArgument(&Argument{
		Switches: []string{"-y"},
	})

	// Illegal
	argv := []string{"-jx"}
	err := p.ParseArgv(argv)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "While parsing value for -j: Cannot convert \"x\" to an integer")

	// Illegal
	argv = []string{"-xj"}
	err = p.ParseArgv(argv)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "The -x switch is valid but not as '-xj'")

	// Illegal
	argv = []string{"-xz"}
	err = p.ParseArgv(argv)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "The -x switch is valid but not as '-xz'")
}

func (s *MySuite) TestParseChoicesString(c *C) {
	v := &TestParseValues{}

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      v,
	}
	p.AddArgument(&Argument{
		Switches:    []string{"--string"},
		Choices: []string{"x", "y", "z"},
	})

	err := p.ParseArgv([]string{"--string", "w"})
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "The possible values for --string are 'x', 'y', and 'z'")

	err = p.ParseArgv([]string{"--string", "x"})
	c.Check(err, IsNil)

}

func (s *MySuite) TestParseChoicesInteger(c *C) {
	v := &TestParseValues{}

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      v,
	}
	p.AddArgument(&Argument{
		Switches:    []string{"--integer"},
		Choices: []string{"1", "22", "333"},
	})

	err := p.ParseArgv([]string{"--integer", "222"})
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "The possible values for --integer are 1, 22, and 333")

	err = p.ParseArgv([]string{"--integer", "22"})
	c.Check(err, IsNil)
}

func (s *MySuite) TestParseEqualsValue(c *C) {
	v := &TestParseValues{}

	p := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      v,
	}
	p.AddArgument(&Argument{
		Switches: []string{"--integer"},
	})
	p.AddArgument(&Argument{
		Switches: []string{"--boolean"},
		Dest: "X",
	})
	p.AddArgument(&Argument{
		Switches: []string{"--string"},
	})

	err := p.ParseArgv([]string{"--integer=222"})
	c.Assert(err, IsNil)
	c.Check(v.Integer, Equals, 222)

	err = p.ParseArgv([]string{"--string=xyz"})
	c.Assert(err, IsNil)
	c.Check(v.String, Equals, "xyz")

	// Error condition
	err = p.ParseArgv([]string{"--="})
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "No such switch: --")

	// Error condition
	err = p.ParseArgv([]string{"--boolean=true"})
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "The --boolean switch does not take a value")
}

func (s *MySuite) TestParseMultipleSwitches(c *C) {
	values := &TestParseValues{}

	p0 := &ArgumentParser{
		Name:             "progname",
		ShortDescription: "This is a simple program",
		Destination:      values,
	}
	p0.AddArgument(&Argument{
		Switches:    []string{"-a", "--b", "-cd", "-w", "--x", "-yz"},
		Dest: "Strings",
	})

	// No string argument passed after subcommand
	argv := []string{
		"-yz", "YZ",
		"--x", "X",
		"-w", "W",
		"-cd", "CD",
		"--b", "B",
		"-a", "A",
	}
	err := p0.ParseArgv(argv)
	c.Assert(err, IsNil)
	c.Check(values.Strings, DeepEquals, []string{"YZ", "X", "W", "CD", "B", "A"})
}
