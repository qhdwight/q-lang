package main

import "fmt"

var (
	strLabelNum  = 0
	loopLabelNum = 0
)

func (node *DefSingleVarNode) Generate(program *Prog) {
	funcSect := program.MainSubSect
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
				if len(operandNode.varName) == 0 {
					operandStackPos := exprScope.AllocVar(operandNode, 4)
					funcSect.Content = append(funcSect.Content,
						fmt.Sprintf("mov dword ptr [rbp - %d], %d", operandStackPos, operandNode.val),
					)
				}
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

func (node *BaseNode) Generate(program *Prog) {
	for _, child := range node.children {
		child.Generate(program)
	}
}

func (node *ImplFuncNode) Generate(program *Prog) {
	program.MainSubSect.Content = append(program.MainSubSect.Content,
		"push rbp",
		"mov rbp, rsp",
		"",
	)
	for _, child := range node.children {
		child.Generate(program)
	}
	program.MainSubSect.Content = append(program.MainSubSect.Content,
		"",
		"pop rbp",
		"ret",
	)
}

func (node *OutNode) Generate(program *Prog) {
	program.FuncSect.Content = append(program.FuncSect.Content,
		fmt.Sprintf("mov eax, %d", 0),
	)
}

func (node *StringLiteralNode) Generate(program *Prog) {
	strLabelNum++
	asmLabel := fmt.Sprintf("string%d", strLabelNum)
	node.label = asmLabel
	msgSubSect := &Sect{
		Label: asmLabel, Content: []string{fmt.Sprintf(`.string "%s\n"`, node.str)},
	}
	program.ConstSect.SubSects = append(program.ConstSect.SubSects, msgSubSect)
}

func (node *LoopNode) Generate(program *Prog) {
	loopLabelNum++
	funcSubSect := program.MainSubSect
	loopScope := NewScope(program.Scope)
	counterPos := loopScope.AllocVar(&DefSingleVarNode{name: "__counter"}, 4)
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

func (node *CallFuncNode) Generate(program *Prog) {
	if node.name == "pln" {
		strNode := node.children[0].(*StringLiteralNode)
		strNode.Generate(program)
		program.FuncSect.Content = append(program.FuncSect.Content,
			fmt.Sprintf("lea rax, [rip + _%s]", strNode.label),
			"mov rsi, rax # Pointer to string",
			fmt.Sprintf("mov rdx, %d # Size", len(strNode.str)+1),
			"mov rax, 0x2000004 # Write",
			"mov rdi, 1 # Standard output",
			"syscall",
		)
	} else if node.name == "i_pln" {
		// https://gcc.godbolt.org/z/3b2P5j
		operand := node.children[0].(*OperandNode)
		operand.Generate(program)
		operandPos := program.Scope.GetVarPos(operand)
		bufferPos := program.Scope.AllocVar(&DefSingleVarNode{name: "__buffer"}, 16)
		program.LibrarySubSect.SubSects = append(program.LibrarySubSect.SubSects, &Sect{
			Label:   "digitToChar",
			Content: []string{"movabs r8, -3689348814741910323"},
		})
		program.LibrarySubSect.SubSects = append(program.LibrarySubSect.SubSects, &Sect{
			Label: "charLoop",
			Content: []string{
				"mov rax, rdi",
				"mul r8",
				"shr rdx, 3", // Shift three right
				"lea eax, [rdx + rdx]",
				"lea eax, [rax + 4*rax]",
				"mov ecx, edi",
				"sub ecx, eax",
				"or cl, 48", // 48 is '0' in ASCII, serves as our base for digit representation
				"mov byte ptr [rsi - 1], cl",
				"dec rsi",
				"cmp rdi, 9", // Loop if we are above 9 (have more digits left)
				"mov rdi, rdx",
				"ja _charLoop",
				"mov rax, rsi",
				"ret",
			},
		})
		program.MainSubSect.Content = append(program.MainSubSect.Content,
			fmt.Sprintf("mov edi, dword ptr [rbp - %d]", operandPos),
			fmt.Sprintf("lea rsi, dword ptr [rbp - %d]", bufferPos),
			"call _digitToChar",
			fmt.Sprintf("mov rdx, %d # Size", 5),
			"mov rax, 0x2000004 # Write",
			"mov rdi, 1 # Standard output",
			"syscall",
		)
	}
}

func (node *ProgNode) Generate(program *Prog) {
	node.BaseNode.Generate(program)
}
