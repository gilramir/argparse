package argparse

// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"time"

	. "gopkg.in/check.v1"
)

type PTestOptions struct {
	Bool1       bool
	Bool2       bool
	Int1        int
	Int2        int
	Int64       int64
	String1     string
	String2     string
	Float1      float64
	Float2      float64
	StringSlice []string

	PosBool1 bool
	PosBool2 bool

	PosInt    int
	PosInt1   int
	PosInt2   int
	PosString string
	PosFloat  float64

	PosBoolSlice     []bool
	PosIntSlice      []int
	PosStringSlice   []string
	PosFloatSlice    []float64
	PosDurationSlice []time.Duration
}

func createPTestParser() (*PTestOptions, *ArgumentParser) {
	opts := &PTestOptions{}
	ap := New(&Command{
		Description: "This is a test program",
		Values:      opts,
	})
	ap.Add(&Argument{
		Switches: []string{"--bool1"},
	})
	ap.Add(&Argument{
		Switches: []string{"--bool2"},
	})
	ap.Add(&Argument{
		Switches: []string{"--int1"},
	})
	ap.Add(&Argument{
		Switches: []string{"--int2"},
	})
	ap.Add(&Argument{
		Switches: []string{"--int64"},
	})
	ap.Add(&Argument{
		Switches: []string{"--string1"},
	})
	ap.Add(&Argument{
		Switches: []string{"--string2"},
	})
	ap.Add(&Argument{
		Switches: []string{"--float1"},
	})
	ap.Add(&Argument{
		Switches: []string{"--float2"},
	})
	return opts, ap
}

// ====================================================== bool

func (s *MySuite) TestRootSwitchesBool1(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--bool1"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Bool1, Equals, true)

	c.Check(ap.Root.Seen["Bool1"], Equals, true)
	c.Check(ap.Root.Seen["Bool2"], Equals, false)
}

func (s *MySuite) TestRootSwitchesBool2(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--bool2", "--bool1"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Bool1, Equals, true)
	c.Check(opts.Bool2, Equals, true)

	c.Check(ap.Root.Seen["Bool1"], Equals, true)
	c.Check(ap.Root.Seen["Bool2"], Equals, true)
}

// ====================================================== string

func (s *MySuite) TestRootSwitchesString1(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--string1", "abc"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.String1, Equals, "abc")

	c.Check(ap.Root.Seen["String1"], Equals, true)
	c.Check(ap.Root.Seen["String2"], Equals, false)
}

func (s *MySuite) TestRootSwitchesString2(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--string2", "xyz", "--string1", "mno"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.String1, Equals, "mno")
	c.Check(opts.String2, Equals, "xyz")

	c.Check(ap.Root.Seen["String1"], Equals, true)
	c.Check(ap.Root.Seen["String2"], Equals, true)
}

func (s *MySuite) TestRootSwitchesString3(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--string1=abc"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.String1, Equals, "abc")

	c.Check(ap.Root.Seen["String1"], Equals, true)
	c.Check(ap.Root.Seen["String2"], Equals, false)
}

// ====================================================== int

func (s *MySuite) TestRootSwitchesInt1(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--int1", "500"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Int1, Equals, 500)

	c.Check(ap.Root.Seen["Int1"], Equals, true)
	c.Check(ap.Root.Seen["Int2"], Equals, false)
}

func (s *MySuite) TestRootSwitchesInt2(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--int2", "-3", "--int1", "77"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Int1, Equals, 77)
	c.Check(opts.Int2, Equals, -3)

	c.Check(ap.Root.Seen["Int1"], Equals, true)
	c.Check(ap.Root.Seen["Int2"], Equals, true)
}

func (s *MySuite) TestRootSwitchesInt3(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--int1=500"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Int1, Equals, 500)

	c.Check(ap.Root.Seen["Int1"], Equals, true)
	c.Check(ap.Root.Seen["Int2"], Equals, false)
}

// ====================================================== int64

func (s *MySuite) TestRootSwitchesInt64Pos(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--int64", "500"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Int64, Equals, int64(500))

	c.Check(ap.Root.Seen["Int64"], Equals, true)
	c.Check(ap.Root.Seen["Int1"], Equals, false)
	c.Check(ap.Root.Seen["Int2"], Equals, false)
}

func (s *MySuite) TestRootSwitchesInt64Neg(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--int64", "-999"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Int64, Equals, int64(-999))

	c.Check(ap.Root.Seen["Int64"], Equals, true)
	c.Check(ap.Root.Seen["Int1"], Equals, false)
	c.Check(ap.Root.Seen["Int2"], Equals, false)
}

// ====================================================== float

func (s *MySuite) TestRootSwitchesFloat1(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--float1", "500.2"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Float1, Equals, 500.2)

	c.Check(ap.Root.Seen["Float1"], Equals, true)
	c.Check(ap.Root.Seen["Float2"], Equals, false)
}

func (s *MySuite) TestRootSwitchesFloat2(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--float2", "-30.0", "--float1", "0.02"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Float1, Equals, 0.02)
	c.Check(opts.Float2, Equals, -30.0)

	c.Check(ap.Root.Seen["Float1"], Equals, true)
	c.Check(ap.Root.Seen["Float2"], Equals, true)
}

func (s *MySuite) TestRootSwitchesFloat3(c *C) {
	opts, ap := createPTestParser()

	argv := []string{"--float1=500"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Float1, Equals, 500.0)

	c.Check(ap.Root.Seen["Float1"], Equals, true)
	c.Check(ap.Root.Seen["Float2"], Equals, false)
}

// ====================================================== bool positional

func (s *MySuite) TestRootPositionalBool1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name: "PosBool1",
	})

	ap.Add(&Argument{
		Name: "PosBool2",
	})

	argv := []string{"false", "true"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosBool1, Equals, false)
	c.Check(opts.PosBool2, Equals, true)

	c.Check(ap.Root.Seen["PosBool1"], Equals, true)
	c.Check(ap.Root.Seen["PosBool2"], Equals, true)
}

// ====================================================== int positional

func (s *MySuite) TestRootPositionalInt1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name: "PosInt",
	})

	argv := []string{"333"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 333)

	c.Check(ap.Root.Seen["PosInt"], Equals, true)
}

// ====================================================== float positional

func (s *MySuite) TestRootPositionalFloat1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name: "PosFloat",
	})

	argv := []string{"400.04"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosFloat, Equals, 400.04)

	c.Check(ap.Root.Seen["PosFloat"], Equals, true)
}

// ====================================================== string positional

func (s *MySuite) TestRootPositionalString1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name: "PosString",
	})

	argv := []string{"foo"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosString, Equals, "foo")

	c.Check(ap.Root.Seen["PosString"], Equals, true)
}

// ====================================================== bool slice positional

func (s *MySuite) TestRootPositionalBoolSlice1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:    "PosBoolSlice",
		NumArgs: 1,
	})

	argv := []string{"true"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosBoolSlice, DeepEquals, []bool{true})

	c.Check(ap.Root.Seen["PosBoolSlice"], Equals, true)
}

func (s *MySuite) TestRootPositionalBoolSliceQuestion0(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosBoolSlice",
		NumArgsGlob: "?",
	})

	argv := []string{}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(len(opts.PosBoolSlice), Equals, 0)

	c.Check(ap.Root.Seen["PosBoolSlice"], Equals, false)
}

func (s *MySuite) TestRootPositionalBoolSliceQuestion1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosBoolSlice",
		NumArgsGlob: "?",
	})

	argv := []string{"true"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosBoolSlice, DeepEquals, []bool{true})

	c.Check(ap.Root.Seen["PosBoolSlice"], Equals, true)
}

func (s *MySuite) TestRootPositionalBoolSliceQuestion2(c *C) {
	_, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosBoolSlice",
		NumArgsGlob: "?",
	})

	argv := []string{"true", "false"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, NotNil)
}

// ====================================================== string slice positional

func (s *MySuite) TestRootPositionalStringSlice1(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name: "PosStringSlice",
	})

	argv := []string{"foo"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"foo"})

	c.Check(ap.Root.Seen["PosStringSlice"], Equals, true)
}

func (s *MySuite) TestRootPositionalStringSliceStar0(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "*",
	})

	argv := []string{}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(len(opts.PosStringSlice), Equals, 0)

	c.Check(ap.Root.Seen["PosStringSlice"], Equals, false)
}

func (s *MySuite) TestRootPositionalStringSliceStar1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "*",
	})

	argv := []string{"z"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"z"})

	c.Check(ap.Root.Seen["PosStringSlice"], Equals, true)
}

func (s *MySuite) TestRootPositionalStringSliceStar2(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "*",
	})

	argv := []string{"a", "b"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"a", "b"})

	c.Check(ap.Root.Seen["PosStringSlice"], Equals, true)
}

// ====================================================== int slice positional

func (s *MySuite) TestRootPositionalIntSlice1(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name: "PosIntSlice",
	})

	argv := []string{"101"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosIntSlice, DeepEquals, []int{101})

	c.Check(ap.Root.Seen["PosIntSlice"], Equals, true)
}

func (s *MySuite) TestRootPositionalIntSlicePlus0(c *C) {
	_, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosIntSlice",
		NumArgsGlob: "+",
	})

	argv := []string{}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, NotNil)
}

func (s *MySuite) TestRootPositionalIntSlicePlus1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosIntSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"101"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosIntSlice, DeepEquals, []int{101})

	c.Check(ap.Root.Seen["PosIntSlice"], Equals, true)
}

func (s *MySuite) TestRootPositionalIntSlicePlus2(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosIntSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"101", "202"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosIntSlice, DeepEquals, []int{101, 202})

	c.Check(ap.Root.Seen["PosIntSlice"], Equals, true)
}

// ====================================================== float slice positional

func (s *MySuite) TestRootPositionalFloatSlice1(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name: "PosFloatSlice",
	})

	argv := []string{"101.2"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosFloatSlice, DeepEquals, []float64{101.2})

	c.Check(ap.Root.Seen["PosFloatSlice"], Equals, true)
}

func (s *MySuite) TestRootPositionalFloatSlice2(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:    "PosFloatSlice",
		NumArgs: 2,
	})

	argv := []string{"101.2", "202.4"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosFloatSlice, DeepEquals, []float64{101.2, 202.4})

	c.Check(ap.Root.Seen["PosFloatSlice"], Equals, true)
}

// ====================================================== time.Duration slice positional

func (s *MySuite) TestRootPositionalTimeDurationSlice1(c *C) {
	opts, ap := createPTestParser()

	// No NumArgs or NumArgsGlob is legal; == 1
	ap.Add(&Argument{
		Name: "PosDurationSlice",
	})

	argv := []string{"1m"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Assert(len(opts.PosDurationSlice), Equals, 1)
	c.Check(opts.PosDurationSlice[0].Seconds(), Equals, 60.0)

	c.Check(ap.Root.Seen["PosDurationSlice"], Equals, true)
}

func (s *MySuite) TestRootPositionalTimeDurationSlice2(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:    "PosDurationSlice",
		NumArgs: 2,
	})

	argv := []string{"2m", "1h"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Assert(len(opts.PosDurationSlice), Equals, 2)
	c.Check(opts.PosDurationSlice[0].Seconds(), Equals, 120.0)
	c.Check(opts.PosDurationSlice[1].Seconds(), Equals, 3600.0)

	c.Check(ap.Root.Seen["PosDurationSlice"], Equals, true)
}

// ====================================================== NumArgsGlob +
func (s *MySuite) TestNumArgsGlobPlusZero(c *C) {
	_, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "+",
	})

	argv := []string{}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, NotNil)
}
func (s *MySuite) TestNumArgsGlobPlusOne(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"a"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"a"})

	c.Check(ap.Root.Seen["PosStringSlice"], Equals, true)
}

func (s *MySuite) TestNumArgsGlobPlusTwo(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"a", "b"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"a", "b"})

	c.Check(ap.Root.Seen["PosStringSlice"], Equals, true)
}

func (s *MySuite) TestNumArgsGlobPlusThree(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"a", "b", "c"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"a", "b", "c"})

	c.Check(ap.Root.Seen["PosStringSlice"], Equals, true)
}

// ====================================================== NumArgsGlob: *

func (s *MySuite) TestRootPositionalStar0(c *C) {
	_, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "*",
	})

	argv := []string{"-a", "-b", "-c"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, NotNil)
}

func (s *MySuite) TestRootPositionalStar1(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "*",
	})

	argv := []string{"--", "-a", "-b", "-c"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"-a", "-b", "-c"})

	c.Check(ap.Root.Seen["PosStringSlice"], Equals, true)
}

// ============================================= NumArgsGlob: ? of string

func (s *MySuite) TestRootPositionalOptionalString0(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosString",
		NumArgsGlob: "?",
	})

	argv := []string{"--bool1"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(ap.Root.Seen["PosString"], Equals, false)
	c.Check(opts.PosString, Equals, "")

	argv = []string{"--bool1", "target"}
	results = ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(ap.Root.Seen["PosString"], Equals, true)
	c.Check(opts.PosString, Equals, "target")
}

// ======================================================

// test switch argument after fixed number of positional arguments

func (s *MySuite) TestSwitchAfterPositional(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:    "PosString",
		NumArgs: 1,
	})

	argv := []string{"x", "--bool1"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.Bool1, Equals, true)
	c.Check(opts.PosString, Equals, "x")
}

// test switch argument after unbounded number of positional arguments
func (s *MySuite) TestSwitchAfterUnboundedPositional(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"x", "--bool1"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(len(opts.PosStringSlice), Equals, 2)
	c.Check(opts.PosStringSlice[0], Equals, "x")
	// "--bool1", as a string, ends up as a positional arg value
	c.Check(opts.PosStringSlice[1], Equals, "--bool1")
}

// ======================================================

// Test positional arguments after a ? option
func (s *MySuite) TestOptionalPositional01(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	argv := []string{"22"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 22)
}

// Test positional arguments after a ? option, with no input
func (s *MySuite) TestOptionalPositional01a(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	argv := []string{}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 0)
}

// ? followed by required
func (s *MySuite) TestOptionalPositional02(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name: "PosInt1",
	})

	argv := []string{"22", "33"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 22)
	c.Check(opts.PosInt1, Equals, 33)
}

// ? followed by required, with only 1 input
func (s *MySuite) TestOptionalPositional02a(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name: "PosInt1",
	})

	argv := []string{"22"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 22)
	c.Check(opts.PosInt1, Equals, 0)
}

// ? followed by required, required
func (s *MySuite) TestOptionalPositional03(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name: "PosInt1",
	})

	ap.Add(&Argument{
		Name: "PosInt2",
	})

	argv := []string{"22", "33", "44"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 22)
	c.Check(opts.PosInt1, Equals, 33)
	c.Check(opts.PosInt2, Equals, 44)
}

// ?? followed by required
func (s *MySuite) TestOptionalPositional04(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name:        "PosInt1",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name: "PosInt2",
	})

	argv := []string{"22", "33", "44"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 22)
	c.Check(opts.PosInt1, Equals, 33)
	c.Check(opts.PosInt2, Equals, 44)
}

// 3 ? arguments
func (s *MySuite) TestOptionalPositional05(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name:        "PosInt1",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name:        "PosInt2",
		NumArgsGlob: "?",
	})

	argv := []string{"22", "33", "44"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 22)
	c.Check(opts.PosInt1, Equals, 33)
	c.Check(opts.PosInt2, Equals, 44)
}

// 2 ? arguments with leftover input
func (s *MySuite) TestOptionalPositional06(c *C) {
	_, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name:        "PosInt1",
		NumArgsGlob: "?",
	})

	argv := []string{"22", "33", "44"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, NotNil)
}

// ? argument followed by *
func (s *MySuite) TestOptionalPositional07(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name:        "PosIntSlice",
		NumArgsGlob: "*",
	})

	argv := []string{"22", "33", "44"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 22)
	c.Check(opts.PosIntSlice, DeepEquals, []int{33, 44})
}

// ? argument followed by +
func (s *MySuite) TestOptionalPositional08(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Name:        "PosInt",
		NumArgsGlob: "?",
	})

	ap.Add(&Argument{
		Name:        "PosIntSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"22", "33", "44"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosInt, Equals, 22)
	c.Check(opts.PosIntSlice, DeepEquals, []int{33, 44})
}

// Test 2 required arguments after a switch argument
func (s *MySuite) TestRequired01(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Switches: []string{"--strings"},
		Dest:     "StringSlice",
		NumArgs:  2,
	})

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"--strings", "22", "33", "44", "55"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(len(opts.StringSlice), Equals, 2)
	c.Check(opts.StringSlice[0], Equals, "22")
	c.Check(opts.StringSlice[1], Equals, "33")
	c.Check(len(opts.PosStringSlice), Equals, 2)
	c.Check(opts.PosStringSlice[0], Equals, "44")
	c.Check(opts.PosStringSlice[1], Equals, "55")
}

// Test 2 required arguments after a switch argument with =
func (s *MySuite) TestRequired02(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Switches: []string{"--strings"},
		Dest:     "StringSlice",
		NumArgs:  2,
	})

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"--strings=22", "33", "44", "55"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(len(opts.StringSlice), Equals, 2)
	c.Check(opts.StringSlice[0], Equals, "22")
	c.Check(opts.StringSlice[1], Equals, "33")
	c.Check(len(opts.PosStringSlice), Equals, 2)
	c.Check(opts.PosStringSlice[0], Equals, "44")
	c.Check(opts.PosStringSlice[1], Equals, "55")
}

// Test 2 required arguments after a switch argument, giving the switch twice
func (s *MySuite) TestRequired03(c *C) {
	opts, ap := createPTestParser()

	ap.Add(&Argument{
		Switches: []string{"--strings"},
		Dest:     "StringSlice",
		NumArgs:  2,
	})

	ap.Add(&Argument{
		Name:        "PosStringSlice",
		NumArgsGlob: "+",
	})

	argv := []string{"--strings", "22", "33", "--strings", "44", "55", "x"}
	results := ap.parseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(len(opts.StringSlice), Equals, 4)
	c.Check(opts.StringSlice[0], Equals, "22")
	c.Check(opts.StringSlice[1], Equals, "33")
	c.Check(opts.StringSlice[2], Equals, "44")
	c.Check(opts.StringSlice[3], Equals, "55")
	c.Check(len(opts.PosStringSlice), Equals, 1)
	c.Check(opts.PosStringSlice[0], Equals, "x")
}
