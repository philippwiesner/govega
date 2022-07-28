package frontend

import (
	"govega/govega/helper"
	"govega/govega/language"
)

type scope struct {
	*helper.HashTable
	name          string
	previousScope *scope
}

func newScope(name string) *scope {
	return &scope{helper.NewHashTable(), name, nil}
}

type SymbolTable struct {
	head *scope
	tail *scope
}

func NewSymbolTable() *SymbolTable {
	globalScope := newScope("global")
	return &SymbolTable{globalScope, globalScope}
}

func (st *SymbolTable) NewScope(name string) {
	newScope := newScope(name)
	old := st.head
	st.head = newScope
	newScope.previousScope = old
}

func (st *SymbolTable) LeaveScope() {
	if st.GetScopeName() != "global" {
		newHead := st.head.previousScope
		st.head = newHead
	}
}

func (st *SymbolTable) GetScopeName() string {
	return st.head.name
}

type Symbol struct {
	name       string
	SymbolType language.IBasicType
	Callable   bool
	Const      bool
}

func NewSymbol(name string, varType language.IBasicType, callable bool, con bool) *Symbol {
	return &Symbol{name, varType, callable, con}
}

func (st *SymbolTable) Add(s *Symbol) {
	st.head.Add(s.name, s)
}

func (st *SymbolTable) LookUp(name string) (entry *Symbol, ok bool) {
	currentScope := st.head
	for {
		result, ok := currentScope.Get(name)
		if !ok {
			if currentScope.previousScope == nil {
				return nil, false
			} else {
				currentScope = currentScope.previousScope
			}
		} else {
			return result.(*Symbol), true
		}
	}
}
