// Package tokens
//
// Defines basic language structures which can be used in the frontend to parse the language
//
// interfaces.go defines external interfaces for tokens

package tokens

// IToken interface to retrieve tag for all token implementations
type IToken interface {
	GetTag() int
}

// NewToken generates new IToken interface based on token struct
func NewToken(tag int) IToken {
	var token IToken = newToken(tag)
	return token
}

// INum interface extends IToken interface for number tokens
type INum interface {
	IToken
	GetValue() int
}

// NewNum generates new INum interface based on num struct
func NewNum(v int) INum {
	var num INum = newNum(v)
	return num
}

// IWord interface extends IToken interface for word tokens
type IWord interface {
	IToken
	GetLexeme() string
}

// NewWord generates new IWord interface based on word struct
func NewWord(l string, t int) IWord {
	var word IWord = newWord(l, t)
	return word
}

// IReal interface extends IToken interface for real number tokens
type IReal interface {
	IToken
	GetValue() float64
}

// NewReal generates new IReal interfaces based on realNumber token
func NewReal(v float64) IReal {
	var realNumber IReal = newReal(v)
	return realNumber
}

// ILiteral interface extends IToken interface for literal tokens
type ILiteral interface {
	IToken
	GetContent() string
}

// NewLiteral generates new ILiteral interface based on literal token
func NewLiteral(c string) ILiteral {
	var literal ILiteral = newLiteral(c)
	return literal
}
