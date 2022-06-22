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
		token, ok := lexer.tokenStream.Remove()
		if !ok {
			t.Fatalf("Error retrieving token from stream")
		}

		if token.GetTokenTag() != tc.want {
			t.Fatalf("Test%d, token should be %v, but is %v", i, tc.want, token.GetTokenTag())
		}

	}
}
