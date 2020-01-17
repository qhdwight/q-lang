package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"q-lang-go/parser"
	"q-lang-go/parser/gen"
)

func main() {
	inputFiles := flag.String("input", "", "Input Q files")
	flag.Parse()
	if *inputFiles == "" {
		panic("No input Q files provided!")
	}
	program := parser.Parse(*inputFiles)

	assembly := new(gen.Program)
	program.Generate(assembly)
	fmt.Println("Output assembly:\n", assembly.ToString())

	assemble(assembly)
}

func assemble(assembly *gen.Program) {
	// TODO currently we write to file first and then call GCC. Find a way to pass it without file as we don't really want to touch the filesystem
	err := ioutil.WriteFile("program.s", []byte(assembly.ToString()), 0644)
	if err != nil {
		panic(err)
	}
	command := exec.Command("gcc", "program.s", "-o", "program")
	fmt.Println("Running command:", command.String())
	// Run GCC to create native binary from assembly code
	// We want combined output since it shows standard output and error.
	// If GCC fails to compile the assembly, it will show us why
	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Println("Error running command:", string(output))
		panic(err)
	}
}
