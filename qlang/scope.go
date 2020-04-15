package main

import "fmt"

var anonNum = 0

type (
	Scope struct {
		vars       map[string]int
		Head, Base int
		Parent     *Scope
	}
)

func (scope *Scope) AllocVar(node Node, size int) int {
	scope.Head += size
	pos := scope.Base + scope.Head
	scope.vars[NodeName(node)] = pos
	return pos
}

func NodeName(node Node) string {
	switch n := node.(type) {
	case *DefSingleVarNode:
		return n.name
	case *OperandNode:
		if len(n.varName) == 0 {
			// Give the operand an anonymous name if it hasn't been set yet
			n.varName = fmt.Sprintf("__anon%d", anonNum)
			anonNum++
		}
		return n.varName
	default:
		panic("Can't get name for this type")
	}
}

func NewScope(parent *Scope) *Scope {
	newScope := &Scope{
		vars:   make(map[string]int),
		Head:   0,
		Parent: parent,
	}
	if parent == nil {
		newScope.Base = 0
	} else {
		newScope.Base = parent.Head
		newScope.Head = parent.Head
		for k, v := range parent.vars {
			newScope.vars[k] = v
		}
	}
	return newScope
}

func (scope *Scope) GetVarPos(node Node) int {
	if pos, has := scope.vars[NodeName(node)]; has {
		return pos
	} else {
		panic("Scope does not contain requested anonymous variable")
	}
}
