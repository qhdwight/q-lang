package main

import "fmt"

var (
	strLabelNum  = 0
	loopLabelNum = 0
)

func (node *DefSingleVarNode) Generate(program *Prog) {
	// TODO: handle types & sizes of variables instead of assuming 4 byte integers
	varStackPos := program.Scope.AllocVar(node, 4)
	program.Scope = NewScope(program.Scope)
	if len(node.children) == 1 {
		operand := node.children[0].(*OperandNode)
		program.CurSect.Content = append(program.CurSect.Content,
			fmt.Sprintf("mov dword ptr [rbp - %d], %d", varStackPos, operand.val),
		)
	} else {
		program.CurSect.Content = append(program.CurSect.Content,
			fmt.Sprintf("mov dword ptr [rbp - %d], 0", varStackPos),
		)
		// "Anonymous" variables that have no names can be allocated in a temporary scope
		for _, child := range node.children {
			if operandNode, isOperand := child.(*OperandNode); isOperand {
				if len(operandNode.varName) == 0 {
					operandStackPos := program.Scope.AllocVar(operandNode, 4)
					program.CurSect.Content = append(program.CurSect.Content,
						fmt.Sprintf("mov dword ptr [rbp - %d], %d", operandStackPos, operandNode.val),
					)
				}
			}
		}
		for nodeIndex, child := range node.children {
			switch child.(type) {
			case *AdditionNode:
				firstOperand, secondOperand := node.children[nodeIndex-1].(*OperandNode), node.children[nodeIndex+1].(*OperandNode)
				program.CurSect.Content = append(program.CurSect.Content,
					fmt.Sprintf("mov eax, dword ptr [rbp - %d]", program.Scope.GetVarPos(firstOperand)),
					fmt.Sprintf("add eax, dword ptr [rbp - %d]", program.Scope.GetVarPos(secondOperand)),
				)
				break
			}
		}
		program.CurSect.Content = append(program.CurSect.Content,
			fmt.Sprintf("mov dword ptr [rbp - %d], eax", varStackPos),
		)
	}
	program.Scope = program.Scope.Parent
}

func (node *BaseNode) Generate(program *Prog) {
	for _, child := range node.children {
		child.Generate(program)
	}
}

func (node *ImplFuncNode) Generate(program *Prog) {
	program.CurSect.Content = append(program.CurSect.Content,
		"push rbp",
		"mov rbp, rsp",
		"sub rsp, 64",
		"",
	)
	for _, child := range node.children {
		child.Generate(program)
	}
	program.CurSect.Content = append(program.CurSect.Content,
		"",
		"add rsp, 64",
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
	program.Scope = NewScope(program.Scope)
	counterPos := program.Scope.AllocVar(&DefSingleVarNode{name: "__counter"}, 4)
	program.CurSect.Content = append(program.CurSect.Content,
		fmt.Sprintf("mov dword ptr [rbp - %d], %d # Counter", counterPos, node.start),
	)
	program.CurSect.Content = append(program.CurSect.Content,
		fmt.Sprintf("_loopCheck%d:", loopLabelNum),
		fmt.Sprintf("cmp dword ptr [rbp - %d], %d", counterPos, node.end),
		fmt.Sprintf("jge _loopContinue%d", loopLabelNum),
		fmt.Sprintf("jmp _loopBody%d", loopLabelNum),
	)
	program.CurSect.Content = append(program.CurSect.Content,
		fmt.Sprintf("_loopBody%d:", loopLabelNum),
	)
	for _, child := range node.children {
		child.Generate(program)
	}
	program.CurSect.Content = append(program.CurSect.Content,
		fmt.Sprintf("mov eax, dword ptr [rbp - %d]", counterPos),
		"inc eax",
		fmt.Sprintf("mov dword ptr [rbp - %d], eax", counterPos),
		fmt.Sprintf("jmp _loopCheck%d", loopLabelNum),
	)
	program.CurSect.Content = append(program.CurSect.Content,
		fmt.Sprintf("_loopContinue%d:", loopLabelNum),
	)
	program.Scope = program.Scope.Parent
}

func (node *CallFuncNode) Generate(program *Prog) {
	if node.name == "pln" {
		strNode := node.children[0].(*StringLiteralNode)
		strNode.Generate(program)
		program.CurSect.Content = append(program.CurSect.Content,
			fmt.Sprintf("lea rax, [rip + _%s]", strNode.label),
			"mov rsi, rax # Pointer to string",
			fmt.Sprintf("mov rdx, %d # Size", len(strNode.str)+1),
			"mov rax, 0x2000004 # Write",
			"mov rdi, 1 # Standard output",
			"syscall",
		)
	} else if node.name == "i_pln" {
		program.Scope = NewScope(program.Scope)
		operand := node.children[0].(*OperandNode)
		operand.Generate(program)
		operandPos := program.Scope.GetVarPos(operand)
		bufferPos := program.Scope.AllocVar(&DefSingleVarNode{name: "__buffer"}, 16)
		// Algorithm in C: https://gcc.godbolt.org/z/3b2P5j
		needsLibrary := true
		for _, subSect := range program.LibrarySubSect.SubSects {
			if subSect.Label == "charLoop" {
				needsLibrary = false
				break
			}
		}
		if needsLibrary {
			program.LibrarySubSect.SubSects = append(program.LibrarySubSect.SubSects, &Sect{
				Label: "digitToChar",
				Content: []string{
					"movabs r8, -3689348814741910323",
					"xor r9, r9",                 // Character count - xor with self cheaply zeroes
					"mov byte ptr [rsi - 1], 10", // Add newline at the end
					"dec rsi",
					"inc r9",
				},
			})
			program.LibrarySubSect.SubSects = append(program.LibrarySubSect.SubSects, &Sect{
				Label: "charLoop",
				Content: []string{ // i32: edi, char[]: rsi
					"movsxd rax, edi", // Cast signed 32 bit to unsigned 64 bit
					"mul r8",
					"shr rdx, 3", // Shift three right
					"lea eax, [rdx + rdx]",
					"lea eax, [rax + 4*rax]",
					"mov ecx, edi",
					"sub ecx, eax",
					"or cl, 48",                  // 48 is '0' in ASCII, serves as our base for digit representation
					"mov byte ptr [rsi - 1], cl", // 1 byte for character
					"dec rsi",
					"inc r9",
					"cmp rdi, 9", // Check if we are above 9 (have more digits left)
					"mov rdi, rdx",
					"ja _charLoop", // Check cmp instruction to loop if we have more digits
					"ret",
				},
			})
			program.Scope = program.Scope.Parent
		}
		program.CurSect.Content = append(program.CurSect.Content,
			fmt.Sprintf("mov edi, dword ptr [rbp - %d] # Integer argument", operandPos),
			fmt.Sprintf("lea rsi, [rbp - %d] # Buffer pointer argument", bufferPos-16),
			"call _digitToChar",
			// fmt.Sprintf("lea rsi, [rbp - %d]", bufferPos),
			"mov rdx, r9 # Size",
			"mov rax, 0x2000004 # Write",
			"mov rdi, 1 # Standard output",
			"syscall", // edi already set earlier
		)
	}
}

func (node *ProgNode) Generate(program *Prog) {
	node.BaseNode.Generate(program)
}
