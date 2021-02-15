// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	. "gopkg.in/check.v1"
)

type APTestOptions struct {
	Bool1   bool
	String1 string
}

func (s *MySuite) TestChooseHelpLongAtRoot(c *C) {

	opts := &APTestOptions{}
	ap := New(&Command{
		Description: "This is a test program",
		Values:      opts,
	})

	argv := []string{"--help"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(results.helpRequested, Equals, true)
	c.Check(results.triggeredCommand, Equals, ap.Root)
}

func (s *MySuite) TestChooseHelpShortAtRoot(c *C) {

	opts := &APTestOptions{}
	ap := New(&Command{
		Description: "This is a test program",
		Values:      opts,
	})

	argv := []string{"-h"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(results.helpRequested, Equals, true)
	c.Check(results.triggeredCommand, Equals, ap.Root)
}

func (s *MySuite) TestChooseHelpCustomAtRoot(c *C) {

	opts := &APTestOptions{}
	ap := New(&Command{
		Description: "This is a test program",
		Values:      opts,
	})
	ap.HelpSwitches = []string{"--custom-help"}

	argv := []string{"--custom-help"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check(results.helpRequested, Equals, true)
	c.Check(results.triggeredCommand, Equals, ap.Root)
}
