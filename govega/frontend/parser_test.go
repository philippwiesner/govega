package frontend

import (
	"govega/govega/language"
	"govega/govega/language/tokens"
	"testing"
)

type mockBlockParser struct {
	parserObject
}

func (mock *mockBlockParser) parseFunctionDeclaration(p Parser) error {
	return nil
}

func (mock *mockBlockParser) parseFunctionReturnType(p Parser, s *Symbol) error {
	return nil
}

func (mock *mockBlockParser) parseScope(p Parser, name string) error {
	return nil
}

func TestParseBlock(t *testing.T) {
	tokenStream := NewTokenStream()

	tokenStream.Add(tokens.NewWord("func", tokens.FUNC), ErrorState{})
	tokenStream.Add(tokens.NewWord("test", tokens.ID), ErrorState{})
	tokenStream.Add(tokens.NewToken('('), ErrorState{})
	tokenStream.Add(tokens.NewToken(')'), ErrorState{})
	tokenStream.Add(language.ReturnType, ErrorState{})
	tokenStream.Add(tokens.NewWord("func", tokens.FUNC), ErrorState{})
	tokenStream.Add(tokens.NewWord("bla", tokens.ID), ErrorState{})
	tokenStream.Add(tokens.NewToken('('), ErrorState{})
	tokenStream.Add(tokens.NewToken(')'), ErrorState{})
	tokenStream.Add(language.ReturnType, ErrorState{})
	tokenStream.Add(tokens.NewToken(tokens.EOF), ErrorState{})

	var mockParser Parser = &mockBlockParser{
		parserObject{
			ErrorState{0, 0, "test", ""},
			tokenStream,
			nil,
			NewSymbolTable(),
		},
	}

	err := mockParser.parse(mockParser)
	if err != nil {
		t.Error(err)
	}

}
