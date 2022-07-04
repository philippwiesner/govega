package frontend

import (
	"errors"
	"fmt"
	"govega/govega/language"
	"govega/govega/language/tokens"
	"reflect"
	"testing"
)

type LiteralTestWant struct {
	start   rune
	end     rune
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
	want *LexicalError
}

func TestNewLexer(t *testing.T) {
	text := "1+1\n."
	want := []rune{'1', '+', '1', '\n', '.'}
	counter := 0
	lexer := NewLexer([]byte(text), "test")
	var err error

	err = lexer.readch()
	for ; ; err = lexer.readch() {
		if err != nil {
			t.Error(err)
		} else if lexer.peek == 0 {
			break
		}
		if lexer.peek != want[counter] {
			t.Fatalf("Char %v on position %v is not %v", lexer.peek, counter, want[counter])
		}
		counter++
	}
}

func TestReadcch(t *testing.T) {
	text := "1+1"
	lexer := NewLexer([]byte(text), "test")

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
		lexer := NewLexer([]byte(tc.in), test)
		err := lexer.readch()
		if err != nil {
			t.Fatal(err)
		}

		err = lexer.scanCombinedTokens('!', '=', language.Ne)
		if err != nil {
			t.Fatal(err)
		}
		token, _ := lexer.tokenStream.Remove()

		if token.GetTokenTag() != tc.want {
			t.Fatalf("%v, token should be %v, but is %v", test, tc.want, token.GetTokenTag())
		}

	}
}

func TestScanLiterals(t *testing.T) {
	tests := []LiteralTest{
		{
			"'my literal'",
			LiteralTestWant{'\'', '\'', tokens.LITERAL, "my literal"},
		},
		{
			"'\\tmy \\nliteral'",
			LiteralTestWant{'\'', '\'', tokens.LITERAL, "\tmy \nliteral"},
		},
		{
			"\"my literal\"",
			LiteralTestWant{'"', '"', tokens.LITERAL, "my literal"},
		},
		{
			"'my \\x3A'",
			LiteralTestWant{'\'', '\'', tokens.LITERAL, "my \x3a"},
		},
		{
			"'my \\123 \\\\'",
			LiteralTestWant{'\'', '\'', tokens.LITERAL, "my \123 \\"},
		},
	}
	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		lexer := NewLexer([]byte(tc.in), test)
		err := lexer.readch()
		if err != nil {
			t.Fatalf("%v: error reading first literal indicator: %v", test, err)
		}

		err = lexer.scanLiterals(lexer.peek)
		if err != nil {
			t.Fatalf("%v: error reading literal: %v", test, err)
		}

		token, _ := lexer.tokenStream.Remove()

		if token.GetTokenTag() != int(tc.want.start) {
			t.Fatalf("%v: Want token to be %v, but got: %v", test, tc.want.start, token.GetTokenTag())
		}

		token, _ = lexer.tokenStream.Remove()
		if token.GetTokenTag() != tc.want.tag {
			t.Fatalf("%v: Want token to be %v, but got: %v", test, tc.want.tag, token.GetTokenTag())
		} else {
			literalToken := token.GetToken().(tokens.ILiteral)
			if literalToken.GetContent() != tc.want.literal {
				t.Fatalf("%v: Want literal to be %v, but got: %v", test, tc.want.literal, literalToken.GetContent())
			}
		}

		token, _ = lexer.tokenStream.Remove()
		if token.GetTokenTag() != int(tc.want.start) {
			t.Fatalf("%v: Want token to be %v, but got: %v", test, tc.want.start, token.GetTokenTag())
		}

	}
}

func TestScanLiteralsFailures(t *testing.T) {
	tests := []LiteralFailureTest{
		{
			"linebreak in literal",
			"'fooBar\n'",
			NewLexicalError(literalNotTerminated, ErrorState{1, 8, "test1", "'fooBar\n"}),
		},
		{
			"invalid escape sequence",
			"fooBar\\-",
			NewLexicalError(invalidEscapeSequence, ErrorState{1, 8, "test2", "fooBar\\-"}),
		},
		{
			"double colon literal escape in single colon literal",
			"'foo\\\"Bar'",
			NewLexicalError(invalidEscapeSequenceLiteral, ErrorState{1, 6, "test3", "'foo\\\""}),
		},
		{
			"single colon literal escape in double colon literal",
			"\"foo\\'Bar\"",
			NewLexicalError(invalidEscapeSequenceLiteral, ErrorState{1, 6, "test4", "\"foo\\'"}),
		},
		{
			"invalid hex escape sequence",
			"'foo\\xAg'",
			NewLexicalError(invalidEscapeSequenceHex, ErrorState{1, 8, "test5", "'foo\\xAg"}),
		},
		{
			"invalid oct escape sequence",
			"'foo\\088Bar'",
			NewLexicalError(invalidEscapeSequenceOct, ErrorState{1, 8, "test6", "'foo\\088"}),
		},
		{
			"no literal terminator",
			"'fooBar",
			NewLexicalError(literalNotTerminated, ErrorState{1, 8, "test7", "'fooBar\x00"}),
		},
	}
	var lexError *LexicalError
	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		lexer := NewLexer([]byte(tc.in), test)
		err := lexer.readch()
		if err != nil {
			t.Fatalf("%v %v:\nerror reading first literal indicator: %v", test, tc.name, err)
		}

		err = lexer.scanLiterals(lexer.peek)
		switch {
		case errors.As(err, &lexError):
			gotErr := err.(*LexicalError)
			if !reflect.DeepEqual(tc.want, gotErr) {
				t.Fatalf("%v %v:\nexpected error msg\n%v\n, but got\n%v", test, tc.name, tc.want, gotErr)
			}
		default:
			t.Fatalf("%v %v:\nUnexpected error: %v\nDebug information: %#v", test, tc.name, err, lexer.errorState)
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
		lexer := NewLexer([]byte(tc.in), test)
		err := lexer.readch()
		if err != nil {
			t.Fatalf("%v: Error reading first digit: %v\nDebug Output: %v", test, err, lexer.errorState)
		}

		err = lexer.scanNumbers()
		if err != nil {
			t.Fatalf("%v: Error reading full digit: %v\nDebug Output: %v", test, err, lexer.errorState)
		}

		tokenBucket, err := lexer.tokenStream.Remove()
		numberToken := tokenBucket.GetToken()
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
		lexer := NewLexer([]byte(tc.in), test)

		err := lexer.readch()
		if err != nil {
			t.Fatalf("%v: Error reading first word char %v\nDebug: %v", test, err, lexer.errorState)
		}

		err = lexer.scanWords()
		if err != nil {
			t.Fatalf("%v: Error reading full word: %v\nDebug Output: %v", test, err, lexer.errorState)
		}

		token, err := lexer.tokenStream.Remove()
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
		lexer := NewLexer([]byte(tc.in), test)

		err := lexer.readch()
		if err != nil {
			t.Error(err)
		}

		err = lexer.scanComments()
		if err != nil {
			t.Error(err)
		}

		err = lexer.readch()
		if err != nil {
			t.Error(err)
		}

		if lexer.peek != tc.want {
			t.Fatalf("%v: Want peek to be: %v, but got: %v", test, tc.want, string(lexer.peek))
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
			"func test() -> bool {\n" +
			"\t    return true;\n" +
			"}",
			4,
			[]interface{}{
				tokens.NewWord("func", tokens.FUNC),
				tokens.NewWord("test", tokens.ID),
				tokens.NewToken('('),
				tokens.NewToken(')'),
				language.ReturnValue,
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
		lexer := NewLexer([]byte(tc.in), test)

		ts, err := lexer.Scan()
		if err != nil {
			t.Error(err)
		}

		if lexer.errorState.lineNumber != tc.lineNumber {
			t.Fatalf("%v: Linecount is not %v, got: %v", test, tc.lineNumber, lexer.errorState.lineNumber)
		}

		if ts.GetCount() != len(tc.want) {
			t.Fatalf("%v: Token count is not %v, got: %v", test, tc.want, ts.GetCount())
		}

		for _, tok := range tc.want {
			tokenInterface, err := ts.Remove()
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
	}

}
