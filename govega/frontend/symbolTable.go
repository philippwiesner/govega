package frontend

import (
	"errors"
	"govega/govega/helper"
	"govega/govega/language"
)

type SymbolTable struct {
	*helper.HashTable
	previousScope *SymbolTable
}

type Symbol struct {
	name     string
	BaseType *language.BasicType
	Callable bool
	Const    bool
}

func NewTable() *SymbolTable {
	return &SymbolTable{helper.NewHashTable(), nil}
}

func (st *SymbolTable) NewScope() *SymbolTable {
	return &SymbolTable{helper.NewHashTable(), st}
}

func (st *SymbolTable) LeaveScope() (symbolTable *SymbolTable, err error) {
	if st.previousScope == nil {
		return nil, errors.New("already last scope")
	}
	return st.previousScope, nil
}

func NewSymbolEntry(name string, baseType *language.BasicType, callable bool, con bool) *Symbol {
	return &Symbol{name, baseType, callable, con}
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
