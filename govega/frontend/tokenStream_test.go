package frontend

import (
	"govega/govega/language/tokens"
	"testing"
)

func TestNewTokenStream(t *testing.T) {
	ts := NewTokenStream()

	ts.Add(tokens.NewToken(3), *new(ErrorState))
	ts.Add(tokens.NewNum(5), *new(ErrorState))

	e, _ := ts.Remove()

	if e.GetTokenTag() != 3 {
		t.Fatalf("Token %d is not tag 3", e.GetTokenTag())
	}

	e, _ = ts.Remove()
	token := e.GetToken()
	if e.GetTokenTag() == tokens.NUM {
		num := token.(tokens.INum)
		if num.GetValue() != 5 {
			t.Fatalf("Token %d has not number 5", token)
		}
	}

}
