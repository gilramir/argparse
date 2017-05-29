package main

import (
	"fmt"
	"os"

	"github.com/gilramir/argparse"
)

type Options struct {
	Filenames []string
}

func (self *Options) Run([]argparse.Destination) error {
	return argparse.ParseError("The CLI syntax is bad")
}

func main() {
	p := &argparse.ArgumentParser{
		Name:             "my_program",
		ShortDescription: "This program takes positional arguments",
		Destination:      &Options{},
	}

	p.AddArgument(&argparse.Argument{
		Name:    "filenames",
		Help:    "The file(s) to look at",
		NumArgs: '+',
	})

	err := p.ParseArgs()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	fmt.Printf("Done.\n")
}
