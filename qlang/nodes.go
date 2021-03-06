package main

const (
	ifKeyword = "iff"
)

var (
	// Allows us to take q-lang code, which is raw text, and convert it into a node-based Go representation
	factory = map[string]func() ParsableNode{
		ifKeyword: func() ParsableNode { return new(IfNode) },
		"pkg":     func() ParsableNode { return new(PckgNode) },
		"def":     func() ParsableNode { return new(DefFuncNode) },
		"imp":     func() ParsableNode { return new(ImplFuncNode) },
		"dat":     func() ParsableNode { return new(DefDatNode) },
		"out":     func() ParsableNode { return new(OutNode) },
		"itr":     func() ParsableNode { return new(LoopNode) },
	}
	OperatorFactory = map[string]func() OperationalNode{
		"+": func() OperationalNode { return new(AddNode) },
		"-": func() OperationalNode { return new(SubtractionNode) },
		"*": func() OperationalNode { return new(MulNode) },
		"/": func() OperationalNode { return new(DivisionNode) },
		"&": func() OperationalNode { return new(AndNode) },
		"|": func() OperationalNode { return new(OrNode) },
		"^": func() OperationalNode { return new(XorNode) },
	}
	// TODO:refactor operators are defined in two locations. Add better system for managing tokens
	tokens = []string{endKeyword, rangeKeyword, delimKeyword, "&&", "&", "||", "|", "{", "}", "(", ")", "->", "+", "-", "*", "/", "'", assignKeyword, "$", "&"}
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
		pckgName  string
		datDefs   map[string]*DefDatNode
		funcImpls map[string]*ImplFuncNode
	}
	PckgNode struct {
		ParseNode
	}
	ParseNode struct {
		BaseNode
	}
	DefFuncNode struct {
		ParseNode
		parameterTypes []string
		retType        string
	}
	ImplFuncNode struct {
		ParseNode
		name           string
		parameterNames []string
	}
	CallFuncNode struct {
		ParseNode
		name   string
		retVar ScopeVar
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
	AddNode struct {
		OperatorNode
	}
	SubtractionNode struct {
		OperatorNode
	}
	MulNode struct {
		OperatorNode
	}
	DivisionNode struct {
		OperatorNode
	}
	AndNode struct {
		OperatorNode
	}
	OrNode struct {
		OperandNode
	}
	XorNode struct {
		OperandNode
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
