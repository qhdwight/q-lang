package node

import (
	"fmt"
	"q-lang-go/parser/util"
)

type IntNode struct {
	DefVarNode
}

type ImplVarNode struct {
	BaseNode
}

type DefVarNode struct {
	ParseBlockNode
}

func (node *IntNode) Parse(scanner *util.Scanner) {

}

func (node *DefVarNode) Parse(scanner *util.Scanner) {
	if scanner.Next(Split) != "{" {
		panic("Expected block for variable declaration!")
	}
	for {
		name := scanner.Next(Split)
		if scanner.Next(Split) != "=" {
			panic("Expected assignment!")
		}
		value := scanner.Next(Split)
		if scanner.Next(Split) != ";" {
			panic("Expected semicolon in variable definition!")
		}
		fmt.Println("Variable named:", name, "with value:", value)
		if nodeFunc, isNode := Factory[name]; isNode {
			childNode := nodeFunc()
			childNode.Parse(scanner)
			childNode.SetParent(node)
			node.Add(childNode)
		}
		if scanner.Next(Split) == "}" {
			break
		}
	}
}
