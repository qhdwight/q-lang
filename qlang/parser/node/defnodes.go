package node

import (
	"fmt"
	"q-lang-go/parser/util"
)

var (
	Factory = map[string]func() ParsableBlockNode{
		"pkg": func() ParsableBlockNode { return new(PackageNode) },
		"def": func() ParsableBlockNode { return new(DefineFuncNode) },
		"imp": func() ParsableBlockNode { return new(ImplFuncNode) },
		"i32": func() ParsableBlockNode { return new(IntNode) },
	}
)

type Node interface {
	Add(child Node)
	SetParent(parent Node)
}

type ParsableBlockNode interface {
	Node
	Parse(scanner *util.Scanner)
	Generate()
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
	ParseBlockNode
	Constants   map[string]Constant
	packageName string
}

type PackageNode struct {
	ParseBlockNode
}

type ParseBlockNode struct {
	BaseNode
}

func (node ParseBlockNode) Generate() {
	panic("implement me")
}

func (node *ProgramNode) Parse(scanner *util.Scanner) {
	nodeName := scanner.Next(util.Split)
	if nodeName != "pkg" {
		panic("Expected package first!")
	}
	node.packageName = scanner.Next(util.Split)
	fmt.Println("Package name: ", node.packageName)
	if scanner.Next(util.Split) == "{" {
		nodeName = scanner.Next(util.Split)
		childNode := Factory[nodeName]()
		childNode.Parse(scanner)
		node.Add(childNode)
		childNode.SetParent(node)
	}
}

func (node *ParseBlockNode) Parse(scanner *util.Scanner) {
	nodeName := scanner.Next(util.Split)
	fmt.Println(nodeName)
	childNode := Factory[nodeName]()
	childNode.Parse(scanner)
	node.Add(childNode)
	childNode.SetParent(node)
}
