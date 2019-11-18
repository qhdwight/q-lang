package node

import (
	"regexp"
	"strings"
)

var (
	varNodeFactory = map[string]func() LineNode{
		"i32": func() LineNode { return new(IntNode) },
	}
	tokenRegex = regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
)

type DefineFunctionNode struct {
	ParsableNode
	parseNode ParseNode
}

func (node *DefineFunctionNode) Add(child Node) {
	node.parseNode.base.Add(child)
}

func (node *DefineFunctionNode) Parse(body, innerBody string, tokens []string, parent Node) {
	node.parseNode.Parse(body, innerBody, tokens, parent)
}

type ImplementFunctionNode struct {
	ParsableNode
	parseNode ParseNode
}

func (node *ImplementFunctionNode) Parse(body, innerBody string, tokens []string, parent Node) {
	node.parseNode.Parse(body, innerBody, tokens, parent)
	lines := strings.Split(innerBody, ";")
	for _, line := range lines {
		tokens := tokenRegex.FindAllString(line, -1)
		firstToken := tokens[0]
		if varNodeFunc, isVar := varNodeFactory[firstToken]; isVar {
			varNode := varNodeFunc()
			node.parseNode.base.Add(varNode)
		}
	}
}
