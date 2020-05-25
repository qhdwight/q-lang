package main

type ScopeVar struct {
	typeName, varName string
	stackPos          int // Relative to base pointer
}

type (
	Scope struct {
		vars       map[Node]ScopeVar
		Head, Base int
		Parent     *Scope
	}
)

func (scope *Scope) Alloc(size int) int {
	scope.Head += size
	return scope.Head
}

func (scope *Scope) BindNamedVar(node *SingleNamedVarNode) ScopeVar {
	size := getSizeOfType(node.typeName)
	pos := scope.Alloc(size)
	scopeVar := ScopeVar{typeName: node.typeName, varName: node.name, stackPos: pos}
	scope.vars[node] = scopeVar
	return scopeVar
}

func (scope *Scope) BindUnnamed(node Node, scopeVar ScopeVar) ScopeVar {
	scope.vars[node] = scopeVar
	return scopeVar
}
func NewScope(parent *Scope) *Scope {
	newScope := &Scope{
		vars:   make(map[Node]ScopeVar),
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
