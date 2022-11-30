package frontend

import (
	"testing"

	"govega/vega/frontend/vega"
	"govega/vega/parser"
)

func createCustomError(v vega.testVegaInterface) error {
	return v.getVega().newLexicalSyntaxError(malformedCode, 1, 4, "custom error message")
}

func TestVega_ErrorType(t *testing.T) {
	v := vega.createTestVega("/path/to/test.vg", []string{"a = 1 + 2;", "b = a + 4;"})

	vErr := createCustomError(v)
	want := malformedCode
	got := parser.GetVErrorType(vErr)

	if got != want {
		t.Fatalf("Excpeted VegaObject Error Type to be %v, but got %v", want, got)
	}
}

func TestVega_ErrorEmptyCode(t *testing.T) {
	v := vega.createTestVega("/path/to/test.vg", []string{})

	vErr := v.getVega().newLexicalSyntaxError(malformedCode, 1, 0, "Error parsing file")

	want := malformedCode
	got := parser.GetVErrorType(vErr)

	if got != want {
		t.Fatalf("Excpeted VegaObject Error Type to be %v, but got %v", want, got)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Code paniced with %v", r)
		}
	}()

	vErr.Error()
}
