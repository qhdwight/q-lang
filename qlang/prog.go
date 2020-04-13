package main

type (
	Scope struct {
		vars       map[Node]int
		Head, Base int
		Parent     *Scope
	}
)

func (scope *Scope) AllocVar(node Node, size int) int {
	pos := scope.Head
	scope.vars[node] = pos
	scope.Head += size
	return scope.Base + pos
}

func NewScope(parent *Scope) *Scope {
	newScope := &Scope{
		vars:   make(map[Node]int),
		Head:   0,
		Parent: parent,
	}
	if parent == nil {
		newScope.Base = 0
	} else {
		newScope.Base = parent.Head
	}
	return newScope
}

func (scope *Scope) GetVarPos(node Node) int {
	_scope := scope
	for {
		if pos, has := _scope.vars[node]; has {
			return _scope.Base + pos
		} else {
			// Try again with parent scope
			_scope = scope.Parent
		}
		if _scope.vars == nil {
			panic("Scope does not contain requested variable")
		}
	}
}
