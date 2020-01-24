package main

import "fmt"

var (
	strLabelNum  = 0
	loopLabelNum = 0
)

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

func (node *OutNode) Generate(program *Program) {
	program.FuncSubSect.Content = append(program.FuncSubSect.Content,
		fmt.Sprintf("mov eax, %d", node.returnValue),
	)
}

func (node *StringLiteralNode) Generate(program *Program) {
	strLabelNum++
	asmLabel := fmt.Sprintf("string%d", strLabelNum)
	node.label = asmLabel
	msgSubSect := &SubSect{
		Label: asmLabel, Content: []string{fmt.Sprintf(`.string "%s\n"`, node.str)},
		Vars: make(map[string]int), Anons: make(map[Node]int),
	}
	program.ConstSect.SubSects = append(program.ConstSect.SubSects, msgSubSect)
}

func (node *LoopNode) Generate(program *Program) {
	loopLabelNum++
	funcSubSect := program.FuncSubSect
	program.FuncStackHead++
	counterPos := program.FuncStackHead
	funcSubSect.Content = append(funcSubSect.Content,
		fmt.Sprintf("mov dword ptr [rbp - %d], %d", counterPos*4, node.start),
	)
	funcSubSect.Content = append(funcSubSect.Content,
		fmt.Sprintf("_loopCheck%d:", loopLabelNum),
		fmt.Sprintf("cmp dword ptr [rbp - %d], %d", counterPos*4, node.end),
		fmt.Sprintf("jge _loopContinue%d", loopLabelNum),
		fmt.Sprintf("jmp _loopBody%d", loopLabelNum),
	)
	funcSubSect.Content = append(funcSubSect.Content,
		fmt.Sprintf("_loopBody%d:", loopLabelNum),
	)
	for _, child := range node.children {
		child.Generate(program)
	}
	funcSubSect.Content = append(funcSubSect.Content,
		fmt.Sprintf("mov eax, dword ptr [rbp - %d]", counterPos*4),
		"add eax, 1",
		fmt.Sprintf("mov dword ptr [rbp - %d], eax", counterPos*4),
		fmt.Sprintf("jmp _loopCheck%d", loopLabelNum),
	)
	funcSubSect.Content = append(funcSubSect.Content,
		fmt.Sprintf("_loopContinue%d:", loopLabelNum),
	)
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
