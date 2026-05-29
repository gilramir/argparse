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

// The help output includes a one-line usage synopsis. Required switches appear
// without brackets; optional ones and the help switch are bracketed.
func (s *MySuite) TestHelpUsageSynopsis(c *C) {
	o := &PTestOptions{}
	ap := New(&Command{Name: "demo", Description: "d", Values: o})
	ap.Add(&Argument{Switches: []string{"--int1", "-c"}, Dest: "Int1", MetaVar: "N"})
	ap.Add(&Argument{Switches: []string{"--string1"}, Required: true})
	ap.Add(&Argument{Name: "PosStringSlice", NumArgsGlob: "+"})

	help := ap.helpString(ap.Root, nil)
	c.Check(strings.Contains(help,
		"usage: demo [--int1 N] --string1 STRING1 [-h] PosStringSlice [...]"),
		Equals, true)
}

// The synopsis shows a <command> placeholder at a command with sub-commands,
// and the full command path for a sub-command.
func (s *MySuite) TestHelpUsageSynopsisSubcommand(c *C) {
	ap := New(&Command{Name: "prog", Values: &reqRootOptions{}})
	ap.Add(&Argument{Switches: []string{"-v", "--verbose"}})
	sub := ap.New(&Command{Name: "open", Values: &reqSubOptions{}})
	sub.Add(&Argument{Switches: []string{"--reason"}})

	rootHelp := ap.helpString(ap.Root, nil)
	c.Check(strings.Contains(rootHelp, "usage: prog [-v] [-h] <command>"), Equals, true)

	subHelp := ap.helpString(sub, []*Command{ap.Root})
	c.Check(strings.Contains(subHelp, "usage: prog open [--reason REASON] [-h]"), Equals, true)
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
