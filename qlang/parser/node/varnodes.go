package node

type LineNode interface {
	ParsableNode
}

type IntNode struct {
	parseNode ParseNode
}

func (node *IntNode) Add(child Node) {
	node.parseNode.Add(child)
}

func (node *IntNode) Parse(body, innerBody string, tokens []string, parent Node) {
}

func (node *IntNode) Generate() {
}




