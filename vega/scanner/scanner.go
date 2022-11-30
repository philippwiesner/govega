// Package scanner
//
// The compilers' frontend is responsible for validating and analysing the source programming language.
//
// Scanner.go implements a lexical Scanner which reads the source source character by character and tries to create language
// token which can later be analysed for correct syntax with the parser.
package scanner

import (
	"bytes"
	"io"
	"math"

	"govega/vega/token"
)

// Token is a token wrapper which also contains information where the token origins in the source
type Token struct {
	token.Token
	token.Location
	Literal string
}

// Scanner implements the Scanner object
type Scanner struct {
	file     token.File
	peek     rune                   // holds the current scanned character
	source   *bytes.Reader          // source to be analysed in memory
	words    map[string]token.Token // collection of keywords and identifiers
	LineFeed string
	line     int
	position int
	eof      bool
}

// NewScanner creates a new Scanner object
func NewScanner(file token.File) *Scanner {
	return &Scanner{
		file:     file,
		peek:     0,
		source:   bytes.NewReader(file.Source),
		words:    token.Keywords,
		line:     1,
		position: 0,
	}
}

func (s *Scanner) getTokenLocation(pos int) token.Location {
	return token.Location{
		FileName: s.file.Name,
		Line:     s.line,
		Position: pos,
		LineFeed: s.LineFeed,
	}
}

// newScannerToken creates a new lexical token
// As the position is calculated during Scanner we need to substract the token length to get the starting location
func (s *Scanner) newScannerToken(tok token.Token, literal string) *Token {
	var (
		tokenPosition int
		word          string
	)
	if s.position-len(tok.String()) < 0 {
		tokenPosition = s.position
	} else {
		tokenPosition = s.position - len(tok.String())
	}
	if tok.IsLiteral() {
		word = literal
	} else {
		word = tok.String()
	}
	return &Token{Token: tok, Literal: word, Location: s.getTokenLocation(tokenPosition)}
}

// unreadch private method to put the last read character back on the source stream (revert previous readch)
func (s *Scanner) unreadch() error {
	if !s.eof {
		if err := s.source.UnreadRune(); err != nil {
			return err
		}
		s.LineFeed = s.LineFeed[:len(s.LineFeed)-1]
		s.position--
	}
	s.peek = 0
	return nil
}

// readch private method to retrieve one character from source stream.
// update errorState for each new read character.
func (s *Scanner) readch() error {
	ch, _, err := s.source.ReadRune()
	if err != nil {
		if err != io.EOF {
			vErr := newScannerSyntaxError(malformedCode, s.getTokenLocation(s.position), "Error parsing file")
			return vErr
		} else {
			s.eof = true
			s.peek = ch
			return nil
		}
	}
	s.LineFeed += string(ch)
	s.position += 1
	s.peek = ch
	return nil
}

// readcch private method to read one character ahead and return true if it matches given character
func (s *Scanner) readcch(char rune) (bool, error) {
	err := s.readch()
	if s.peek != char {
		return false, err
	}
	s.peek = 0
	return true, err
}

// scanCombinedTokens private method to scan token which are combined from two single tokes, e.g. == or !=
func (s *Scanner) scanCombinedTokens(sch rune, tag token.Token) (*Token, error) {
	fch := s.peek
	ok, err := s.readcch(sch)
	if err != nil {
		return nil, err
	}
	if ok {
		return s.newScannerToken(tag, ""), nil
	} else {
		var tok token.Token
		switch fch {
		case '!':
			tok = token.LNOT
		case '=':
			tok = token.ASSIGN
		case '<':
			tok = token.LESS
		case '>':
			tok = token.GREATER
		case '|':
			tok = token.BIT_OR
		case '&':
			tok = token.BIT_AND
		}
		return s.newScannerToken(tok, ""), nil
	}
}

func toLower(sign rune) rune { return sign + ('a' - 'A') }
func isUpper(sign rune) bool { return sign >= 'A' && sign <= 'Z' }
func isHex(sign rune) bool   { return (sign >= '0' && sign <= '9') || (sign >= 'a' && sign <= 'f') }
func isOct(sign rune) bool   { return sign >= '0' && sign <= '7' }
func isDigit(sign rune) bool { return sign >= '0' && sign <= '9' }
func hexToDec(sign rune) rune {
	if isDigit(sign) {
		sign -= '0'
	} else {
		sign -= 'a'
		sign += 10
	}
	return sign
}

func (s *Scanner) scanEncoded() (rune, error) {
	var (
		sign rune
		err  error
	)
	switch s.peek {
	case 'x':
		for i := 1; i >= 0; i-- {
			if err = s.readch(); err != nil {
				return sign, err
			}
			if isUpper(s.peek) {
				s.peek = toLower(s.peek)
			}
			if !isHex(s.peek) {
				err = newScannerSyntaxError(invalidEscapeSequenceHexadecimal, s.getTokenLocation(s.position), "Invalid hexadecimal literal. Must contain two digits between 00-FF")
			}
			s.peek = hexToDec(s.peek)
			sign += s.peek * int32(math.Pow(16, float64(i)))
		}
	case 'u':
		for i := 3; i >= 0; i-- {
			if err = s.readch(); err != nil {
				return sign, err
			}
			if isUpper(s.peek) {
				s.peek = toLower(s.peek)
			}
			if !isHex(s.peek) {
				err = newScannerSyntaxError(invalidEscapeSequenceUnicode, s.getTokenLocation(s.position), "Invalid unicode literal. Must contain four digits between 0000-FFFF")
			}
			s.peek = hexToDec(s.peek)
			sign += s.peek * int32(math.Pow(16, float64(i)))
		}
	case 'o':
		if err = s.readch(); err != nil {
			return sign, err
		}
		if s.peek >= '0' && s.peek <= '3' {
			s.peek = hexToDec(s.peek)
			sign = s.peek * int32(math.Pow(8, 2))
			for i := 1; i >= 0; i-- {
				if err = s.readch(); err != nil {
					return sign, err
				}
				if !isOct(s.peek) {
					err = newScannerSyntaxError(invalidEscapeSequenceOctal, s.getTokenLocation(s.position), "Invalid octal literal. Must contain three digits between 000-377")
				}
				s.peek = hexToDec(s.peek)
				sign += s.peek * int32(math.Pow(8, float64(i)))
			}
		} else {
			err = newScannerSyntaxError(invalidEscapeSequenceOctal, s.getTokenLocation(s.position), "Invalid octal literal. Must contain three digits between 000-377")
		}
	}
	return sign, err
}

func (s *Scanner) scanSigns(tok token.Token) (rune, error) {
	var (
		sign rune
		err  error
	)
	switch s.peek {
	case '\n', 0:
		err = newScannerSyntaxError(literalNotTerminated, s.getTokenLocation(s.position), "String literal not terminated")
	case '\\':
		err = s.readch()
		if err != nil {
			return sign, err
		}
		switch s.peek {
		case 'b':
			sign = '\b'
		case 'f':
			sign = '\f'
		case 'n':
			sign = '\n'
		case 'r':
			sign = '\r'
		case 't':
			sign = '\t'
		case 'v':
			sign = '\v'
		case '\\':
			sign = '\\'
		case '"':
			if tok == token.STRING {
				sign = '"'
			} else {
				err = newScannerSyntaxError(invalidEscapeSequence, s.getTokenLocation(s.position), "Invalid escape sequence in literal")
			}
		case '\'':
			if tok == token.CHAR {
				sign = '\''
			} else {
				err = newScannerSyntaxError(invalidEscapeSequence, s.getTokenLocation(s.position), "Invalid escape sequence in literal")
			}
		case 'x', 'u', 'o':
			sign, err = s.scanEncoded()
		default:
			err = newScannerSyntaxError(invalidEscapeSequence, s.getTokenLocation(s.position), "Invalid escape sequence in literal")
		}
	default:
		sign = s.peek
	}
	return sign, err
}

func (s *Scanner) scanChar() (*Token, error) {
	var (
		literal string
	)
	literal += string(s.peek)
	if err := s.readch(); err != nil {
		return nil, err
	}
	if sign, err := s.scanSigns(token.CHAR); err != nil {
		return nil, err
	} else {
		literal += string(sign)
	}
	if err := s.readch(); err != nil {
		return nil, err
	}
	if s.peek != '\'' {
		vErr := newScannerSyntaxError(literalNotTerminated, s.getTokenLocation(s.position), "Char can only have one character")
		return nil, vErr
	}
	literal += string(s.peek)
	return s.newScannerToken(token.CHAR, literal), nil
}

func (s *Scanner) scanString() (*Token, error) {
	var (
		literal string
		err     error
	)
	literal += string(s.peek)
	if err = s.readch(); err != nil {
		return nil, err
	}
	for ; s.peek != '"' && err == nil; err = s.readch() {
		if sign, err := s.scanSigns(token.STRING); err != nil {
			return nil, err
		} else {
			literal += string(sign)
		}
	}
	if err != nil {
		return nil, err
	}
	literal += string(s.peek)
	return s.newScannerToken(token.STRING, literal), err
}

// scanNumbers private method to scan integer and floating point numbers
func (s *Scanner) scanNumbers() (*Token, error) {
	var (
		err   error
		digit string
	)
	for ; isDigit(s.peek) && err == nil; err = s.readch() {
		digit += string(s.peek)
	}
	if err != nil {
		return nil, err
	}
	if s.peek != '.' {
		return s.newScannerToken(token.INT, digit), nil
	} else {
		digit += string(s.peek)
		err = s.readch()
		for ; isDigit(s.peek) && err == nil; err = s.readch() {
			digit += string(s.peek)
		}
		if err != nil {
			return nil, err
		}
		return s.newScannerToken(token.FLOAT, digit), nil
	}
}

// scanWords private method to scan keywords and identifiers. new identifier are registered in words hashtable to be
// recognized later
func (s *Scanner) scanWords() (*Token, error) {
	var (
		word string
		err  error
	)
	for ; ((s.peek > 64 && s.peek < 91) || (s.peek > 96 && s.peek < 123) || (s.peek > 47 && s.peek < 58)) && err == nil; err = s.readch() {
		word += string(s.peek)
	}
	if err != nil {
		return nil, err
	}
	return s.newScannerToken(token.Lookup(word), word), nil
}

// scanComments private method which skips all single and multi-line comments
func (s *Scanner) scanComments() (*Token, error) {
	var (
		err error
		ok  bool
	)
	if err = s.readch(); err != nil {
		return nil, err
	}
	if s.peek == '/' {
		for ; s.peek != '\n' && s.peek != 0 && err == nil; err = s.readch() {
		}
		if err != nil {
			return nil, err
		}
		if err = s.unreadch(); err != nil {
			return nil, err
		}
	} else if s.peek == '*' {
		for ; err == nil; err = s.readch() {
			if s.peek == '\n' {
				s.line++
			} else if s.peek == '*' {
				ok, err = s.readcch('/')
				if ok || err != nil {
					break
				}
			}
		}
		if err != nil {
			return nil, err
		}
	} else {
		return s.newScannerToken(token.DIV, ""), nil
	}
	return nil, nil
}

// Scan public method to scan the actual source and return a tokenStream with all scanned token
func (s *Scanner) Scan() (*Token, error) {
	err := s.readch()
	for ; err == nil; err = s.readch() {
		if s.peek == 0 {
			return s.newScannerToken(token.EOF, ""), nil
		}
		switch {
		// skip line breaks
		case s.peek == '\n':
			tok := s.newScannerToken(token.LINEBREAK, "")
			s.LineFeed = ""
			s.position = 0
			s.line++
			return tok, nil
		// skip comments
		case s.peek == '/':
			var tok *Token
			tok, err = s.scanComments()
			if tok != nil {
				return tok, err
			}
		// skip whitespaces
		case s.peek == ' ', s.peek == '\t', s.peek == '\v', s.peek == '\r':
		// read token ! or !=
		case s.peek == '!':
			return s.scanCombinedTokens('=', token.NE)
		// read token = or ==
		case s.peek == '=':
			return s.scanCombinedTokens('=', token.EQ)
		// read token < or <=
		case s.peek == '<':
			return s.scanCombinedTokens('=', token.LE)
		// read token > or >=
		case s.peek == '>':
			return s.scanCombinedTokens('=', token.GE)
		// read token & or &&
		case s.peek == '&':
			return s.scanCombinedTokens('&', token.LAND)
		// read token | or ||
		case s.peek == '|':
			return s.scanCombinedTokens('|', token.LOR)
		// read char
		case s.peek == '\'':
			return s.scanChar()
		// read string
		case s.peek == '"':
			return s.scanString()
		// read numbers
		case s.peek > 47 && s.peek < 58:
			tok, err := s.scanNumbers()
			if err != nil {
				return nil, err
			}
			err = s.unreadch()
			return tok, err
		// read words
		case (s.peek > 64 && s.peek < 91) || (s.peek > 96 && s.peek < 123):
			tok, err := s.scanWords()
			if err != nil {
				return nil, err
			}
			err = s.unreadch()
			return tok, err
			// read +
		case s.peek == '+':
			return s.newScannerToken(token.ADD, ""), nil
			// read -
		case s.peek == '-':
			return s.newScannerToken(token.SUB, ""), nil
			// read *
		case s.peek == '*':
			return s.newScannerToken(token.MUL, ""), nil
			// read {
		case s.peek == '{':
			return s.newScannerToken(token.LCBRACKET, ""), nil
			// read }
		case s.peek == '}':
			return s.newScannerToken(token.RCBRACKET, ""), nil
			// read [
		case s.peek == '[':
			return s.newScannerToken(token.LSBRACKET, ""), nil
			// read ]
		case s.peek == ']':
			return s.newScannerToken(token.RSBRACKET, ""), nil
			// read (
		case s.peek == '(':
			return s.newScannerToken(token.LBRACKET, ""), nil
			// read )
		case s.peek == ')':
			return s.newScannerToken(token.RBRACKET, ""), nil
			// read ;
		case s.peek == ';':
			return s.newScannerToken(token.DELIMITER, ""), nil
			// read :
		case s.peek == ':':
			return s.newScannerToken(token.COLON, ""), nil
			// read ,
		case s.peek == ',':
			return s.newScannerToken(token.COMMA, ""), nil
		// token not in alphabet
		default:
			return nil, newScannerSyntaxError(invalidCharacter, s.getTokenLocation(s.position), "Invalid character")
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, err
}
