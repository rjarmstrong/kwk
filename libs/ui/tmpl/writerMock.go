package tmpl

type MockWriter struct {
	RenderCalledWith []interface{}
}

func (t *MockWriter) Render(templateName string, data interface{}) {
	t.RenderCalledWith = []interface{}{templateName, data}
}
