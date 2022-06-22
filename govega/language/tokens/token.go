// Package language
//
// Defines basic language structures which can be used in the frontend to parse the language
//
// tokens.go defines language tokens which are defined by a unique ID for each tokens

package tokens

const (
	EQ           int = iota + 256 // ==
	LE                            // <=
	GE                            // >=
	NE                            // !=
	CONST                         // const
	FUNC                          // func
	WHILE                         // while
	IF                            // if
	ELIF                          // elif
	ELSE                          // else
	RETURN_VALUE                  // ->
	RETURN                        // return
	PASS                          // pass
	CONTINUE                      // continue
	BREAK                         // break
	TRUE                          // true
	FALSE                         // false
	NOT                           // not
	AND                           // and
	BOOL_AND                      // &&
	OR                            // or
	BOOL_OR                       // ||
	INDEX                         // [i]
	ID                            // identifier
	BASIC                         // basic data type (e.g. int, char)
	FUNCTION                      // function identifier
	TYPE                          // non-basic data types (e.g. string, array)
	NUM                           // normal numbers (int)
	REAL                          // real numbers (floating point)
	LITERAL                       // everything enclosed in '' or ""
)

// Token struct to represent simple basic language tokens
type token struct {
	tag int
}

// NewToken is the constructor for a new tokens
func newToken(t int) *token {
	return &token{tag: t}
}

// GetTag public getter method for getting the tag
func (t *token) GetTag() int {
	return t.tag
}

// Num is a numeric tokens
type num struct {
	*token
	value int // tokens value
}

// NewNum is the constructor for a new numeric tokens
//
// The tokens tag is being set to NUM
func newNum(v int) *num {
	return &num{newToken(NUM), v}
}

// GetValue public getter method for retrieving the numeric tokens value
func (n *num) GetValue() int {
	return n.value
}

// Word is a word tokens
type word struct {
	*token
	lexeme string // word of the language (more than one character)
}

// NewWord is the constructor for a new word
//
// The tag depends on the word, language keywords have tags defined, identifier and function identifier get
// a special tag. The differentiation is being made in the lexer
func newWord(l string, t int) *word {
	return &word{newToken(t), l}
}

// GetLexeme public getter method for lexeme
func (w *word) GetLexeme() string {
	return w.lexeme
}

// RealNumber is a floating point number tokens
type realNumber struct {
	*token
	value float64 // floating point number value
}

// NewReal is the constructor for a new real number
func newReal(v float64) *realNumber {
	return &realNumber{newToken(REAL), v}
}

// GetValue public getter method for the floating point number value
func (r *realNumber) GetValue() float64 {
	return r.value
}

// Literal is a literal tokens (everything enclosed in '' or "")
type literal struct {
	*token
	content []rune // content between '' or ""
}

// NewLiteral is the constructor for a new literal tokens
func newLiteral(c []rune) *literal {
	return &literal{newToken(LITERAL), c}
}

// GetContent public getter method for the literal content
func (l *literal) GetContent() []rune {
	return l.content
}
