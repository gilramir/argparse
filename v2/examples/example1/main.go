// Copyright (c) 2021 by Gilbert Ramirez <gram@alumni.rice.edu>

package main

import (
	"fmt"
	"time"

	"github.com/gilramir/argparse/v2"
)

type MyOptions struct {
	Count      int
	Expiration time.Duration
	Verbose    bool
	Names      []string
}

func main() {
	opts := &MyOptions{}
	ap := argparse.New(&argparse.Command{
		Name:        "Example 1",
		Description: "This is an example program",
		Values:      opts,
                Epilog: `This is the first line.
This is the second line. We will then skip a line.

And now 1 line was skipped. The end.`,
	})

	// These are switch arguments
	ap.Add(&argparse.Argument{
		Switches: []string{"--count"},
		MetaVar:  "N",
		Help:     "How many items",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--expiration", "-x"},
		Help:     "How long: #(h|m|s|ms|us|ns)",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"-v", "--verbose"},
		Help:     "Set verbose mode",
	})

	// This is a positional argument
	ap.Add(&argparse.Argument{
		Name: "names",
		Help: "Some names passed into the program",
		// We require one or more names
		NumArgsGlob: "+",
	})

	// The library handles errors, and -h/--help
	ap.Parse()

	fmt.Printf("Verbose is %t\n", opts.Verbose)
	fmt.Printf("Count is %d\n", opts.Count)

	if ap.Root.Seen["Expiration"] {
		fmt.Printf("Expiration: %s\n", opts.Expiration.String())
	}

	fmt.Printf("Number of names: %d\n", len(opts.Names))
	for i := 0; i < len(opts.Names); i++ {
		fmt.Printf("%d. %s\n", i+1, opts.Names[i])
	}
}
