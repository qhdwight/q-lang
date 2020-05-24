package main

import "fmt"

// TODO:refactor use sub-sections more instead of inlining labels in content

func (node *LoopNode) Generate(prog *Prog) {
	startVar, endVar := genOperandVar(prog, node.start), genOperandVar(prog, node.end)
	loopLabelNum++
	prog.Scope = NewScope(prog.Scope)
	counterPos := prog.Scope.Alloc(4)
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("mov eax, dword ptr [rbp - %d]", startVar.stackPos),
		fmt.Sprintf("mov dword ptr [rbp - %d], eax # Counter", counterPos),
	)
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("_loopCheck%d:", loopLabelNum),
		fmt.Sprintf("mov eax, dword ptr [rbp - %d]", counterPos),
		fmt.Sprintf("cmp eax, dword ptr [rbp - %d]", endVar.stackPos),
		fmt.Sprintf("jge _loopContinue%d", loopLabelNum),
		fmt.Sprintf("jmp _loopBody%d", loopLabelNum),
	)
	prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("_loopBody%d:", loopLabelNum))
	for _, child := range node.children {
		child.Generate(prog)
	}
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("mov eax, dword ptr [rbp - %d]", counterPos),
		"inc eax",
		fmt.Sprintf("mov dword ptr [rbp - %d], eax", counterPos),
		fmt.Sprintf("jmp _loopCheck%d", loopLabelNum),
	)
	prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("_loopContinue%d:", loopLabelNum))
	prog.Scope = prog.Scope.Parent
}

func (node *IfNode) Generate(prog *Prog) {
	v1, v2 := genOperandVar(prog, node.o1), genOperandVar(prog, node.o2)
	falseLabelNum := ifLabelNum
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("mov eax, dword ptr [rbp - %d]", v1.stackPos),
		fmt.Sprintf("cmp eax, dword ptr [rbp - %d]", v2.stackPos),
		fmt.Sprintf("jne _iff%d", falseLabelNum),
	)
	prog.Scope = NewScope(prog.Scope)
	for _, child := range node.t.children {
		child.Generate(prog)
	}
	prog.Scope = prog.Scope.Parent
	ifLabelNum++
	escapeLabelNum := ifLabelNum
	prog.CurSect.Content = append(prog.CurSect.Content,
		fmt.Sprintf("jmp _iff%d", escapeLabelNum),
		fmt.Sprintf("_iff%d:", falseLabelNum),
	)
	prog.Scope = NewScope(prog.Scope)
	for _, child := range node.f.children {
		child.Generate(prog)
	}
	prog.Scope = prog.Scope.Parent
	prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("_iff%d:", escapeLabelNum))
}
