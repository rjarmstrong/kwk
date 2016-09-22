package gui

import (
	"io"
)

type TemplateWriter struct {
	io.Writer
}

func NewTemplateWriter(w io.Writer) *TemplateWriter {
	return &TemplateWriter{Writer:w}
}

func (w *TemplateWriter) Render(templateName string, data interface{})  {
	if t := Templates[templateName]; t != nil {
		Templates[templateName].Execute(w.Writer, data)
	} else {
		panic("Template not found: " + templateName)
	}
}
