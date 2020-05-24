package main

import "fmt"

func (node *CallFuncNode) Generate(prog *Prog) {
	if node.name == "wln" {
		node.writeLine(prog)
	} else if node.name == "rln" {
		node.readLine(prog)
	} else {
		funcImpl := progNode.funcImpls[node.name]
		funcDef := funcImpl.Parent.(*DefFuncNode)
		// Allocate space for return variable in caller function stack
		// Place the start pointer of return variable into rdi so function can fill
		retStackPos := prog.Scope.Alloc(getSizeOfType(funcDef.retType))
		node.retVar = ScopeVar{typeName: funcDef.retType, stackPos: retStackPos}
		prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("lea rdi, dword ptr [rbp - %d] # Return stack pointer", retStackPos))
		for paraIdx, paraType := range funcDef.parameterTypes {
			paraVar := genOperandVar(prog, node.children[paraIdx].(*OperandNode))
			paraSize := getSizeOfType(paraType)
			// Copy parameters from source frame to destination frame
			for i := 0; i < paraSize; i += 4 {
				prog.CurSect.Content = append(prog.CurSect.Content,
					fmt.Sprintf("mov eax, dword ptr [rbp - %d]", paraVar.stackPos+i),
					fmt.Sprintf("mov dword ptr [rsp - %d], eax # Copy parameter %s to called function stack frame",
						paraIdx*paraSize+i+16+4+8, funcImpl.parameterNames[paraIdx]), // 16 bytes padding for call operation, 4 for dword size, 8 for return stack position pointer (rdi)
				)
			}
		}
		prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("call _%s", node.name))
	}
}

func (node *OutNode) Generate(prog *Prog) {
	exprRes := node.genExpr(prog)
	funcImpl := node.Parent.(*ImplFuncNode)
	if funcImpl.name == "main" {
		prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("mov eax, dword ptr [rbp - %d] # Program exit code", exprRes.stackPos))
	} else {
		prog.CurSect.Content = append(prog.CurSect.Content, "mov rdi, qword ptr [rbp - 8] # Caller return variable stack position")
		funcDef := funcImpl.Parent.(*DefFuncNode)
		retSize := getSizeOfType(funcDef.retType)
		for i := 0; i < retSize; i += 4 {
			prog.CurSect.Content = append(prog.CurSect.Content,
				fmt.Sprintf("mov eax, dword ptr [rbp - %d]", exprRes.stackPos-i),
				fmt.Sprintf("mov dword ptr [rdi + %d], eax", i), // rdi points to base of return variable
			)
		}
	}
}

func (node *ImplFuncNode) Generate(prog *Prog) {
	// TODO:warning detect stack size properly instead of subtracting constant
	prog.Scope = NewScope(prog.Scope)
	if node.name != "main" { // Main function does not use our protocol of placing return variable pointer in rdi
		prog.Scope.Alloc(8)
	}
	prog.CurSect = &Sect{Label: node.name}
	funcDef := node.Parent.(*DefFuncNode)
	for i := 0; i < len(funcDef.parameterTypes); i++ {
		parameterVar := &SingleNamedVarNode{typeName: funcDef.parameterTypes[i], name: node.parameterNames[i]}
		prog.Scope.BindNamedVar(parameterVar)
	}
	prog.CurSect.Content = append(prog.CurSect.Content,
		"push rbp",
		"mov rbp, rsp",
		"sub rsp, 1024",
	)
	if node.name != "main" { // Save pointer to return variable, since system calls may need rdi
		prog.CurSect.Content = append(prog.CurSect.Content, "mov qword ptr [rbp - 8], rdi")
	}
	prog.CurSect.Content = append(prog.CurSect.Content, "")
	for _, child := range node.children {
		child.Generate(prog)
	}
	prog.CurSect.Content = append(prog.CurSect.Content,
		"",
		"add rsp, 1024",
		"pop rbp",
		"ret",
	)
	prog.FuncSect.SubSects = append(prog.FuncSect.SubSects, prog.CurSect)

	prog.Scope = prog.Scope.Parent
}
