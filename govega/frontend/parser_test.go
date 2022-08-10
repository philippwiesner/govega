package frontend

import (
	"testing"
)

type mockBlockParser struct {
	parserObject
}

func (mock *mockBlockParser) parseFunctionDeclaration(p Parser) error {
	return nil
}

func (mock *mockBlockParser) parseFunctionReturnType(p Parser) error {
	return nil
}

func (mock *mockBlockParser) parseScope(p Parser) error {
	return nil
}

func TestParseBlock(t *testing.T) {
	code := "fun test()"
	lexer := NewLexer([]byte(code), "test")
	tokenStream, _ := lexer.Scan()

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
