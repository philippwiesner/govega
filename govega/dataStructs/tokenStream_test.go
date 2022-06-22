package dataStructs

import (
	"govega/govega/language/tokens"
	"testing"
)

func TestNewTokenStream(t *testing.T) {
	ts := NewTokenStream()

	ts.Add(tokens.NewToken(3), 5)
	ts.Add(tokens.NewNum(5), 6)

	e, ok := ts.Remove()
	if !ok {
		t.Fatalf("Error removing element")
	}

	if e.GetTokenTag() != 3 {
		t.Fatalf("Token %d is not tag 3", e.GetTokenTag())
	}

	e, ok = ts.Remove()
	if !ok {
		t.Fatalf("Error removing element")
	}
	token := e.GetToken()
	if e.GetTokenTag() == tokens.NUM {
		num := token.(tokens.INum)
		if num.GetValue() != 5 {
			t.Fatalf("Token %d has not number 5", token)
		}
	}

}
