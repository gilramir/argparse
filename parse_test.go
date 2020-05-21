// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	. "gopkg.in/check.v1"
)

type PTestOptions struct {
	Bool1	bool
	Bool2	bool
	Int1	int
	Int2	int
	String1	string
	String2	string
	Float1	float64
	Float2	float64

	PosBool1	bool
	PosBool2	bool

	PosInt		int
	PosString	string
	PosFloat	float64

	PosBoolSlice	[]bool
	PosIntSlice	[]int
	PosStringSlice	[]string
	PosFloatSlice	[]float64
}

func createPTestParser() (*PTestOptions,*ArgumentParser) {
	opts := &PTestOptions{}
	ap := New(&Command{
		Description:	"This is a test program",
		Values:		opts,
	})
	ap.Add(&Argument{
		Switches:	[]string{"--bool1"},
	})
	ap.Add(&Argument{
		Switches:	[]string{"--bool2"},
	})
	ap.Add(&Argument{
		Switches:	[]string{"--int1"},
	})
	ap.Add(&Argument{
		Switches:	[]string{"--int2"},
	})
	ap.Add(&Argument{
		Switches:	[]string{"--string1"},
	})
	ap.Add(&Argument{
		Switches:	[]string{"--string2"},
	})
	ap.Add(&Argument{
		Switches:	[]string{"--float1"},
	})
	ap.Add(&Argument{
		Switches:	[]string{"--float2"},
	})
	return opts, ap
}

// ====================================================== bool

func (s *MySuite) TestRootSwitchesBool1(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--bool1"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.Bool1, Equals, true )

	c.Check( ap.Root.Seen["Bool1"], Equals, true )
	c.Check( ap.Root.Seen["Bool2"], Equals, false )
}

func (s *MySuite) TestRootSwitchesBool2(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--bool2", "--bool1"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.Bool1, Equals, true )
	c.Check( opts.Bool2, Equals, true )

	c.Check( ap.Root.Seen["Bool1"], Equals, true )
	c.Check( ap.Root.Seen["Bool2"], Equals, true )
}

// ====================================================== string

func (s *MySuite) TestRootSwitchesString1(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--string1", "abc"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.String1, Equals, "abc" )

	c.Check( ap.Root.Seen["String1"], Equals, true )
	c.Check( ap.Root.Seen["String2"], Equals, false )
}

func (s *MySuite) TestRootSwitchesString2(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--string2", "xyz", "--string1", "mno"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.String1, Equals, "mno" )
	c.Check( opts.String2, Equals, "xyz" )

	c.Check( ap.Root.Seen["String1"], Equals, true )
	c.Check( ap.Root.Seen["String2"], Equals, true )
}

func (s *MySuite) TestRootSwitchesString3(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--string1=abc"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.String1, Equals, "abc" )

	c.Check( ap.Root.Seen["String1"], Equals, true )
	c.Check( ap.Root.Seen["String2"], Equals, false )
}

// ====================================================== int

func (s *MySuite) TestRootSwitchesInt1(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--int1", "500"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.Int1, Equals, 500 )

	c.Check( ap.Root.Seen["Int1"], Equals, true )
	c.Check( ap.Root.Seen["Int2"], Equals, false )
}

func (s *MySuite) TestRootSwitchesInt2(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--int2", "-3", "--int1", "77"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.Int1, Equals, 77 )
	c.Check( opts.Int2, Equals, -3 )

	c.Check( ap.Root.Seen["Int1"], Equals, true )
	c.Check( ap.Root.Seen["Int2"], Equals, true )
}

func (s *MySuite) TestRootSwitchesInt3(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--int1=500"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.Int1, Equals, 500 )

	c.Check( ap.Root.Seen["Int1"], Equals, true )
	c.Check( ap.Root.Seen["Int2"], Equals, false )
}

// ====================================================== float

func (s *MySuite) TestRootSwitchesFloat1(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--float1", "500.2"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.Float1, Equals, 500.2 )

	c.Check( ap.Root.Seen["Float1"], Equals, true )
	c.Check( ap.Root.Seen["Float2"], Equals, false )
}

func (s *MySuite) TestRootSwitchesFloat2(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--float2", "-30.0", "--float1", "0.02"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.Float1, Equals, 0.02 )
	c.Check( opts.Float2, Equals, -30.0 )

	c.Check( ap.Root.Seen["Float1"], Equals, true )
	c.Check( ap.Root.Seen["Float2"], Equals, true )
}

func (s *MySuite) TestRootSwitchesFloat3(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--float1=500"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.Float1, Equals, 500.0 )

	c.Check( ap.Root.Seen["Float1"], Equals, true )
	c.Check( ap.Root.Seen["Float2"], Equals, false )
}

// ====================================================== bool positional

func (s *MySuite) TestRootPositionalBool1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosBool1",
	})

	ap.Add(&Argument{
		Name:	"PosBool2",
	})

	argv := []string{"false", "true"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosBool1, Equals, false )
	c.Check( opts.PosBool2, Equals, true )

	c.Check( ap.Root.Seen["PosBool1"], Equals, true )
	c.Check( ap.Root.Seen["PosBool2"], Equals, true )
}

// ====================================================== int positional

func (s *MySuite) TestRootPositionalInt1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosInt",
	})

	argv := []string{"333"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosInt, Equals, 333 )

	c.Check( ap.Root.Seen["PosInt"], Equals, true )
}

// ====================================================== float positional

func (s *MySuite) TestRootPositionalFloat1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosFloat",
	})

	argv := []string{"400.04"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosFloat, Equals, 400.04 )

	c.Check( ap.Root.Seen["PosFloat"], Equals, true )
}


// ====================================================== string positional

func (s *MySuite) TestRootPositionalString1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosString",
	})

	argv := []string{"foo"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosString, Equals, "foo" )

	c.Check( ap.Root.Seen["PosString"], Equals, true )
}

// ====================================================== bool slice positional

func (s *MySuite) TestRootPositionalBoolSlice1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosBoolSlice",
		NumArgs: 1,
	})

	argv := []string{"true"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosBoolSlice, DeepEquals, []bool{true} )

	c.Check( ap.Root.Seen["PosBoolSlice"], Equals, true )
}

func (s *MySuite) TestRootPositionalBoolSliceQuestion0(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosBoolSlice",
		NumArgsGlob: "?",
	})

	argv := []string{}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( len(opts.PosBoolSlice), Equals, 0 )

	c.Check( ap.Root.Seen["PosBoolSlice"], Equals, false )
}

func (s *MySuite) TestRootPositionalBoolSliceQuestion1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosBoolSlice",
		NumArgsGlob: "?",
	})

	argv := []string{"true"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosBoolSlice, DeepEquals, []bool{true} )

	c.Check( ap.Root.Seen["PosBoolSlice"], Equals, true )
}

func (s *MySuite) TestRootPositionalBoolSliceQuestion2(c *C) {
	_, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosBoolSlice",
		NumArgsGlob: "?",
	})

	argv := []string{"true", "false"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, NotNil)
}

// ====================================================== string slice positional

func (s *MySuite) TestRootPositionalStringSlice1(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:	"PosStringSlice",
	})

	argv := []string{"foo"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosStringSlice, DeepEquals, []string{"foo"} )

	c.Check( ap.Root.Seen["PosStringSlice"], Equals, true )
}

func (s *MySuite) TestRootPositionalStringSliceStar0(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:	"PosStringSlice",
		NumArgsGlob: "*",
	})

	argv := []string{}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( len(opts.PosStringSlice), Equals, 0 )

	c.Check( ap.Root.Seen["PosStringSlice"], Equals, false )
}

func (s *MySuite) TestRootPositionalStringSliceStar1(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:	"PosStringSlice",
		NumArgsGlob: "*",
	})

	argv := []string{"z"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosStringSlice, DeepEquals, []string{"z"} )

	c.Check( ap.Root.Seen["PosStringSlice"], Equals, true )
}

func (s *MySuite) TestRootPositionalStringSliceStar2(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:	"PosStringSlice",
		NumArgsGlob: "*",
	})

	argv := []string{"a", "b"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosStringSlice, DeepEquals, []string{"a", "b"} )

	c.Check( ap.Root.Seen["PosStringSlice"], Equals, true )
}

// ====================================================== int slice positional

func (s *MySuite) TestRootPositionalIntSlice1(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:	"PosIntSlice",
	})

	argv := []string{"101"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosIntSlice, DeepEquals, []int{101} )

	c.Check( ap.Root.Seen["PosIntSlice"], Equals, true )
}

func (s *MySuite) TestRootPositionalIntSlicePlus0(c *C) {
	_, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosIntSlice",
		NumArgsGlob: "+",
	})

	argv := []string{}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, NotNil)
}

func (s *MySuite) TestRootPositionalIntSlicePlus1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosIntSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"101"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosIntSlice, DeepEquals, []int{101} )

	c.Check( ap.Root.Seen["PosIntSlice"], Equals, true )
}

func (s *MySuite) TestRootPositionalIntSlicePlus2(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosIntSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"101", "202"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosIntSlice, DeepEquals, []int{101, 202} )

	c.Check( ap.Root.Seen["PosIntSlice"], Equals, true )
}

// ====================================================== float slice positional

func (s *MySuite) TestRootPositionalFloatSlice1(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:	"PosFloatSlice",
	})

	argv := []string{"101.2"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosFloatSlice, DeepEquals, []float64{101.2} )

	c.Check( ap.Root.Seen["PosFloatSlice"], Equals, true )
}

func (s *MySuite) TestRootPositionalFloatSlice2(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:	"PosFloatSlice",
		NumArgs: 2,
	})

	argv := []string{"101.2", "202.4"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosFloatSlice, DeepEquals, []float64{101.2, 202.4} )

	c.Check( ap.Root.Seen["PosFloatSlice"], Equals, true )
}

// ====================================================== NumArgsGlob +
func (s *MySuite) TestNumArgsGlobPlusZero(c *C) {
	_, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:		"PosStringSlice",
		NumArgsGlob:	"+",
	})

	argv := []string{}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, NotNil)
}
func (s *MySuite) TestNumArgsGlobPlusOne(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:		"PosStringSlice",
		NumArgsGlob:	"+",
	})

	argv := []string{"a" }
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosStringSlice, DeepEquals, []string{"a"} )

	c.Check( ap.Root.Seen["PosStringSlice"], Equals, true )
}

func (s *MySuite) TestNumArgsGlobPlusTwo(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:		"PosStringSlice",
		NumArgsGlob:	"+",
	})

	argv := []string{"a", "b" }
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosStringSlice, DeepEquals, []string{"a", "b"} )

	c.Check( ap.Root.Seen["PosStringSlice"], Equals, true )
}

func (s *MySuite) TestNumArgsGlobPlusThree(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name:		"PosStringSlice",
		NumArgsGlob:	"+",
	})

	argv := []string{"a", "b", "c" }
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.PosStringSlice, DeepEquals, []string{"a", "b", "c"} )

	c.Check( ap.Root.Seen["PosStringSlice"], Equals, true )
}

/*
import (
	"bytes"

	. "gopkg.in/check.v1"
)
*/
/*
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
*/
