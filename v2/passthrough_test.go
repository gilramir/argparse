package argparse

// Copyright (c) 2026 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	. "gopkg.in/check.v1"
)

type ptRootOptions struct{}

type ptSubOptions struct {
	Rest []string
}

// "--" sends the whole remainder to the pass-through argument, even when there
// are no other positional arguments, and even tokens beginning with "-".
func (s *MySuite) TestPassThroughDashDash(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Name: "PosStringSlice", PassThrough: true})

	results := ap.parseArgv([]string{"--", "a", "-b", "--c"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"a", "-b", "--c"})
	c.Check(ap.Root.Seen["PosStringSlice"], Equals, true)
}

// A non-dash first token also engages the pass-through, capturing the rest raw.
func (s *MySuite) TestPassThroughNoDashDash(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Name: "PosStringSlice", PassThrough: true})

	results := ap.parseArgv([]string{"sub", "-x", "--y"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"sub", "-x", "--y"})
}

// The program's own switches parse normally until the pass-through engages.
func (s *MySuite) TestPassThroughAfterKnownSwitch(c *C) {
	opts, ap := createPTestParser()
	ap.Add(&Argument{Name: "PosStringSlice", PassThrough: true})

	results := ap.parseArgv([]string{"--bool1", "cmd", "--bool2", "-z"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.Bool1, Equals, true)
	// --bool2 is part of the captured remainder, not parsed.
	c.Check(opts.Bool2, Equals, false)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"cmd", "--bool2", "-z"})
}

// A preceding positional is filled first; the pass-through captures the rest.
func (s *MySuite) TestPassThroughAfterPositional(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Name: "PosString", NumArgs: 1})
	ap.Add(&Argument{Name: "PosStringSlice", PassThrough: true})

	results := ap.parseArgv([]string{"first", "-x", "--y"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosString, Equals, "first")
	c.Check(opts.PosStringSlice, DeepEquals, []string{"-x", "--y"})
}

// A "--" appearing inside the captured remainder is taken literally.
func (s *MySuite) TestPassThroughLiteralDashDash(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Name: "PosStringSlice", PassThrough: true})

	results := ap.parseArgv([]string{"--", "a", "--", "b"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"a", "--", "b"})
}

// A token matching one of the program's own switches is captured raw, not parsed.
func (s *MySuite) TestPassThroughCapturesOwnSwitch(c *C) {
	opts, ap := createPTestParser()
	ap.Add(&Argument{Name: "PosStringSlice", PassThrough: true})

	results := ap.parseArgv([]string{"--", "--bool1"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.Bool1, Equals, false)
	c.Check(opts.PosStringSlice, DeepEquals, []string{"--bool1"})
}

// Empty pass-through is fine (zero or more), and not marked as Seen.
func (s *MySuite) TestPassThroughEmpty(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Name: "PosStringSlice", PassThrough: true})

	results := ap.parseArgv([]string{"--"})
	c.Assert(results.parseError, IsNil)
	c.Check(len(opts.PosStringSlice), Equals, 0)
	c.Check(ap.Root.Seen["PosStringSlice"], Equals, false)

	results = ap.parseArgv([]string{})
	c.Assert(results.parseError, IsNil)
	c.Check(len(opts.PosStringSlice), Equals, 0)
}

// Pass-through works on a sub-command.
func (s *MySuite) TestPassThroughSubcommand(c *C) {
	ap := New(&Command{Name: "prog", Values: &ptRootOptions{}})
	sub := ap.New(&Command{Name: "run", Values: &ptSubOptions{}})
	sub.Add(&Argument{Name: "Rest", PassThrough: true})

	results := ap.parseArgv([]string{"run", "--", "cmd", "-flag"})
	c.Assert(results.parseError, IsNil)
	c.Check(results.triggeredCommand.Values.(*ptSubOptions).Rest,
		DeepEquals, []string{"cmd", "-flag"})
}

// PassThrough must be a positional []string, must be last, and there is only one.
func (s *MySuite) TestPassThroughValidation(c *C) {
	// Not a positional.
	c.Assert(func() {
		opts := &PTestOptions{}
		ap := New(&Command{Description: "x", Values: opts})
		ap.Add(&Argument{Switches: []string{"--bool1"}, PassThrough: true})
	}, PanicMatches, "PassThrough can only be set on a positional.*")

	// Not a []string.
	c.Assert(func() {
		opts := &PTestOptions{}
		ap := New(&Command{Description: "x", Values: opts})
		ap.Add(&Argument{Name: "PosIntSlice", PassThrough: true})
	}, PanicMatches, "PassThrough argument.*must have a \\[\\]string.*")

	// Nothing may follow it (it behaves like "*").
	c.Assert(func() {
		opts := &PTestOptions{}
		ap := New(&Command{Description: "x", Values: opts})
		ap.Add(&Argument{Name: "PosStringSlice", PassThrough: true})
		ap.Add(&Argument{Name: "PosString"})
	}, PanicMatches, "Cannot add a positional argument after.*")
}
