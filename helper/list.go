package helper

import "errors"

type Node struct {
	data interface{}
	last *Node
	next *Node
}

func newNode(data interface{}) *Node {
	n := Node{data, nil, nil}
	return &n
}

type List struct {
	tail  *Node
	head  *Node
	count int
}

func NewList() *List {
	l := List{nil, nil, 0}
	return &l
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

type Stack struct {
	*List
}

func NewStack() *Stack {
	l := NewList()
	s := Stack{l}
	return &s
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

type Queue struct {
	*List
}

func NewQueue() *Queue {
	l := NewList()
	q := Queue{l}
	return &q
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
