package main

import (
	"fmt"
	"os"

	"github.com/gilramir/argparse"
)

type Options struct {
	Filename        []string
        OtherArguments  []string
}

func (self *Options) Run([]argparse.Destination) error {
        fmt.Printf("Filenames: %v\n", self.Filename)
        fmt.Printf("Other Arguments: %v\n", self.OtherArguments)
	return nil
}

func main() {
	p := &argparse.ArgumentParser{
		Name:             "my_program",
		ShortDescription: "This program takes positional arguments",
		Destination:      &Options{},
	}

	p.AddArgument(&argparse.Argument{
		Name: "filename",
		Help: "The file to look at. Can be given more than once.",
	})

	p.AddArgument(&argparse.Argument{
		ParseCommand: argparse.PassThrough,
		String:        "--",
		Dest:         "OtherArguments",
	})

	err := p.ParseArgs()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

