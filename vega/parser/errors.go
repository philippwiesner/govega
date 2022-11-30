package parser

import (
	"fmt"

	"govega/vega/scanner"
	"govega/vega/token"
)

const (
	unexpectedEOF scanner.VErrorType = "UnexpectedEOF"
	invalidSyntax scanner.VErrorType = "InvalidSyntax"
)

type VParserError struct {
	scanner.VScanError
	lineFeed string
	token    token.Token
}

func newParserSyntaxErrorObject(etype scanner.VErrorType, token *scanner.Token, msg string, lineFeed string) *VParserError {
	return &VParserError{
		VScanError: scanner.VScanError{
			VError: scanner.VError{
				Location: token.Location,
				Class:    scanner.SyntaxError,
				Type:     etype,
			},
			Message: msg,
		},
		lineFeed: lineFeed,
		token:    token.Token,
	}
}

func newParserSyntaxError(etype scanner.VErrorType, token *scanner.Token, msg string, line string) scanner.Verror {
	var err scanner.Verror = newParserSyntaxErrorObject(etype, token, msg, line)
	return err
}

func (v *VParserError) String() string {
	errString := fmt.Sprintf(`Error in: %v
%v -> %v: at line %v position %v
%v
%v
`, v.GetFile(), v.Class, v.Type, v.Line, v.Position, v.lineFeed, v.Message)
	return errString
}

func (v *VParserError) Error() string {
	return v.String()
}

func GetVErrorType(err error) scanner.VErrorType {
	switch e := err.(type) {
	case scanner.Verror:
		return e.GetErrorType()
	default:
		return "GoError"
	}
}
