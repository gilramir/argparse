// Copyright (c) 2020 by Gilbert Ramirez <gram@alumni.rice.edu>

package main

import (
	"fmt"
	"os"

	"github.com/gilramir/argparse"
)

type RootOptions struct {
	Debug bool
	Verbose	bool
}

type OpenOptions struct {
	RootOptions

	Name	string
	Reason	string
}

type CloseOptions struct {
	RootOptions

	Name	string
	Reason	string
}

func main() {
	ap := argparse.NewArgumentParser( &argparse.Program{
		ShortDescription:	"This is an example program",
		Values:			&RootOptions{},
	})

	ap.AddArgument(&argparse.Bool{
		Switches:	[]string{"--debug"},
		Help:		"Set debug mode",
	})

	ap.AddArgument(&argparse.Bool{
		Switches:	[]string{"-v", "--verbose"},
		Help:		"Set verbose mode",
	})

	// open
	open_ap := ap.AddSubCommand( &argparse.SubCommand {
		Name:			"open",
		ShortDescription:	"Open something",
		Function:		do_open,
	})

	open_ap.AddArgument(&argparse.String{
		Switches:	[]string{"-r", "--reason"},
		Help:		"Why you are opening this",
	})

	open_ap.AddArgument(&argparse.String{
		Name:		"name",
		Help:		"The thing you are opening",
	})

	// close
/*
	close_ap := ap.AddParser( &argparse.ArgumentParser{
		Name:			"close",
		ShortDescription:	"Close something",
		Function:		do_close,
	})

	close_ap.AddArgument(&argparse.Argument{
		Switches:	[]string{"-r", "--reason"},
		Help:		"Why you are closing this",
	})

	close_ap.AddArgument(&argparse.Argument{
		Name:		"name",
		Help:		"The thing you are closing",
	})
*/

	err := ap.ExecuteOsArgs()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func do_open( opts_ OptionValues ) error {

	opts := opts_.(*OpenOptions)

	fmt.Printf("(open) Verbose is %v\n", opts.Verbose)
	fmt.Printf("(open) Debug is %v\n", opts.Debug)
	fmt.Printf("(open) Reason is %s\n", opts.Reason)
	fmt.Printf("(open) name is %s\n", opts.Name)
}
