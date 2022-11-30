package scanner

import (
	"fmt"

	"govega/vega/token"
)

type VErrorClass string
type VErrorType string

const (
	SyntaxError VErrorClass = "SyntaxError"
)

const (
	malformedCode                    VErrorType = "MalformedCode"
	literalNotTerminated             VErrorType = "LiteralNotTerminated"
	invalidCharacter                 VErrorType = "InvalidCharacter"
	invalidEscapeSequence            VErrorType = "InvalidEscapeSequence"
	invalidEscapeSequenceHexadecimal VErrorType = "InvalidEscapeSequenceHexadecimal"
	invalidEscapeSequenceOctal       VErrorType = "InvalidEscapeSequenceOctal"
	invalidEscapeSequenceUnicode     VErrorType = "InvalidEscapeSequenceUnicode"
)

type Verror interface {
	error
	GetErrorType() VErrorType
	GetErrorClass() VErrorClass
	GetMessage() string
}

type VError struct {
	token.Location
	Class VErrorClass
	Type  VErrorType
}

func (v *VError) GetFile() string {
	return fmt.Sprintf("\033[36m%v\033[0m", v.FileName)
}

func newError(class VErrorClass, etype VErrorType, loc token.Location) VError {
	return VError{
		Location: loc,
		Class:    class,
		Type:     etype,
	}
}

func (v *VError) GetErrorClass() VErrorClass {
	return v.Class
}

func (v *VError) GetErrorType() VErrorType {
	return v.Type
}

type VScanError struct {
	VError
	Message string
}

func newScannerError(etype VErrorType, loc token.Location, msg string) *VScanError {
	err := newError(SyntaxError, etype, loc)
	return &VScanError{
		VError:  err,
		Message: msg,
	}
}

func newScannerSyntaxError(etype VErrorType, loc token.Location, msg string) Verror {
	var err Verror = newScannerError(etype, loc, msg)
	return err
}

func (v *VScanError) GetMessage() string {
	return v.Message
}

func (v *VScanError) String() string {
	errString := fmt.Sprintf(`Error in: %v
%v -> %v: at line %v position %v
%v
%v
`, v.GetFile(), v.Class, v.Type, v.Line, v.Position, v.LineFeed, v.Message)
	return errString
}

func (v *VScanError) Error() string {
	return v.String()
}
