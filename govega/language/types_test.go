package language

import (
	"reflect"
	"testing"
)

func TestNewArray(t *testing.T) {
	a1 := NewArray(IntType, 4)
	if a1.GetLexeme() != "[]" {
		t.Fatalf("Want Array1 lexeme to be [], got: %v", a1.GetLexeme())
	}

	if a1.GetWidth() != 16 {
		t.Fatalf("Want Array1 width 16, got: %v", a1.GetWidth())
	}

	if a1.GetType() != IntType {
		t.Fatalf("Want Array1 type INT, got: %v", a1.GetType())
	}

	a2 := NewArray(a1, 5)

	if a2.GetLexeme() != "[]" {
		t.Fatalf("Want Array2 lexeme to be [], got: %v", a2.GetLexeme())
	}

	if a2.GetWidth() != 80 {
		t.Fatalf("Want Array2 width 16, got: %v", a2.GetWidth())
	}

	if a2.GetType() != IntType {
		t.Fatalf("Want Array2 type INT, got: %v", a2.GetType())
	}

	if !reflect.DeepEqual(a2.GetDimensions(), []int{4, 5}) {
		t.Fatalf("Want Array2 dimensions to be {4, 5}, got: %v", a2.GetDimensions())
	}

	a3 := NewArray(a2, 2)

	if a3.GetWidth() != 160 {
		t.Fatalf("Want Array3 width 160, got: %v", a3.GetWidth())
	}

	if !reflect.DeepEqual(a3.GetDimensions(), []int{4, 5, 2}) {
		t.Fatalf("Want Array3 dimensions to be {4, 5, 2}, got: %v", a2.GetDimensions())
	}
}

func TestNewString(t *testing.T) {
	s1 := NewString(10)
	if s1.GetLexeme() != "[]" {
		t.Fatalf("Want String1 lexeme to be [], got: %v", s1.GetLexeme())
	}

	if s1.GetWidth() != 80 {
		t.Fatalf("Want String1 width 16, got: %v", s1.GetWidth())
	}

	if s1.GetType() != CharType {
		t.Fatalf("Want String1 type char, got: %v", s1.GetType())
	}
}
