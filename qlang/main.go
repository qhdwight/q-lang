package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	outputProgram *string
)

func main() {
	inputFiles := flag.String("input", "program.qq", "Input Q files")
	outputProgram = flag.String("output", "program", "Output program file")
	flag.Parse()
	prog := Parse(*inputFiles)

	asm := NewProg()
	prog.Generate(asm)
	fmt.Println("Output asm:\n", asm.ToString())

	assemble(asm)
}

func assemble(assembly *Prog) {
	// TODO currently we write to file first and then call GCC. Find a way to pass it without file as we don't really want to touch the filesystem
	asmFilePath := fmt.Sprintf("%s.s", *outputProgram)
	_ = os.MkdirAll(filepath.Dir(asmFilePath), 0644)
	err := ioutil.WriteFile(asmFilePath, []byte(assembly.ToString()), 0644)
	if err != nil {
		panic(err)
	}
	command := exec.Command("gcc", asmFilePath, "-o", *outputProgram)
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
