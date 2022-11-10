package frontend_test

import (
	"testing"

	. "govega/vega/frontend"
)

func TestParser_ParseError(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			"Misspelled first mandatory keyword",
			"fonc",
			"Missing 'func' at 'fonc'",
		},
		{
			"Function keyword and EOF",
			"func",
			"Unexpected End Of File",
		},
		{
			"No function identifier",
			"func #",
			"Mismatched input '#', expected <identifier>",
		},
		{
			"No function param definition",
			"func test#",
			"Mismatched input '#', expected '('",
		},
		{
			"No function param type or empty function",
			"func test(#",
			"Mismatched input '#', expected <terminal_variable_type> or ')'",
		},
		{
			"No function param name",
			"func test(int #",
			"Mismatched input '#', expected '[' or <identifier>",
		},
		{
			"No closing array definition",
			"func test(int [#",
			"Extraneous input '#', expected ']'",
		},
		{
			"No array name or new array declaration",
			"func test(int []#",
			"Mismatched input '#', expected '[' or <identifier>",
		},
		{
			"No closing bracket or next param operator",
			"func test(int []a int",
			"Mismatched input 'int', expected ',' or ')'",
		},
		{
			"Empty param type",
			"func test(int []a, )",
			"Mismatched input ')', expected <terminal_variable_type>",
		},
		{
			"Invalid variable type",
			"func test(int []a, int b) #",
			"Mismatched input '#', expected <variable_type>",
		},
		{
			"No function body",
			"func test(int []a, int b) int #",
			"Mismatched input '#', expected '{'",
		},
		{
			"Missing function body statement",
			"func test(int []a, int b) int { #",
			"Mismatched input '#', expected 'pass;' or <statement>",
		},
		{
			"Empty function body statement",
			"func test(int []a, int b) int { }",
			"Mismatched input '}', expected 'pass;' or <statement>",
		},
		{
			"Missing line delimiter after pass statement",
			"func test(int []a, int b) int { pass }",
			"Mismatched input '}', expected ';' or line break",
		},
		{
			"Missing line delimiter after continue statement",
			"func test(int []a, int b) int { continue #",
			"Mismatched input '#', expected ';' or line break",
		},
		{
			"Missing line delimiter after break statement",
			"func test(int []a, int b) int { break #",
			"Mismatched input '#', expected ';' or line break",
		},
		{
			"Pass and following statements",
			"func test(int []a, int b) int { break; pass;",
			"Mismatched input 'pass', expected <statement> or '}'",
		},
		{
			"No pass in non empty scope",
			"func test(int []a, int b) int { pass; break;",
			"Mismatched input 'break', expected '}'",
		},
		{
			"No valid expression",
			"func test(int []a, int b) int { if #",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing conditional body",
			"func test(int []a, int b) int { if true #",
			"Mismatched input '#', expected '{'",
		},
		{
			"No closing scope or following statement",
			"func test(int []a, int b) int { if true { pass; } #",
			"Mismatched input '#', expected <statement> or '}'",
		},
		{
			"Missing elif conditional",
			"func test(int []a, int b) int { if true { pass; } elif #",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing second elif conditional",
			"func test(int []a, int b) int { if true { pass; } elif true { break; } elif #",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing conditional else statement body",
			"func test(int []a, int b) int { if true { pass; } else #",
			"Mismatched input '#', expected '{'",
		},
		{
			"Missing while conditional",
			"func test(int []a, int b) int { while #",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing return expression",
			"func test(int []a, int b) int { return #",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing return delimiter",
			"func test(int []a, int b) int { return true #",
			"Mismatched input '#', expected ';' or line break",
		},
		{
			"No valid const variable type",
			"func test(int []a, int b) int { const foo",
			"Mismatched input 'foo', expected <variable_type>",
		},
		{
			"No identifier or array",
			"func test(int []a, int b) int { int #",
			"Mismatched input '#', expected <identifier> or '['",
		},
		{
			"No int in array declaration",
			"func test(int []a, int b) int { int[#",
			"Mismatched input '#', expected <INT>",
		},
		{
			"No closing array bracket",
			"func test(int []a, int b) int { int[8#",
			"Mismatched input '#', expected ']'",
		},
		{
			"No int in two-dim array declaration",
			"func test(int []a, int b) int { int[8][#",
			"Mismatched input '#', expected <INT>",
		},
		{
			"No int in array declaration",
			"func test(int []a, int b) int { int[8][#",
			"Mismatched input '#', expected <INT>",
		},
		{
			"No new array dimension or identifier after array",
			"func test(int []a, int b) int { int[8]#",
			"Mismatched input '#', expected <identifier> or '['",
		},
		{
			"No comma or assignment",
			"func test(int []a, int b) int { int[8] a #",
			"Mismatched input '#', expected ';' or line break",
		},
		{
			"Missing delimiter after declaration assignment",
			"func test(int []a, int b) int { int[8] a = 5#",
			"Mismatched input '#', expected ';' or line break",
		},
		{
			"Missing funcCall, array, comma or assignment",
			"func test(int []a, int b) int { a#",
			"Mismatched input '#', expected '(', '[', ',' or '='",
		},
		{
			"Missing funcCall parameter",
			"func test(int []a, int b) int { a(#",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing funcCall comma or closing bracket",
			"func test(int []a, int b) int { a(b c",
			"Mismatched input 'c', expected ',' or ')'",
		},
		{
			"Missing funcCall closing bracket",
			"func test(int []a, int b) int { a(b, c#",
			"Mismatched input '#', expected ',' or ')'",
		},
		{
			"No assigment after funccall",
			"func test(int []a, int b) int { a(b, c)=",
			"Mismatched input '=', expected ';' or line break",
		},
		{
			"No expression in array assignment",
			"func test(int []a, int b) int { a[#",
			"Mismatched input '#', expected <unary>",
		},
		{
			"No expression in two dim array assignment",
			"func test(int []a, int b) int { a[b][#",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing assignment expression",
			"func test(int []a, int b) int { a[b] = #",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing delimiter",
			"func test(int []a, int b) int { a[b] = 5#",
			"Mismatched input '#', expected ';' or line break",
		},
		{
			"Missing unary in function call",
			"func test(int []a, int b) int { a = b( #",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing coma or closing bracket",
			"func test(int []a, int b) int { a = b(a#",
			"Mismatched input '#', expected ',' or ')'",
		},
		{
			"Missing unary after open bracket",
			"func test(int []a, int b) int { a = (#",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Missing closing bracket pair",
			"func test(int []a, int b) int { a = (b+(c+d)#",
			"Mismatched input '#', expected ')'",
		},
		{
			"Incomplete array declaration",
			"func test(int []a, int b) int { a = [#",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Incomplete array declaration, second element",
			"func test(int []a, int b) int { a = [b#",
			"Mismatched input '#', expected ',' or ']'",
		},
		{
			"Incomplete array access",
			"func test(int []a, int b) int { a = b[#",
			"Mismatched input '#', expected <unary>",
		},
		{
			"Incomplete array access, closing array bracket",
			"func test(int []a, int b) int { a = b[4#",
			"Mismatched input '#', expected ']'",
		},
		{
			"String literal not terminated",
			"func test(int []a, int b) int { a = 'fooBar",
			"String literal not terminated",
		},
		{
			"Invalid excape sequence",
			"func test(int []a, int b) int { a = '\\Fd'",
			"Invalid escape sequence",
		},
	}

	for i, tc := range tests {

		testNumber := i + 1

		vega := NewVega("/path/to/test.vg")
		lexer := vega.NewLexer([]byte(tc.in))
		parser := vega.NewParser(lexer)
		parseErr := parser.Parse(parser)

		if parseErr == nil {
			t.Fatalf("Test%d: Expected error, got nil", testNumber)
		}

		switch err := parseErr.(type) {
		case IVError:
			if err.GetMessage() != tc.want {
				t.Fatalf("Test%d: %v:\n\n%v\n\nExpected error message to be:\n\t%v\nbut got:\n\t%v", testNumber, tc.name, tc.in, tc.want, err.GetMessage())
			}
		case error:
			t.Fatalf("Test%d: Expected IVError, but got %v", testNumber, err)

		}

	}
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{
			"Full Code test",
			`


/* This is a multiline comment
which spans over mutiple lines */

// This is a single line comment

// reserved array method
func fooBar(int[] a, 

bool f) 


int
{
	a = 1 + 6+ f(4+6) + a[3]
	b = true == not false != false or false and true
	const int i; const int a
	const int i
	const int g

	return 1


}

func main() int {
	int[5] a = [1, 2, 4, 5, 6 + 8]
	char c = 'g'
	str s = '\xFF Hello World'
	bool a = fooBar()
	if c == 'g' and a {
		while true {
			if c == 'g' {
				continue
			} else {
				break
			}
		}
	} elif false {
		pass
	} else {
		float x = 0.5
	}

	return 0
}
`,
		},
	}

	for i, tc := range tests {

		testNumber := i + 1

		vega := NewVega("/path/to/test.vg")
		lexer := vega.NewLexer([]byte(tc.in))
		parser := vega.NewParser(lexer)
		parseErr := parser.Parse(parser)

		if parseErr != nil {
			t.Fatalf("Test%d: %v: Expected no error, but got:\n\n%v", testNumber, tc.name, parseErr)
		}

	}
}
