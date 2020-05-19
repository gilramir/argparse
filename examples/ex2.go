package main

import (
	"fmt"
	"os"

	"github.com/gilramir/argparse"
)

type Options struct {
	Pattern   string
	Filenames []string
}

func (self *Options) Run([]argparse.Destination) error {
        fmt.Printf("Pattern: %s\n", self.Pattern)
        fmt.Printf("Filenames: %v\n", self.Filenames)
	return nil
}

func main() {
	p := &argparse.ArgumentParser{
		Name:             "my_program",
		ShortDescription: "This program takes positional arguments",
		Destination:      &Options{},
	}

	p.AddArgument(&argparse.Argument{
		Name: "pattern",
		Help: "The pattern to look for",
	})

	p.AddArgument(&argparse.Argument{
		Name:    "filenames",
		Help:    "The file(s) to look at",
		NumArgs: '+',
	})

	err := p.ParseArgs()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
