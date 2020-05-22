package main

var (
	// Allows us to take q-lang code, which is raw text, and convert it into a node-based Go representation
	factory = map[string]func() ParsableNode{
		"iff": func() ParsableNode { return new(IfNode) },
		"pkg": func() ParsableNode { return new(PckgNode) },
		"def": func() ParsableNode { return new(DefineFuncNode) },
		"imp": func() ParsableNode { return new(ImplFuncNode) },
		"dat": func() ParsableNode { return new(DefDatNode) },
		"out": func() ParsableNode { return new(OutNode) },
		"for": func() ParsableNode { return new(LoopNode) },
	}
	OperatorFactory = map[string]func() OperationalNode{
		"+": func() OperationalNode { return new(AdditionNode) },
		"-": func() OperationalNode { return new(SubtractionNode) },
		"*": func() OperationalNode { return new(MultiplicationNode) },
		"/": func() OperationalNode { return new(DivisionNode) },
	}
	// TODO:refactor operators are defined in two locations. Add better system for managing tokens
	tokens = []string{";", "..", ",", "&&", "||", "{", "}", "(", ")", "->", "+", "-", "*", "/", "'", ":=", "=", "$", "&"}
)

type (
	Node interface {
		Add(child Node)
		GetChildren() []Node
		SetParent(parent Node)
		Generate(prog *Prog) // Create assembly text
	}
	ParsableNode interface {
		Node
		Parse(scanner *Scanner)
	}
	OperationalNode interface {
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
		start, end *OperandNode
	}
	ProgNode struct {
		ParseNode
		pckgName string
		datDefs  map[string]*DefDatNode
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
	NamedVarsNode struct {
		ParseNode
		typeName string
	}
	SingleNamedVarNode struct {
		ExpressionNode
		name, typeName string
	}
	AssignmentNode struct {
		ParseNode
		accessor string // Left-hand side
	}
	OperandNode struct {
		BaseNode
		typeName   string // Always
		accessor   string // If reference to variable
		literalVal string // If storing literal (int)
	}
	PropFill struct {
		ParseNode
		propDef *DefDatPropNode
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

	/* Control flow */

	IfNode struct {
		ParseNode
		o1, o2 *OperandNode
		t, f   *BaseNode
	}

	/* Data system */

	DefDatNode struct {
		ParseNode
		name string
		size int
	}

	DefDatPropNode struct {
		ParseNode
		name, typeName string
		offset         int
	}
)
