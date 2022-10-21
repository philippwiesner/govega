package frontend

import (
	"bytes"
	"fmt"
	"govega/vega/language"
	"govega/vega/language/tokens"
	"reflect"
	"testing"
)

type testLexerInterface interface {
	Lexer
	getPeek() rune
	getLine() int
	unreadch() error
	readch() error
	readcch(char rune) (bool, error)
	scanCombinedTokens(fch rune, sch rune, word tokens.IWord) (*lexicalToken, error)
	scanLiterals(indicator rune) (*lexicalToken, error)
	scanNumbers() (*lexicalToken, error)
	scanWords() (*lexicalToken, error)
	scanComments() (*lexicalToken, error)
	scan() (*lexicalToken, error)
}

type testLexer struct {
	lexer
}

func newTestLexer(vegaLines []string, inputCode []byte) testLexerInterface {
	v := createTestVega("/path/to/test.vg", vegaLines)
	var lexer testLexerInterface = &testLexer{
		lexer{
			vega:     v.getVega(),
			peek:     0,
			code:     bytes.NewReader(inputCode),
			words:    language.KeyWords,
			lineFeed: "",
			line:     1,
			position: 0,
		},
	}
	return lexer
}

func (t *testLexer) getPeek() rune {
	return t.lexer.peek
}

func (t *testLexer) getLine() int {
	return t.lexer.line
}

type LiteralTestWant struct {
	tag     int
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
	lexer := newTestLexer([]string{}, []byte(text))
	var err error

	err = lexer.readch()
	for ; ; err = lexer.readch() {
		if err != nil {
			t.Error(err)
		} else if lexer.getPeek() == 0 {
			break
		}
		if lexer.getPeek() != want[counter] {
			t.Fatalf("Char %v on position %v is not %v", lexer.getPeek(), counter, want[counter])
		}
		counter++
	}
}

func TestReadcch(t *testing.T) {
	text := "1+1"
	lexer := newTestLexer([]string{}, []byte(text))

	err := lexer.readch()
	if err != nil {
		t.Fatal(err)
	}

	ok, err := lexer.readcch('+')
	if !ok {
		t.Fatalf("Next character not '+'")
	}

}

func TestCombinedTokens(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"!=", tokens.NE},
		{"!-", '!'},
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		lexer := newTestLexer([]string{}, []byte(tc.in))
		err := lexer.readch()
		if err != nil {
			t.Fatal(err)
		}

		token, err := lexer.scanCombinedTokens('!', '=', language.Ne)
		if err != nil {
			t.Fatal(err)
		}

		if token.GetTag() != tc.want {
			t.Fatalf("%v, token should be %v, but is %v", test, tc.want, token.GetTag())
		}

	}
}

func TestScanLiterals(t *testing.T) {
	tests := []LiteralTest{
		{
			"'my literal'",
			LiteralTestWant{tokens.LITERAL, "'my literal'"},
		},
		{
			"'\\tmy \\nliteral'",
			LiteralTestWant{tokens.LITERAL, "'\tmy \nliteral'"},
		},
		{
			"\"my literal\"",
			LiteralTestWant{tokens.LITERAL, "\"my literal\""},
		},
		{
			"'my \\x3A'",
			LiteralTestWant{tokens.LITERAL, "'my \x3a'"},
		},
		{
			"'my \\123 \\\\'",
			LiteralTestWant{tokens.LITERAL, "'my \123 \\'"},
		},
	}
	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		lexer := newTestLexer([]string{}, []byte(tc.in))
		err := lexer.readch()
		if err != nil {
			t.Fatalf("%v: error reading first literal indicator: %v", test, err)
		}

		token, err := lexer.scanLiterals(lexer.getPeek())
		if err != nil {
			t.Fatalf("%v: error reading literal: %v", test, err)
		}

		if token.GetTag() != tc.want.tag {
			t.Fatalf("%v: Want token to be %v, but got: %v", test, tc.want.tag, token.GetTag())
		} else {
			literalToken := token.GetToken().(tokens.ILiteral)
			if literalToken.GetContent() != tc.want.literal {
				t.Fatalf("%v: Want literal to be %v, but got: %v", test, tc.want.literal, literalToken.GetContent())
			}
		}

	}
}

func TestScanLiteralsFailures(t *testing.T) {
	tests := []LiteralFailureTest{
		{
			"linebreak in literal",
			"'fooBar\n'",
			literalNotTerminated,
		},
		{
			"invalid escape sequence",
			"fooBar\\-",
			invalidEscapeSequence,
		},
		{
			"double colon literal escape in single colon literal",
			"'foo\\\"Bar'",
			invalidEscapeSequenceLiteral,
		},
		{
			"single colon literal escape in double colon literal",
			"\"foo\\'Bar\"",
			invalidEscapeSequenceLiteral,
		},
		{
			"invalid hex escape sequence",
			"'foo\\xAg'",
			invalidEscapeSequenceHexadecimal,
		},
		{
			"invalid oct escape sequence",
			"'foo\\088Bar'",
			invalidEscapeSequenceOctal,
		},
		{
			"no literal terminator",
			"'fooBar",
			literalNotTerminated,
		},
	}
	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		lexer := newTestLexer([]string{}, []byte(tc.in))
		err := lexer.readch()
		if err != nil {
			t.Fatalf("%v %v:\nerror reading first literal indicator: %v", test, tc.name, err)
		}

		_, err = lexer.scanLiterals(lexer.getPeek())
		gotErrType := GetVErrorType(err)
		if gotErrType != tc.want {
			t.Fatalf("%v %v: Expected vega Error type to be %v, but got %v", test, tc.name, tc.want, gotErrType)
		}

	}
}

func TestScanNumbers(t *testing.T) {
	tests := []struct {
		in   string
		want interface{}
	}{
		{"123", tokens.NewNum(123)}, {"12,3", tokens.NewNum(12)}, {"12.3", tokens.NewReal(12.3)},
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		lexer := newTestLexer([]string{}, []byte(tc.in))
		err := lexer.readch()
		if err != nil {
			t.Fatalf("%v: Error reading first digit:\n%v", test, err)
		}

		token, err := lexer.scanNumbers()
		if err != nil {
			t.Fatalf("%v: Error reading full digit:\n%v", test, err)
		}

		numberToken := token.GetToken()
		switch number := numberToken.(type) {
		case tokens.INum:
			wantNumber := tc.want.(tokens.INum)
			if number.GetValue() != wantNumber.GetValue() {
				t.Fatalf("%v: Want number to be %d, but got %d", test, wantNumber.GetValue(), number.GetValue())
			}
		case tokens.IReal:
			wantNumber := tc.want.(tokens.IReal)
			if number.GetValue() != wantNumber.GetValue() {
				t.Fatalf("%v: Want number to be %f, but got %f", test, wantNumber.GetValue(), number.GetValue())
			}
		default:
			t.Fatalf("%v: Number Token not valid: %#v", test, numberToken)
		}

	}
}

func TestScanWords(t *testing.T) {
	tests := []struct {
		in   string
		want tokens.IWord
	}{
		{"while", tokens.NewWord("while", tokens.WHILE)},
		{"var1", tokens.NewWord("var1", tokens.ID)},
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		lexer := newTestLexer([]string{}, []byte(tc.in))

		err := lexer.readch()
		if err != nil {
			t.Fatalf("%v: Error reading first word char:\n%v", test, err)
		}

		token, err := lexer.scanWords()
		if err != nil {
			t.Fatalf("%v: Error reading full word:\n%v", test, err)
		}

		word := token.GetToken().(tokens.IWord)
		if !reflect.DeepEqual(word, tc.want) {
			t.Fatalf("%v: Word: %#v is not %#v", test, word, tc.want)
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
		lexer := newTestLexer([]string{}, []byte(tc.in))

		err := lexer.readch()
		if err != nil {
			t.Error(err)
		}

		_, err = lexer.scanComments()
		if err != nil {
			t.Error(err)
		}

		err = lexer.readch()
		if err != nil {
			t.Error(err)
		}

		if lexer.getPeek() != tc.want {
			t.Fatalf("%v: Want peek to be: %v, but got: %v", test, tc.want, string(lexer.getPeek()))
		}
	}
}

func TestLexer_Scan(t *testing.T) {
	tests := []struct {
		in         string
		lineNumber int
		want       []interface{}
	}{
		{"// test function\n" +
			"func test() bool {\n" +
			"\t    return true;\n" +
			"}",
			4,
			[]interface{}{
				tokens.NewWord("func", tokens.FUNC),
				tokens.NewWord("test", tokens.ID),
				tokens.NewToken('('),
				tokens.NewToken(')'),
				language.BoolType,
				tokens.NewToken('{'),
				tokens.NewWord("return", tokens.RETURN),
				tokens.NewWord("true", tokens.TRUE),
				tokens.NewToken(';'),
				tokens.NewToken('}'),
				tokens.NewToken(tokens.EOF),
			}},
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		lexer := newTestLexer([]string{}, []byte(tc.in))

		for _, tok := range tc.want {
			tokenInterface, err := lexer.scan()
			if err != nil {
				t.Error(err)
			}
			switch token := tokenInterface.GetToken().(type) {
			case tokens.IWord:
				wantToken := tok.(tokens.IWord)
				if token.GetTag() != wantToken.GetTag() || token.GetLexeme() != wantToken.GetLexeme() {
					t.Fatalf("%v: Want token %#v, but got %#v", test, wantToken, token)
				}
			case tokens.IToken:
				wantToken := tok.(tokens.IToken)
				if token.GetTag() != wantToken.GetTag() {
					t.Fatalf("%v: Want token %#v, but got %#v", test, wantToken, token)
				}
			default:
				t.Fatalf("%v: Unexpected token: %#v", test, token)
			}
		}

		if lexer.getLine() != tc.lineNumber {
			t.Fatalf("%v: Linecount is not %v, got: %v", test, tc.lineNumber, lexer.getLine())
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
		lexer := newTestLexer([]string{}, []byte(tc))
		token, err := lexer.scan()

		if err != nil {
			t.Error(err)
		}

		if token == nil {
			t.Fatalf("%v: Expected Token to be filled, but is empty", test)
		}

	}

}

func TestLexer_ScanError(t *testing.T) {
	var err error
	input := []string{"// test doc string\n", "a = 'a\\-"}
	inputCode := input[0] + input[1]
	testVega := &vega{
		file:      "/path/to/test.vg",
		codeLines: input,
	}
	wantError := testVega.newLexicalSyntaxError(invalidEscapeSequence, 2, 8, "Invalid escape sequence")

	lexer := newTestLexer([]string{}, []byte(inputCode))

	for ; err == nil; _, err = lexer.scan() {
	}

	if err == nil {
		t.Fatalf("Expected error output")
	}

	if wantError.Error() != err.Error() {
		t.Fatalf("Expected:\n---\n%v\n---\nbut got:\n---\n%v\n---", wantError, err)
	}
}
