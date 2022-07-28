// Package frontend
//
// The compilers' frontend is responsible for validating and analysing the source programming language.
//
// lexer.go implements a lexical scanner which reads the source code character by character and tries to create language
// tokens which can later be analysed for correct syntax with the parser.
package frontend

import (
	"bytes"
	"govega/govega/helper"
	"govega/govega/language"
	"govega/govega/language/tokens"
	"io"
)

// lexer implements the lexer object
type lexer struct {
	errorState  ErrorState        // gives information when an error occures during lexical scanning
	peek        rune              // holds the current scanned character
	code        *bytes.Reader     // code to be analysed in memory
	tokenStream *TokenStream      // stream of tokens returned by lexer
	words       *helper.HashTable // collection of keywords and identifiers
}

// NewLexer creates a new lexer object
func NewLexer(code []byte, fileName string) *lexer {
	return &lexer{ErrorState{1, 0, fileName, ""},
		0,
		bytes.NewReader(code),
		NewTokenStream(),
		language.KeyWords,
	}
}

// unreadch private method to put the last read character back on the code stream (revert previous readch)
func (l *lexer) unreadch() error {
	err := l.code.UnreadRune()
	if err != nil {
		return err
	}
	l.errorState.position--
	l.errorState.lineFeed = l.errorState.lineFeed[:len(l.errorState.lineFeed)-1]
	l.peek = 0
	return nil
}

// readch private method to retrieve one character from code stream.
// update errorState for each new read character.
func (l *lexer) readch() error {
	ch, _, err := l.code.ReadRune()
	l.errorState.position++
	l.errorState.lineFeed += string(ch)
	if err != nil {
		if err != io.EOF {
			return NewLexicalError(malformedCode, l.errorState)
		} else {
			l.peek = ch
			return nil
		}
	}
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
func (l *lexer) scanCombinedTokens(fch rune, sch rune, word tokens.IWord) error {
	ok, err := l.readcch(sch)
	if err != nil {
		return err
	}
	if ok {
		l.tokenStream.Add(word, l.errorState)
	} else {
		l.tokenStream.Add(tokens.NewToken(int(fch)), l.errorState)
	}
	return nil
}

// scanLiterals private method to scan all types of literals. Currently, supports only ASCII.
func (l *lexer) scanLiterals(indicator rune) error {
	var (
		literal string
		char    rune
		err     error
	)
	l.tokenStream.Add(tokens.NewToken(int(indicator)), l.errorState)
	err = l.readch()
	for ; l.peek != indicator && err == nil; err = l.readch() {
		if l.peek == '\n' || l.peek == 0 {
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
			case 'u':
				unicode := ""
				for i := 0; i < 4; i++ {
					err = l.readch()
					if err != nil {
						return NewLexicalError(invalidEscapeSequenceUnicode, l.errorState)
					}
					if l.peek > 64 && l.peek < 91 {
						l.peek = l.peek + 32
					}
					unicode = unicode + string(l.peek)
				}
				unicodeLookup, ok := language.EscapeUnicodeLiterals.Get(unicode)
				if !ok {
					return NewLexicalError(invalidEscapeSequenceUnicode, l.errorState)
				}
				char = unicodeLookup.(rune)
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
	return nil
}

// scanNumbers private method to scan integer and floating point numbers
func (l *lexer) scanNumbers() error {
	var (
		err   error
		value int
	)
	for ; l.peek > 47 && l.peek < 58 && err == nil; err = l.readch() {
		value = 10*value + (int(l.peek) - '0')
	}
	if err != nil {
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
			realNumber = realNumber + (float64(l.peek)-'0')/fraction
			fraction *= 10
		}
		if err != nil {
			return err
		}
		l.tokenStream.Add(tokens.NewReal(realNumber), l.errorState)
	}
	return nil
}

// scanWords private method to scan keywords and identifiers. new identifier are registered in words hashtable to be
// recognized later
func (l *lexer) scanWords() error {
	var (
		word string
		err  error
	)
	for ; ((l.peek > 64 && l.peek < 91) || (l.peek > 96 && l.peek < 123) || (l.peek > 47 && l.peek < 58)) && err == nil; err = l.readch() {
		word += string(l.peek)
	}
	if err != nil {
		return err
	}
	lookup, ok := l.words.Get(word)
	if ok {
		l.tokenStream.Add(lookup.(tokens.IWord), l.errorState)
		return nil
	}
	identifier := tokens.NewWord(word, tokens.ID)
	l.words.Add(word, identifier)
	l.tokenStream.Add(identifier, l.errorState)
	return nil
}

// scanComments private method which skips all single and multi-line comments
func (l *lexer) scanComments() error {
	var (
		err error
		ok  bool
	)
	if err = l.readch(); err != nil {
		return err
	}
	if l.peek == '/' {
		for ; l.peek != '\n' && l.peek != 0 && err == nil; err = l.readch() {
		}
		if err != nil {
			return err
		}
		if err = l.unreadch(); err != nil {
			return err
		}
	} else if l.peek == '*' {
		for ; err == nil; err = l.readch() {
			if l.peek == '\n' {
				l.errorState.lineNumber++
			} else if l.peek == '*' {
				ok, err = l.readcch('/')
				if ok || err != nil {
					break
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Scan public method to scan the actual source code and return a tokenStream with all scanned tokens
func (l *lexer) Scan() (ts *TokenStream, e error) {
	err := l.readch()
	for ; err == nil; err = l.readch() {
		if l.peek == 0 {
			l.tokenStream.Add(tokens.NewToken(tokens.EOF), l.errorState)
			break
		}
		switch {
		// skip line breaks
		case l.peek == '\n':
			l.errorState.lineFeed = ""
			l.errorState.position = 0
			l.errorState.lineNumber++
		// skip comments
		case l.peek == '/':
			err = l.scanComments()
		// skip whitespaces
		case l.peek == ' ', l.peek == '\t', l.peek == '\v', l.peek == '\r':
		// read token ! or !=
		case l.peek == '!':
			err = l.scanCombinedTokens(l.peek, '=', language.Ne)
		// read token ! or !=
		case l.peek == '=':
			err = l.scanCombinedTokens(l.peek, '=', language.Eq)
		// read token < or <=
		case l.peek == '<':
			err = l.scanCombinedTokens(l.peek, '=', language.Le)
		// read token > or >=
		case l.peek == '>':
			err = l.scanCombinedTokens(l.peek, '=', language.Ge)
		// read token & or &&
		case l.peek == '&':
			err = l.scanCombinedTokens(l.peek, '&', language.BoolAnd)
		// read token | or ||
		case l.peek == '|':
			err = l.scanCombinedTokens(l.peek, '|', language.BoolOr)
		// read token ->
		case l.peek == '-':
			err = l.scanCombinedTokens(l.peek, '>', language.ReturnType)
		// read literals encapsulated in '' or ""
		case l.peek == '\'', l.peek == '"':
			err = l.scanLiterals(l.peek)
		// read numbers
		case l.peek > 47 && l.peek < 58:
			if err = l.scanNumbers(); err != nil {
				return nil, err
			}
			err = l.unreadch()
		// read words
		case (l.peek > 64 && l.peek < 91) || (l.peek > 96 && l.peek < 123):
			if err = l.scanWords(); err != nil {
				return nil, err
			}
			err = l.unreadch()
		// read every token left
		default:
			l.tokenStream.Add(tokens.NewToken(int(l.peek)), l.errorState)
		}
	}
	if err != nil {
		return nil, err
	}
	return l.tokenStream, nil
}
