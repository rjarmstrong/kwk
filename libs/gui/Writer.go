package gui


func NewWriter(templates map[string]Template) *Writer {
	return &Writer{templates:templates}
}

type Writer struct {
	templates map[string]Template
}

type Template func(input interface{})

var templates = map[string]Template{}

func (w *Writer) PrintWithTemplate(templateName string, input interface{}) {
	templates[templateName](input)
}
