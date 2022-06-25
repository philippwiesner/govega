package frontend

type LexicalError struct {
	Msg  string
	Line int
	Pos  int
}

type LexicalErrorChar struct {
	LexicalError
	Char rune
}

func (l *LexicalError) Error() string {
	return l.Msg
}
