package parser

import (
	"fmt"
	"io/ioutil"
	"q-lang-go/parser/node"
	"strings"
)

func Parse(fileName string) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	code := string(bytes)
	getProgram(code)
}

func getProgram(code string) *node.ProgramNode {
	code = strings.ReplaceAll(code, "\n", "")
	program := new(node.ProgramNode)
	recurseNodes(code, &program.Base, 0)
	return program
}

func recurseNodes(code string, parent node.CodeNode, depth int) {
	level, blockInfoStart, blockStart := 0, 0, 0
	for i := 0; i < len(code); i++ {
		c := code[i]
		if c == '{' {
			if level == 0 {
				blockStart = i + 1
			}
			level++
		} else if c == '}' {
			level--
			if level == 0 {
				blockWithInfo := code[blockInfoStart:i]
				blockWithInfo = strings.TrimSpace(blockWithInfo)
				tokens := strings.FieldsFunc(blockWithInfo, func(c rune) bool {
					return c == ' '
				})
				nodeName := tokens[0]
				blockContentStart := blockStart
				fmt.Println(nodeName)
				child := node.Factory[nodeName]()
				child.Parse()
				parent.Add(child)
				recurseNodes(code[blockContentStart:i], child, depth+1)
				blockInfoStart = i + 1
			}
		}
	}
}
