// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>
package argparse

import (
	. "gopkg.in/check.v1"
)

type CTestCommon struct {
	Verbose	bool
	Debug	bool
}

type CTestOptionsRoot struct {
	CTestCommon
	Int1	int
	String1	string
}

type CTestOptionsSubA struct {
	CTestOptionsRoot
	String2 string
}

type CTestOptionsSubB struct {
	CTestOptionsRoot
	String2 string
}

// These are put togther in a single struct only to make
// the createCTestParser function friendlier to use.
type CTestOptions struct {
	root CTestOptionsRoot
	a CTestOptionsSubA
	b CTestOptionsSubB
}

func createCTestParser() (*CTestOptions, *ArgumentParser, *Command, *Command) {
	opts := &CTestOptions{}
	ap := New(&Command{
		Description:	"This is a test program",
		Values:		&opts.root,
	})
	ap.Add(&Argument{
		Switches:	[]string{"--verbose", "-v"},
		Inherit:	true,
	})
	ap.Add(&Argument{
		Switches:	[]string{"--debug", "-d"},
		//Inherit:	true,
	})
	ap.Add(&Argument{
		Switches:	[]string{"--int1"},
	})
	ap.Add(&Argument{
		Switches:	[]string{"--string1"},
	})

	subA := ap.New(&Command{
		Name:	"sub-a",
		Description: "Sub-command A",
		Values:		&opts.a,
	})
	/*
	subA.Add(&Argument{
		Switches:	[]string{"--verbose", "-v"},
	})
	subA.Add(&Argument{
		Switches:	[]string{"--debug", "-d"},
	})
	*/
	subA.Add(&Argument{
		Switches:	[]string{"--string2"},
	})

	subB := ap.New(&Command{
		Name:	"sub-b",
		Description: "Sub-command B",
		Values:		&opts.b,
	})
	/*
	subB.Add(&Argument{
		Switches:	[]string{"--verbose", "-v"},
	})
	subB.Add(&Argument{
		Switches:	[]string{"--debug", "-d"},
	})
	*/
	subB.Add(&Argument{
		Switches:	[]string{"--string2"},
	})

	return opts, ap, subA, subB
}

// ====================================================== bool

func (s *MySuite) TestSubCommandRoot(c *C) {
	opts, ap, _, _ := createCTestParser()

	argv := []string{"--verbose"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)
	c.Check( opts.root.Verbose, Equals, true )

	c.Check( ap.Root.Seen["Verbose"], Equals, true )
}

func (s *MySuite) TestSubCommandSubAPre(c *C) {
	opts, ap, suba, _ := createCTestParser()

	// --verbose is before the sub-command "sub-a"
	argv := []string{"--verbose", "sub-a", "--string2", "foo"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)

	c.Check( ap.Root.CommandSeen["sub-a"], Equals, true )

	c.Check( opts.a.Verbose, Equals, true )
	c.Check( opts.a.String2, Equals, "foo" )

	c.Check( suba.Seen["Verbose"], Equals, true )
	c.Check( suba.Seen["String2"], Equals, true )
}

func (s *MySuite) TestSubCommandSubAPost(c *C) {
	opts, ap, suba, _ := createCTestParser()

	// --verbose is after the sub-command "sub-a"
	argv := []string{"sub-a", "--verbose", "--string2", "foo"}
	results := ap.ParseArgv(argv)

	c.Assert(results.parseError, IsNil)

	c.Check( ap.Root.CommandSeen["sub-a"], Equals, true )

	c.Check( opts.a.Verbose, Equals, true )
	c.Check( opts.a.String2, Equals, "foo" )

	c.Check( suba.Seen["Verbose"], Equals, true )
	c.Check( suba.Seen["String2"], Equals, true )
}
