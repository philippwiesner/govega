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
