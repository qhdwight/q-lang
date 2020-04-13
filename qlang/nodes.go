package main

var (
	// Allows us to take q-lang code, which is raw text, and convert it into a Go representation
	factory = map[string]func() ParsableNode{
		"pkg": func() ParsableNode { return new(PckgNode) },
		"def": func() ParsableNode { return new(DefineFuncNode) },
		"imp": func() ParsableNode { return new(ImplFuncNode) },
		"i32": func() ParsableNode { return new(DefIntNode) },
		"out": func() ParsableNode { return new(OutNode) },
		"for": func() ParsableNode { return new(LoopNode) },
	}
	OperatorFactory = map[string]func() OperableNode{
		"+": func() OperableNode { return new(AdditionNode) },
		"-": func() OperableNode { return new(SubtractionNode) },
		"*": func() OperableNode { return new(MultiplicationNode) },
		"/": func() OperableNode { return new(DivisionNode) },
	}
	// TODO operators are defined in two locations. Add better system for managing tokens
	tokens = []string{";", "..", ",", "&&", "||", "{", "}", "(", ")", "->", "+", "-", "*", "/", "'"}
)

type (
	Node interface {
		Add(child Node)
		GetChildren() []Node
		SetParent(parent Node)
		Generate(program *Program) // Create assembly string
	}
	ParsableNode interface {
		Node
		Parse(scanner *Scanner)
	}
	OperableNode interface {
		Node
	}
)

func (node *BaseNode) Add(child Node) {
	node.children = append(node.children, child)
}

func (node *BaseNode) GetChildren() []Node {
	return node.children
}

func (node *BaseNode) SetParent(parent Node) {
	node.Parent = parent
}

type (
	BaseNode struct {
		children []Node
		Parent   Node
	}
	ExpressionNode struct {
		BaseNode
	}
	LoopNode struct {
		ParseNode
		start, end int
	}
	ProgNode struct {
		ParseNode
		pckgName string
	}
	PckgNode struct {
		ParseNode
	}
	ParseNode struct {
		BaseNode
	}
	DefineFuncNode struct {
		ParseNode
	}
	ImplFuncNode struct {
		ParseNode
	}
	CallFuncNode struct {
		ParseNode
		name string
	}
	OutNode struct {
		ParseNode
	}
	StringLiteralNode struct {
		ParseNode
		str, label string
	}
	DefIntNode struct {
		DefVarNode
	}
	DefVarNode struct {
		ParseNode
	}
	DefSingleVarNode struct {
		ParseNode
		name string
	}
	OperandNode struct {
		BaseNode
		val int
	}
	OperatorNode struct {
		ExpressionNode
	}
	AdditionNode struct {
		OperatorNode
	}
	SubtractionNode struct {
		OperatorNode
	}
	MultiplicationNode struct {
		OperatorNode
	}
	DivisionNode struct {
		OperatorNode
	}
)
