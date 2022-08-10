package frontend

import (
	"fmt"
	"testing"
)

type mockObject struct {
	object
}

func (mo *mockObject) getValue() int {
	return 6
}

func (mo *mockObject) print(of printObjectInterface) {
	fmt.Printf("Hello %v: %v\n", mo.name, of.getValue())
}

func Test(t *testing.T) {
	var o printObjectInterface = NewObject(5, "Wilma")
	print(o)

	var mo printObjectInterface = &mockObject{object{name: "paul"}}
	print(mo)
	mo.print(mo)
}
