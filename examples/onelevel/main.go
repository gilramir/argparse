// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package main

import (
	"fmt"

	"github.com/gilramir/argparse"
)

type MyOptions struct {
	Debug bool
	Verbose	bool
	Names []string
}

func main() {
	opts := &MyOptions{}
	ap := argparse.New(&argparse.Command{
		Description:	"This is an example program",
		Values:		opts,
	})
/*
	ap.Add(&argparse.Argument{
		Switches:	[]string{"--debug"},
		Help:		"Set debug mode",
	})

	ap.Add(&argparse.Argument{
		Switches:	[]string{"-v", "--verbose"},
		Help:		"Set verbose mode",
	})
	*/
/*
	ap.Add(&argparse.String{
		Name:		"names",
		Help:		"Some names passed into the program",
		NumArgsGlob:	'+',
	})
*/
	// The library handles errors, and -h/--help
	ap.Parse()
	/*
	opts_ := ap.Parse()
	opts := opts_.(*MyOptions)
	*/

	fmt.Printf("Verbose is %v\n", opts.Verbose)
	fmt.Printf("Debug is %v\n", opts.Debug)
	fmt.Printf("Number of names: %d\n", len(opts.Names))

	for i := 0 ; i < len(opts.Names) ; i++ {
		fmt.Printf("%d. %s\n", i+1, opts.Names[i])
	}
}
