package node

import (
	"fmt"
	"q-lang-go/parser/util"
	"strings"
	"unicode"
)

var (
	Factory = map[string]func() ParsableBlockNode{
		"pkg": func() ParsableBlockNode { return new(PackageNode) },
		"def": func() ParsableBlockNode { return new(DefineFuncNode) },
		"imp": func() ParsableBlockNode { return new(ImplFuncNode) },
		"i32": func() ParsableBlockNode { return new(IntNode) },
	}
	tokens = []string{";", ",", "&&", "||", "{", "}", "(", ")"}
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

func (node *ParseBlockNode) Parse(scanner *util.Scanner) {
	nodeName := scanner.Next(Split)
	fmt.Println(nodeName)
	node.parseNextChild(nodeName, scanner)
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
