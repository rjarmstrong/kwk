package tmpl

import (
	"io"
)

// StdWriter is the default template writer.
type StdWriter struct {
	io.Writer
}

func NewWriter(w io.Writer) Writer {
	return &StdWriter{Writer: w}
}

func (w *StdWriter) Render(templateName string, data interface{}) {
	if t := Templates[templateName]; t != nil {
		Templates[templateName].Execute(w.Writer, data)
	} else {
		panic("Template not found: " + templateName)
	}
}
