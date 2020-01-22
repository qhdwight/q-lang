package node

import (
	"fmt"
	"q-lang-go/parser/gen"
	"q-lang-go/parser/util"
	"strconv"
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
	i int
	n string
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
				i, _ := strconv.Atoi(value)
				childNode = &OperandNode{i: i, n: name}
			}
			childNode.SetParent(node)
			node.Add(childNode)
			fmt.Println("Variable node:", name, "with value:", value)
		}
	}
}

func (node *IntNode) Generate(program *gen.Program) {
	for _, child := range node.children {
		if operandNode, isOperand := child.(*OperandNode); isOperand {
			program.FuncStackHead += 4
			program.FuncSubSection.Content = append(program.FuncSubSection.Content,
				fmt.Sprintf("mov dword ptr [rbp - %d], %d", program.FuncStackHead, operandNode.i),
			)
			program.FuncSubSection.Variables[operandNode.n] = program.FuncStackHead
		}
	}
	for nodeIndex, child := range node.children {
		if _, isOperator := child.(*AdditionNode); isOperator {
			firstOperand, secondOperand := node.children[nodeIndex-1].(*OperandNode), node.children[nodeIndex+1].(*OperandNode)
			variables := program.FuncSubSection.Variables
			program.FuncSubSection.Content = append(program.FuncSubSection.Content,
				fmt.Sprintf("mov eax, dword ptr [rbp - %d]", variables[firstOperand.n]),
				fmt.Sprintf("add eax, dword ptr [rbp - %d]", variables[secondOperand.n]),
			)
		}
	}
	program.FuncSubSection.Content = append(program.FuncSubSection.Content,
	)
}

func (node *DefVarNode) Parse(scanner *util.Scanner) {

}

func (node *DefVarNode) Generate(program *gen.Program) {
	//for _, s := range program.FuncSection.SubSections {
	//}
}
