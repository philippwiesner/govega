package dataStructs

import (
	"govega/govega/helper"
	"govega/govega/language/tokens"
)

type TokenBucket struct {
	token tokens.IToken
	line  int
}

func (tb *TokenBucket) GetToken() tokens.IToken {
	return tb.token
}

func (tb *TokenBucket) GetLine() int {
	return tb.line
}

func (tb *TokenBucket) GetTokenTag() int {
	return tb.GetToken().GetTag()
}

type TokenStream struct {
	*helper.Queue
}

func NewTokenStream() *TokenStream {
	return &TokenStream{helper.NewQueue()}
}

func (ts *TokenStream) Add(token tokens.IToken, line int) {
	ts.Queue.Add(&TokenBucket{token, line})
}

func (ts *TokenStream) Remove() (tokenBucket *TokenBucket, ok bool) {
	data, ok := ts.Queue.Remove()
	if !ok {
		return nil, ok
	}
	return data.(*TokenBucket), ok
}

func (ts *TokenStream) Top() *TokenBucket {
	data := ts.Queue.Top()
	return data.(*TokenBucket)
}
