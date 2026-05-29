package argparse

// Copyright (c) 2026 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	. "gopkg.in/check.v1"
)

// upperString is a custom type whose UnmarshalText upper-cases the input.
type upperString string

func (u *upperString) UnmarshalText(text []byte) error {
	*u = upperString(strings.ToUpper(string(text)))
	return nil
}

// evenInt is a custom type whose UnmarshalText rejects odd numbers.
type evenInt int

func (e *evenInt) UnmarshalText(text []byte) error {
	n, err := strconv.Atoi(string(text))
	if err != nil {
		return err
	}
	if n%2 != 0 {
		return fmt.Errorf("%d is not even", n)
	}
	*e = evenInt(n)
	return nil
}

type customOptions struct {
	Name  upperString
	Even  evenInt
	Addr  net.IP
	Names []upperString
}

func (s *MySuite) TestCustomScalarType(c *C) {
	opts := &customOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--name"}})

	results := ap.parseArgv([]string{"--name", "hello"})
	c.Assert(results.parseError, IsNil)
	c.Check(string(opts.Name), Equals, "HELLO")
}

func (s *MySuite) TestCustomScalarTypeError(c *C) {
	opts := &customOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--even"}})

	// Odd value -> UnmarshalText returns an error, wrapped by CannotParseValueFmt.
	results := ap.parseArgv([]string{"--even", "3"})
	c.Assert(results.parseError, NotNil)
	c.Check(strings.Contains(results.parseError.Error(), "is not even"), Equals, true)

	// Even value -> ok
	results = ap.parseArgv([]string{"--even", "4"})
	c.Assert(results.parseError, IsNil)
	c.Check(int(opts.Even), Equals, 4)
}

// net.IP from the standard library implements encoding.TextUnmarshaler.
func (s *MySuite) TestCustomStdlibType(c *C) {
	opts := &customOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{Switches: []string{"--addr"}})

	results := ap.parseArgv([]string{"--addr", "10.0.0.7"})
	c.Assert(results.parseError, IsNil)
	c.Check(opts.Addr.String(), Equals, "10.0.0.7")
}

func (s *MySuite) TestCustomSliceType(c *C) {
	opts := &customOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	ap.Add(&Argument{
		Name:        "Names",
		NumArgsGlob: "+",
	})

	results := ap.parseArgv([]string{"foo", "bar"})
	c.Assert(results.parseError, IsNil)
	c.Assert(len(opts.Names), Equals, 2)
	c.Check(string(opts.Names[0]), Equals, "FOO")
	c.Check(string(opts.Names[1]), Equals, "BAR")
}

// Choices is not supported for custom types; setting it panics at Add().
func (s *MySuite) TestCustomTypeChoicesPanics(c *C) {
	opts := &customOptions{}
	ap := New(&Command{Description: "x", Values: opts})
	c.Assert(func() {
		ap.Add(&Argument{
			Switches: []string{"--name"},
			Choices:  []string{"A", "B"},
		})
	}, PanicMatches, ".*not supported for a custom.*")
}
