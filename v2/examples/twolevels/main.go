// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package main

import (
	"fmt"

	"github.com/gilramir/argparse/v2"
)

type RootOptions struct {
	Debug   bool
	Verbose bool
}

type OpenOptions struct {
	RootOptions

	Name   string
	Reason string
}

type CloseOptions struct {
	RootOptions

	Name   string
	Reason string
}

func DoOpen(cmd *argparse.Command, values argparse.Values) error {
	opts := values.(*OpenOptions)

	fmt.Printf("Open: Verbose is %v, Seen=%v\n", opts.Verbose,
		cmd.Seen["Verbose"])
	fmt.Printf("Open: Debug is %v, Seen=%v\n", opts.Debug,
		cmd.Seen["Debug"])
	fmt.Printf("Open: Reason is %s, Seen=%v\n", opts.Reason,
		cmd.Seen["Reason"])
	fmt.Printf("Open: Name is %v\n", opts.Name)

	return nil
}

func DoClose(cmd *argparse.Command, values argparse.Values) error {
	opts := values.(*CloseOptions)

	fmt.Printf("Close: Verbose is %v, Seen=%v\n", opts.Verbose,
		cmd.Seen["Verbose"])
	fmt.Printf("Close: Debug is %v, Seen=%v\n", opts.Debug,
		cmd.Seen["Debug"])
	fmt.Printf("Close: Reason is %s, Seen=%v\n", opts.Reason,
		cmd.Seen["Reason"])
	fmt.Printf("Close: Name is %v\n", opts.Name)

	return nil
}

func main() {
	opts := &RootOptions{}
	ap := argparse.New(&argparse.Command{
		Description: "This is an example program",
		Values:      opts,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"--debug"},
		Help:     "Set debug mode",
		Inherit:  true,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"-v", "--verbose"},
		Help:     "Set verbose mode",
		Inherit:  true,
	})

	// open
	open_ap := ap.New(&argparse.Command{
		Name:        "open",
		Description: "Open something",
		Function:    DoOpen,
		Values:      &OpenOptions{},
	})

	open_ap.Add(&argparse.Argument{
		Switches: []string{"-r", "--reason"},
		Help:     "Why you are opening this",
	})

	open_ap.Add(&argparse.Argument{
		Name: "name",
		Help: "The thing you are opening",
	})

	// close
	close_ap := ap.New(&argparse.Command{
		Name:        "close",
		Description: "Close something",
		Function:    DoClose,
		Values:      &CloseOptions{},
	})

	close_ap.Add(&argparse.Argument{
		Switches: []string{"-r", "--reason"},
		Help:     "Why you are closing this",
	})

	close_ap.Add(&argparse.Argument{
		Name: "name",
		Help: "The thing you are closing",
	})

	ap.Parse()

	fmt.Println("After Parse")
	fmt.Printf("root: Verbose is %v\n", opts.Verbose)
	fmt.Printf("root: Debug is %v\n", opts.Debug)
}
