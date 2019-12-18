package node

import (
	"fmt"
	"q-lang-go/parser/util"
)

type DefineFuncNode struct {
	ParseBlockNode
}

type ImplFuncNode struct {
	ParseBlockNode
}

func (node *DefineFuncNode) Parse(scanner *util.Scanner) {
	parameterType := scanner.Next(Split)
	if scanner.Next(Split) != "->" {
		panic("Expected return type!")
	}
	fmt.Println("Parameter type:", parameterType)
	returnType := scanner.Next(Split)
	if scanner.Next(Split) != "{" {
		panic("Expected body for function definition!")
	}
	fmt.Println("Return type:", returnType)
	if scanner.Next(Split) != "imp" {
		panic("Expected function implementation!")
	}
	implNode := new(ImplFuncNode)
	node.parseAndAdd(implNode, scanner)
}

func (node *ImplFuncNode) Parse(scanner *util.Scanner) {
	funcName := scanner.Next(Split)
	if scanner.Next(Split) != "{" {
		panic("Expected block in function implementation!")
	}
	nextToken := scanner.Next(Split)
	node.parseNextChild(nextToken, scanner)
	fmt.Println("Func name: ", funcName)
}

func (node *ImplFuncNode) Generate() string {
	return `_main:
	push rbp
	mov	rbp, rsp`
}

func (node *BaseNode) parseNextChild(nextToken string, scanner *util.Scanner) {
	if nodeFunc, isNode := Factory[nextToken]; isNode {
		childNode := nodeFunc()
		node.parseAndAdd(childNode, scanner)
	} else {
		panic("Unrecognized phrase")
	}
}

func (node *BaseNode) parseAndAdd(childNode ParsableBlockNode, scanner *util.Scanner) {
	childNode.Parse(scanner)
	childNode.SetParent(node)
	node.Add(childNode)
}
