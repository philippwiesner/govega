// Package helper
//
// provides basic data structures to build more complex data structures used in the compiler code architecture
//
// list.go implements basic double linked list implementations for a queue and a stack

package helper

import "errors"

type Node struct {
	data interface{}
	last *Node
	next *Node
}

func newNode(data interface{}) *Node {
	return &Node{data: data}
}

func (n *Node) GetData() interface{} {
	return n.data
}

type List struct {
	tail  *Node
	head  *Node
	count int
}

func NewList() *List {
	return &List{}
}

func (l *List) increment() {
	l.count += 1
}

func (l *List) decrement() {
	l.count -= 1
}

func (l *List) IsEmpty() bool {
	return l.tail == nil && l.head == nil && l.count == 0
}

func (l *List) GetCount() int {
	return l.count
}

func (l *List) getHead() interface{} {
	if l.head == nil {
		return nil
	} else {
		return l.head.GetData()
	}
}

func (l *List) getTail() interface{} {
	if l.tail == nil {
		return nil
	} else {
		return l.tail.GetData()
	}
}

type Stack struct {
	*List
}

func NewStack() *Stack {
	return &Stack{List: NewList()}
}

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

func (s *Stack) Pop() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("empty Stack")
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

func (s *Stack) Top() interface{} {
	return s.getHead()
}

func (s *Stack) Bottom() interface{} {
	return s.getTail()
}

type Queue struct {
	*List
}

func NewQueue() *Queue {
	return &Queue{List: NewList()}
}

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

func (q *Queue) Remove() (interface{}, error) {
	if q.IsEmpty() {
		return nil, errors.New("empty Queue")
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

func (q *Queue) Top() interface{} {
	return q.getHead()
}

func (q *Queue) Bottom() interface{} {
	return q.getTail()
}
