package argparse

// Copyright (c) 2026 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	. "gopkg.in/check.v1"
)

type numOptions struct {
	I8  int8
	I16 int16
	I32 int32
	U   uint
	U8  uint8
	U64 uint64
	F32 float32

	I8s  []int8
	Us   []uint
	F32s []float32
}

func (s *MySuite) TestScalarSignedSubInts(c *C) {
	opts := &numOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--i8"}})
	ap.Add(&Argument{Switches: []string{"--i16"}})
	ap.Add(&Argument{Switches: []string{"--i32"}})

	results := ap.parseArgv([]string{"--i8", "-12", "--i16", "0x100", "--i32", "70000"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.I8, Equals, int8(-12))
	c.Check(opts.I16, Equals, int16(256))
	c.Check(opts.I32, Equals, int32(70000))
}

func (s *MySuite) TestScalarUnsignedInts(c *C) {
	opts := &numOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--u"}})
	ap.Add(&Argument{Switches: []string{"--u8"}})
	ap.Add(&Argument{Switches: []string{"--u64"}})

	results := ap.parseArgv([]string{"--u", "42", "--u8", "0xff", "--u64", "18446744073709551615"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.U, Equals, uint(42))
	c.Check(opts.U8, Equals, uint8(255))
	c.Check(opts.U64, Equals, uint64(18446744073709551615))
}

func (s *MySuite) TestUnsignedRejectsNegative(c *C) {
	opts := &numOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--u"}})

	results := ap.parseArgv([]string{"--u", "-1"})
	c.Assert(results.parseError, NotNil)
}

func (s *MySuite) TestScalarFloat32(c *C) {
	opts := &numOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--f32"}})

	results := ap.parseArgv([]string{"--f32", "1.5"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.F32, Equals, float32(1.5))
}

func (s *MySuite) TestSliceNumericTypes(c *C) {
	opts := &numOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--i8s"}, NumArgs: 2})
	ap.Add(&Argument{Switches: []string{"--us"}, NumArgs: 2})
	ap.Add(&Argument{Switches: []string{"--f32s"}, NumArgs: 2})

	results := ap.parseArgv([]string{
		"--i8s", "1", "2",
		"--us", "10", "20",
		"--f32s", "1.5", "2.5",
	})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.I8s, DeepEquals, []int8{1, 2})
	c.Check(opts.Us, DeepEquals, []uint{10, 20})
	c.Check(opts.F32s, DeepEquals, []float32{1.5, 2.5})
}

func (s *MySuite) TestUnsignedChoices(c *C) {
	opts := &numOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--u"}, Choices: []uint{1, 2, 3}})

	results := ap.parseArgv([]string{"--u", "2"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.U, Equals, uint(2))

	results = ap.parseArgv([]string{"--u", "5"})
	c.Assert(results.parseError, NotNil)
}
