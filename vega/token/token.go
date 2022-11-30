// Package token
//
// Defines basic language structures which can be used in the frontend to parse the language
//
// token.go defines language token which are defined by a unique IDENT for each token

package token

type Token int

// token tag constants for identifying token.
const (
	ILLEGAL Token = iota
	EOF           // end of file

	// literals
	literal_start
	IDENT  // identifier
	INT    // integer 123
	FLOAT  // float 123.45
	CHAR   // 'a'
	STRING // "hello world"
	literal_end

	operator_start
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	BIT_OR  // |
	BIT_AND // &

	LAND // &&
	LOR  // ||
	LNOT // !

	EQ      // ==
	LE      // <=
	GE      // >=
	NE      // !=
	LESS    // <
	GREATER // >
	ASSIGN  // =

	LCBRACKET // {
	RCBRACKET // }
	LSBRACKET // [
	RSBRACKET // ]
	LBRACKET  // (
	RBRACKET  // )
	COMMA     // ,

	COLON     // :
	LINEBREAK // \n
	DELIMITER // ;
	operator_end

	keyword_start
	AND
	BREAK
	CASE
	CONST
	CONTINUE
	DEFAULT
	ELIF
	ELSE
	FALSE
	FUNC
	IF
	NOT
	OR
	PASS
	RETURN
	SWITCH
	TRUE
	VAR
	WHILE
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	CHAR:   "CHAR",
	STRING: "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	BIT_OR:  "|",
	BIT_AND: "&",

	LAND: "&&",
	LOR:  "||",
	LNOT: "!",

	EQ:      "==",
	LE:      "<=",
	GE:      ">=",
	NE:      "!=",
	LESS:    "<",
	GREATER: ">",
	ASSIGN:  "=",

	LCBRACKET: "{",
	RCBRACKET: "}",
	LSBRACKET: "[",
	RSBRACKET: "]",
	LBRACKET:  "(",
	RBRACKET:  ")",
	COMMA:     ",",

	COLON:     ":",
	LINEBREAK: "\n",
	DELIMITER: ";",

	AND:      "and",
	BREAK:    "break",
	CASE:     "case",
	CONST:    "const",
	CONTINUE: "continue",
	DEFAULT:  "default",
	ELIF:     "elif",
	ELSE:     "else",
	FALSE:    "false",
	FUNC:     "func",
	IF:       "if",
	NOT:      "not",
	OR:       "or",
	PASS:     "pass",
	RETURN:   "return",
	SWITCH:   "switch",
	TRUE:     "true",
	VAR:      "var",
	WHILE:    "while",
}

// String print token as string
func (t Token) String() string {
	return tokens[t]
}

func (t Token) IsLiteral() bool {
	return t > literal_start && t < literal_end
}

var Keywords = initKeywords()

func initKeywords() map[string]Token {
	keywords := make(map[string]Token, keyword_end-(keyword_start+1))
	for i := keyword_start + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
	return keywords
}

func IsKeyword(name string) bool {
	_, ok := Keywords[name]
	return ok
}

func IsIdentifier(name string) bool {
	if name == "" || IsKeyword(name) {
		return false
	}
	return true
}

func Lookup(name string) Token {
	if !IsIdentifier(name) {
		return Keywords[name]
	} else {
		return IDENT
	}
}
