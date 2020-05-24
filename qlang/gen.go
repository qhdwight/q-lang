package main

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	strLabelNum  = 0
	loopLabelNum = 0
	ifLabelNum   = 0
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
	// Function call
	if len(operand.children) == 1 {
		if funcCall, isFuncCall := operand.children[0].(*CallFuncNode); isFuncCall {
			funcCall.Generate(prog)
			return funcCall.retVar
		}
	}
	// Existing variable in scope
	if len(operand.accessor) > 0 {
		boundVar := getBoundVar(prog, operand.accessor)
		if len(operand.typeName) == 0 {
			operand.typeName = boundVar.typeName
		}
		return boundVar
	}
	// New int
	if operand.typeName == uintKeyword {
		pos := prog.Scope.Alloc(getSizeOfType(uintKeyword))
		operandScopeVar := ScopeVar{typeName: uintKeyword, stackPos: pos}
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
		// New data
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
	if name == uintKeyword {
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
	// for i := 0; i < size; i++ {
	// 	prog.CurSect.Content = append(prog.CurSect.Content,
	// 		fmt.Sprintf("mov dl, byte ptr [rbp - %d]", src.stackPos-i),
	// 		fmt.Sprintf("mov byte ptr [rbp - %d], dl", dst.stackPos-i),
	// 	)
	// }
	for i := 0; i < size; i += 4 {
		prog.CurSect.Content = append(prog.CurSect.Content,
			fmt.Sprintf("mov edx, dword ptr [rbp - %d]", src.stackPos-i),
			fmt.Sprintf("mov dword ptr [rbp - %d], edx", dst.stackPos-i),
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
		genOperation := func(opName string) {
			rhsOperand := node.children[nodeIndex+1].(*OperandNode)
			rhsVar := genOperandVar(prog, rhsOperand)
			prog.CurSect.Content = append(prog.CurSect.Content, fmt.Sprintf("%s eax, dword ptr [rbp - %d]", opName, rhsVar.stackPos))
		}
		switch child.(type) {
		case *AddNode:
			genOperation("add")
			break
		case *MulNode:
			genOperation("imul")
			break
		case *SubtractionNode:
			genOperation("sub")
			break
		case *DivisionNode:
			panic("Division not supported yet!")
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
	if len(node.children) > 0 {
		node.genAssignmentPhrase(prog, namedVar)
	}
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

func (node *StringLiteralNode) Generate(prog *Prog) {
	strLabelNum++
	asmLabel := fmt.Sprintf("string%d", strLabelNum)
	node.label = asmLabel
	msgSubSect := &Sect{
		Label: asmLabel, Content: []string{fmt.Sprintf(`.string "%s\n"`, node.str)},
	}
	prog.ConstSect.SubSects = append(prog.ConstSect.SubSects, msgSubSect)
}

func libraryLabelExists(prog *Prog, labelName string) bool {
	has := false
	for _, subSect := range prog.LibrarySubSect.SubSects {
		if subSect.Label == labelName {
			has = true
			break
		}
	}
	return has
}

func (node *ProgNode) Generate(prog *Prog) {
	node.BaseNode.Generate(prog)
}
