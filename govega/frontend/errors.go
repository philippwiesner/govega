package frontend

import "fmt"

const (
	malformedCode                = "error parsing file, malformed"
	unexpectedEOF                = "unexpected EOF"
	literalNotTerminated         = "string literal not terminated"
	invalidEscapeSequence        = "invalid escape squence"
	invalidEscapeSequenceLiteral = "invalid escape sequence in literal"
	invalidEscapeSequenceHex     = "invalid hexadecimal literal. must contain two digits between 00-FF"
	invalidEscapeSequenceOct     = "invalid octal literal. must contain three digits between 000-377"
	invalidEscapeSequenceUnicode = "invalid unicode literal. must contain four digits between 0000-FFFF"
	invalidSyntax                = "invalid Syntax"
	alreadyDefined               = "identifier already defined"
)

type ErrorState struct {
	lineNumber int
	position   int
	fileName   string
	lineFeed   string
}

type LexicalError struct {
	msg string
	ErrorState
}

func NewLexicalError(msg string, state ErrorState) *LexicalError {
	return &LexicalError{msg, state}
}

func (l *LexicalError) Error() string {
	return fmt.Sprintf("%v:%d:%d: %v:\n\n\t%v", l.fileName, l.lineNumber, l.position, l.msg, l.lineFeed)
}

func (l LexicalError) String() string {
	return l.Error()
}

type SyntaxError struct {
	msg string
	ErrorState
}

func NewSyntaxError(msg string, state ErrorState) *SyntaxError {
	return &SyntaxError{msg, state}
}

func (s *SyntaxError) Error() string {
	return fmt.Sprintf("%v:%d: Syntax Error: %v", s.fileName, s.lineNumber, s.msg)
}

func (s *SyntaxError) String() string {
	return s.Error()
}

type DeclarationError struct {
	msg string
	ErrorState
}

func NewDeclarationError(msg string, state ErrorState) *DeclarationError {
	return &DeclarationError{msg, state}
}

func (s *DeclarationError) Error() string {
	return fmt.Sprintf("%v:%d: Declaration Error: %v", s.fileName, s.lineNumber, s.msg)
}

func (s *DeclarationError) String() string {
	return s.Error()
}
