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
	parameterType := scanner.Next(util.Split)
	if scanner.Next(util.Split) != "->" {
		panic("Expected return type!")
	}
	fmt.Println("Parameter type: ", parameterType)
	returnType := scanner.Next(util.Split)
	if scanner.Next(util.Split) != "{" {
		panic("Expected block!")
	}
	fmt.Println("Return type: ", returnType)
	if scanner.Next(util.Split) != "imp" {
		panic("Expected function implementation!")
	}
	implNode := new(ImplFuncNode)
	implNode.Parse(scanner)
	implNode.SetParent(node)
	node.Add(implNode)
}

func (node *ImplFuncNode) Parse(scanner *util.Scanner) {
	funcName := scanner.Next(util.Split)
	if scanner.Next(util.Split) != "{" {
		return
	}
	fmt.Println("Func name: ", funcName)
}
