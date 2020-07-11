// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package main

import (
	"fmt"

	"github.com/gilramir/argparse/v2"
)

type RootOptions struct {
	Debug   bool
	Verbose bool
	Reason  string
}

type OpenOptions struct {
	RootOptions

	Name string
}

type CloseOptions struct {
	RootOptions

	Name string
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

func build_argparse() *argparse.ArgumentParser {
	opts := &RootOptions{
		Reason: "default-reason-at-root",
	}
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

	build_argparse_open(ap)
	build_argparse_close(ap)
	return ap
}

func build_argparse_open(parent_ap *argparse.ArgumentParser) {
	opts := &OpenOptions{
		RootOptions{
			Reason: "default-reason-at-open",
		},
		/*Name:*/ "",
	}
	// open
	ap := parent_ap.New(&argparse.Command{
		Name:        "open",
		Description: "Open something",
		Function:    DoOpen,
		Values:      opts,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"-r", "--reason"},
		Help:     "Why you are opening this",
	})

	ap.Add(&argparse.Argument{
		Name: "name",
		Help: "The thing you are opening",
	})
}

func build_argparse_close(parent_ap *argparse.ArgumentParser) {

	opts := &CloseOptions{
		RootOptions{
			Reason: "default-reason-at-close",
		},
		/*Name:*/ "",
	}

	// close
	ap := parent_ap.New(&argparse.Command{
		Name:        "close",
		Description: "Close something",
		Function:    DoClose,
		Values:      opts,
	})

	ap.Add(&argparse.Argument{
		Switches: []string{"-r", "--reason"},
		Help:     "Why you are closing this",
	})

	ap.Add(&argparse.Argument{
		Name: "name",
		Help: "The thing you are closing",
	})
}

func main() {

	ap := build_argparse()
	ap.Parse()
}
