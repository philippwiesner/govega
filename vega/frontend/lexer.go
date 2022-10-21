// Package frontend
//
// The compilers' frontend is responsible for validating and analysing the source programming language.
//
// lexer.go implements a lexical scanner which reads the source code character by character and tries to create language
// tokens which can later be analysed for correct syntax with the parser.
package frontend

import (
	"bytes"
	"io"

	"govega/vega/helper"
	"govega/vega/language"
	"govega/vega/language/tokens"
)

// tokenLocation stores line and position of token occurrence
type tokenLocation struct {
	line     int
	position int
}

// lexicalToken is a token wrapper which also contains information where the token origins in the code
type lexicalToken struct {
	token tokens.IToken
	tokenLocation
}

// GetToken getter method for the tokens.IToken interface
func (t *lexicalToken) GetToken() tokens.IToken {
	return t.token
}

func (t *lexicalToken) GetLocation() (line int, position int) {
	return t.line, t.position
}

// GetTag getter method to retrieve the token tag
func (t *lexicalToken) GetTag() int {
	return t.token.GetTag()
}

// lexer implements the lexer object
type lexer struct {
	*vega
	peek     rune             // holds the current scanned character
	code     *bytes.Reader    // code to be analysed in memory
	words    helper.HashTable // collection of keywords and identifiers
	lineFeed string
	line     int
	position int
	eof      bool
}

// NewLexer creates a new lexer object
func (v *vega) NewLexer(code []byte) Lexer {
	var lexer Lexer = &lexer{
		vega:     v,
		peek:     0,
		code:     bytes.NewReader(code),
		words:    language.KeyWords,
		lineFeed: "",
		line:     1,
		position: 0,
	}
	return lexer
}

func (l *lexer) getLineFeed() string {
	return l.lineFeed
}

// newLexicalToken creates a new lexical token
// As the position is calculated during lexer we need to substract the token length to get the starting location
func (l *lexer) newLexicalToken(token tokens.IToken) *lexicalToken {
	var tokenPosition int
	if l.position-len(token.String()) < 0 {
		tokenPosition = l.position
	} else {
		tokenPosition = l.position - len(token.String())
	}
	loc := tokenLocation{line: l.line, position: tokenPosition}
	return &lexicalToken{token: token, tokenLocation: loc}
}

// unreadch private method to put the last read character back on the code stream (revert previous readch)
func (l *lexer) unreadch() error {
	if !l.eof {
		err := l.code.UnreadRune()
		if err != nil {
			return err
		}
		l.lineFeed = l.lineFeed[:len(l.lineFeed)-1]
		l.position--
	}
	l.peek = 0
	return nil
}

// readch private method to retrieve one character from code stream.
// update errorState for each new read character.
func (l *lexer) readch() error {
	ch, _, err := l.code.ReadRune()
	if err != nil {
		if err != io.EOF {
			l.codeLines = append(l.codeLines, l.lineFeed)
			vErr := l.newLexicalSyntaxError(malformedCode, l.line, l.position, "Error parsing file")
			return vErr
		} else {
			l.eof = true
			l.peek = ch
			return nil
		}
	}
	l.lineFeed += string(ch)
	l.position += 1
	l.peek = ch
	return nil
}

// readcch private method to read one character ahead and return true if it matches given character
func (l *lexer) readcch(char rune) (bool, error) {
	err := l.readch()
	if l.peek != char {
		return false, err
	}
	l.peek = 0
	return true, err
}

// scanCombinedTokens private method to scan tokens which are combined from two single tokes, e.g. == or !=
func (l *lexer) scanCombinedTokens(fch rune, sch rune, word tokens.IWord) (*lexicalToken, error) {
	ok, err := l.readcch(sch)
	if err != nil {
		return nil, err
	}
	if ok {
		return l.newLexicalToken(word), nil
	} else {
		return l.newLexicalToken(tokens.NewToken(int(fch))), nil
	}
}

// scanLiterals private method to scan all types of literals.
func (l *lexer) scanLiterals(indicator rune) (*lexicalToken, error) {
	var (
		literal string
		char    rune
		err     error
	)
	literal = string(l.peek)
	err = l.readch()
	for ; l.peek != indicator && err == nil; err = l.readch() {
		if l.peek == '\n' || l.peek == 0 {
			l.codeLines = append(l.codeLines, l.lineFeed)
			vErr := l.newLexicalSyntaxError(literalNotTerminated, l.line, l.position, "String literal not terminated")
			return nil, vErr
		}
		if l.peek == '\\' {
			err = l.readch()
			if err != nil {
				l.codeLines = append(l.codeLines, l.lineFeed)
				vErr := l.newLexicalSyntaxError(invalidEscapeSequence, l.line, l.position, "Invalid escape sequence")
				return nil, vErr
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
					l.codeLines = append(l.codeLines, l.lineFeed)
					vErr := l.newLexicalSyntaxError(invalidEscapeSequenceLiteral, l.line, l.position, "Invalid escape sequence in literal")
					return nil, vErr
				}
			case '\'':
				if indicator == '\'' {
					char = '\''
				} else {
					l.codeLines = append(l.codeLines, l.lineFeed)
					vErr := l.newLexicalSyntaxError(invalidEscapeSequenceLiteral, l.line, l.position, "Invalid escape sequence in literal")
					return nil, vErr
				}
			case 'x':
				hex := ""
				for i := 0; i < 2; i++ {
					err = l.readch()
					if err != nil {
						l.codeLines = append(l.codeLines, l.lineFeed)
						vErr := l.newLexicalSyntaxError(invalidEscapeSequenceHexadecimal, l.line, l.position, "Invalid hexadecimal literal. Must contain two digits between 00-FF")
						return nil, vErr
					}
					// transform uppercase hex in lowercase for lookup
					if l.peek > 64 && l.peek < 91 {
						l.peek = l.peek + 32
					}
					hex = hex + string(l.peek)
				}
				hexLookup, ok := language.EscapeHexaLiterals.Get(hex)
				if !ok {
					l.codeLines = append(l.codeLines, l.lineFeed)
					vErr := l.newLexicalSyntaxError(invalidEscapeSequenceHexadecimal, l.line, l.position, "Invalid hexadecimal literal. Must contain two digits between 00-FF")
					return nil, vErr
				}
				char = hexLookup.(rune)
			case 'u':
				unicode := ""
				for i := 0; i < 4; i++ {
					err = l.readch()
					if err != nil {
						l.codeLines = append(l.codeLines, l.lineFeed)
						vErr := l.newLexicalSyntaxError(invalidEscapeSequenceUnicode, l.line, l.position, "Invalid unicode literal. Must contain four digits between 0000-FFFF")
						return nil, vErr
					}
					if l.peek > 64 && l.peek < 91 {
						l.peek = l.peek + 32
					}
					unicode = unicode + string(l.peek)
				}
				unicodeLookup, ok := language.EscapeUnicodeLiterals.Get(unicode)
				if !ok {
					l.codeLines = append(l.codeLines, l.lineFeed)
					vErr := l.newLexicalSyntaxError(invalidEscapeSequenceUnicode, l.line, l.position, "Invalid unicode literal. Must contain four digits between 0000-FFFF")
					return nil, vErr
				}
				char = unicodeLookup.(rune)
			case '0', '1', '2', '3':
				oct := string(l.peek)
				for i := 0; i < 2; i++ {
					err = l.readch()
					if err != nil {
						l.codeLines = append(l.codeLines, l.lineFeed)
						vErr := l.newLexicalSyntaxError(invalidEscapeSequenceOctal, l.line, l.position, "Invalid octal literal. Must contain three digits between 000-377")
						return nil, vErr
					}
					oct = oct + string(l.peek)
				}
				octLookup, ok := language.EscapeOctalLiterals.Get(oct)
				if !ok {
					l.codeLines = append(l.codeLines, l.lineFeed)
					vErr := l.newLexicalSyntaxError(invalidEscapeSequenceOctal, l.line, l.position, "Invalid octal literal. Must contain three digits between 000-377")
					return nil, vErr
				}
				char = octLookup.(rune)
			default:
				l.codeLines = append(l.codeLines, l.lineFeed)
				vErr := l.newLexicalSyntaxError(invalidEscapeSequence, l.line, l.position, "Invalid escape sequence")
				return nil, vErr
			}
		} else {
			char = l.peek
		}
		literal = literal + string(char)
	}
	if err != nil {
		l.codeLines = append(l.codeLines, l.lineFeed)
		vErr := l.newLexicalSyntaxError(literalNotTerminated, l.line, l.position, "String literal not terminated")
		return nil, vErr
	}
	literal = literal + string(l.peek)
	return l.newLexicalToken(tokens.NewLiteral(literal)), nil
}

// scanNumbers private method to scan integer and floating point numbers
func (l *lexer) scanNumbers() (*lexicalToken, error) {
	var (
		err   error
		value int
	)
	for ; l.peek > 47 && l.peek < 58 && err == nil; err = l.readch() {
		value = 10*value + (int(l.peek) - '0')
	}
	if err != nil {
		return nil, err
	}
	if l.peek != '.' {
		return l.newLexicalToken(tokens.NewNum(value)), nil
	} else {
		realNumber := float64(value)
		fraction := float64(10)
		err = l.readch()
		for ; l.peek > 47 && l.peek < 58 && err == nil; err = l.readch() {
			realNumber = realNumber + (float64(l.peek)-'0')/fraction
			fraction *= 10
		}
		if err != nil {
			return nil, err
		}
		return l.newLexicalToken(tokens.NewReal(realNumber)), nil
	}
}

// scanWords private method to scan keywords and identifiers. new identifier are registered in words hashtable to be
// recognized later
func (l *lexer) scanWords() (*lexicalToken, error) {
	var (
		word string
		err  error
	)
	for ; ((l.peek > 64 && l.peek < 91) || (l.peek > 96 && l.peek < 123) || (l.peek > 47 && l.peek < 58)) && err == nil; err = l.readch() {
		word += string(l.peek)
	}
	if err != nil {
		return nil, err
	}
	lookup, ok := l.words.Get(word)
	if ok {
		return l.newLexicalToken(lookup.(tokens.IWord)), nil
	}
	identifier := tokens.NewWord(word, tokens.ID)
	l.words.Add(word, identifier)
	return l.newLexicalToken(identifier), nil
}

// scanComments private method which skips all single and multi-line comments
func (l *lexer) scanComments() (*lexicalToken, error) {
	var (
		err error
		ok  bool
	)
	if err = l.readch(); err != nil {
		return nil, err
	}
	if l.peek == '/' {
		for ; l.peek != '\n' && l.peek != 0 && err == nil; err = l.readch() {
		}
		if err != nil {
			return nil, err
		}
		if err = l.unreadch(); err != nil {
			return nil, err
		}
	} else if l.peek == '*' {
		for ; err == nil; err = l.readch() {
			if l.peek == '\n' {
				l.line++
			} else if l.peek == '*' {
				ok, err = l.readcch('/')
				if ok || err != nil {
					break
				}
			}
		}
		if err != nil {
			return nil, err
		}
	} else {
		return l.newLexicalToken(tokens.NewToken('/')), nil
	}
	return nil, nil
}

// scan public method to scan the actual source code and return a tokenStream with all scanned tokens
func (l *lexer) scan() (*lexicalToken, error) {
	err := l.readch()
	for ; err == nil; err = l.readch() {
		if l.peek == 0 {
			return l.newLexicalToken(tokens.NewToken(tokens.EOF)), nil
		}
		switch {
		// skip line breaks
		case l.peek == '\n':
			l.codeLines = append(l.codeLines, l.lineFeed)
			l.lineFeed = ""
			l.position = 0
			l.line++
		// skip comments
		case l.peek == '/':
			var token *lexicalToken
			token, err = l.scanComments()
			if token != nil {
				return token, err
			}
		// skip whitespaces
		case l.peek == ' ', l.peek == '\t', l.peek == '\v', l.peek == '\r':
		// read token ! or !=
		case l.peek == '!':
			return l.scanCombinedTokens(l.peek, '=', language.Ne)
		// read token = or ==
		case l.peek == '=':
			return l.scanCombinedTokens(l.peek, '=', language.Eq)
		// read token < or <=
		case l.peek == '<':
			return l.scanCombinedTokens(l.peek, '=', language.Le)
		// read token > or >=
		case l.peek == '>':
			return l.scanCombinedTokens(l.peek, '=', language.Ge)
		// read token & or &&
		case l.peek == '&':
			return l.scanCombinedTokens(l.peek, '&', language.BoolAnd)
		// read token | or ||
		case l.peek == '|':
			return l.scanCombinedTokens(l.peek, '|', language.BoolOr)
		// read literals encapsulated in '' or ""
		case l.peek == '\'', l.peek == '"':
			return l.scanLiterals(l.peek)
		// read numbers
		case l.peek > 47 && l.peek < 58:
			tok, err := l.scanNumbers()
			if err != nil {
				return nil, err
			}
			err = l.unreadch()
			return tok, err
		// read words
		case (l.peek > 64 && l.peek < 91) || (l.peek > 96 && l.peek < 123):
			tok, err := l.scanWords()
			if err != nil {
				return nil, err
			}
			err = l.unreadch()
			return tok, err
		// read every token left
		default:
			return l.newLexicalToken(tokens.NewToken(int(l.peek))), nil
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, err
}
