package frontend

import (
	"bytes"
	"fmt"
	"govega/govega/dataStructs"
	"govega/govega/helper"
	"govega/govega/language"
	"govega/govega/language/tokens"
	"io"
	"unicode"
)

type lexer struct {
	line        int
	peek        rune
	code        *bytes.Reader
	tokenStream *dataStructs.TokenStream
	words       *helper.HashTable
}

func NewLexer(code []byte) *lexer {
	return &lexer{
		0,
		0,
		bytes.NewReader(code),
		dataStructs.NewTokenStream(),
		language.KeyWords,
	}
}

func (l *lexer) readch() error {
	ch, _, err := l.code.ReadRune()
	if err != nil {
		if err != io.EOF {
			return fmt.Errorf("lexical error: Error reading character %v. Error is: %v", ch, err)
		} else {
			l.peek = ch
			return err
		}
	}
	l.peek = ch
	return nil
}

func (l *lexer) readcch(char rune) (bool, error) {
	err := l.readch()
	if l.peek != char {
		return false, err
	}
	l.peek = 0
	return true, err
}

func (l *lexer) scanCombinedTokens(fch rune, sch rune, word tokens.IWord) error {
	if l.peek == fch {
		ok, err := l.readcch(sch)
		if err != nil {
			return err
		}
		if ok {
			l.tokenStream.Add(word, l.line)
		} else {
			l.tokenStream.Add(tokens.NewToken(int(fch)), l.line)
		}
	}
	return nil
}

func (l *lexer) scanLiterals(indicator rune) error {
	var (
		literal []rune
		err     error
	)
	if l.peek == indicator {
		l.tokenStream.Add(tokens.NewToken(int(indicator)), l.line)
		for ; l.peek != indicator && err == nil; err = l.readch() {
			literal = append(literal, l.peek)
		}
		if err != nil {
			return err
		}
		l.tokenStream.Add(tokens.NewLiteral(literal), l.line)
		if err = l.readch(); err != nil {
			return err
		}
		l.tokenStream.Add(tokens.NewToken(int(indicator)), l.line)
	}
	return nil
}

func (l *lexer) scanNumbers() error {
	var err error
	if unicode.IsDigit(l.peek) {
		value := 0
		for ; unicode.IsDigit(l.peek) && err == nil; err = l.readch() {
			value = 10*value + int(l.peek)
		}
		if err != nil {
			return err
		}
		if l.peek != '.' {
			l.tokenStream.Add(tokens.NewNum(value), l.line)
			return nil
		} else {
			realNumber := float64(value)
			fraction := float64(10)
			for ; unicode.IsDigit(l.peek) && err == nil; err = l.readch() {
				realNumber = realNumber + float64(l.peek)/fraction
				fraction *= 10
			}
			if err != nil {
				return err
			}
			l.tokenStream.Add(tokens.NewReal(realNumber), l.line)
		}
	}
	return nil
}
