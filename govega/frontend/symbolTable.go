package frontend

import (
	"govega/govega/helper"
	"govega/govega/language"
)

// scope defines a variable lookup scope
type scope struct {
	*helper.HashTable        // HashTable store all defined variables in the current scope
	name              string // give the scope a name
	previousScope     *scope // link to the previous scope for outer scope lookups
}

// newScope is the internal method for creating a new scope
func newScope(name string) *scope {
	return &scope{helper.NewHashTable(), name, nil}
}

// SymbolTable defines a lookup table for symbols. The table can store multiple tables linkes as a stack to lookup data
// from previous tables and remove unneeded tables when the scope is left.
type SymbolTable struct {
	head *scope
	tail *scope
}

// NewSymbolTable creates a new SymbolTable and adds a global scope
func NewSymbolTable() *SymbolTable {
	globalScope := newScope("global")
	return &SymbolTable{globalScope, globalScope}
}

// NewScope adds a new scope on top of the SymbolTables last scope
func (st *SymbolTable) NewScope(name string) {
	newScope := newScope(name)
	old := st.head
	st.head = newScope
	newScope.previousScope = old
}

// LeaveScope removed the current scope from the table
func (st *SymbolTable) LeaveScope() {
	if st.getScopeName() != "global" {
		newHead := st.head.previousScope
		st.head = newHead
	}
}

// getScopeName returnes the name of the current scope (mainly for debugging and testing)
func (st *SymbolTable) getScopeName() string {
	return st.head.name
}

// Symbol is stored in the symbol table
type Symbol struct {
	name       string              // symbol (identifier) name to be looked up with
	SymbolType language.IBasicType // identifier data tybe
	Callable   bool                // flag if identifier is callable (function declaration)
	Const      bool                // flag if identifier is a constant
}

// NewSymbol creates a new Symbol
func NewSymbol(name string, varType language.IBasicType, callable bool, con bool) *Symbol {
	return &Symbol{name, varType, callable, con}
}

// Add adds a new symbol to the current scope of the SymbolTable
func (st *SymbolTable) Add(s *Symbol) {
	st.head.Add(s.name, s)
}

// Lookup searches if the given symbol can already be found in the current or any previous scopes. Returns also a
// boolean if a symbol is being found or not.
func (st *SymbolTable) Lookup(name string) (entry *Symbol, ok bool) {
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
