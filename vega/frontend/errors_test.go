package frontend

import (
	"testing"
)

func createCustomError(v testVegaInterface) error {
	return v.getVega().newLexicalSyntaxError(malformedCode, 1, 4, "custom error message")
}

func TestVega_ErrorType(t *testing.T) {
	v := createTestVega("/path/to/test.vg", []string{"a = 1 + 2;", "b = a + 4;"})

	vErr := createCustomError(v)
	want := malformedCode
	got := GetVErrorType(vErr)

	if got != want {
		t.Fatalf("Excpeted vega Error Type to be %v, but got %v", want, got)
	}
}

func TestVega_ErrorEmptyCode(t *testing.T) {
	v := createTestVega("/path/to/test.vg", []string{})

	vErr := v.getVega().newLexicalSyntaxError(malformedCode, 1, 0, "Error parsing file")

	want := malformedCode
	got := GetVErrorType(vErr)

	if got != want {
		t.Fatalf("Excpeted vega Error Type to be %v, but got %v", want, got)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Code paniced with %v", r)
		}
	}()

	vErr.Error()
}
