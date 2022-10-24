package frontend

import (
	"testing"

	"govega/vega/language/tokens"
)

type testParserInterface interface {
	Parser
	getCurrentToken() *lexicalToken
	getNextToken() *lexicalToken
	readToken() error
	getToken() (*lexicalToken, error)
	lookAHead(tag int) bool
	matchToken(tag int) bool
}

type testParser struct {
	parser
}

func newTestParser(inputCode []byte) testParserInterface {
	v := &vega{
		file:      "/path/to/test.vg",
		codeLines: []string{},
	}
	l := v.NewLexer(inputCode)
	var parser testParserInterface = &testParser{
		parser{
			vega:  v,
			lexer: l,
		},
	}
	return parser
}

func (t *testParser) getCurrentToken() *lexicalToken {
	return t.currentToken
}

func (t *testParser) getNextToken() *lexicalToken {
	return t.nextToken
}

func TestParserObject_getLexerError(t *testing.T) {
	parser := newTestParser([]byte("'blubb\\H'"))
	err := parser.readToken()
	vErr := err.(IVError)

	if vErr.GetErrorType() != invalidEscapeSequence {
		t.Fatalf("Expected %v, got %v", invalidEscapeSequence, vErr.GetErrorType())
	}
}

func TestParserObject_getEOFToken(t *testing.T) {
	parser := newTestParser([]byte{})

	err := parser.readToken()

	if err != nil {
		t.Error(err)
	}

	gotCurrent := parser.getCurrentToken().GetTag()
	gotNext := parser.getNextToken().GetTag()
	want := tokens.EOF

	if gotCurrent != want && gotNext != want {
		t.Fatalf("Expected EOF token %v, but got current: %v and next: %v", want, gotCurrent, gotNext)
	}

}

func TestParserObject_lookAHead(t *testing.T) {
	parser := newTestParser([]byte("while if"))

	err := parser.readToken()

	if err != nil {
		t.Error(err)
	}

	boolean := parser.lookAHead(tokens.IF)

	if !boolean {
		t.Fatalf("Token is not if token, got %v", parser.getNextToken().GetToken().String())
	}

	want := tokens.IF

	if parser.getNextToken().GetTag() != want {
		t.Fatalf("Token is not %v token, got %v", want, parser.getNextToken().GetToken().String())
	}

	got := parser.getNextToken().GetTag()

	if want != got && parser.getNextToken() != nil {
		t.Fatalf("Expected token to be %v, but got %v", want, got)
	}

}
