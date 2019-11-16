package node

var (
	Factory = map[string]func() CodeNode{
		"pkg": func() CodeNode { return new(PackageNode) },
		"def": func() CodeNode { return new(DefineFunctionNode) },
		"imp": func() CodeNode { return new(ImplementFunctionNode) },
	}
)

type Node interface {
	Add(child Node)
}

type CodeNode interface {
	Node
	Parse()
	Generate()
}

type BaseNode struct {
	Node
	children []Node
	Parent   Node
}

func (node *BaseNode) Generate() {

}

func (node *BaseNode) Parse() {
}

func (node *BaseNode) Add(child Node) {
	node.children = append(node.children, child)
}

type ProgramNode struct {
	CodeNode
	Base BaseNode
}

type PackageNode struct {
	CodeNode
	parseNode ParseNode
}

func (node *PackageNode) Parse() {
}

type DefineFunctionNode struct {
	CodeNode
	parseNode ParseNode
}

func (node *DefineFunctionNode) Parse() {
}

type ImplementFunctionNode struct {
	CodeNode
	parseNode ParseNode
}

func (node *ImplementFunctionNode) Parse() {
}

type ParseNode struct {
	base   BaseNode
	body   string
	tokens []string
}
