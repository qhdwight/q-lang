package node

import (
	"fmt"
	"q-lang-go/parser/gen"
	"q-lang-go/parser/util"
	"strconv"
)

var (
	labelNum = 0
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

type OutNode struct {
	ParseNode
	returnValue int
}

type StringLiteralNode struct {
	ParseNode
	str, label string
}

func (node *StringLiteralNode) Parse(scanner *util.Scanner) {
	node.str = scanner.Next(StrSplit)
	fmt.Printf("String literal: '%s'\n", node.str)
}

func (node *StringLiteralNode) Generate(program *gen.Program) {
	labelNum++
	asmLabel := "string" + strconv.Itoa(labelNum)
	node.label = asmLabel
	msgSubSection := &gen.SubSection{
		Label:     asmLabel,
		Content:   []string{fmt.Sprintf(`.string "%s\n"`, node.str)},
		Variables: make(map[string]int),
	}
	program.ConstSections.SubSections = append(program.ConstSections.SubSections, msgSubSection)
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
		funcSubSection := program.FuncSubSection
		strNode := node.children[0].(*StringLiteralNode)
		strNode.Generate(program)
		funcSubSection.Content = append(funcSubSection.Content,
			fmt.Sprintf("lea rax, [rip + _%s]", strNode.label),
			"mov rsi, rax # Pointer to string",
			fmt.Sprintf("mov rdx, %d # Size", len(strNode.str)+1),
			"mov rax, 0x2000004 # Write",
			"mov rdi, 1 # Standard output",
			"syscall",
		)
		//program.ConstSections.SubSections.
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
	for {
		nextToken := scanner.Next(Split)
		node.parseNextChild(nextToken, scanner)
		if nextToken == "out" {
			break
		}
		fmt.Println("Func name:", funcName)
	}
}

func (node *ImplFuncNode) Generate(program *gen.Program) {
	content := &[]string{
		"push rbp",
		"mov rbp, rsp",
		"",
	}
	funcSubSection := &gen.SubSection{Label: "main", Content: *content, Variables: make(map[string]int)}
	program.FuncSubSection = funcSubSection
	program.FuncStackHead = 0
	program.FuncSection.SubSections = append(program.FuncSection.SubSections, funcSubSection)
	for _, child := range node.children {
		child.Generate(program)
	}
	funcSubSection.Content = append(funcSubSection.Content,
		"",
		"pop rbp",
		"ret",
	)
}

func (node *OutNode) Parse(scanner *util.Scanner) {
	for {
		var err error
		node.returnValue, err = strconv.Atoi(scanner.Next(Split))
		if err != nil {
			panic(err)
		}
		nextToken := scanner.Next(Split)
		if nextToken == ";" {
			break
		}
	}
}

func (node *OutNode) Generate(program *gen.Program) {
	program.FuncSubSection.Content = append(program.FuncSubSection.Content,
		fmt.Sprintf("mov eax, %d", node.returnValue),
	)
}

func (node *BaseNode) parseNextChild(nextToken string, scanner *util.Scanner) {
	if nodeFunc, isNode := Factory[nextToken]; isNode {
		childNode := nodeFunc()
		node.parseAndAdd(childNode, scanner)
	} else {
		callNode := &CallFuncNode{name: nextToken}
		fmt.Println("Function Call:", callNode.name)
		node.parseAndAdd(callNode, scanner)
	}
}

func (node *BaseNode) parseAndAdd(childNode ParsableNode, scanner *util.Scanner) {
	childNode.Parse(scanner)
	childNode.SetParent(node)
	node.Add(childNode)
}
