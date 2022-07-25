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

// newBasicType is the constructor for new a BasicType
//
// Takes the name of the type, tag and memory storage space.
func newBasicType(varType string, tag int, width int) *BasicType {
	return &BasicType{tokens.NewWord(varType, tag), width}
}

// GetWidth public getter method for the memory storage space
func (t *BasicType) GetWidth() int {
	return t.width
}

// Define basic data types: integer, float, char and bool
var (
	IntType   = NewBasicType("int", tokens.BASIC, 4)
	FloatType = NewBasicType("float", tokens.BASIC, 8)
	CharType  = NewBasicType("char", tokens.BASIC, 8)
	BoolType  = NewBasicType("bool", tokens.BASIC, 1)
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

// newArray is the constructor for new array types.
//
// When creating multidimensional arrays the type of the given array is being used and the total size to be allocated
// is multiplicated by the new and the old size. All array dimensions are stored in a list
func newArray(t interface{}, s int) *ArrayType {
	arr := new(ArrayType)
	arr.size = s
	switch v := t.(type) {
	case *BasicType:
		arr.BasicType = newBasicType("[]", tokens.INDEX, s*v.width)
		arr.arrayType = v
		arr.dimensions = []int{s}
	case *ArrayType:
		arr.BasicType = newBasicType(v.GetLexeme(), v.GetTag(), s*v.GetWidth())
		arr.arrayType = v.arrayType
		arr.dimensions = append(v.dimensions, s)
	}
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

// newString is the constructor for creating new strings with an array of the type CHAR
func newString(s int) *StringType {
	return &StringType{newArray(CharType, s)}
}
