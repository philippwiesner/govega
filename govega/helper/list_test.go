package helper

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewList(t *testing.T) {
	list := newList()
	if !list.IsEmpty() {
		t.Fatalf("List not empty, got: %v", list)
	}
}

func TestNewStack(t *testing.T) {
	stack := NewStack()
	if !stack.IsEmpty() {
		t.Fatalf("Stack not empty, got: %v", stack)
	}
}

func TestNewQueue(t *testing.T) {
	queue := NewQueue()
	if !queue.IsEmpty() {
		t.Fatalf("Queue not empty, got: %v", queue)
	}
}

func TestNode_GetData(t *testing.T) {
	tests := []struct {
		in   interface{}
		want interface{}
	}{
		{1, 1},
		{struct {
			a int
			b int
		}{1, 2}, struct {
			a int
			b int
		}{1, 2}},
	}
	for i, tc := range tests {
		n := newNode(tc.in)
		got := n.GetData()
		if tc.want != got {
			t.Fatalf("test %d: expected: %v, got: %v", i+1, tc.want, got)
		}
	}
}

func TestStack(t *testing.T) {
	tests := []struct {
		in   []interface{}
		want []interface{}
	}{
		{[]interface{}{1, 2, 3}, []interface{}{3, 2, 1}},
		{[]interface{}{"a", "b", "c"}, []interface{}{"c", "b", "a"}},
	}

	for i, tc := range tests {
		stack := NewStack()
		for _, el := range tc.in {
			stack.Push(el)
		}
		var got []interface{}
		for !stack.IsEmpty() {
			data, err := stack.Pop()
			if err == EmptyStackError {
				t.Error(err)
			}
			got = append(got, data)
		}
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("test %d: expected: %v, got: %v", i+1, tc.want, got)
		}
	}
}

func TestStackEmptyPop(t *testing.T) {
	stack := NewStack()
	_, err := stack.Pop()
	if !errors.Is(err, EmptyStackError) {
		t.Fatalf("error is not %v, but %v", EmptyStackError, err)
	}
}

func TestQueue(t *testing.T) {
	tests := []struct {
		in   []interface{}
		want []interface{}
	}{
		{[]interface{}{1, 2, 3}, []interface{}{1, 2, 3}},
		{[]interface{}{"a", "b", "c"}, []interface{}{"a", "b", "c"}},
	}

	for i, tc := range tests {
		queue := NewQueue()
		for _, el := range tc.in {
			queue.Add(el)
		}
		var got []interface{}
		for !queue.IsEmpty() {
			data, err := queue.Remove()
			if err == EmptyQueueError {
				t.Error(err)
			}
			got = append(got, data)
		}
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("test %d: expected: %v, got: %v", i+1, tc.want, got)
		}
	}
}

func TestQueueEmptyRemove(t *testing.T) {
	stack := NewQueue()
	_, err := stack.Remove()
	if !errors.Is(err, EmptyQueueError) {
		t.Fatalf("error is not %v, but %v", EmptyQueueError, err)
	}
}

func TestBla(t *testing.T) {
	s := NewStack()
	s.Top().GetLast().GetLast().GetNext().GetData()
}
