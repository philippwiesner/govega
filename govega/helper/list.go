// Package helper
//
// provides basic data structures to build more complex data structures used in the compiler code architecture
//
// list.go implements basic double linked list implementations for a queue and a stack

package helper

import "fmt"

// Node is a basic list element which contains some data and has two links to the next and last Node
type Node struct {
	data interface{}
	last *Node
	next *Node
}

// newNode is the constructor for creating a new Node with data
func newNode(data interface{}) *Node {
	return &Node{data: data}
}

// GetData getter method for returning the data
func (n *Node) GetData() interface{} {
	return n.data
}

// list is a basic List which stores nodes in a double linked list
type list struct {
	tail  *Node
	head  *Node
	count int
}

// newList is the constructor for a new list
func newList() *list {
	return &list{}
}

// increment increases the local list element counter
func (l *list) increment() {
	l.count += 1
}

// decrement decreases the local list element counter
func (l *list) decrement() {
	l.count -= 1
}

// IsEmpty checks if the list is empty or not
func (l *list) IsEmpty() bool {
	return l.tail == nil && l.head == nil && l.count == 0
}

// GetCount Getter method for returning the counter
func (l *list) GetCount() int {
	return l.count
}

// getHead private method for getting the first element in list
func (l *list) getHead() interface{} {
	if l.head == nil {
		return nil
	} else {
		return l.head.GetData()
	}
}

// getTail private method for getting the last element in list
func (l *list) getTail() interface{} {
	if l.tail == nil {
		return nil
	} else {
		return l.tail.GetData()
	}
}

// Stack is a double-linked list where elements are added on top and retrieved from the top
type Stack struct {
	*list
}

// NewStack is the constructor for a new Stack
func NewStack() *Stack {
	return &Stack{newList()}
}

// Push public method for adding new elements on top of the Stack
func (s *Stack) Push(data interface{}) {
	node := newNode(data)
	if s.IsEmpty() {
		s.head = node
		s.tail = node
	} else {
		old := s.head
		s.head = node
		old.next = node
		node.last = old
	}
	s.increment()
}

// Pop public method for retrieving elements from the top of the Stack
func (s *Stack) Pop() (result interface{}, err error) {
	if s.IsEmpty() {
		return nil, fmt.Errorf("stack pop: %w", EmptyStackError)
	}

	node := s.head
	if node == s.tail {
		s.head = nil
		s.tail = nil
	} else {
		s.head = node.last
		s.head.next = nil
	}
	data := node.data
	node.next = nil
	node.last = nil
	node = nil
	s.decrement()
	return data, nil
}

// Top public getter method for getting the top element of the Stack
func (s *Stack) Top() interface{} {
	return s.getHead()
}

// Bottom public getter method for getting the oldest element in the Stack
func (s *Stack) Bottom() interface{} {
	return s.getTail()
}

// Queue is a double-linked list where new elements are added at the bottom of the list and retrieved from the top
type Queue struct {
	*list
}

// NewQueue is the constructor for the Queue
func NewQueue() *Queue {
	return &Queue{newList()}
}

// Add public method for adding new elements at the end of the Queue
func (q *Queue) Add(data interface{}) {
	node := newNode(data)
	if q.IsEmpty() {
		q.tail = node
		q.head = node
	} else {
		old := q.tail
		q.tail = node
		old.last = node
		node.next = old
	}
	q.increment()
}

// Remove public method for retrieving elements from the front of the Queue
func (q *Queue) Remove() (result interface{}, err error) {
	if q.IsEmpty() {
		return nil, fmt.Errorf("queue remove: %w", EmptyQueueError)
	}

	node := q.head
	if node == q.tail {
		q.tail = nil
		q.head = nil
	} else {
		q.head = node.last
		q.head.next = nil
	}
	data := node.data
	node.next = nil
	node.last = nil
	q.decrement()
	return data, nil
}

// Top public getter method for getting the oldest element in the Queue
func (q *Queue) Top() interface{} {
	return q.getHead()
}

// Bottom public getter method for getting the newest element in the Queue
func (q *Queue) Bottom() interface{} {
	return q.getTail()
}
