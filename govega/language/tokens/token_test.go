package tokens

import (
	"fmt"
	"testing"
)

func TestNewToken(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{'{', '{'},
		{NUM, NUM},
	}

	for i, tc := range tests {
		tk := NewToken(tc.in)
		got := tk.GetTag()
		if got != tc.want {
			t.Fatalf("Test %v: Want %v, but got: %v", i+1, tc.want, tc.in)
		}
	}

	t1 := NewToken('+')
	t2 := NewToken(OR)

	if t1.GetTag() == t2.GetTag() {
		t.Fatalf("Tag should not be equal between: %v and %v", t1, t2)
	}
}

func TestNewNum(t *testing.T) {
	n1 := NewNum(5)
	if n1.GetTag() != NUM && n1.GetValue() != 5 {
		t.Fatalf("Want Tag: %v, but got: %v, Want Value: %v, but is %v", NUM, n1.GetTag(), 5, n1.GetValue())
	}
}

func TestNewReal(t *testing.T) {
	n := NewNum(5)
	r := NewReal(5.0)
	if n.GetTag() == r.GetTag() {
		t.Fatalf("Num and Real are incomparable types")
	}
}

func Test(t *testing.T) {
	a := NewToken('g')
	b := NewWord("böib", ID)
	a.GetTag()
	b.GetTag()
}

func TestInterface(t *testing.T) {
	to := NewToken(NUM)
	fmt.Println(to.GetTag())
	to2 := NewWord("blubb", FUNCTION)
	fmt.Println(to2.GetTag())
	fmt.Println(to2.GetLexeme())
}
