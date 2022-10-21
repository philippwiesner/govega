// Package language
//
// Defines basic language structures which can be used in the frontend to parse the language
//
// vocabulary.go defines language vocabulary and functions to check if special characters are in the language alphabet

package language

import (
	"govega/vega/helper"
	"govega/vega/language/tokens"
)

// define special combined tokens, keywords and special escaped characters to be used by the lexer
var (
	Eq                    = tokens.NewWord("==", tokens.EQ)
	Ne                    = tokens.NewWord("!=", tokens.NE)
	Le                    = tokens.NewWord("<=", tokens.LE)
	Ge                    = tokens.NewWord(">=", tokens.GE)
	BoolAnd               = tokens.NewWord("&&", tokens.BOOLAND)
	BoolOr                = tokens.NewWord("||", tokens.BOOLOR)
	KeyWords              = initKeyWords()
	EscapeHexaLiterals    = initHexadecimalChars()
	EscapeOctalLiterals   = initOctalChars()
	EscapeUnicodeLiterals = initUnicodeChars()
)

// initKeyWords creates a new lookup Hashtable containing all the keywords of the language
func initKeyWords() helper.HashTable {
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

// initHexadecimalChars creates a lookup Hashtable for all hexadecimal escaped characters.
//
// valid hexadecimal escape sequences are \x00 - \xff. Uppercase will automatically be converted to lowercase.
func initHexadecimalChars() helper.HashTable {
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

// initOctalChars creates a lookup Hashtable for all octal escaped characters.
//
// valid octal escape sequences are \o000 - \0377.
func initOctalChars() helper.HashTable {
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

// initUnicodeChars creates a lookup Hashtable for all unicode escaped characters.
//
// valid unicode escape sequences are \u0000 - \uffff. Uppercase will automatically be converted to lowercase.
func initUnicodeChars() helper.HashTable {
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
