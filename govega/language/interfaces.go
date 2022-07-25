package language

import "govega/govega/language/tokens"

type IBasicType interface {
	tokens.IWord
	GetWidth() int
}

func NewBasicType(varType string, tag int, width int) IBasicType {
	var newType IBasicType = newBasicType(varType, tag, width)
	return newType
}

type IArrayType interface {
	IBasicType
	GetSize() int
	GetType() *BasicType
	GetDimensions() []int
}

func NewArray(t interface{}, s int) IArrayType {
	var newArray IArrayType = newArray(t, s)
	return newArray
}

type IStringType interface {
	IArrayType
}

func NewString(s int) IStringType {
	var newString IStringType = newString(s)
	return newString
}
