package language

import (
	"govega/govega/language/tokens"
	"testing"
)

func TestKeyWords(t *testing.T) {
	result, ok := KeyWords.Get("int")
	if !ok {
		t.Fatalf("element int not found")
	}
	integer := result.(*BasicType)

	if integer.GetTag() != tokens.BASIC {
		t.Fatalf("Word integer is not of type BASIC, got: %v", integer.GetTag())
	}

	if integer.GetWidth() != 4 {
		t.Fatalf("Word integer has not width of 4, got: %v", integer.GetWidth())
	}

	result, ok = KeyWords.Get("continue")
	if !ok {
		t.Fatalf("element continue not found")
	}

	word := result.(tokens.IWord)

	if word.GetTag() != tokens.CONTINUE {
		t.Fatalf("continue keyWord is not of type CONTINUE, got: %v", word.GetTag())
	}
	if word.GetLexeme() != "continue" {
		t.Fatalf("continue keyWord has not lexeme continue, got: %v", word.GetLexeme())
	}

}

func TestHexaLiterals(t *testing.T) {
	tests := []struct {
		in   string
		want rune
	}{
		{"00", rune(00)},
		{"0a", rune(10)},
		{"ff", rune(255)},
	}

	if EscapeHexaLiterals.BucketCount != 256 {
		t.Fatalf("Table should have 256 entries")
	}

	for i, tc := range tests {
		hex, ok := EscapeHexaLiterals.Get(tc.in)
		if !ok {
			t.Fatalf("test%d: Element %v not found", i, tc.in)
		}
		hexRune := hex.(rune)
		if hexRune != tc.want {
			t.Fatalf("test%d: Want %v, but got %v", i, tc.want, hexRune)
		}
	}
}

func TestOctaLiterals(t *testing.T) {
	tests := []struct {
		in   string
		want rune
	}{
		{"000", rune(000)},
		{"123", rune(83)},
		{"377", rune(255)},
	}

	if EscapeOctalLiterals.BucketCount != 256 {
		t.Fatalf("Table should have 256 entries")
	}

	for i, tc := range tests {
		oct, ok := EscapeOctalLiterals.Get(tc.in)
		if !ok {
			t.Fatalf("test%d: Element %v not found", i, tc.in)
		}
		octRune := oct.(rune)
		if octRune != tc.want {
			t.Fatalf("test%d: Want %v, but got %v", i, tc.want, octRune)
		}
	}
}
