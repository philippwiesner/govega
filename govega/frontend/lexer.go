package frontend

import (
	"bytes"
	"govega/govega/dataStructs"
	"govega/govega/helper"
	"govega/govega/language"
	"govega/govega/language/tokens"
	"io"
	"unicode"
)

type lexer struct {
	line        int
	pos         int
	peek        rune
	code        *bytes.Reader
	tokenStream *dataStructs.TokenStream
	words       *helper.HashTable
}

func NewLexer(code []byte) *lexer {
	return &lexer{
		0,
		0,
		0,
		bytes.NewReader(code),
		dataStructs.NewTokenStream(),
		language.KeyWords,
	}
}

func (l *lexer) readch() error {
	ch, _, err := l.code.ReadRune()
	l.pos++
	if err != nil {
		if err != io.EOF {
			return &LexicalError{
				"Error parsing file, probably malformed",
				l.line,
				l.pos,
			}
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
		literal string
		char    rune
		err     error
	)
	if l.peek == indicator {
		l.tokenStream.Add(tokens.NewToken(int(indicator)), l.line)
		err = l.readch()
		for ; l.peek != indicator && err == nil; err = l.readch() {
			if l.peek == '\n' {
				return &LexicalError{
					"string literal not terminated",
					l.line,
					l.pos,
				}
			}
			if l.peek == '\\' {
				err = l.readch()
				if err != nil {
					return &LexicalError{
						"invalid escape sequence",
						l.line,
						l.pos,
					}
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
					if indicator != '"' {
						char = '"'
					} else {
						return &LexicalErrorChar{
							LexicalError{"invalid escape sequence in literal",
								l.line,
								l.pos},
							l.peek,
						}
					}
				case '\'':
					if indicator != '\'' {
						char = '\''
					} else {
						return &LexicalErrorChar{
							LexicalError{"invalid escape sequence in literal",
								l.line,
								l.pos},
							l.peek,
						}
					}
				case 'x':
					hex := ""
					for i := 0; i < 2; i++ {
						err = l.readch()
						if err != nil {
							return &LexicalError{
								"invalid escape sequence",
								l.line,
								l.pos,
							}
						}
						// transform uppercase hex in lowercase for lookup
						if l.peek > 64 && l.peek < 91 {
							l.peek = l.peek + 32
						}
						hex = hex + string(l.peek)
					}
					hexLookup, ok := language.EscapeHexaLiterals.Get(hex)
					if !ok {
						return &LexicalErrorChar{
							LexicalError{"invalid hexadecimal literal. must contain two digits between 00-FF",
								l.line,
								l.pos},
							l.peek,
						}
					}
					char = hexLookup.(rune)
				case '0', '1', '2', '3':
					oct := string(l.peek)
					for i := 0; i < 2; i++ {
						err = l.readch()
						if err != nil {
							return &LexicalError{
								"invalid escape sequence",
								l.line,
								l.pos,
							}
						}
						oct = oct + string(l.peek)
					}
					octLookup, ok := language.EscapeOctalLiterals.Get(oct)
					if !ok {
						return &LexicalErrorChar{
							LexicalError{"invalid octal literal. must contain three digits between 000-377",
								l.line,
								l.pos},
							l.peek,
						}
					}
					char = octLookup.(rune)
				default:
					return &LexicalError{
						"invalid escape sequence",
						l.line,
						l.pos,
					}
				}
				l.peek = 0
			} else {
				char = l.peek
			}
			literal = literal + string(char)
		}
		if err != nil {
			return &LexicalError{
				"string literal not terminated",
				l.line,
				l.pos,
			}
		}
		l.tokenStream.Add(tokens.NewLiteral(literal), l.line)
		l.tokenStream.Add(tokens.NewToken(int(l.peek)), l.line)
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

func (l *lexer) Scan() (ts *dataStructs.TokenStream, e error) {
	for err := *new(error); ; err = l.readch() {
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		l.tokenStream.Add(tokens.NewToken(int(l.peek)), l.line)
	}
	return l.tokenStream, nil
}
