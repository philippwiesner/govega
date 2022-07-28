// Package frontend
//
// The compilers' frontend is responsible for validating and analysing the source programming language.
//
// tokenStream.go defines a stream of language tokens with is created during the lexical analyses and disected by the
// parser to check for the correct language grammar
package frontend

import (
	"govega/govega/language/tokens"
	"io"
)

// TokenBucket stores an tokens.IToken interface as a token the code line the token occures
type TokenBucket struct {
	token      tokens.IToken
	errorState ErrorState
	last       *TokenBucket
}

// newTokenBucket creates a new bucket for a token
func newTokenBucket(token tokens.IToken, state ErrorState) *TokenBucket {
	return &TokenBucket{token: token, errorState: state}
}

// GetToken getter method for the tokens.IToken interface
func (tb *TokenBucket) GetToken() tokens.IToken {
	return tb.token
}

// GetErrorState getter method for retrieving error information during failure
func (tb *TokenBucket) GetErrorState() *ErrorState {
	return &tb.errorState
}

// GetTokenTag getter method to retrieve the token tag
func (tb *TokenBucket) GetTokenTag() int {
	return tb.GetToken().GetTag()
}

// GetTokenLine getter method to retrieve token occurence
func (tb *TokenBucket) GetTokenLine() int {
	return tb.errorState.lineNumber
}

// TokenStream implements a token stream
type TokenStream struct {
	tail *TokenBucket
	head *TokenBucket
}

// NewTokenStream generates a new TokenStream
func NewTokenStream() *TokenStream {
	return &TokenStream{}
}

func (ts *TokenStream) IsEmpty() bool {
	return ts.tail == nil && ts.head == nil
}

func (ts *TokenStream) GetTail() *TokenBucket {
	return ts.tail
}

func (ts *TokenStream) GetHead() *TokenBucket {
	return ts.head
}

// Add a new token on top of the token stream
func (ts *TokenStream) Add(token tokens.IToken, state ErrorState) {
	newToken := newTokenBucket(token, state)
	if ts.IsEmpty() {
		ts.tail = newToken
		ts.head = newToken
	} else {
		old := ts.tail
		ts.tail = newToken
		old.last = newToken
	}
}

// Remove the top element from the token stream
func (ts *TokenStream) Remove() (TokenBucket *TokenBucket, err error) {
	if ts.IsEmpty() {
		return nil, io.EOF
	}
	head := ts.head
	ts.head = head.last
	head.last = nil
	return head, nil
}
