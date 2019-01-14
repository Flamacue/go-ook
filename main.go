package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/flamacue/go-ook/compiler"
	"github.com/flamacue/go-ook/machine"
)

func main() {
	// TODO: Add help message instead of exit 1
	if len(os.Args) < 2 {
		fmt.Println("Please specify a file to execute.")
		os.Exit(1)
	}

	code, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		handleError(err)
	}
	compiler := compiler.New(string(code))
	instructions, err := compiler.Compile()
	if err != nil {
		handleError(err)
	}

	m := machine.New(instructions, os.Stdin, os.Stdout)
	m.Execute()
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}
