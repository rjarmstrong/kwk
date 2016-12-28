package tmpl

type WriterMock struct {
	RenderCalledWith []interface{}
}

func (t *WriterMock) Render(templateName string, data interface{}) {
	t.RenderCalledWith = []interface{}{templateName, data}
}
