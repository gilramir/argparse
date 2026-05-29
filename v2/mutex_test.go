package argparse

// Copyright (c) 2026 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	. "gopkg.in/check.v1"
)

// At most one of a non-required group may be given.
func (s *MySuite) TestMutexAtMostOne(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.AddMutuallyExclusive(false,
		&Argument{Switches: []string{"--bool1"}},
		&Argument{Switches: []string{"--bool2"}})

	// Neither given: ok (not required).
	results := ap.parseArgv([]string{})
	c.Assert(results.parseError, IsNil)

	// One given: ok.
	results = ap.parseArgv([]string{"--bool1"})
	c.Assert(results.parseError, IsNil)

	// Both given: error.
	results = ap.parseArgv([]string{"--bool1", "--bool2"})
	c.Assert(results.parseError, NotNil)
	c.Check(results.parseError.Error(), Equals,
		"Only one of --bool1, --bool2 may be given")
}

// A required group needs exactly one.
func (s *MySuite) TestMutexRequiredOne(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.AddMutuallyExclusive(true,
		&Argument{Switches: []string{"--int1"}},
		&Argument{Switches: []string{"--int2"}})

	// None given: error.
	results := ap.parseArgv([]string{})
	c.Assert(results.parseError, NotNil)
	c.Check(results.parseError.Error(), Equals,
		"One of --int1, --int2 is required")

	// Exactly one: ok.
	results = ap.parseArgv([]string{"--int1", "5"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.Int1, Equals, 5)

	// Both: error.
	results = ap.parseArgv([]string{"--int1", "5", "--int2", "6"})
	c.Assert(results.parseError, NotNil)
}

// A positional argument cannot be placed in a group.
func (s *MySuite) TestMutexPositionalPanics(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	c.Assert(func() {
		ap.AddMutuallyExclusive(false,
			&Argument{Switches: []string{"--bool1"}},
			&Argument{Name: "PosString"})
	}, PanicMatches, "Cannot put positional argument.*")
}

// An individual group member must not set Required.
func (s *MySuite) TestMutexRequiredMemberPanics(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	c.Assert(func() {
		ap.AddMutuallyExclusive(false,
			&Argument{Switches: []string{"--bool1"}, Required: true},
			&Argument{Switches: []string{"--bool2"}})
	}, PanicMatches, ".*must not set Required.*")
}

// A group needs at least two members.
func (s *MySuite) TestMutexTooFewPanics(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	c.Assert(func() {
		ap.AddMutuallyExclusive(false, &Argument{Switches: []string{"--bool1"}})
	}, PanicMatches, ".*at least two.*")
}
