package frontend

import (
	"govega/govega/language"
	"testing"
)

func TestNewScope(t *testing.T) {
	inMain := []*Symbol{
		NewSymbolEntry("var1", language.IntType, false, false),
		NewSymbolEntry("var2", language.CharType, false, true),
	}
	inSub := []*Symbol{
		NewSymbolEntry("var3", language.FloatType, false, false),
		NewSymbolEntry("var4", language.NewString(6), true, false),
	}

	table := NewScope("main", nil)

	if table.GetName() != "main" {
		t.Fatalf("Name of scope not main, got: %v", table.GetName())
	}

	for _, s := range inMain {
		table.NewEntry(s)
	}

	var1, ok := table.LookUp("var1")
	if !ok {
		t.Fatalf("Element var1 not found")
	}

	if var1.SymbolType != language.IntType || var1.Callable != false || var1.Const != false {
		t.Fatalf("var1 not as it should be, got: %v", var1)
	}

	table = NewScope("Sub", table)

	if table.GetName() != "Sub" {
		t.Fatalf("Name of scope not sub, got: %v", table.GetName())
	}

	for _, s := range inSub {
		table.NewEntry(s)
	}

	var2, ok := table.LookUp("var2")
	if !ok {
		t.Fatalf("Element var2 not found")
	}

	if var2.SymbolType != language.CharType || var2.Callable != false || var2.Const != true {
		t.Fatalf("var2 not as it should be, got: %v", var2)
	}

	var4, ok := table.LookUp("var4")
	if !ok {
		t.Fatalf("Element var4 not found")
	}

	stringType := var4.SymbolType.(language.IStringType)

	if stringType.GetLexeme() != "[]" || stringType.GetType() != language.CharType || var4.Callable != true || var4.Const != false {
		t.Fatalf("var4 not as it should be, got: %v", var4)
	}

	// leave sub scope, and try to lookup again var in main scope and see if element cannot be found in sub scope
	table, err := table.LeaveScope()
	if err != nil {
		t.Errorf("TestNewScope: %v", err)
	}

	if table.GetName() != "main" {
		t.Fatalf("Name of scope not main, got: %v", table.GetName())
	}

	var1, ok = table.LookUp("var1")
	if !ok {
		t.Fatalf("Element var1 not found")
	}

	if var1.SymbolType != language.IntType || var1.Callable != false || var1.Const != false {
		t.Fatalf("var1 not as it should be, got: %v", var1)
	}

	var4, ok = table.LookUp("var4")
	if ok {
		t.Fatalf("Element var4 found, but shouldn't as scope left")
	}

}
