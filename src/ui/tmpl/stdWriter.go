package tmpl

import (
	"bitbucket.com/sharingmachine/types/errs"
	"io"
	"bitbucket.com/sharingmachine/kwkcli/src/models"
)

/*
StdWriter is the default template writer.
*/
type StdWriter struct {
	io.Writer
}

func NewWriter(w io.Writer) Writer {
	return &StdWriter{Writer: w}
}

func (w *StdWriter) Render(templateName string, data interface{}) {
	t := Templates[templateName]
	if t != nil {
		Templates[templateName].Execute(w.Writer, data)
		return
	}
	w.HandleErr(errs.New(errs.CodePrinterNotFound, "`%s` template not found.", templateName))
}

func (w *StdWriter) Out(templateName string, data interface{}) {
	if t := Printers[templateName]; t != nil {
		Printers[templateName](w.Writer, data)
	} else {
		w.HandleErr(errs.New(errs.CodePrinterNotFound, "`%s` printer not found.", templateName))
	}
}
