package main

import (
	"fmt"
	"runtime"
)

const (
	writeSyscall = iota
	readSyscall
)

func getSyscall(id int) string {
	if runtime.GOOS == "linux" {
		switch id {
		case writeSyscall:
			return "0x1"
		case readSyscall:
			return "0x0"
		}
	} else if runtime.GOOS == "darwin" {
		switch id {
		case writeSyscall:
			return "0x2000004"
		case readSyscall:
			return "0x2000003"
		}
	}
	panic("Unsupported OS")
}

func (node *CallFuncNode) writeLine(prog *Prog) {
	if strNode, isStr := node.children[0].(*StringLiteralNode); isStr {
		strNode.Generate(prog)
		prog.CurSect.Content = append(prog.CurSect.Content,
			fmt.Sprintf("lea rax, [rip + _%s]", strNode.label),
			"mov rsi, rax # Pointer to string",
		fmt.Sprintf("mov rdx, %d # Size", len(strNode.str)+1),
			fmt.Sprintf("mov rax, %s # Write system call identifier", getSyscall(writeSyscall)),
			"mov rdi, 1 # Standard output",
			"syscall",
	)
	} else {
		typeName := genOperandVar(prog, node.children[0].(*OperandNode)).typeName
		if typeName == uintKeyword {
			node.writeUInt(prog)
		} else if typeName == floatKeyword {
			node.writeFloat(prog)
		} else {
			panic("Type not supported for printing")
		}
	}
}

func (node *CallFuncNode) writeFloat(prog *Prog) {
	// https://godbolt.org/z/UpTbKK
	if !libraryLabelExists(prog, "floatToAscii") {
		prog.LibrarySubSect.SubSects = append(prog.LibrarySubSect.SubSects, &Sect{
			Label: "floatToAscii",
			Content: []string{
				"xorps xmm1, xmm1",
				"comiss xmm1, xmm0",
				"movaps xmm1, xmm0",
				"jbe _L13",
				"mulss xmm1, dword ptr _LC1[rip]",
				"mov ecx, 10000",
				"cvttss2si eax, xmm1",
				"movaps xmm1, xmm0",
				"xorps xmm1, xmmword PTR _LC2[rip]",
				"cdq",
				"idiv ecx",
				"cvttss2si ecx, xmm1",
				"mov eax, edx",
				"jmp _L4",
			},
		}, &Sect{
			Label: "L13",
			Content: []string{
				"mulss xmm1, dword ptr _LC3[rip]",
				"mov ecx, 10000",
				"cvttss2si eax, xmm1",
				"cdq",
				"idiv ecx",
				"cvttss2si ecx, xmm0",
				"mov eax, edx",
			},
		}, &Sect{
			Label: "L4",
			Content: []string{
				"mov byte ptr [rdi+15], 10",
				"xor esi, esi",
				"mov r8d, 10",
			},
		}, &Sect{
			Label: "L5",
			Content: []string{
				"xor edx, edx",
				"div r8w",
				"add edx, 48",
				"mov byte ptr [rdi+14+rsi], dl",
				"dec rsi",
				"cmp rsi, -4",
				"jne _L5",
				"mov byte ptr [rdi+10], 46",
				"lea rsi, [rdi+10]",
				"mov r8d, 6",
				"mov edi, 10",
			},
		}, &Sect{
			Label: "L7",
			Content: []string{
				"test ecx, ecx",
				"jle _L6",
				"mov eax, ecx",
				"dec rsi",
				"inc r8d",
				"cdq",
				"idiv edi",
				"add edx, 48",
				"mov ecx, eax",
				"mov byte ptr [rsi], dl",
				"jmp _L7",
			},
		}, &Sect{
			Label: "L6",
			Content: []string{
				"xorps xmm1, xmm1",
				"comiss xmm1, xmm0",
				"jbe _L1",
				"mov byte ptr [rsi-1], 45",
			},
		}, &Sect{
			Label: "L1",
			Content: []string{
				"mov eax, r8d",
				"ret",
				"_LC1:",
				".long -971227136",
				"_LC2:",
				".long -2147483648",
				".long 0",
				".long 0",
				".long 0",
				"_LC3:",
				".long 1176256512",
			},
		})
	}
	prog.Scope = NewScope(prog.Scope)
	operandPos := node.genExpr(prog).stackPos
	bufferPos := prog.Scope.Alloc(16)
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("movups xmm0, xmmword ptr [rbp - %d] # Float argument", operandPos),
		fmt.Sprintf("lea rdi, [rbp - %d] # Buffer pointer argument", bufferPos),
		"call _floatToAscii",
		fmt.Sprintf("lea rsi, [rbp - %d] # Buffer pointer argument", bufferPos-16),
		"sub rsi, rax",
		"mov rdx, rax # Size",
		fmt.Sprintf("mov rax, %s # Write system call identifier", getSyscall(writeSyscall)),
		"mov rdi, 1 # Standard output file descriptor",
		"syscall",
	)
	prog.Scope = prog.Scope.Parent
}

func (node *CallFuncNode) writeUInt(prog *Prog) {
	prog.Scope = NewScope(prog.Scope)
	operandPos := node.genExpr(prog).stackPos
	bufferPos := prog.Scope.Alloc(16)
	// Algorithm in C: https://gcc.godbolt.org/z/3b2P5j
	if !libraryLabelExists(prog, "uintToAscii") {
		prog.LibrarySubSect.SubSects = append(prog.LibrarySubSect.SubSects, &Sect{
			Label: "uintToAscii",
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
	}
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("mov edi, dword ptr [rbp - %d] # Integer argument", operandPos),
		fmt.Sprintf("lea rsi, [rbp - %d] # Buffer pointer argument", bufferPos-16),
		"call _uintToAscii",
		// fmt.Sprintf("lea rsi, [rbp - %d]", bufferPos),
		"mov rdx, r9 # Size",
		fmt.Sprintf("mov rax, %s # Write system call identifier", getSyscall(writeSyscall)),
		"mov rdi, 1 # Standard output file descriptor",
		"syscall", // edi already set earlier
	)
	prog.Scope = prog.Scope.Parent
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
		fmt.Sprintf("mov rax, %s # Write system call identifier", getSyscall(readSyscall)),
		"mov rdi, 0 # Standard input file descriptor",
		"syscall",

		"mov esi, eax # Length of characters",
		fmt.Sprintf("lea rdi, [rbp - %d] # Pointer to ASCII buffer", bufferPos),
		"call _asciiToInt32",

		fmt.Sprintf("mov dword ptr [rbp - %d], eax", ref.stackPos),
	)
}
