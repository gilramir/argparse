package argparse

// Copyright (c) 2026 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"strings"

	. "gopkg.in/check.v1"
)

// The help output shows a non-zero default value and the list of choices.
func (s *MySuite) TestHelpShowsDefaultAndChoices(c *C) {
	opts := &PTestOptions{}
	opts.Int1 = 5 // a non-zero default

	ap := New(&Command{Name: "prog", Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--int1"}, Help: "count"})
	ap.Add(&Argument{
		Switches: []string{"--string1"},
		Help:     "mode",
		Choices:  []string{"a", "b", "c"},
	})

	help := ap.helpString(ap.Root, nil)
	c.Check(strings.Contains(help, "(default: 5)"), Equals, true)
	c.Check(strings.Contains(help, "(choices: a, b, c)"), Equals, true)
}

// Zero-valued defaults are not shown (no noise like "(default: 0)").
func (s *MySuite) TestHelpHidesZeroDefault(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Name: "prog", Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--int2"}, Help: "count"})
	ap.Add(&Argument{Switches: []string{"--bool1"}, Help: "flag"})

	help := ap.helpString(ap.Root, nil)
	c.Check(strings.Contains(help, "(default:"), Equals, false)
}
