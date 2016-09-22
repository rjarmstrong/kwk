package gui

type TemplateWriterMock struct {
     RenderCalledWith []interface{}
}

func (t *TemplateWriterMock) Render(templateName string, data interface{}) {
     t.RenderCalledWith = []interface{}{templateName, data}
}
