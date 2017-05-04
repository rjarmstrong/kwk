package ui

import (
	"io"
	"bitbucket.com/sharingmachine/types/errs"
	"bitbucket.com/sharingmachine/kwkcli/src/models"
)

func NewErrWriter(w io.Writer) io.Writer {
	return &ErrWriter{Writer:w}
}

type ErrWriter struct {
	io.Writer
}

func (er ErrWriter) Write(p []byte) (n int, err error) {
	return er.Writer.Write(p)
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

func (er *ErrWriter) handleKwkError(e *errs.Error) {
	switch e.Code {
	case errs.CodeInvalidArgument:
		models.Debug("Unhandled err, requires mapping to types.Error:")
		models.LogErr(e)
		return
	case errs.CodeNotAuthenticated:
		er.Render("api:not-authenticated", nil)
		return
	case errs.CodeNotFound:
		er.Render("api:not-found", e)
		return
	case errs.CodeAlreadyExists:
		er.Render("api:exists", nil)
		return
	case errs.CodePermissionDenied:
		er.Render("api:denied", nil)
		return
	case errs.CodeNotImplemented:
		er.Render("api:not-implemented", nil)
		return
	case errs.CodeInternalError:
		models.LogErr(e)
		er.Render("api:error", nil)
		return
	case errs.CodeNotAvailable:
		er.Render("api:not-available", nil)
		return
	}
	o := getMessageOverride(e.Code)
	if o != "" {
		e.Message = o
	}
	er.Render("validation:one-line", e.Message)
}

var overrides = map[errs.ErrCode]string{
	errs.CodeInvalidPassword: "Password must have one upper, one lower and one numeric",
	errs.CodeInvalidUsername: "Username must be bl",
	errs.CodeEmailTaken:      "That email has been taken",
}

func getMessageOverride(code errs.ErrCode) string {
	return overrides[code]
}
