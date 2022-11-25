package frontend

import (
	"fmt"

	"govega/vega/language/tokens"
)

type VErrorClass string
type VErrorType string

const (
	syntaxError VErrorClass = "SyntaxError"
)

const (
	malformedCode                    VErrorType = "MalformedCode"
	unexpectedEOF                    VErrorType = "UnexpectedEOF"
	literalNotTerminated             VErrorType = "LiteralNotTerminated"
	invalidCharacter                 VErrorType = "InvalidCharacter"
	invalidEscapeSequence            VErrorType = "InvalidEscapeSequence"
	invalidEscapeSequenceLiteral     VErrorType = "InvalidEscapeSequenceLiteral"
	invalidEscapeSequenceHexadecimal VErrorType = "InvalidEscapeSequenceHexadecimal"
	invalidEscapeSequenceOctal       VErrorType = "InvalidEscapeSequenceOctal"
	invalidEscapeSequenceUnicode     VErrorType = "InvalidEscapeSequenceUnicode"
	invalidSyntax                    VErrorType = "InvalidSyntax"
)

type IVError interface {
	error
	GetErrorType() VErrorType
	GetErrorClass() VErrorClass
	GetMessage() string
}

type vError struct {
	*vega
	class     VErrorClass
	errorType VErrorType
	line      int
}

func (v *vError) getErrorLines() string {
	var errorLines string
	if len(v.codeLines) == 0 {
		errorLines = ""
	} else {
		errorLines = v.codeLines[v.line-1]
	}
	return errorLines
}

func (v *vError) getFile() string {
	return fmt.Sprintf("\033[36m%v\033[0m", v.file)
}

func (v *vega) newError(class VErrorClass, etype VErrorType, line int) *vError {
	return &vError{
		vega:      v,
		class:     class,
		errorType: etype,
		line:      line,
	}
}

func (v *vError) GetErrorClass() VErrorClass {
	return v.class
}

func (v *vError) GetErrorType() VErrorType {
	return v.errorType
}

type vLexerError struct {
	vError
	message  string
	position int
}

func (v *vega) newLexicalSyntaxErrorObject(etype VErrorType, line int, pos int, msg string) *vLexerError {
	vegaError := *v.newError(syntaxError, etype, line)
	return &vLexerError{
		vError:   vegaError,
		message:  msg,
		position: pos,
	}
}

func (v *vega) newLexicalSyntaxError(etype VErrorType, line int, pos int, msg string) IVError {
	var vErr IVError = v.newLexicalSyntaxErrorObject(etype, line, pos, msg)
	return vErr
}

func (v *vLexerError) GetMessage() string {
	return v.message
}

func (v *vLexerError) String() string {
	errString := fmt.Sprintf(`Error in: %v
%v -> %v: at line %v position %v
%v
%v
`, v.getFile(), v.class, v.errorType, v.line, v.position, v.getErrorLines(), v.message)
	return errString
}

func (v *vLexerError) Error() string {
	return v.String()
}

type vParserError struct {
	vLexerError
	lineFeed string
	token    tokens.IToken
}

func (v *vega) newParserSyntaxErrorObject(etype VErrorType, token *lexicalToken, msg string, lineFeed string) *vParserError {
	line, position := token.GetLocation()
	return &vParserError{
		vLexerError: *v.newLexicalSyntaxErrorObject(etype, line, position, msg),
		lineFeed:    lineFeed,
		token:       token.GetToken(),
	}
}

func (v *vega) newParserSyntaxError(etype VErrorType, token *lexicalToken, msg string, line string) IVError {
	var vErr IVError = v.newParserSyntaxErrorObject(etype, token, msg, line)
	return vErr
}

func (v *vParserError) String() string {
	errString := fmt.Sprintf(`Error in: %v
%v -> %v: at line %v position %v
%v
%v
`, v.getFile(), v.class, v.errorType, v.line, v.position, v.lineFeed, v.message)
	return errString
}

func (v *vParserError) Error() string {
	return v.String()
}

func GetVErrorType(err error) VErrorType {
	switch e := err.(type) {
	case IVError:
		return e.GetErrorType()
	default:
		return "GoError"
	}
}
