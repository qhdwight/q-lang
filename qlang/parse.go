package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
)

var (
	// TODO:refactor avoid global variable
	progNode *ProgNode
)

func getProg(code string) *ProgNode {
	progNode = &ProgNode{datDefs: make(map[string]*DefDatNode)}
	progNode.Parse(NewScanner(code))
	return progNode
}

func Parse(fileName string) *ProgNode {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	code := string(bytes)
	return getProg(code)
}

func (node *ParseNode) Parse(scanner *Scanner) {
	nodeName := scanner.Next(Split)
	node.parseNextStatementNode(nodeName, scanner)
}

func (node *StringLiteralNode) Parse(scanner *Scanner) {
	node.str = scanner.Next(StrSplit)
	fmt.Printf("String literal: '%s'\n", node.str)
}

func (node *CallFuncNode) Parse(scanner *Scanner) {
	fmt.Println("Function Call:", node.name)
	if node.name == "wln" {
		if scanner.PeekAdvanceIf(Split, func(str string) bool { return str == "'" }) {
			literalNode := new(StringLiteralNode)
			node.parseAndAdd(literalNode, scanner)
			if scanner.Next(Split) != ";" {
				panic("Expecting semicolon to end string literal!")
			}
		} else {
			node.parseExpr(scanner, func(token string) bool { return token == ";" })
		}
	} else if node.name == "rln" {
		if scanner.Next(Split) != "$" {
			panic("Expecting reference to 32-bit integer!")
		}
		operandNode := parseOperand(scanner, scanner.Next(Split))
		node.children = append(node.children, operandNode)
		operandNode.Parent = operandNode
		if scanner.Next(Split) != ";" {
			panic("Expecting semicolon to end function call!")
		}
	} else {
		panic(fmt.Sprintf("Unrecognized function call %s", node.name))
	}
}

func (node *NamedVarsNode) Parse(scanner *Scanner) {
	if scanner.Next(Split) != "{" {
		panic("Expected block for variable declaration!")
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		name := nextToken
		varNode := &SingleNamedVarNode{name: name, typeName: node.typeName}
		node.children = append(node.children, varNode)
		if scanner.Next(Split) != ":=" {
			panic("Expected assignment to named variable!")
		}
		varNode.parseExpr(scanner, func(token string) bool { return token == ";" })
		fmt.Println("Variable node:", name)
	}
}

func (node *BaseNode) parseExpr(scanner *Scanner, tokenBreakCond func(token string) bool) {
	for {
		nextToken := scanner.Next(Split)
		if tokenBreakCond(nextToken) {
			break
		}
		var child Node
		if operatorFunc, isOperator := OperatorFactory[nextToken]; isOperator {
			child = operatorFunc()
		} else {
			operandNode := parseOperand(scanner, nextToken)
			child = operandNode
		}
		node.Add(child)
		child.SetParent(node)
	}
}

func getImmediatePropDef(datDef *DefDatNode, name string) *DefDatPropNode {
	for _, child := range datDef.children {
		propDef := child.(*DefDatPropNode)
		if propDef.name == name {
			return propDef
		}
	}
	panic("No property found with name")
}

func parseOperand(scanner *Scanner, strVal string) *OperandNode {
	operandNode := &OperandNode{}
	// i32 literal
	_, err := strconv.Atoi(strVal)
	if err == nil {
		operandNode.typeName = "i32"
		operandNode.literalVal = strVal
		return operandNode
	}
	// Data constructor
	for _, datDef := range progNode.datDefs {
		if strVal == datDef.name {
			operandNode.typeName = datDef.name
			if scanner.Next(Split) != "{" {
				panic("Expected block for data property fills")
			}
			for {
				nextToken := scanner.Next(Split)
				if nextToken == "}" {
					break
				}
				propName := nextToken
				if scanner.Next(Split) != "=" {
					panic("Expected assignment!")
				}
				propDef := getImmediatePropDef(datDef, propName)
				fillNode := &PropFill{propDef: propDef}
				operandNode.parseAndAdd(fillNode, scanner)
			}
			return operandNode
		}
	}
	// Accessor to existing variable
	operandNode.accessor = strVal
	return operandNode
}

func (node *PropFill) Parse(scanner *Scanner) {
	node.parseExpr(scanner, func(token string) bool { return token == "," })
}

func (node *ProgNode) Parse(scanner *Scanner) {
	nextToken := scanner.Next(Split)
	if nextToken != "pkg" {
		panic("Expected package first!")
	}
	node.pckgName = scanner.Next(Split)
	fmt.Println("Package name:", node.pckgName)
	if scanner.Next(Split) != "{" {
		panic("Expected block for package!")
	}
	for {
		nextToken = scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		childNode := factory[nextToken]()
		node.parseAndAdd(childNode, scanner)
	}
}

func (node *LoopNode) Parse(scanner *Scanner) {
	if scanner.Next(Split) != "range" {
		panic("Expected range")
	}
	start, err := strconv.Atoi(scanner.Next(Split))
	node.start = start
	if err != nil {
		panic(err)
	}
	if scanner.Next(Split) != ".." {
		panic("Expected comma")
	}
	end, err := strconv.Atoi(scanner.Next(Split))
	node.end = end
	if err != nil {
		panic(err)
	}
	fmt.Println("Loops with range", start, "to", end)
	if scanner.Next(Split) != "{" {
		panic("Expected block")
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		node.parseNextStatementNode(nextToken, scanner)
	}
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
	fmt.Println("Function name:", funcName)
	if scanner.Next(Split) != "{" {
		panic("Expected block in function implementation!")
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		node.parseNextStatementNode(nextToken, scanner)
		if nextToken == "out" {
			break
		}
	}
}

func (node *OutNode) Parse(scanner *Scanner) {
	// TODO:warning implement
	scanner.Next(Split)
	scanner.Next(Split)
}

func (node *DefDatNode) Parse(scanner *Scanner) {
	node.name = scanner.Next(Split)
	progNode.datDefs[node.name] = node
	fmt.Println("Data structure with name:", node.name)
	if scanner.Next(Split) != "{" {
		panic("Expected block in data definition!")
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		typeName := nextToken
		if scanner.Next(Split) != "{" {
			panic("Expected block after type for property definition!")
		}
		for {
			nextToken := scanner.Next(Split)
			if nextToken == "}" {
				break
			}
			defPropNode := &DefDatPropNode{typeName: typeName, name: nextToken}
			node.children = append(node.children, defPropNode)
			defPropNode.Parent = node
			if scanner.Next(Split) != ";" {
				panic("Expected semicolon to end property definition!")
			}
		}
	}
}

// TODO:refactor remove next token and use peek instead
func (node *BaseNode) parseNextStatementNode(nextToken string, scanner *Scanner) {
	var childNode ParsableNode
	if nodeFunc, isNode := factory[nextToken]; isNode {
		childNode = nodeFunc()
	} else if nextToken == "i32" {
		childNode = &NamedVarsNode{typeName: nextToken}
	} else if datDef, isDatDef := progNode.datDefs[nextToken]; isDatDef {
		childNode = &NamedVarsNode{typeName: datDef.name}
	} else if scanner.PeekAdvanceIf(Split, func(str string) bool { return str == "=" }) {
		childNode = &AssignmentNode{accessor: nextToken}
	} else {
		childNode = &CallFuncNode{name: nextToken}
	}
	node.parseAndAdd(childNode, scanner)
}

func (node *IfNode) Parse(scanner *Scanner) {
	node.o1 = parseOperand(scanner, scanner.Next(Split))
	if scanner.Next(Split) != "=" {
		panic("Expected comparison!")
	}
	node.o2 = parseOperand(scanner, scanner.Next(Split))
	if scanner.Next(Split) != "{" {
		panic("Expected block for if statement!")
	}
	node.t, node.f = new(BaseNode), new(BaseNode)
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		node.t.parseNextStatementNode(nextToken, scanner)
	}
	if !scanner.PeekAdvanceIf(Split, func(str string) bool { return str == "els" }) {
		return
	}
	if scanner.Next(Split) != "{" {
		panic("Expected block for else statement!")
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		node.f.parseNextStatementNode(nextToken, scanner)
	}
}

func (node *AssignmentNode) Parse(scanner *Scanner) {
	fmt.Println("Assignment to:", node.accessor)
	node.parseExpr(scanner, func(token string) bool { return token == ";" })
}

func (node *BaseNode) parseAndAdd(childNode ParsableNode, scanner *Scanner) {
	childNode.Parse(scanner)
	childNode.SetParent(node)
	node.Add(childNode)
}
