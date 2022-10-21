package frontend

type testVegaInterface interface {
	Vega
	getVega() *vega
}

type testVega struct {
	vega
}

func createTestVega(path string, lines []string) testVegaInterface {
	var v testVegaInterface = &testVega{
		vega{
			file:      path,
			codeLines: lines,
		},
	}
	return v
}

func (t *testVega) getFile() string {
	return t.vega.file
}

func (t *testVega) getLines() []string {
	return t.vega.codeLines
}
