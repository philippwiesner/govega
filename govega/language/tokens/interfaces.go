package tokens

type IToken interface {
	GetTag() int
}

func NewToken(tag int) IToken {
	var token IToken = newToken(tag)
	return token
}

type INum interface {
	IToken
	GetValue() int
}

func NewNum(v int) INum {
	var num INum = newNum(v)
	return num
}

type IWord interface {
	IToken
	GetLexeme() string
}

func NewWord(l string, t int) IWord {
	var word IWord = newWord(l, t)
	return word
}

type IReal interface {
	IToken
	GetValue() float64
}

func NewReal(v float64) IReal {
	var realNumber IReal = newReal(v)
	return realNumber
}

type ILiteral interface {
	IToken
	GetContent() []rune
}

func NewLiteral(c []rune) ILiteral {
	var literal ILiteral = newLiteral(c)
	return literal
}
