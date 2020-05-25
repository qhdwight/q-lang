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

const (
	uintKeyword   = "u32"
	floatKeyword  = "f32"
	assignKeyword = "<-"
	endKeyword    = ";"
	rangeKeyword  = "..."
	delimKeyword  = ","
)

func getProg(code string) *ProgNode {
	progNode = &ProgNode{datDefs: map[string]*DefDatNode{}, funcImpls: map[string]*ImplFuncNode{}}
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
	parseNextStatementNode(nodeName, node, scanner)
}

func (node *StringLiteralNode) Parse(scanner *Scanner) {
	node.str = scanner.Next(StrSplit)
	fmt.Printf("String literal: '%s'\n", node.str)
}

func (node *CallFuncNode) Parse(scanner *Scanner) {
	fmt.Println("Function Call:", node.name)
	switch node.name {
	case "wln":
		if scanner.PeekAdvanceIf(Split, func(str string) bool { return str == "'" }) {
			literalNode := new(StringLiteralNode)
			parseAndAdd(literalNode, node, scanner)
			if scanner.Next(Split) != endKeyword {
				panic("Expecting semicolon to end string literal!")
			}
		} else {
			node.parseExpr(scanner, func(token string) bool { return token == endKeyword })
		}
	case "rln":
		if scanner.Next(Split) != "$" {
			panic("Expecting reference to 32-bit integer!")
		}
		operandNode := parseOperand(scanner, scanner.Next(Split))
		node.children = append(node.children, operandNode)
		operandNode.Parent = operandNode
		if scanner.Next(Split) != endKeyword {
			panic("Expecting semicolon to end function call!")
		}
	default:
		if _, isFuncImpl := progNode.funcImpls[node.name]; isFuncImpl {
			for {
				nextToken, amount := scanner.Peek(Split)
				if nextToken == delimKeyword {
					scanner.SkipBy(amount)
					continue
				}
				if nextToken == endKeyword {
					break
				}
				scanner.Skip(Split)
				operand := parseOperand(scanner, nextToken)
				node.children = append(node.children, operand)
				operand.Parent = node
			}
		} else {
			panic(fmt.Sprintf("Unrecognized function call %s", node.name))
		}
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
		nextToken = scanner.Next(Split)
		if nextToken != assignKeyword {
			panic("Expected assignment")
		}
		varNode.parseExpr(scanner, func(token string) bool { return token == endKeyword })
		fmt.Println("Variable node:", name)
	}
}

func (node *BaseNode) parseExpr(scanner *Scanner, tokenBreakCond func(token string) bool) {
	for {
		nextToken := scanner.Next(Split)
		if tokenBreakCond(nextToken) {
			break
		}
		if nextToken == "?" {
			if !tokenBreakCond(scanner.Next(Split)) {
				panic("Expected end after ?")
			}
			return
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
	// Unsigned 32-bit integer literal
	_, err := strconv.Atoi(strVal)
	if err == nil {
		operandNode.typeName, operandNode.literalVal = uintKeyword, strVal
		return operandNode
	}
	// 32-bit float literal
	_, err = strconv.ParseFloat(strVal, 64)
	if err == nil {
		operandNode.typeName, operandNode.literalVal = floatKeyword, strVal
		return operandNode
	}
	// Data constructor
	if datDef, isDat := progNode.datDefs[strVal]; isDat {
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
			if scanner.Next(Split) != assignKeyword {
				panic("Expected assignment!")
			}
			propDef := getImmediatePropDef(datDef, propName)
			fillNode := &PropFill{propDef: propDef}
			parseAndAdd(fillNode, operandNode, scanner)
		}
		return operandNode
	}
	// Function call
	if funcImpl, isCall := progNode.funcImpls[strVal]; isCall {
		callNode := &CallFuncNode{name: strVal}
		operandNode.typeName = funcImpl.Parent.(*DefFuncNode).retType
		parseAndAdd(callNode, operandNode, scanner)
	}
	// Accessor to existing variable
	operandNode.accessor = strVal
	return operandNode
}

func (node *PropFill) Parse(scanner *Scanner) {
	node.parseExpr(scanner, func(token string) bool { return token == endKeyword })
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
		parseAndAdd(childNode, node, scanner)
	}
}

func (node *LoopNode) Parse(scanner *Scanner) {
	node.start = parseOperand(scanner, scanner.Next(Split))
	if scanner.Next(Split) != rangeKeyword {
		panic("Expected comma")
	}
	node.end = parseOperand(scanner, scanner.Next(Split))
	fmt.Println("Loop")
	if scanner.Next(Split) != "{" {
		panic("Expected block")
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		parseNextStatementNode(nextToken, node, scanner)
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

func (node *DefFuncNode) Parse(scanner *Scanner) {
	for {
		nextToken := scanner.Next(Split)
		if nextToken == delimKeyword {
			continue
		}
		if nextToken == "->" {
			break
		}
		parameterType := nextToken
		if parameterType == "nil" {
			continue
		}
		node.parameterTypes = append(node.parameterTypes, parameterType)
		fmt.Println("Parameter type:", parameterType)
	}
	node.retType = scanner.Next(Split)
	if scanner.Next(Split) != "{" {
		panic("Expected body for function definition!")
	}
	fmt.Println("Return type:", node.retType)
	if scanner.Next(Split) != "imp" {
		panic("Expected function implementation!")
	}
	implNode := new(ImplFuncNode)
	parseAndAdd(implNode, node, scanner)
	if scanner.Next(Split) != "}" {
		panic("Expected close to block!")
	}
}

func (node *ImplFuncNode) Parse(scanner *Scanner) {
	node.name = scanner.Next(Split)
	progNode.funcImpls[node.name] = node
	fmt.Println("Function name:", node.name)
	for {
		nextToken := scanner.Next(Split)
		if nextToken == delimKeyword {
			continue
		}
		if nextToken == "{" {
			break
		}
		node.parameterNames = append(node.parameterNames, nextToken)
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		parseNextStatementNode(nextToken, node, scanner)
	}
}

func (node *OutNode) Parse(scanner *Scanner) {
	node.parseExpr(scanner, func(token string) bool { return token == endKeyword })
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
			if scanner.Next(Split) != endKeyword {
				panic("Expected semicolon to end property definition!")
			}
		}
	}
}

// TODO:refactor remove next token and use peek instead
func parseNextStatementNode(nextToken string, parent Node, scanner *Scanner) {
	if nextToken == endKeyword {
		return
	}
	var childNode ParsableNode
	if nodeFunc, isNode := factory[nextToken]; isNode {
		childNode = nodeFunc()
	} else if nextToken == uintKeyword || nextToken == floatKeyword {
		childNode = &NamedVarsNode{typeName: nextToken}
	} else if datDef, isDatDef := progNode.datDefs[nextToken]; isDatDef {
		childNode = &NamedVarsNode{typeName: datDef.name}
	} else if scanner.PeekAdvanceIf(Split, func(str string) bool { return str == assignKeyword }) {
		childNode = &AssignmentNode{accessor: nextToken}
	} else {
		childNode = &CallFuncNode{name: nextToken}
	}
	parseAndAdd(childNode, parent, scanner)
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
		parseNextStatementNode(nextToken, node.t, scanner)
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
		parseNextStatementNode(nextToken, node.f, scanner)
	}
}

func (node *AssignmentNode) Parse(scanner *Scanner) {
	fmt.Println("Assignment to:", node.accessor)
	node.parseExpr(scanner, func(token string) bool { return token == endKeyword })
}

func parseAndAdd(childNode ParsableNode, parent Node, scanner *Scanner) {
	childNode.Parse(scanner)
	childNode.SetParent(parent)
	parent.Add(childNode)
}
