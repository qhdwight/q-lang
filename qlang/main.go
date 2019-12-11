package main

import (
	"flag"
	"q-lang-go/parser"
)

func main() {
	inputFiles := flag.String("input", "", "Input Q files")
	flag.Parse()
	if *inputFiles == "" {
		panic("No input Q files provided!")
	}
	parser.Parse(*inputFiles)

	//fmt.Println("[" + strings.Join(util.Tokenize("x =  2 + && a", []string{
	//	"+", "-", "*", "/", "&&", "||",
	//}), "|") + "]")
}
