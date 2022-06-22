// Package language
//
// Defines basic language structures which can be used in the frontend to parse the language
//
// types.go implements data types which can be recognised in lexical analysis

package language

import "govega/govega/language/tokens"

// BasicType represents simple or basic variable types like integers, floating point numbers, chars or boolean values
type BasicType struct {
	tokens.IWord     // keyWord describing type (int, float, char, bool)
	width        int // size of memory storage space
}

// NewBasicType is the constructor for new a BasicType
//
// Takes the name of the type, tag and memory storage space.
func NewBasicType(varType string, tag int, width int) *BasicType {
	return &BasicType{tokens.NewWord(varType, tag), width}
}

// GetWidth public getter method for the memory storage space
func (t *BasicType) GetWidth() int {
	return t.width
}

// Define basic data types: integer, float, char and bool
var (
	intType   = NewBasicType("int", tokens.BASIC, 4)
	floatType = NewBasicType("float", tokens.BASIC, 8)
	charType  = NewBasicType("char", tokens.BASIC, 8)
	boolType  = NewBasicType("bool", tokens.BASIC, 1)
)

// ArrayType is the data type describing an array
//
// Each array inherits the BasicType to describe the array as a dedicated type, further a size for the array, a data
// type for the data stored inside the array and the dimensions of the array for multidimensional arrays exists
type ArrayType struct {
	*BasicType
	size       int
	arrayType  *BasicType
	dimensions []int
}

// NewArray is the constructor for a new one dimensional array
func NewArray(t *BasicType, s int) *ArrayType {
	arr := new(ArrayType)
	arr.BasicType = NewBasicType("[]", tokens.INDEX, s*t.width)
	arr.size = s
	arr.arrayType = t
	arr.dimensions = []int{s}
	return arr
}

// NewArrayArray is the constructor for multidimensional arrays
//
// A new ArrayType is being created from the lexeme, tags and types of the old array but the memory space is
// multiplicated with the new array size. Further the dimensions list is being extended with the new dimensions
func NewArrayArray(a *ArrayType, s int) *ArrayType {
	arr := new(ArrayType)
	arr.BasicType = NewBasicType(a.GetLexeme(), a.GetTag(), s*a.GetWidth())
	arr.size = s
	arr.arrayType = a.arrayType
	arr.dimensions = append(a.dimensions, s)
	return arr
}

// GetSize public getter method for getting the array size
func (a *ArrayType) GetSize() int {
	return a.size
}

// GetType public getter method for getting the array type
func (a *ArrayType) GetType() *BasicType {
	return a.arrayType
}

// GetDimensions public getter method for getting the array dimensions
func (a *ArrayType) GetDimensions() []int {
	return a.dimensions
}

// StringType is basically a special array just for characters
type StringType struct {
	*ArrayType
}

// NewString is the constructor for creating new strings with an array of the type CHAR
func NewString(s int) *StringType {
	return &StringType{NewArray(charType, s)}
}
