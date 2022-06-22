package helper

import (
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
			if !err {
				t.Errorf(`Error: #{err}`)
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
	want := false
	_, got := stack.Pop()
	if want != got {
		t.Fatalf("Want %v, got: %v", want, got)
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
			if !err {
				t.Errorf(`Error: #{err}`)
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
	want := false
	_, got := stack.Remove()
	if want != got {
		t.Fatalf("Want %v, got: %v", want, got)
	}
}
