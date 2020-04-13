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
	tokens = []string{";", ",", "&&", "||", "{", "}", "(", ")", "->", "+", "-", "*", "/", "'"}
)

type Node interface {
	Add(child Node)
	GetChildren() []Node
	SetParent(parent Node)
	Generate(program *Program) // Create assembly string
}

type ParsableNode interface {
	Node
	Parse(scanner *Scanner)
}

type OperableNode interface {
	Node
}

type SimplifiableNode interface {
	Node
}

type BaseNode struct {
	children []Node
	Parent   Node
}

type ExpressionNode struct {
	BaseNode
}

func (node *BaseNode) Add(child Node) {
	node.children = append(node.children, child)
}

func (node *BaseNode) GetChildren() []Node {
	return node.children
}

func (node *BaseNode) SetParent(parent Node) {
	node.Parent = parent
}

type LoopNode struct {
	ParseNode
	start, end int
}

type ProgNode struct {
	ParseNode
	pckgName string
}

type PckgNode struct {
	ParseNode
}

type ParseNode struct {
	BaseNode
}

type DefineFuncNode struct {
	ParseNode
}

type ImplFuncNode struct {
	ParseNode
}

type CallFuncNode struct {
	ParseNode
	name string
}

type ArgumentNode struct {
	ParseNode
}

type OutNode struct {
	ParseNode
}

type StringLiteralNode struct {
	ParseNode
	str, label string
}

type DefIntNode struct {
	DefVarNode
}

type ImplVarNode struct {
	BaseNode
}

type DefVarNode struct {
	ParseNode
}

type DefSingleIntNode struct {
	ParseNode
	name string
}

type IntOperandNode struct {
	BaseNode
	val int
}

type IntOperatorNode struct {
	ExpressionNode
}

type AdditionNode struct {
	IntOperatorNode
}

type SubtractionNode struct {
	IntOperatorNode
}

type MultiplicationNode struct {
	IntOperatorNode
}

type DivisionNode struct {
	IntOperatorNode
}
