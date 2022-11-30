package utils

import (
	"testing"

	types2 "govega/vega/types"
)

func TestNewScope(t *testing.T) {
	inMain := []*Symbol{
		NewSymbol("var1", types2.IntType, false, false),
		NewSymbol("var2", types2.CharType, false, true),
	}
	inSub := []*Symbol{
		NewSymbol("var3", types2.FloatType, false, false),
		NewSymbol("var4", types2.NewString(6), true, false),
	}

	table := NewSymbolTable()
	table.NewScope("main")

	if table.getScopeName() != "main" {
		t.Fatalf("Name of scope not main, got: %v", table.getScopeName())
	}

	for _, s := range inMain {
		table.Add(s)
	}

	var1, ok := table.Lookup("var1")
	if !ok {
		t.Fatalf("Element var1 not found")
	}

	if var1.SymbolType != types2.IntType || var1.Callable != false || var1.Const != false {
		t.Fatalf("var1 not as it should be, got: %v", var1)
	}

	table.NewScope("Sub")

	if table.getScopeName() != "Sub" {
		t.Fatalf("Name of scope not sub, got: %v", table.getScopeName())
	}

	for _, s := range inSub {
		table.Add(s)
	}

	var2, ok := table.Lookup("var2")
	if !ok {
		t.Fatalf("Element var2 not found")
	}

	if var2.SymbolType != types2.CharType || var2.Callable != false || var2.Const != true {
		t.Fatalf("var2 not as it should be, got: %v", var2)
	}

	var4, ok := table.Lookup("var4")
	if !ok {
		t.Fatalf("Element var4 not found")
	}

	if var4.SymbolType.(*types2.StringType) == types2.NewString(6) || var4.Callable != true || var4.Const != false {
		t.Fatalf("var4 not as it should be, got: %v", var4)
	}

	// leave sub scope, and try to lookup again var in main scope and see if element cannot be found in sub scope
	table.LeaveScope()

	if table.getScopeName() != "main" {
		t.Fatalf("Name of scope not main, got: %v", table.getScopeName())
	}

	var1, ok = table.Lookup("var1")
	if !ok {
		t.Fatalf("Element var1 not found")
	}

	if var1.SymbolType != types2.IntType || var1.Callable != false || var1.Const != false {
		t.Fatalf("var1 not as it should be, got: %v", var1)
	}

	var4, ok = table.Lookup("var4")
	if ok {
		t.Fatalf("Element var4 found, but shouldn't as scope left")
	}

}
