package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"unicode"
)

func getProgram(code string) *ProgNode {
	program := new(ProgNode)
	program.Parse(NewScanner(code))
	return program
}

func Parse(fileName string) *ProgNode {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	code := string(bytes)
	return getProgram(code)
}

func (node *ParseNode) Parse(scanner *Scanner) {
	nodeName := scanner.Next(Split)
	fmt.Println(nodeName)
	node.parseNextChild(nodeName, scanner)
}

func (node *StringLiteralNode) Parse(scanner *Scanner) {
	node.str = scanner.Next(StrSplit)
	fmt.Printf("String literal: '%s'\n", node.str)
}

func (node *CallFuncNode) Parse(scanner *Scanner) {
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "'" {
			literalNode := new(StringLiteralNode)
			node.parseAndAdd(literalNode, scanner)
		} else if nextToken == ";" {
			break
		} else {
			node.Add(parseOperand(nextToken))
		}
	}
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
		varNode := &DefSingleVarNode{name: name}
		node.children = append(node.children, varNode)
		if scanner.Next(Split) != "=" {
			panic("Expected assignment!")
		}
		for {
			nextToken = scanner.Next(Split)
			if nextToken == ";" {
				break
			}
			var childNode Node
			if operatorFunc, isOperator := OperatorFactory[nextToken]; isOperator {
				childNode = operatorFunc()
			} else {
				operandNode := parseOperand(nextToken)
				childNode = operandNode
			}
			childNode.SetParent(varNode)
			varNode.Add(childNode)
		}
		fmt.Println("Variable node:", name)
	}
}

func parseOperand(strVal string) *OperandNode {
	operandNode := &OperandNode{}
	val, err := strconv.Atoi(strVal)
	if err == nil {
		operandNode.val = val
	} else {
		// Reference to existing variable
		operandNode.varName = strVal
	}
	return operandNode
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
	childNode := factory[nodeName]()
	node.parseAndAdd(childNode, scanner)
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
		node.parseNextChild(nextToken, scanner)
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
	fmt.Println("Func name:", funcName)
	if scanner.Next(Split) != "{" {
		panic("Expected block in function implementation!")
	}
	for {
		nextToken := scanner.Next(Split)
		if nextToken == "}" {
			break
		}
		node.parseNextChild(nextToken, scanner)
		if nextToken == "out" {
			break
		}
	}
}

func (node *OutNode) Parse(scanner *Scanner) {
	// TODO: implement
}

func (node *BaseNode) parseNextChild(nextToken string, scanner *Scanner) {
	if nodeFunc, isNode := factory[nextToken]; isNode {
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
