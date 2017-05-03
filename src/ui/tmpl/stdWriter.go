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

func (w *StdWriter) HandleErr(e error) {
	ce, ok := e.(*errs.Error)
	if ok {
		w.handleKwkError(ce)
		return
	}
	models.Debug("%+v", e)
	models.LogErr(e)
}

func (w *StdWriter) handleKwkError(e *errs.Error) {
	switch e.Code {
	case errs.CodeInvalidArgument:
		models.Debug("Unhandled err, requires mapping to types.Error:")
		models.LogErr(e)
		return
	case errs.CodeNotAuthenticated:
		w.Render("api:not-authenticated", nil)
		return
	case errs.CodeNotFound:
		w.Render("api:not-found", e)
		return
	case errs.CodeAlreadyExists:
		w.Render("api:exists", nil)
		return
	case errs.CodePermissionDenied:
		w.Render("api:denied", nil)
		return
	case errs.CodeNotImplemented:
		w.Render("api:not-implemented", nil)
		return
	case errs.CodeInternalError:
		models.LogErr(e)
		w.Render("api:error", nil)
		return
	case errs.CodeNotAvailable:
		w.Render("api:not-available", nil)
		return
	}
	o := getMessageOverride(e.Code)
	if o != "" {
		e.Message = o
	}
	w.Render("validation:one-line", e.Message)
}

var overrides = map[errs.ErrCode]string{
	errs.CodeInvalidPassword: "Password must have one upper, one lower and one numeric",
	errs.CodeInvalidUsername: "Username must be bl",
	errs.CodeEmailTaken:      "That email has been taken",
}

func getMessageOverride(code errs.ErrCode) string {
	return overrides[code]
}
