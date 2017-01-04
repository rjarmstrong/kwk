package tmpl

type WriterMock struct {
	RenderCalledWith []interface{}
	RenderCallCount int
}

func (t *WriterMock) Render(templateName string, data interface{}) {
	t.RenderCallCount += 1
	t.RenderCalledWith = []interface{}{templateName, data}
}

