package main

import "fmt"

var (
	strLabelNum  = 0
	loopLabelNum = 0
)

func (node *DefSingleVarNode) Generate(program *Program) {
	funcSect := program.FuncSubSect
	// TODO: handle types & sizes of variables instead of assuming 4 byte integers
	varStackPos := program.Scope.AllocVar(node, 4)
	if len(node.children) == 1 {
		funcSect.Content = append(funcSect.Content,
			fmt.Sprintf("mov dword ptr [rbp - %d], %d", varStackPos, node.children[0].(*OperandNode).val),
		)
	} else {
		funcSect.Content = append(funcSect.Content,
			fmt.Sprintf("mov dword ptr [rbp - %d], 0", varStackPos),
		)
		// "Anonymous" variables that have no names can be allocated in a temporary scope
		exprScope := NewScope(program.Scope)
		for _, child := range node.children {
			if operandNode, isOperand := child.(*OperandNode); isOperand {
				// TODO: operand ref to existing var on stack
				operandStackPos := exprScope.AllocVar(operandNode, 4)
				funcSect.Content = append(funcSect.Content,
					fmt.Sprintf("mov dword ptr [rbp - %d], %d", operandStackPos, operandNode.val),
				)
			}
		}
		for nodeIndex, child := range node.children {
			switch child.(type) {
			case *AdditionNode:
				firstOperand, secondOperand := node.children[nodeIndex-1].(*OperandNode), node.children[nodeIndex+1].(*OperandNode)
				funcSect.Content = append(funcSect.Content,
					fmt.Sprintf("mov eax, dword ptr [rbp - %d]", exprScope.GetVarPos(firstOperand)),
					fmt.Sprintf("add eax, dword ptr [rbp - %d]", exprScope.GetVarPos(secondOperand)),
				)
				break
			}
		}
		funcSect.Content = append(funcSect.Content,
			fmt.Sprintf("mov dword ptr [rbp - %d], eax", varStackPos),
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
	}
	program.FuncSubSect = funcSubSect
	program.Scope = NewScope(nil)
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
		fmt.Sprintf("mov eax, %d", 0),
	)
}

func (node *StringLiteralNode) Generate(program *Program) {
	strLabelNum++
	asmLabel := fmt.Sprintf("string%d", strLabelNum)
	node.label = asmLabel
	msgSubSect := &SubSect{
		Label: asmLabel, Content: []string{fmt.Sprintf(`.string "%s\n"`, node.str)},
	}
	program.ConstSect.SubSects = append(program.ConstSect.SubSects, msgSubSect)
}

func (node *LoopNode) Generate(program *Program) {
	loopLabelNum++
	funcSubSect := program.FuncSubSect
	loopScope := NewScope(program.Scope)
	counterPos := loopScope.AllocVar(&BaseNode{}, 4)
	funcSubSect.Content = append(funcSubSect.Content,
		fmt.Sprintf("mov dword ptr [rbp - %d], %d", counterPos, node.start),
	)
	funcSubSect.Content = append(funcSubSect.Content,
		fmt.Sprintf("_loopCheck%d:", loopLabelNum),
		fmt.Sprintf("cmp dword ptr [rbp - %d], %d", counterPos, node.end),
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
		fmt.Sprintf("mov eax, dword ptr [rbp - %d]", counterPos),
		"add eax, 1",
		fmt.Sprintf("mov dword ptr [rbp - %d], eax", counterPos),
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
