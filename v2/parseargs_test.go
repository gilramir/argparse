package argparse

// Copyright (c) 2026 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"bytes"
	"strings"

	. "gopkg.in/check.v1"
)

func (s *MySuite) TestParseArgsSuccess(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Name: "prog", Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--int1"}})

	cmd, err := ap.ParseArgs([]string{"--int1", "7"})
	c.Assert(err, IsNil)
	c.Check(cmd, Equals, ap.Root)
	c.Check(opts.Int1, Equals, 7)
}

// A bad command-line returns an error instead of exiting the process.
func (s *MySuite) TestParseArgsError(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Name: "prog", Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--int1"}})

	_, err := ap.ParseArgs([]string{"--nope"})
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "No such switch: --nope")
}

// A help request returns ErrHelp and writes the help text to Stdout.
func (s *MySuite) TestParseArgsHelp(c *C) {
	var buf bytes.Buffer
	opts := &PTestOptions{}
	ap := New(&Command{Name: "prog", Description: "this is prog", Values: opts})
	ap.Stdout = &buf
	ap.Add(&Argument{Switches: []string{"--int1"}})

	cmd, err := ap.ParseArgs([]string{"--help"})
	c.Assert(err, Equals, ErrHelp)
	c.Check(cmd, Equals, ap.Root)
	c.Check(strings.Contains(buf.String(), "this is prog"), Equals, true)
}

// ParseArgs returns the triggered sub-command without running its Function.
func (s *MySuite) TestParseArgsSubcommand(c *C) {
	ran := false
	ap := New(&Command{Name: "prog", Values: &reqRootOptions{}})
	sub := ap.New(&Command{
		Name:   "open",
		Values: &reqSubOptions{},
		Function: func(cmd *Command, values Values) error {
			ran = true
			return nil
		},
	})
	sub.Add(&Argument{Switches: []string{"--reason"}})

	cmd, err := ap.ParseArgs([]string{"open", "--reason", "y"})
	c.Assert(err, IsNil)
	c.Check(cmd.Name, Equals, "open")
	c.Check(cmd.Values.(*reqSubOptions).Reason, Equals, "y")
	// ParseArgs must not run the callback.
	c.Check(ran, Equals, false)
}
