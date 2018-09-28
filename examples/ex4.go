package main

import (
	"fmt"
	"os"

	"github.com/gilramir/argparse"
)

type Options struct {
	Pattern     string
	Filename    []string
}

func (self *Options) Run([]argparse.Destination) error {
        fmt.Printf("Pattern: %s\n", self.Pattern)
        fmt.Printf("Filenames: %v\n", self.Filename)
	return nil
}

func main() {
	p := &argparse.ArgumentParser{
		Name:             "my_program",
		ShortDescription: "This program takes positional arguments",
		Destination:      &Options{},
	}

	p.AddArgument(&argparse.Argument{
		Switches: []string{"-p", "--pattern"},
		Help: "The pattern to look for",
	})

	p.AddArgument(&argparse.Argument{
		Switches: []string{"-f", "--filename"},
		Help: "The file to look at. Can be given more than once.",
	})

	err := p.ParseArgs()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
