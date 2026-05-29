package argparse

// Copyright (c) 2026 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"bytes"
	"strings"

	. "gopkg.in/check.v1"
)

// When Version is set, --version prints it and ParseArgs returns ErrVersion.
func (s *MySuite) TestVersionSwitch(c *C) {
	var buf bytes.Buffer
	opts := &PTestOptions{}
	ap := New(&Command{Name: "prog", Values: opts})
	ap.Version = "prog 1.2.3"
	ap.Stdout = &buf

	cmd, err := ap.ParseArgs([]string{"--version"})
	c.Assert(err, Equals, ErrVersion)
	c.Check(cmd, Equals, ap.Root)
	c.Check(strings.TrimSpace(buf.String()), Equals, "prog 1.2.3")
}

// Without a Version set, --version is just an unknown switch.
func (s *MySuite) TestVersionSwitchNotConfigured(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Name: "prog", Values: opts})

	_, err := ap.ParseArgs([]string{"--version"})
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "No such switch: --version")
}

// The version switches can be customized.
func (s *MySuite) TestVersionSwitchCustom(c *C) {
	var buf bytes.Buffer
	opts := &PTestOptions{}
	ap := New(&Command{Name: "prog", Values: opts})
	ap.Version = "9"
	ap.VersionSwitches = []string{"-V", "--version"}
	ap.Stdout = &buf

	_, err := ap.ParseArgs([]string{"-V"})
	c.Assert(err, Equals, ErrVersion)
	c.Check(strings.TrimSpace(buf.String()), Equals, "9")
}

// The version switch appears in the help output and synopsis when configured.
func (s *MySuite) TestVersionShownInHelp(c *C) {
	opts := &PTestOptions{}
	ap := New(&Command{Name: "prog", Description: "d", Values: opts})
	ap.Version = "1.0"

	help := ap.helpString(ap.Root, nil)
	c.Check(strings.Contains(help, "--version"), Equals, true)
	c.Check(strings.Contains(help, "Show the version and exit"), Equals, true)
	c.Check(strings.Contains(help, "[--version]"), Equals, true)
}
