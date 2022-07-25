package frontend

import (
	"errors"
	"govega/govega/helper"
	"govega/govega/language"
)

type SymbolTable struct {
	*helper.HashTable
	name          string
	previousScope *SymbolTable
}

func (st *SymbolTable) GetName() string {
	return st.name
}

type Symbol struct {
	name       string
	SymbolType language.IBasicType
	Callable   bool
	Const      bool
}

func NewScope(name string, st *SymbolTable) *SymbolTable {
	return &SymbolTable{helper.NewHashTable(), name, st}
}

func (st *SymbolTable) LeaveScope() (symbolTable *SymbolTable, err error) {
	if st.previousScope == nil {
		return nil, errors.New("already last scope")
	}
	return st.previousScope, nil
}

func NewSymbolEntry(name string, varType language.IBasicType, callable bool, con bool) *Symbol {
	return &Symbol{name, varType, callable, con}
}

func (st *SymbolTable) NewEntry(entry *Symbol) {
	st.Add(entry.name, entry)
}

func (st *SymbolTable) LookUp(name string) (entry *Symbol, ok bool) {
	table := st
	for {
		result, ok := table.Get(name)
		if !ok {
			if table.previousScope == nil {
				return nil, false
			} else {
				table = table.previousScope
			}
		} else {
			return result.(*Symbol), true
		}
	}
}
