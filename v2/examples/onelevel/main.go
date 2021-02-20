// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package main

import (
	"fmt"
	"time"

	"github.com/gilramir/argparse/v2"
)

type MyOptions struct {
	Debug    bool
	Duration time.Duration
	Verbose  bool
	Names    []string
	N        int
}

func main() {
	opts := &MyOptions{}
	ap := argparse.New(&argparse.Command{
		Description: "This is an example program",
		Values:      opts,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--debug"},
		Help:     "Set debug mode",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--duration"},
		MetaVar:  "#(h|m|s|ms|us|ns)",
		Help:     "How long do you want to run?",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"-v", "--verbose"},
		Help:     "Set verbose mode",
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"-n"},
		Help:     "Number",
		Choices:  []int{1, 2, 5},
	})

	ap.Add(&argparse.Argument{
		Name: "names",
		Help: `Some names passed into the program. This is
                an example of a very long help message so that the word wrap
                can be checked.`,
		NumArgsGlob: "+",
	})

	// The library handles errors, and -h/--help
	ap.Parse()

	fmt.Printf("Verbose is %v\n", opts.Verbose)
	fmt.Printf("Debug is %v\n", opts.Debug)
	fmt.Printf("N is %v\n", opts.N)

	if ap.Root.Seen["Duration"] {
		fmt.Printf("Duration: %s\n", opts.Duration.String())
	}

	fmt.Printf("Number of names: %d\n", len(opts.Names))
	for i := 0; i < len(opts.Names); i++ {
		fmt.Printf("%d. %s\n", i+1, opts.Names[i])
	}
}
