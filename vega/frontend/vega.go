package frontend

type vega struct {
	file      string
	codeLines []string
}

func NewVega(filePath string) Vega {
	var lines []string
	var v Vega = &vega{
		file:      filePath,
		codeLines: lines,
	}
	return v
}

func (v *vega) getVega() *vega {
	return v
}
