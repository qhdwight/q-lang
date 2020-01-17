package node

import (
	"fmt"
	"q-lang-go/parser/gen"
	"q-lang-go/parser/util"
)

var (
	OperatorFactory = map[string]func() OperableNode{
		"+": func() OperableNode { return new(AdditionNode) },
		"-": func() OperableNode { return new(SubtractionNode) },
		"*": func() OperableNode { return new(MultiplicationNode) },
		"/": func() OperableNode { return new(DivisionNode) },
	}
)

type IntNode struct {
	DefVarNode
}

type ImplVarNode struct {
	BaseNode
}

type DefVarNode struct {
	ParseNode
}

type OperandNode struct {
	BaseNode
}

type OperatorNode struct {
	BaseNode
}

type AdditionNode struct {
	OperatorNode
}

type SubtractionNode struct {
	OperatorNode
}

type MultiplicationNode struct {
	OperatorNode
}

type DivisionNode struct {
	OperatorNode
}

type OperableNode interface {
	Node
}

func (node *IntNode) Parse(scanner *util.Scanner) {
	if scanner.Next(Split) != "{" {
		panic("Expected block for variable declaration!")
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		name := nextToken
		if scanner.Next(Split) != "=" {
			panic("Expected assignment!")
		}
		for {
			nextToken = scanner.Next(Split)
			if nextToken == ";" {
				break
			}
			value := nextToken
			var childNode Node
			// TODO post-processing needs to be done on nodes to make a tree which can be systematically simplified
			if operatorFunc, isOperator := OperatorFactory[value]; isOperator {
				childNode = operatorFunc()
			} else {
				childNode = new(OperandNode)
			}
			childNode.SetParent(node)
			node.Add(childNode)
			fmt.Println("Variable node:", name, "with value:", value)
		}
	}
}

func (node *DefVarNode) Parse(scanner *util.Scanner) {

}

func (node *DefVarNode) Generate(program *gen.Program) {
	//for _, s := range program.FuncSection.SubSections {
	//}
}
