package frontend

import "fmt"

type objectInterface interface {
	getValue() int
	getName() string
}

type printObjectInterface interface {
	objectInterface
	print(of printObjectInterface)
}

type object struct {
	value int
	name  string
}

type ObjectFactory struct{}

func NewObject(v int, n string) *object {
	return &object{v, n}
}

func (o *object) getValue() int {
	return o.value
}

func (o *object) getName() string {
	return o.name
}

func (o *object) print(of printObjectInterface) {
}

func print(of printObjectInterface) {
	fmt.Printf("Hello %v, your number is: %v\n", of.getName(), of.getValue())
}
