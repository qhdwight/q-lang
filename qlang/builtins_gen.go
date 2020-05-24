package main

import "fmt"

func (node *CallFuncNode) writeLine(prog *Prog) {
	if strNode, isStr := node.children[0].(*StringLiteralNode); isStr {
		strNode.Generate(prog)
		prog.CurSect.Content = append(prog.CurSect.Content,
			fmt.Sprintf("lea rax, [rip + _%s]", strNode.label),
			"mov rsi, rax # Pointer to string",
			fmt.Sprintf("mov rdx, %d # Size", len(strNode.str)+1),
			"mov rax, 0x2000004 # Write",
			"mov rdi, 1 # Standard output",
			"syscall",
		)
	} else {
		prog.Scope = NewScope(prog.Scope)
		operandPos := node.genExpr(prog).stackPos
		bufferPos := prog.Scope.Alloc(16)
		// Algorithm in C: https://gcc.godbolt.org/z/3b2P5j
		if !libraryLabelExists(prog, "charLoop") {
			prog.LibrarySubSect.SubSects = append(prog.LibrarySubSect.SubSects, &Sect{
				Label: "digitToChar",
				Content: []string{
					"movabs r8, -3689348814741910323",
					"xor r9, r9",                 // Character count - xor with self cheaply zeroes
					"mov byte ptr [rsi - 1], 10", // Add newline at the end
					"dec rsi",
					"inc r9",
				},
			})
			prog.LibrarySubSect.SubSects = append(prog.LibrarySubSect.SubSects, &Sect{
				Label: "charLoop",
				Content: []string{ // u32: edi, u8[]: rsi
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
			prog.Scope = prog.Scope.Parent
		}
		prog.CurSect.Content = append(prog.CurSect.Content,
			fmt.Sprintf("mov edi, dword ptr [rbp - %d] # Integer argument", operandPos),
			fmt.Sprintf("lea rsi, [rbp - %d] # Buffer pointer argument", bufferPos-16),
			"call _digitToChar",
			// fmt.Sprintf("lea rsi, [rbp - %d]", bufferPos),
			"mov rdx, r9 # Size",
			"mov rax, 0x2000004 # Write system call identifier",
			"mov rdi, 1 # Standard output file descriptor",
			"syscall", // edi already set earlier
		)
	}
}

func (node *CallFuncNode) readLine(prog *Prog) {
	// Algorithm in C: https://godbolt.org/z/qKpeBB
	bufferPos := prog.Scope.Alloc(16)
	ref := genOperandVar(prog, node.children[0].(*OperandNode))
	if !libraryLabelExists(prog, "asciiToInt32") {
		prog.LibrarySubSect.SubSects = append(prog.LibrarySubSect.SubSects, &Sect{
			Label: "asciiToInt32",
			Content: []string{
				"xor edx, edx",
				"xor eax, eax",
				"dec esi",
				"_whileBody:",
				"cmp esi, edx",
				"jle _whileEnd",
				"imul eax, eax, 10",
				"movsx ecx, byte ptr [rdi+rdx]",
				"inc rdx",
				"lea eax, [rax-48+rcx]",
				"jmp _whileBody",
				"_whileEnd:",
				"ret",
			},
		})
	}
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("lea rsi, [rbp - %d] # Pointer to ASCII buffer", bufferPos),
		fmt.Sprintf("mov rdx, %d # Size", 16),
		"mov rax, 0x2000003 # Read system call identifier",
		"mov rdi, 0 # Standard input file descriptor",
		"syscall",

		"mov esi, eax # Length of characters",
		fmt.Sprintf("lea rdi, [rbp - %d] # Pointer to ASCII buffer", bufferPos),
		"call _asciiToInt32",

		fmt.Sprintf("mov dword ptr [rbp - %d], eax", ref.stackPos),
	)
}
