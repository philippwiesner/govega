package language

import (
	"govega/govega/helper"
	"govega/govega/language/tokens"
)

var (
	ReturnValue = tokens.NewWord("->", tokens.RETURNVALUE)
	Eq          = tokens.NewWord("<=", tokens.EQ)
	Ne          = tokens.NewWord("!=", tokens.NE)
	Le          = tokens.NewWord("<=", tokens.LE)
	Ge          = tokens.NewWord(">=", tokens.GE)
	BoolAnd     = tokens.NewWord("&&", tokens.BOOLAND)
	BoolOr      = tokens.NewWord("||", tokens.BOOLOR)
	KeyWords    = initKeyWords()
)

func initKeyWords() *helper.HashTable {
	basicTypes := []*BasicType{intType, floatType, charType, boolType}
	vocabulary := []tokens.IWord{
		tokens.NewWord("str", tokens.TYPE),
		tokens.NewWord("true", tokens.TRUE),
		tokens.NewWord("false", tokens.FALSE),
		tokens.NewWord("func", tokens.FUNC),
		tokens.NewWord("const", tokens.CONST),
		tokens.NewWord("return", tokens.RETURN),
		tokens.NewWord("while", tokens.WHILE),
		tokens.NewWord("break", tokens.BREAK),
		tokens.NewWord("continue", tokens.CONTINUE),
		tokens.NewWord("pass", tokens.PASS),
		tokens.NewWord("if", tokens.IF),
		tokens.NewWord("elif", tokens.ELIF),
		tokens.NewWord("else", tokens.ELSE),
		tokens.NewWord("and", tokens.AND),
		tokens.NewWord("or", tokens.OR),
		tokens.NewWord("not", tokens.NOT),
	}

	table := helper.NewHashTable()

	for _, b := range basicTypes {
		table.Add(b.GetLexeme(), b)
	}

	for _, v := range vocabulary {
		table.Add(v.GetLexeme(), v)
	}

	return table
}
