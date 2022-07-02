package frontend

import (
	"bytes"
	"govega/govega/helper"
	"govega/govega/language"
	"govega/govega/language/tokens"
	"io"
)

type lexer struct {
	errorState  ErrorState
	peek        rune
	code        *bytes.Reader
	tokenStream *TokenStream
	words       *helper.HashTable
}

func NewLexer(code []byte, fileName string) *lexer {
	return &lexer{ErrorState{1, 0, fileName, ""},
		0,
		bytes.NewReader(code),
		NewTokenStream(),
		language.KeyWords,
	}
}

func (l *lexer) readch() error {
	ch, _, err := l.code.ReadRune()
	l.errorState.position++
	l.errorState.lineFeed += string(ch)
	if err != nil {
		if err != io.EOF {
			return NewLexicalError(malformedCode, l.errorState)
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
			l.tokenStream.Add(word, l.errorState)
		} else {
			l.tokenStream.Add(tokens.NewToken(int(fch)), l.errorState)
		}
	}
	return nil
}

func (l *lexer) scanLiterals(indicator rune) error {
	var (
		literal string
		char    rune
		err     error
	)
	if l.peek == indicator {
		l.tokenStream.Add(tokens.NewToken(int(indicator)), l.errorState)
		err = l.readch()
		for ; l.peek != indicator && err == nil; err = l.readch() {
			if l.peek == '\n' {
				return NewLexicalError(literalNotTerminated, l.errorState)
			}
			if l.peek == '\\' {
				err = l.readch()
				if err != nil {
					return NewLexicalError(invalidEscapeSequence, l.errorState)
				}
				switch l.peek {
				case 'b':
					char = '\b'
				case 'f':
					char = '\f'
				case 'n':
					char = '\n'
				case 'r':
					char = '\r'
				case 't':
					char = '\t'
				case 'v':
					char = '\v'
				case '\\':
					char = '\\'
				case '"':
					if indicator == '"' {
						char = '"'
					} else {
						return NewLexicalError(invalidEscapeSequenceLiteral, l.errorState)
					}
				case '\'':
					if indicator == '\'' {
						char = '\''
					} else {
						return NewLexicalError(invalidEscapeSequenceLiteral, l.errorState)
					}
				case 'x':
					hex := ""
					for i := 0; i < 2; i++ {
						err = l.readch()
						if err != nil {
							return NewLexicalError(invalidEscapeSequenceHex, l.errorState)
						}
						// transform uppercase hex in lowercase for lookup
						if l.peek > 64 && l.peek < 91 {
							l.peek = l.peek + 32
						}
						hex = hex + string(l.peek)
					}
					hexLookup, ok := language.EscapeHexaLiterals.Get(hex)
					if !ok {
						return NewLexicalError(invalidEscapeSequenceHex, l.errorState)
					}
					char = hexLookup.(rune)
				case '0', '1', '2', '3':
					oct := string(l.peek)
					for i := 0; i < 2; i++ {
						err = l.readch()
						if err != nil {
							return NewLexicalError(invalidEscapeSequenceOct, l.errorState)
						}
						oct = oct + string(l.peek)
					}
					octLookup, ok := language.EscapeOctalLiterals.Get(oct)
					if !ok {
						return NewLexicalError(invalidEscapeSequenceOct, l.errorState)
					}
					char = octLookup.(rune)
				default:
					return NewLexicalError(invalidEscapeSequence, l.errorState)
				}
				l.peek = 0
			} else {
				char = l.peek
			}
			literal = literal + string(char)
		}
		if err != nil {
			return NewLexicalError(literalNotTerminated, l.errorState)
		}
		l.tokenStream.Add(tokens.NewLiteral(literal), l.errorState)
		l.tokenStream.Add(tokens.NewToken(int(l.peek)), l.errorState)
	}
	return nil
}

func (l *lexer) scanNumbers() error {
	var err error
	if l.peek > 47 && l.peek < 58 {
		value := 0
		for ; l.peek > 47 && l.peek < 58 && err == nil; err = l.readch() {
			value = 10*value + (int(l.peek) - 48)
		}
		if err == io.EOF {
			l.tokenStream.Add(tokens.NewNum(value), l.errorState)
			l.tokenStream.Add(tokens.NewToken(tokens.EOF), l.errorState)
			return err
		} else if err != nil {
			return err
		}
		if l.peek != '.' {
			l.tokenStream.Add(tokens.NewNum(value), l.errorState)
			return nil
		} else {
			realNumber := float64(value)
			fraction := float64(10)
			err = l.readch()
			for ; l.peek > 47 && l.peek < 58 && err == nil; err = l.readch() {
				realNumber = realNumber + (float64(l.peek)-48)/fraction
				fraction *= 10
			}
			if err == io.EOF {
				l.tokenStream.Add(tokens.NewReal(realNumber), l.errorState)
				l.tokenStream.Add(tokens.NewToken(tokens.EOF), l.errorState)
				return err
			} else if err != nil {
				return err
			}
			l.tokenStream.Add(tokens.NewReal(realNumber), l.errorState)
		}
	}
	return nil
}

func (l *lexer) Scan() (ts *TokenStream, e error) {
	for err := *new(error); ; err = l.readch() {
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		if l.peek == '\n' {
			l.errorState.lineFeed = ""
			l.errorState.position = 0
			continue
		}
		l.tokenStream.Add(tokens.NewToken(int(l.peek)), l.errorState)
	}
	return l.tokenStream, nil
}
