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
	recurseNodes(code, program, 0)
	return program
}

func recurseNodes(code string, parent node.Node, depth int) {
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
				blockWithInfo := code[blockInfoStart : i+1]
				blockWithInfo = strings.TrimSpace(blockWithInfo)
				tokens := strings.FieldsFunc(blockWithInfo, func(c rune) bool {
					return c == ' '
				})
				nodeName := tokens[0]
				fmt.Println(nodeName)
				blockContentStart := blockStart
				blockContents := code[blockContentStart-blockInfoStart : i]
				nodeFunc, exists := node.Factory[nodeName]
				if exists {
					child := nodeFunc()
					child.Parse(blockWithInfo, blockContents, tokens, parent)
					parent.Add(child)
					//switch v := child.(type) {
					//case node.DeclarationNode:
					//	fmt.Println(nodeName)
					//default:
					//	recurseNodes(code[blockContentStart:i], v, depth+1)
					//}
					recurseNodes(code[blockContentStart:i], child, depth+1)
					blockInfoStart = i + 1
				} else {
					fmt.Printf("Unrecognized keyword: %s\n", nodeName)
				}
			}
		}
	}
}
