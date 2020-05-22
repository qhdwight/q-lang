package main

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	strLabelNum  = 0
	loopLabelNum = 0
)

func getPropDef(datDef *DefDatNode, split []string) *DefDatPropNode {
	// TODO:warning recursion for depth greater than one
	for _, child := range datDef.children {
		propDefNode := child.(*DefDatPropNode)
		if propDefNode.name == split[1] {
			return propDefNode
		}
	}
	panic("Property not found")
}

func genOperandVar(prog *Prog, operand *OperandNode) ScopeVar {
	// Existing variable in scope
	if len(operand.accessor) > 0 {
		boundVar := getBoundVar(prog, operand.accessor)
		if len(operand.typeName) == 0 {
			operand.typeName = boundVar.typeName
		}
		return boundVar
	}
	// New int
	if operand.typeName == "i32" {
		pos := prog.Scope.Alloc(getSizeOfType("i32"))
		operandScopeVar := ScopeVar{typeName: "i32", stackPos: pos}
		if len(operand.literalVal) > 0 {
			i, err := strconv.Atoi(operand.literalVal)
			if err == nil {
				prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("mov dword ptr [rbp - %d], %d # Integer literal", pos, i))
			}
		}
		return operandScopeVar
	}
	// New data
	datDef := progNode.datDefs[operand.typeName]
	basePos := prog.Scope.Alloc(datDef.size)
	for _, child := range operand.children {
		fillProp := child.(*PropFill)
		propPos := fillProp.genExpr(prog)
		genStackCopy(prog, propPos, ScopeVar{typeName: fillProp.propDef.typeName, stackPos: basePos - fillProp.propDef.offset})
	}
	return ScopeVar{typeName: operand.typeName, stackPos: basePos}
}

func getBoundVar(prog *Prog, accessor string) ScopeVar {
	for _, scopeVar := range prog.Scope.vars {
		if strings.HasPrefix(accessor, scopeVar.varName) {
			split := strings.Split(accessor, ".")
			if len(split) == 1 {
				return scopeVar
			} else {
				datDef := progNode.datDefs[scopeVar.typeName]
				propDef := getPropDef(datDef, split)
				return ScopeVar{typeName: propDef.typeName, stackPos: scopeVar.stackPos - propDef.offset}
			}
		}
	}
	panic(fmt.Sprintf("Variable %s not found in scope!", accessor))
}

func getSizeOfType(name string) int {
	if name == "i32" {
		return 4
	} else {
		return progNode.datDefs[name].size
	}
}

func genStackCopy(prog *Prog, src, dst ScopeVar) {
	if src.typeName != dst.typeName {
		panic("Moving different types!")
	}
	size := getSizeOfType(dst.typeName)
	prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("# Copy %+v to %+v", src, dst))
	// TODO:performance use memory copy equivalent
	for i := 0; i < size; i++ {
		prog.CurSect.Content = append(prog.CurSect.Content,
			fmt.Sprintf("mov dl, byte ptr [rbp - %d]", src.stackPos-i),
			fmt.Sprintf("mov byte ptr [rbp - %d], dl", dst.stackPos-i),
		)
	}
}

func (node *BaseNode) genExpr(prog *Prog) ScopeVar {
	child0 := node.children[0].(*OperandNode)
	var0 := genOperandVar(prog, child0)
	exprResVar := genOperandVar(prog, &OperandNode{typeName: child0.typeName}) // Create clone
	genStackCopy(prog, var0, exprResVar)
	prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("# Expression base %+v", exprResVar))
	if len(node.children) == 1 {
		return exprResVar
	}
	// Multiple operands in expression with operators in between
	prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("mov eax, dword ptr [rbp - %d]", exprResVar.stackPos))
	for nodeIndex, child := range node.children {
		switch child.(type) {
		case *AdditionNode:
			rhsOperand := node.children[nodeIndex+1].(*OperandNode)
			rhsVar := genOperandVar(prog, rhsOperand)
			prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("add eax, dword ptr [rbp - %d]", rhsVar.stackPos))
			break
		}
	}
	prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("mov dword ptr [rbp - %d], eax", exprResVar.stackPos))
	return exprResVar
}

func (node *AssignmentNode) Generate(prog *Prog) {
	dstVar := getBoundVar(prog, node.accessor)
	node.genAssignmentPhrase(prog, dstVar)
}

func (node *SingleNamedVarNode) Generate(prog *Prog) {
	namedVar := prog.Scope.BindNamedVar(node)
	node.genAssignmentPhrase(prog, namedVar)
}

func (node *BaseNode) genAssignmentPhrase(prog *Prog, dstVar ScopeVar) {
	prog.Scope = NewScope(prog.Scope)
	exprResPos := node.genExpr(prog)
	genStackCopy(prog, exprResPos, dstVar)
	prog.Scope = prog.Scope.Parent
}

func (node *DefDatNode) Generate(*Prog) {
	offset := 0
	for _, child := range node.children {
		propNode := child.(*DefDatPropNode)
		propSize := getSizeOfType(propNode.typeName)
		propNode.offset = offset
		offset += propSize
	}
	node.size = offset
}

func (node *BaseNode) Generate(prog *Prog) {
	for _, child := range node.children {
		child.Generate(prog)
	}
}

func (node *ImplFuncNode) Generate(prog *Prog) {
	// TODO:warning detect stack size properly instead of subtracting constant
	prog.CurSect.Content = append(prog.CurSect.Content,
		"push rbp",
		"mov rbp, rsp",
		"sub rsp, 64",
		"",
	)
	for _, child := range node.children {
		child.Generate(prog)
	}
	prog.CurSect.Content = append(prog.CurSect.Content,
		"",
		"add rsp, 64",
		"pop rbp",
		"ret",
	)
}

func (node *OutNode) Generate(prog *Prog) {
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("mov eax, %d", 0),
	)
}

func (node *StringLiteralNode) Generate(prog *Prog) {
	strLabelNum++
	asmLabel := fmt.Sprintf("string%d", strLabelNum)
	node.label = asmLabel
	msgSubSect := &Sect{
		Label: asmLabel, Content: []string{fmt.Sprintf(`.string "%s\n"`, node.str)},
	}
	prog.ConstSect.SubSects = append(prog.ConstSect.SubSects, msgSubSect)
}

func (node *LoopNode) Generate(prog *Prog) {
	loopLabelNum++
	prog.Scope = NewScope(prog.Scope)
	counterPos := prog.Scope.Alloc(4)
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("mov dword ptr [rbp - %d], %d # Counter", counterPos, node.start),
	)
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("_loopCheck%d:", loopLabelNum),
		fmt.Sprintf("cmp dword ptr [rbp - %d], %d", counterPos, node.end),
		fmt.Sprintf("jge _loopContinue%d", loopLabelNum),
		fmt.Sprintf("jmp _loopBody%d", loopLabelNum),
	)
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("_loopBody%d:", loopLabelNum),
	)
	for _, child := range node.children {
		child.Generate(prog)
	}
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("mov eax, dword ptr [rbp - %d]", counterPos),
		"inc eax",
		fmt.Sprintf("mov dword ptr [rbp - %d], eax", counterPos),
		fmt.Sprintf("jmp _loopCheck%d", loopLabelNum),
	)
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("_loopContinue%d:", loopLabelNum),
	)
	prog.Scope = prog.Scope.Parent
}

func (node *CallFuncNode) Generate(prog *Prog) {
	if node.name == "pln" {
		if strNode, isStr := node.children[0].(*StringLiteralNode); isStr {
			strNode.Generate(prog)
			prog.CurSect.Content = append(prog.CurSect.Content,
				fmt.Sprintf("lea rax, [rip + _%s]", strNode.label),
				"mov rsi, rax # Pointer to string",
				fmt.Sprintf("mov rdx, %d # Size", len(strNode.str)+1),
				"mov rax, 1 # Write",
				"mov rdi, 1 # Standard output",
				"syscall",
			)
		} else {
			prog.Scope = NewScope(prog.Scope)
			operandPos := node.genExpr(prog).stackPos
			bufferPos := prog.Scope.Alloc(16)
			// Algorithm in C: https://gcc.godbolt.org/z/3b2P5j
			needsLibrary := true
			for _, subSect := range prog.LibrarySubSect.SubSects {
				if subSect.Label == "charLoop" {
					needsLibrary = false
					break
				}
			}
			if needsLibrary {
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
				prog.Scope = prog.Scope.Parent
			}
			prog.CurSect.Content = append(prog.CurSect.Content,
				fmt.Sprintf("mov edi, dword ptr [rbp - %d] # Integer argument", operandPos),
				fmt.Sprintf("lea rsi, [rbp - %d] # Buffer pointer argument", bufferPos-16),
				"call _digitToChar",
				// fmt.Sprintf("lea rsi, [rbp - %d]", bufferPos),
				"mov rdx, r9 # Size",
				"mov rax, 1 # Write",
				"mov rdi, 1 # Standard output",
				"syscall", // edi already set earlier
			)
		}
	}
}

func (node *ProgNode) Generate(prog *Prog) {
	node.BaseNode.Generate(prog)
}
