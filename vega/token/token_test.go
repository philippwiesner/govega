package token

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewToken(t *testing.T) {
	tests := []struct {
		in   Token
		want Token
	}{
		{LCBRACKET, LCBRACKET},
		{INT, INT},
	}

	for i, tc := range tests {
		tk := tc.in
		if tk != tc.want {
			t.Fatalf("Test %v: Want %v, but got: %v", i+1, tc.want, tc.in)
		}
	}

	t1 := ADD
	t2 := OR

	if t1 == t2 {
		t.Fatalf("Tag should not be equal between: %v and %v", t1, t2)
	}
}

func newTestFile(source []byte) File {
	return File{
		Name:   "/path/to/test.vg",
		Source: source,
	}
}

func TestNewFile(t *testing.T) {
	tests := []struct {
		in   string
		want []byte
	}{
		{
			"", []byte{},
		},
		{
			"vega file", []byte("vega file"),
		},
	}

	for i, tc := range tests {
		test := fmt.Sprintf("test%d", i+1)
		file := newTestFile([]byte(tc.in))

		if file.Name != "/path/to/test.vg" && !reflect.DeepEqual(file.Source, tc.want) {
			t.Fatalf("%v: Want source to be %v, but got %v", test, tc.want, file.Source)
		}

	}

}
