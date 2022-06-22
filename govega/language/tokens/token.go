// Package tokens
//
// Defines basic language structures which can be used in the frontend to parse the language
//
// tokens.go defines language tokens which are defined by a unique ID for each token

package tokens

// token tag constants for identifying tokens. numbering start at 256 as the integers 0-255 represent runes (chars)
const (
	EQ          int = iota + 256 // ==
	LE                           // <=
	GE                           // >=
	NE                           // !=
	CONST                        // const
	FUNC                         // func
	WHILE                        // while
	IF                           // if
	ELIF                         // elif
	ELSE                         // else
	RETURNVALUE                  // ->
	RETURN                       // return
	PASS                         // pass
	CONTINUE                     // continue
	BREAK                        // break
	TRUE                         // true
	FALSE                        // false
	NOT                          // not
	AND                          // and
	BOOLAND                      // &&
	OR                           // or
	BOOLOR                       // ||
	INDEX                        // [i]
	ID                           // identifier
	BASIC                        // basic data type (e.g. int, char)
	FUNCTION                     // function identifier
	TYPE                         // non-basic data types (e.g. string, array)
	NUM                          // normal numbers (int)
	REAL                         // real numbers (floating point)
	LITERAL                      // everything enclosed in '' or ""
)

// token struct represents simple basic language tokens identified by an integer number
type token struct {
	tag int
}

// newToken is the constructor for a new token
func newToken(t int) *token {
	return &token{tag: t}
}

// GetTag public getter method for getting the tag
func (t *token) GetTag() int {
	return t.tag
}

// num is a numeric tokens
type num struct {
	*token
	value int // tokens value
}

// newNum is the constructor for a new numeric tokens
//
// The token tag is being set to NUM
func newNum(v int) *num {
	return &num{newToken(NUM), v}
}

// GetValue public getter method for retrieving the numeric tokens value
func (n *num) GetValue() int {
	return n.value
}

// word is a word token
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

// realNumber is a floating point number token
type realNumber struct {
	*token
	value float64 // floating point number value
}

// newReal is the constructor for a new real number
func newReal(v float64) *realNumber {
	return &realNumber{newToken(REAL), v}
}

// GetValue public getter method for the floating point number value
func (r *realNumber) GetValue() float64 {
	return r.value
}

// literal is a literal tokens (everything enclosed in '' or "")
type literal struct {
	*token
	content []rune // content between '' or ""
}

// newLiteral is the constructor for a new literal tokens
func newLiteral(c []rune) *literal {
	return &literal{newToken(LITERAL), c}
}

// GetContent public getter method for the literal content
func (l *literal) GetContent() []rune {
	return l.content
}
