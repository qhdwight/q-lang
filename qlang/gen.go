package main

import "fmt"

func (node *DefSingleIntNode) Generate(program *Program) {
	funcSect := program.FuncSubSect
	program.FuncStackHead++
	varStackPos := program.FuncStackHead
	anons := funcSect.Anons
	if len(node.children) == 1 {
		funcSect.Content = append(funcSect.Content,
			fmt.Sprintf("mov dword ptr [rbp - %d], %d", varStackPos*4, node.children[0].(*IntOperandNode).i),
		)
	} else {
		funcSect.Content = append(funcSect.Content,
			fmt.Sprintf("mov dword ptr [rbp - %d], 0", varStackPos*4),
		)
		for _, child := range node.children {
			if operandNode, isOperand := child.(*IntOperandNode); isOperand {
				program.FuncStackHead++
				operandStackPos := program.FuncStackHead
				funcSect.Content = append(funcSect.Content,
					fmt.Sprintf("mov dword ptr [rbp - %d], %d", operandStackPos*4, operandNode.i),
				)
				anons[operandNode] = operandStackPos
			}
		}
		for nodeIndex, child := range node.children {
			switch child.(type) {
			case *AdditionNode:
				firstOperand, secondOperand := node.children[nodeIndex-1].(*IntOperandNode), node.children[nodeIndex+1].(*IntOperandNode)
				funcSect.Content = append(funcSect.Content,
					fmt.Sprintf("mov eax, dword ptr [rbp - %d]", anons[firstOperand]*4),
					fmt.Sprintf("add eax, dword ptr [rbp - %d]", anons[secondOperand]*4),
				)
				break
			}
		}
		funcSect.Content = append(funcSect.Content,
			fmt.Sprintf("mov dword ptr [rbp - %d], eax", varStackPos*4),
		)
	}
}

func (node *BaseNode) Generate(program *Program) {
	for _, child := range node.children {
		child.Generate(program)
	}
}

func (node *ProgNode) Generate(program *Program) {
	program.ConstSect = &Sect{
		Decorators: []string{"data"},
	}
	program.FuncSect = &Sect{
		Decorators: []string{"text", "intel_syntax noprefix", "globl _main"},
	}
	node.BaseNode.Generate(program)
}
