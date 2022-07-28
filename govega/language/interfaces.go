// Package language
//
// Defines basic language structures which can be used in the frontend to parse the language
//
// interfaces.go defines external interfaces for data types

package language

import "govega/govega/language/tokens"

// IBasicType interface for BasicType
type IBasicType interface {
	tokens.IWord   // IWord interface from word token
	GetWidth() int // get type memory size
}

// NewBasicType generates new IBasicType interface for BasicType
func NewBasicType(varType string, tag int, width int) IBasicType {
	var newType IBasicType = newBasicType(varType, tag, width)
	return newType
}

// IArrayType interface for ArrayType
type IArrayType interface {
	IBasicType
	GetSize() int
	GetType() *BasicType
	GetDimensions() []int
}

// NewArray generates IArrayType interface for ArrayType
func NewArray(t interface{}, s int) IArrayType {
	var newArray IArrayType = newArray(t, s)
	return newArray
}

// IStringType interface for StringType
type IStringType interface {
	IArrayType
}

// NewString generates IStringType interface for StringType
func NewString(s int) IStringType {
	var newString IStringType = newString(s)
	return newString
}
