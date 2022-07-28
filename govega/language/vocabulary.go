package language

import (
	"govega/govega/helper"
	"govega/govega/language/tokens"
)

var (
	ReturnType            = tokens.NewWord("->", tokens.RETURNTYPE)
	Eq                    = tokens.NewWord("<=", tokens.EQ)
	Ne                    = tokens.NewWord("!=", tokens.NE)
	Le                    = tokens.NewWord("<=", tokens.LE)
	Ge                    = tokens.NewWord(">=", tokens.GE)
	BoolAnd               = tokens.NewWord("&&", tokens.BOOLAND)
	BoolOr                = tokens.NewWord("||", tokens.BOOLOR)
	KeyWords              = initKeyWords()
	EscapeHexaLiterals    = initHexaLiterals()
	EscapeOctalLiterals   = initOctalLiterals()
	EscapeUnicodeLiterals = initUnicodeLiterals()
)

func initKeyWords() *helper.HashTable {
	basicTypes := []IBasicType{IntType, FloatType, CharType, BoolType}
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

func initHexaLiterals() *helper.HashTable {
	table := helper.NewHashTable()
	alphabet := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			hexKey := string(alphabet[i]) + string(alphabet[j])
			hexValue := rune(i*16 + j)
			table.Add(hexKey, hexValue)
		}
	}
	return table
}

func initOctalLiterals() *helper.HashTable {
	table := helper.NewHashTable()
	alphabet := []rune{'0', '1', '2', '3', '4', '5', '6', '7'}
	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			for k := 0; k < 8; k++ {
				octKey := string(alphabet[i]) + string(alphabet[j]) + string(alphabet[k])
				octValue := rune(i*8*8 + j*8 + k)
				table.Add(octKey, octValue)
			}
		}
	}
	return table
}

func initUnicodeLiterals() *helper.HashTable {
	table := helper.NewHashTable()
	alphabet := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			for k := 0; k < 16; k++ {
				for l := 0; l < 16; l++ {
					uniKey := string(alphabet[i]) + string(alphabet[j]) + string(alphabet[k]) + string(alphabet[l])
					uniValue := rune(i*16*16*16 + j*16*16 + k*16 + l)
					table.Add(uniKey, uniValue)
				}
			}
		}
	}
	return table
}
