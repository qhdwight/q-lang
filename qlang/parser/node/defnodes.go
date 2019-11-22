package node

var (
	Factory = map[string]func() ParsableNode{
		"pkg": func() ParsableNode { return new(PackageNode) },
		"def": func() ParsableNode { return new(DefineFuncNode) },
		"imp": func() ParsableNode { return new(ImplFuncNode) },
		"i32": func() ParsableNode { return new(IntNode) },
	}
)

type Node interface {
	Add(child Node)
	SetParent(parent Node)
}

type ParsableNode interface {
	Node
	Parse(body, innerBody string, tokens []string, parent Node)
	Generate()
}

type DeclarationNode interface {
	ParsableNode
}

type BaseNode struct {
	children []Node
	Parent   Node
}

func (node *BaseNode) Add(child Node) {
	node.children = append(node.children, child)
}

func (node *BaseNode) SetParent(parent Node) {
	node.Parent = parent
}

type ProgramNode struct {
	ParseNode
	Constants map[string]Constant
}

type PackageNode struct {
	ParseNode
}

type ParseNode struct {
	BaseNode
	body, innerBody string
	tokens          []string
}

func (node *ParseNode) Generate() {
	panic("implement me")
}

func (node *ParseNode) Parse(body, innerBody string, tokens []string, parent Node) {
	node.SetParent(parent)
	node.body = body
	node.innerBody = innerBody
	node.tokens = tokens
}
