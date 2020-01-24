package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
)

var (
	labelNum = 0
)

func (node *StringLiteralNode) Parse(scanner *Scanner) {
	node.str = scanner.Next(StrSplit)
	fmt.Printf("String literal: '%s'\n", node.str)
}

func (node *StringLiteralNode) Generate(program *Program) {
	labelNum++
	asmLabel := fmt.Sprintf("string%d", labelNum)
	node.label = asmLabel
	msgSubSect := &SubSect{
		Label: asmLabel, Content: []string{fmt.Sprintf(`.string "%s\n"`, node.str)},
		Vars: make(map[string]int), Anons: make(map[Node]int),
	}
	program.ConstSect.SubSects = append(program.ConstSect.SubSects, msgSubSect)
}

func (node *CallFuncNode) Parse(scanner *Scanner) {
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

func Parse(fileName string) *ProgNode {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	code := string(bytes)
	return getProgram(code)
}

func getProgram(code string) *ProgNode {
	program := new(ProgNode)
	program.Parse(NewScanner(code))
	return program
}

func (node *DefIntNode) Parse(scanner *Scanner) {
	if scanner.Next(Split) != "{" {
		panic("Expected block for variable declaration!")
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		name := nextToken
		varNode := &DefSingleIntNode{name: name}
		node.children = append(node.children, varNode)
		if scanner.Next(Split) != "=" {
			panic("Expected assignment!")
		}
		for {
			nextToken = scanner.Next(Split)
			if nextToken == ";" {
				break
			}
			value := nextToken
			var childNode Node
			// TODO post-processing needs to be done on nodes to make a tree which can be systematically simplified
			if operatorFunc, isOperator := OperatorFactory[value]; isOperator {
				childNode = operatorFunc()
			} else {
				i, _ := strconv.Atoi(value)
				childNode = &IntOperandNode{i: i}
			}
			childNode.SetParent(varNode)
			varNode.Add(childNode)
			fmt.Println("Variable node:", name, "with value:", value)
		}
	}
}

func (node *ProgNode) Parse(scanner *Scanner) {
	nodeName := scanner.Next(Split)
	if nodeName != "pkg" {
		panic("Expected package first!")
	}
	node.pckgName = scanner.Next(Split)
	fmt.Println("Package name:", node.pckgName)
	if scanner.Next(Split) != "{" {
		panic("Expected block for package!")
	}
	nodeName = scanner.Next(Split)
	childNode := Factory[nodeName]()
	node.parseAndAdd(childNode, scanner)
}

func (node *ParseNode) Parse(scanner *Scanner) {
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
	var word string
	skipLength := 0
	for {
		// Cut out all spaces in the beginning
		blankLength := 0
		for ; blankLength < len(rest); blankLength++ {
			if !unicode.IsSpace(rune(rest[blankLength])) {
				break
			}
		}
		word = rest[blankLength:]
		skipLength += blankLength
		// Skip line comment if present
		if strings.HasPrefix(word, "#") {
			commentLength := strings.Index(word, "\n")
			rest = word[commentLength:]
			skipLength += commentLength
		} else {
			break
		}
	}
	// If we start with a token, extract it
	for _, token := range tokens {
		if strings.HasPrefix(word, token) {
			return token, len(token) + skipLength
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
	return word[:wordLength], skipLength + wordLength
}

func (node *CallFuncNode) Generate(program *Program) {
	if node.name == "pln" {
		funcSubSect := program.FuncSubSect
		strNode := node.children[0].(*StringLiteralNode)
		strNode.Generate(program)
		funcSubSect.Content = append(funcSubSect.Content,
			fmt.Sprintf("lea rax, [rip + _%s]", strNode.label),
			"mov rsi, rax # Pointer to string",
			fmt.Sprintf("mov rdx, %d # Size", len(strNode.str)+1),
			"mov rax, 0x2000004 # Write",
			"mov rdi, 1 # Standard output",
			"syscall",
		)
		//program.ConstSect.SubSects.
	}
}

func (node *DefineFuncNode) Parse(scanner *Scanner) {
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

func (node *ImplFuncNode) Parse(scanner *Scanner) {
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

func (node *ImplFuncNode) Generate(program *Program) {
	content := &[]string{
		"push rbp",
		"mov rbp, rsp",
		"",
	}
	funcSubSect := &SubSect{
		Label: "main", Content: *content,
		Vars: make(map[string]int), Anons: make(map[Node]int),
	}
	program.FuncSubSect = funcSubSect
	program.FuncStackHead = 0
	program.FuncSect.SubSects = append(program.FuncSect.SubSects, funcSubSect)
	for _, child := range node.children {
		child.Generate(program)
	}
	funcSubSect.Content = append(funcSubSect.Content,
		"",
		"pop rbp",
		"ret",
	)
}

func (node *OutNode) Parse(scanner *Scanner) {
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

func (node *OutNode) Generate(program *Program) {
	program.FuncSubSect.Content = append(program.FuncSubSect.Content,
		fmt.Sprintf("mov eax, %d", node.returnValue),
	)
}

func (node *BaseNode) parseNextChild(nextToken string, scanner *Scanner) {
	if nodeFunc, isNode := Factory[nextToken]; isNode {
		childNode := nodeFunc()
		node.parseAndAdd(childNode, scanner)
	} else {
		callNode := &CallFuncNode{name: nextToken}
		fmt.Println("Function Call:", callNode.name)
		node.parseAndAdd(callNode, scanner)
	}
}

func (node *BaseNode) parseAndAdd(childNode ParsableNode, scanner *Scanner) {
	childNode.Parse(scanner)
	childNode.SetParent(node)
	node.Add(childNode)
}
