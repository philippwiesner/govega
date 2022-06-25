package frontend

import (
	"govega/govega/language"
	"govega/govega/language/tokens"
	"io"
	"testing"
)

func TestNewLexer(t *testing.T) {
	text := "1+1\n."
	want := []rune{'1', '+', '1', '\n', '.'}
	counter := 0
	lexer := NewLexer([]byte(text))

	for {
		if err := lexer.readch(); err != nil {
			if err != nil {
				if err == io.EOF {
					break
				} else {
					t.Fatal(err)
				}
			}
		}
		if lexer.peek != want[counter] {
			t.Fatalf("Char %v on position %v is not %v", lexer.peek, counter, want[counter])
		}
		counter++
	}

}

func TestReadcch(t *testing.T) {
	text := "1+1"
	lexer := NewLexer([]byte(text))

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
		lexer := NewLexer([]byte(tc.in))
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
			t.Fatalf("Test%d, token should be %v, but is %v", i, tc.want, token.GetTokenTag())
		}

	}
}

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
		lexer := NewLexer([]byte(tc.in))
		err := lexer.readch()
		if err != nil {
			t.Fatalf("test%d: error reading first literal indicator: %v", i, err)
		}

		err = lexer.scanLiterals(lexer.peek)
		if err != nil {
			t.Fatalf("test%d: error reading literal: %v", i, err)
		}

		token, _ := lexer.tokenStream.Remove()

		if token.GetTokenTag() != int(tc.want.start) {
			t.Fatalf("test%d: Want token to be %v, but got: %v", i, tc.want.start, token.GetTokenTag())
		}

		token, _ = lexer.tokenStream.Remove()
		if token.GetTokenTag() != tc.want.tag {
			t.Fatalf("test%d: Want token to be %v, but got: %v", i, tc.want.tag, token.GetTokenTag())
		} else {
			literalToken := token.GetToken().(tokens.ILiteral)
			if literalToken.GetContent() != tc.want.literal {
				t.Fatalf("test%d: Want literal to be %v, but got: %v", i, tc.want.literal, literalToken.GetContent())
			}
		}

		token, _ = lexer.tokenStream.Remove()
		if token.GetTokenTag() != int(tc.want.start) {
			t.Fatalf("test%d: Want token to be %v, but got: %v", i, tc.want.start, token.GetTokenTag())
		}

	}
}
