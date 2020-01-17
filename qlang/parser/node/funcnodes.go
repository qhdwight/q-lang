package node

import (
	"fmt"
	"q-lang-go/parser/gen"
	"q-lang-go/parser/util"
)

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

type StringLiteralNode struct {
	ParseNode
	str string
}

func (node *StringLiteralNode) Parse(scanner *util.Scanner) {
	node.str = scanner.Next(StrSplit)
	fmt.Println("String literal:", "'"+node.str+"'")
}

func (node *StringLiteralNode) Generate(program *gen.Program) {
	program.CurrentSubSection.Content = append(program.CurrentSubSection.Content, `.string "`+node.str+`\n"`)
}

func (node *CallFuncNode) Parse(scanner *util.Scanner) {
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "'" {
			literalNode := new(StringLiteralNode)
			node.parseAndAdd(literalNode, scanner)
		} else if nextToken == ";" {
			break
		}
	}
}

func (node *CallFuncNode) Generate(program *gen.Program) {
	if node.name == "pln" {
		subSection := &gen.SubSection{
			Name: "_message:",
		}
		funcSubSection := program.CurrentSubSection
		program.ConstantsSection.SubSections = append(program.ConstantsSection.SubSections, subSection)
		program.CurrentSubSection = subSection
		node.children[0].Generate(program)
		funcSubSection.Content = append(funcSubSection.Content,
			"lea rax, [rip + _message]",
			"mov rsi, rax # Pointer to string",
			"mov rdx, 13 # Size",
			"mov rax, 0x2000004 # Write",
			"mov rdi, 1 # Standard output",
			"syscall",
		)
		//program.ConstantsSection.SubSections.
	}
}

func (node *DefineFuncNode) Parse(scanner *util.Scanner) {
	parameterType := scanner.Next(Split)
	if scanner.Next(Split) != "->" {
		panic("Expected return type!")
	}
	fmt.Println("Parameter type:", parameterType)
	returnType := scanner.Next(Split)
	if scanner.Next(Split) != "{" {
		panic("Expected body for function definition!")
	}
	fmt.Println("Return type:", returnType)
	if scanner.Next(Split) != "imp" {
		panic("Expected function implementation!")
	}
	implNode := new(ImplFuncNode)
	node.parseAndAdd(implNode, scanner)
}

func (node *ImplFuncNode) Parse(scanner *util.Scanner) {
	funcName := scanner.Next(Split)
	if scanner.Next(Split) != "{" {
		panic("Expected block in function implementation!")
	}
	nextToken := scanner.Next(Split)
	node.parseNextChild(nextToken, scanner)
	fmt.Println("Func name:", funcName)
}

func (node *ImplFuncNode) Generate(program *gen.Program) {
	content := &[]string{
		"push rbp",
		"mov rbp, rsp",
	}
	funcSubSection := &gen.SubSection{Name: "_main:", Content: *content}
	program.FuncSection.SubSections = append(program.FuncSection.SubSections, funcSubSection)
	program.CurrentSubSection = funcSubSection
	for _, child := range node.children {
		child.Generate(program)
	}
	funcSubSection.Content = append(funcSubSection.Content,
		"pop rbp",
		"ret",
	)
}

func (node *BaseNode) parseNextChild(nextToken string, scanner *util.Scanner) {
	if nodeFunc, isNode := Factory[nextToken]; isNode {
		childNode := nodeFunc()
		node.parseAndAdd(childNode, scanner)
	} else {
		callNode := &CallFuncNode{name: nextToken}
		fmt.Println("Calling Function:", callNode)
		node.parseAndAdd(callNode, scanner)
	}
}

func (node *BaseNode) parseAndAdd(childNode ParsableNode, scanner *util.Scanner) {
	childNode.Parse(scanner)
	childNode.SetParent(node)
	node.Add(childNode)
}
