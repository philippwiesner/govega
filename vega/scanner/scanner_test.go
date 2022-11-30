package scanner

import (
	"fmt"
	"testing"

	"govega/vega/token"
)

func newTestFile(source []byte) token.File {
	return token.File{
		Name:   "/path/to/test.vg",
		Source: source,
	}
}

type LiteralTestWant struct {
	token.Token
	literal string
}

type LiteralTest struct {
	in   string
	want LiteralTestWant
}

type LiteralFailureTest struct {
	name string
	in   string
	want VErrorType
}

func TestNewLexer(t *testing.T) {
	text := "1+1\n."
	want := []rune{'1', '+', '1', '\n', '.'}
	counter := 0
	scanner := NewScanner(newTestFile([]byte(text)))
	var err error

	err = scanner.readch()
	for ; ; err = scanner.readch() {
		if err != nil {
			t.Error(err)
		} else if scanner.peek == 0 {
			break
		}
		if scanner.peek != want[counter] {
			t.Fatalf("Char %v on position %v is not %v", scanner.peek, counter, want[counter])
		}
		counter++
	}
}

func TestReadcch(t *testing.T) {
	text := "1+1"
	scanner := NewScanner(newTestFile([]byte(text)))

	err := scanner.readch()
	if err != nil {
		t.Fatal(err)
	}

	ok, err := scanner.readcch('+')
	if !ok {
		t.Fatalf("Next character not '+'")
	}

}

func TestCombinedTokens(t *testing.T) {
	tests := []struct {
		in   string
		want token.Token
	}{
		{"!=", token.NE},
		{"!-", token.LNOT},
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		scanner := NewScanner(newTestFile([]byte(tc.in)))
		err := scanner.readch()
		if err != nil {
			t.Fatal(err)
		}

		tok, err := scanner.scanCombinedTokens('=', token.NE)
		if err != nil {
			t.Fatal(err)
		}

		if tok.Token != tc.want {
			t.Fatalf("%v, token should be %v, but is %v", test, tc.want, tok.Token)
		}

	}
}

func TestScanLiterals(t *testing.T) {
	var tok *Token
	tests := []LiteralTest{
		{
			"\"my literal\"",
			LiteralTestWant{token.STRING, "\"my literal\""},
		},
		{
			"\"\\tmy \\nliteral\"",
			LiteralTestWant{token.STRING, "\"\tmy \nliteral\""},
		},
		{
			"\"my literal\"",
			LiteralTestWant{token.STRING, "\"my literal\""},
		},
		{
			"'\\x3A'",
			LiteralTestWant{token.CHAR, "'\x3a'"},
		},
		{
			"'A'",
			LiteralTestWant{token.CHAR, "'A'"},
		},
		{
			"\"my \\o123 \\\\\"",
			LiteralTestWant{token.STRING, "\"my \123 \\\""},
		},
	}
	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		scanner := NewScanner(newTestFile([]byte(tc.in)))
		err := scanner.readch()
		if err != nil {
			t.Fatalf("%v: error reading first literal indicator: %v", test, err)
		}

		switch scanner.peek {
		case '\'':
			tok, err = scanner.scanChar()
		// read string
		case '"':
			tok, err = scanner.scanString()
		}
		if err != nil {
			t.Fatalf("%v: error reading literal: %v", test, err)
		}

		if tok.Token != tc.want.Token {
			t.Fatalf("%v: Want token to be %v, but got: %v", test, tc.want.Token, tok.Token)
		} else {
			if tok.Literal != tc.want.literal {
				t.Fatalf("%v: Want literal to be %v, but got: %v", test, tc.want.literal, tok.Literal)
			}
		}

	}
}

func TestScanLiteralsFailures(t *testing.T) {
	tests := []LiteralFailureTest{
		{
			"linebreak in literal",
			"\"fooBar\n\"",
			literalNotTerminated,
		},
		{
			"invalid escape sequence",
			"\"fooBar\\-\"",
			invalidEscapeSequence,
		},
		{
			"single colon literal escape in double colon literal",
			"\"foo\\'Bar\"",
			invalidEscapeSequence,
		},
		{
			"invalid hex escape sequence",
			"\"foo\\xAg\"",
			invalidEscapeSequenceHexadecimal,
		},
		{
			"invalid oct escape sequence",
			"\"foo\\o088Bar\"",
			invalidEscapeSequenceOctal,
		},
		{
			"no literal terminator",
			"\"fooBar",
			literalNotTerminated,
		},
		{
			"more than one sign in character literal",
			"'ab'",
			literalNotTerminated,
		},
	}
	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		scanner := NewScanner(newTestFile([]byte(tc.in)))
		err := scanner.readch()
		if err != nil {
			t.Fatalf("%v %v:\nerror reading first literal indicator: %v", test, tc.name, err)
		}

		switch scanner.peek {
		case '\'':
			_, err = scanner.scanChar()
		// read string
		case '"':
			_, err = scanner.scanString()
		}

		switch err.(type) {
		case Verror:
			gotErrType := err.(Verror).GetErrorType()
			if gotErrType != tc.want {
				t.Fatalf("%v %v: Expected vega Error type to be %v, but got %v", test, tc.name, tc.want, gotErrType)
			}
		default:
			t.Fatalf("%v %v: Excpected Verror, but got %v", test, tc.name, err)
		}

	}
}

func TestScanNumbers(t *testing.T) {
	tests := []struct {
		in   string
		want Token
	}{
		{"123", Token{
			token.INT,
			token.Location{},
			"123",
		}}, {"12,3", Token{
			token.INT,
			token.Location{},
			"12",
		}}, {"12.3", Token{
			token.FLOAT,
			token.Location{},
			"12.3",
		}},
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		scanner := NewScanner(newTestFile([]byte(tc.in)))
		err := scanner.readch()
		if err != nil {
			t.Fatalf("%v: Error reading first digit:\n%v", test, err)
		}

		tok, err := scanner.scanNumbers()
		if err != nil {
			t.Fatalf("%v: Error reading full digit:\n%v", test, err)
		}

		if tok.Token != tc.want.Token || tok.Literal != tc.want.Literal {
			t.Fatalf("%v: Want number to be %v, but got %v", test, tc.want.Literal, tok.Literal)
		}

	}
}

func TestScanWords(t *testing.T) {
	tests := []struct {
		in   string
		want Token
	}{
		{"while", Token{
			token.WHILE,
			token.Location{},
			"while",
		}},
		{"var1", Token{
			token.IDENT,
			token.Location{},
			"var1",
		}},
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		scanner := NewScanner(newTestFile([]byte(tc.in)))

		err := scanner.readch()
		if err != nil {
			t.Fatalf("%v: Error reading first word char:\n%v", test, err)
		}

		tok, err := scanner.scanWords()
		if err != nil {
			t.Fatalf("%v: Error reading full word:\n%v", test, err)
		}

		if tok.Token != tc.want.Token || tok.Literal != tc.want.Literal {
			t.Fatalf("%v: Word: %#v is not %#v", test, tok, tc.want)
		}

	}

}

func TestLexer_scanComments(t *testing.T) {
	tests := []struct {
		in   string
		want rune
	}{
		{"// this is a test comment\n", '\n'},
		{"/* this\nis\na\nmulti-line\ncomment\n*/x", 'x'},
	}
	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		scanner := NewScanner(newTestFile([]byte(tc.in)))

		err := scanner.readch()
		if err != nil {
			t.Error(err)
		}

		_, err = scanner.scanComments()
		if err != nil {
			t.Error(err)
		}

		err = scanner.readch()
		if err != nil {
			t.Error(err)
		}

		if scanner.peek != tc.want {
			t.Fatalf("%v: Want peek to be: %v, but got: %v", test, tc.want, string(scanner.peek))
		}
	}
}

func TestLexer_Scan(t *testing.T) {
	tests := []struct {
		in         string
		lineNumber int
		want       []Token
	}{
		{"// test function\n" +
			"func test() bool {\n" +
			"\t    return true;\n" +
			"}",
			4,
			[]Token{
				{
					Token:    token.LINEBREAK,
					Location: token.Location{},
					Literal:  "\n",
				},
				{
					Token:    token.FUNC,
					Location: token.Location{},
					Literal:  "func",
				},
				{
					Token:    token.IDENT,
					Location: token.Location{},
					Literal:  "test",
				},
				{
					Token:    token.LBRACKET,
					Location: token.Location{},
					Literal:  "(",
				},
				{
					Token:    token.RBRACKET,
					Location: token.Location{},
					Literal:  ")",
				},
				{
					Token:    token.IDENT,
					Location: token.Location{},
					Literal:  "bool",
				},
				{
					Token:    token.LCBRACKET,
					Location: token.Location{},
					Literal:  "{",
				},
				{
					Token:    token.LINEBREAK,
					Location: token.Location{},
					Literal:  "\n",
				},
				{
					Token:    token.RETURN,
					Location: token.Location{},
					Literal:  "return",
				},
				{
					Token:    token.TRUE,
					Location: token.Location{},
					Literal:  "true",
				},
				{
					Token:    token.DELIMITER,
					Location: token.Location{},
					Literal:  ";",
				},
				{
					Token:    token.LINEBREAK,
					Location: token.Location{},
					Literal:  "\n",
				},
				{
					Token:    token.RCBRACKET,
					Location: token.Location{},
					Literal:  "}",
				},
				{
					Token:    token.EOF,
					Location: token.Location{},
					Literal:  "EOF",
				},
			}},
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		scanner := NewScanner(newTestFile([]byte(tc.in)))

		for _, wantTok := range tc.want {
			tok, err := scanner.Scan()
			if err != nil {
				t.Error(err)
			}

			if tok.Token != wantTok.Token || tok.Literal != wantTok.Literal {
				t.Fatalf("%v: Want token %#v, but got %#v", test, wantTok, tok)
			}

		}

		if scanner.line != tc.lineNumber {
			t.Fatalf("%v: Linecount is not %v, got: %v", test, tc.lineNumber, scanner.line)
		}
	}
}

func TestLexer_ScanNonCode(t *testing.T) {
	tests := []string{
		"blubb",
		"lorem ipsum +=\n hello world",
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		scanner := NewScanner(newTestFile([]byte(tc)))
		tok, err := scanner.Scan()

		if err != nil {
			t.Error(err)
		}

		if tok == nil {
			t.Fatalf("%v: Expected Token to be filled, but is empty", test)
		}

	}

}

func TestLexer_ScanError(t *testing.T) {
	var err error
	input := []string{"// test doc string\n", "a = 'a\\-"}
	inputCode := input[0] + input[1]

	scanner := NewScanner(newTestFile([]byte(inputCode)))
	wantLocation := token.Location{
		FileName: "/path/to/test.vg",
		Line:     2,
		Position: 8,
		LineFeed: "a = 'a\\-",
	}
	wantError := newScannerSyntaxError(invalidEscapeSequence, wantLocation, "Invalid escape sequence")

	for ; err == nil; _, err = scanner.Scan() {
	}

	if err == nil {
		t.Fatalf("Expected error output")
	}

	if wantError.Error() != err.Error() {
		t.Fatalf("Expected:\n---\n%v\n---\nbut got:\n---\n%v\n---", wantError, err)
	}
}
