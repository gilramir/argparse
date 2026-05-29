package argparse

// Copyright (c) 2026 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	. "gopkg.in/check.v1"
)

type reqRootOptions struct {
	Verbose bool
}

type reqSubOptions struct {
	reqRootOptions
	Reason string
}

func (s *MySuite) TestRequiredSwitchGiven(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--int1"}, Required: true})

	results := ap.parseArgv([]string{"--int1", "5"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.Int1, Equals, 5)
}

func (s *MySuite) TestRequiredSwitchMissing(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--int1"}, Required: true})

	results := ap.parseArgv([]string{})
	c.Assert(results.parseError, NotNil)
	c.Check(results.parseError.Error(), Equals, "Missing required switch: --int1")
}

// A required boolean switch must still be present.
func (s *MySuite) TestRequiredBoolSwitchMissing(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--bool1"}, Required: true})

	results := ap.parseArgv([]string{})
	c.Assert(results.parseError, NotNil)
	c.Check(results.parseError.Error(), Equals, "Missing required switch: --bool1")
}

// Required is reported using the full switch set in PrettyName.
func (s *MySuite) TestRequiredSwitchPrettyName(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"-x", "--expiration"}, Dest: "String1", Required: true})

	results := ap.parseArgv([]string{})
	c.Assert(results.parseError, NotNil)
	c.Check(results.parseError.Error(), Equals, "Missing required switch: -x/--expiration")
}

// A required switch on a sub-command must be given when that sub-command runs.
func (s *MySuite) TestRequiredSwitchOnSubcommand(c *C) {
	ap := New(&Command{Description: "x", Values: &reqRootOptions{}})
	sub := ap.New(&Command{
		Name:   "open",
		Values: &reqSubOptions{},
	})
	sub.Add(&Argument{Switches: []string{"-r", "--reason"}, Required: true})

	// Missing on the sub-command -> error
	results := ap.parseArgv([]string{"open"})
	c.Assert(results.parseError, NotNil)
	c.Check(results.parseError.Error(), Equals, "Missing required switch: -r/--reason")

	// Given on the sub-command -> ok
	results = ap.parseArgv([]string{"open", "--reason", "because"})
	c.Assert(results.parseError, IsNil)
	c.Check(results.triggeredCommand.Values.(*reqSubOptions).Reason, Equals, "because")
}

// Setting Required on a positional argument is a programming error.
func (s *MySuite) TestRequiredPositionalPanics(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	c.Assert(func() {
		ap.Add(&Argument{Name: "PosString", Required: true})
	}, PanicMatches, "Cannot set Required on positional argument.*")
}
