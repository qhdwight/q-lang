package node

import (
	"fmt"
	"q-lang-go/parser/gen"
	"q-lang-go/parser/util"
	"strings"
	"unicode"
)

var (
	// Allows us to take q-lang code, which is raw text, and convert it into Go structs
	Factory = map[string]func() ParsableNode{
		"pkg": func() ParsableNode { return new(PackageNode) },
		"def": func() ParsableNode { return new(DefineFuncNode) },
		"imp": func() ParsableNode { return new(ImplFuncNode) },
		"i32": func() ParsableNode { return new(IntNode) },
	}
	// TODO operators are defined in two locations. Add better system for managing tokens
	tokens = []string{";", ",", "&&", "||", "{", "}", "(", ")", "->", "+", "-", "*", "/", "'"}
)

type Node interface {
	Add(child Node)
	GetChildren() []Node
	SetParent(parent Node)
	Generate(program *gen.Program) // Create assembly string
}

type ParsableNode interface {
	Node
	Parse(scanner *util.Scanner)
}

type BaseNode struct {
	children []Node
	Parent   Node
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

type ProgramNode struct {
	ParseNode
	constants   map[string]Constant
	packageName string
}

type PackageNode struct {
	ParseNode
}

type ParseNode struct {
	BaseNode
}

func (node *BaseNode) Generate(program *gen.Program) {
	for _, child := range node.children {
		child.Generate(program)
	}
}

func (node *ProgramNode) Generate(program *gen.Program) {
	program.ConstantsSection = &gen.Section{
		Decorators: []string{".data"},
	}
	program.FuncSection = &gen.Section{
		Decorators: []string{".text", ".intel_syntax noprefix", ".globl _main"},
	}
	node.BaseNode.Generate(program)
}

func (node *ProgramNode) Parse(scanner *util.Scanner) {
	nodeName := scanner.Next(Split)
	if nodeName != "pkg" {
		panic("Expected package first!")
	}
	node.packageName = scanner.Next(Split)
	fmt.Println("Package name:", node.packageName)
	if scanner.Next(Split) != "{" {
		panic("Expected block for package!")
	}
	nodeName = scanner.Next(Split)
	childNode := Factory[nodeName]()
	node.parseAndAdd(childNode, scanner)
}

func (node *ParseNode) Parse(scanner *util.Scanner) {
	nodeName := scanner.Next(Split)
	fmt.Println(nodeName)
	node.parseNextChild(nodeName, scanner)
}

func StrSplit(rest string) (string, int) {
	for i := 0; i < len(rest); i++ {
		if rune(rest[i]) == '\'' {
			return rest[:i], i + 1
		}
	}
	panic("String literal is not terminated!")
}

func Split(rest string) (string, int) {
	// Cut out all spaces in the beginning
	blankLength := 0
	for ; blankLength < len(rest); blankLength++ {
		if !unicode.IsSpace(rune(rest[blankLength])) {
			break
		}
	}
	word := rest[blankLength:]
	// If we start with a token, extract it
	for _, token := range tokens {
		if strings.HasPrefix(word, token) {
			return token, len(token) + blankLength
		}
	}
	wordLength := 0
	// Extend word until we find a space or a token
findEnd:
	for ; wordLength < len(word); wordLength++ {
		if unicode.IsSpace(rune(word[wordLength])) {
			break
		}
		for _, token := range tokens {
			if strings.HasPrefix(word[wordLength:], token) {
				break findEnd
			}
		}
	}
	return word[:wordLength], blankLength + wordLength
}
