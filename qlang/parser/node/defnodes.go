package node

var (
	Factory = map[string]func() ParsableNode{
		"pkg": func() ParsableNode { return new(PackageNode) },
		"def": func() ParsableNode { return new(DefineFunctionNode) },
		"imp": func() ParsableNode { return new(ImplementFunctionNode) },
	}
)

type Node interface {
	Add(child Node)
}

type ParsableNode interface {
	Node
	Parse(body, innerBody string, tokens []string, parent Node)
	Generate()
}

type BaseNode struct {
	Node
	children []Node
	Parent   Node
}

func (node *BaseNode) Add(child Node) {
	node.children = append(node.children, child)
}

type ProgramNode struct {
	ParsableNode
	Base      BaseNode
	Constants map[string]Constant
}

type PackageNode struct {
	ParsableNode
	parseNode ParseNode
}

func (node *PackageNode) Add(child Node) {
	node.parseNode.base.Add(child)
}

func (node *PackageNode) Parse(body, innerBody string, tokens []string, parent Node) {
	node.parseNode.Parse(body, innerBody, tokens, parent)
}

type ParseNode struct {
	ParsableNode
	base            BaseNode
	body, innerBody string
	tokens          []string
}

func (node *ParseNode) Parse(body, innerBody string, tokens []string, parent Node) {
	node.base.Parent = parent
	node.body = body
	node.innerBody = innerBody
	node.tokens = tokens
}
