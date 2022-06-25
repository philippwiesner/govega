// Package helper
//
// provides basic data structures to build more complex data structures used in the compiler code architecture
//
// errors.go defines custom errors used in the helper package

package helper

import "errors"

var (
	EmptyQueueError = errors.New("queue is empty")
	EmptyStackError = errors.New("stack is empty")
)
