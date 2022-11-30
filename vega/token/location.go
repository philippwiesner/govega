package token

// Location stores line and position of token occurrence
type Location struct {
	FileName string
	Line     int
	Position int
	LineFeed string
}

type File struct {
	Name   string
	Source []byte
}

func readFile(name string) []byte {
	var code []byte // fill me with actual reading
	return code
}

func NewFile(name string) File {
	return File{
		Name:   name,
		Source: readFile(name),
	}
}
