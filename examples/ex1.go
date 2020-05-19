package main

import (
	"fmt"
	"os"

	"github.com/gilramir/argparse"
)

func main() {
	p := &argparse.RootArgumentParser{
	})

	p := &argparse.ArgumentParser{
		Name:             "my_program",
		ShortDescription: "This program takes no arguments",
	}
	err := p.ParseArgs()
	fmt.Printf("Parsing is finished\n")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}
