// Package frontend
//
// The compilers' frontend is responsible for validating and analysing the source programming language.
//
// tokenStream.go defines a stream of language tokens with is created during the lexical analyses and disected by the
// parser to check for the correct language grammar
package frontend

import (
	"fmt"
	"govega/govega/helper"
	"govega/govega/language/tokens"
)

// TokenBucket stores an tokens.IToken interface as a token the code line the token occures
type TokenBucket struct {
	token      tokens.IToken
	errorState ErrorState
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

// TokenStream implements a token stream
type TokenStream struct {
	*helper.Queue
}

// NewTokenStream generates a new TokenStream
func NewTokenStream() *TokenStream {
	return &TokenStream{helper.NewQueue()}
}

// Add overwrites helper.Queue Add method to add a new token and its line of occurense to the token stream
func (ts *TokenStream) Add(token tokens.IToken, state ErrorState) {
	ts.Queue.Add(&TokenBucket{token, state})
}

// Remove overwrites helper.Queue Remove method to remove the top element from the token stream
func (ts *TokenStream) Remove() (tokenBucket *TokenBucket, err error) {
	data, err := ts.Queue.Remove()
	if err != nil {
		return nil, fmt.Errorf("tokenstream remove: %w", err)
	}
	return data.(*TokenBucket), nil
}

// Top overwrites helper.Queue Top method to lookup on the top of the token stream
func (ts *TokenStream) Top() *TokenBucket {
	data := ts.Queue.Top()
	return data.(*TokenBucket)
}
