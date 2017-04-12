package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"io"
	"fmt"
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
	w.HandleErr(models.ErrOneLine(models.Code_PrinterNotFound, "`%s` template not found.", templateName))
}

func (w *StdWriter) Out(templateName string, data interface{}) {
	if t := Printers[templateName]; t != nil {
		Printers[templateName](w.Writer, data)
	} else {
		w.HandleErr(models.ErrOneLine(models.Code_PrinterNotFound, "`%s` printer not found.", templateName))
	}
}

// Should be able to handle all types of errors not limited to grpc.rpcError and models.ClientErr
func (w *StdWriter) HandleErr(e error) {
	code := grpc.Code(e)
	ce, ok := e.(*models.ClientErr)
	if ok {
		w.handleClientError(ce)
		return
	}

	switch code {
	case codes.InvalidArgument:
		log.Error("Unhandled err, requires mapping to client err.:", e)
		panic(e)
	case codes.Unauthenticated:
		w.Render("api:not-authenticated", nil)
	case codes.NotFound:
		w.Render("api:not-found", e)
	case codes.AlreadyExists:
		w.Render("api:exists", nil)
	case codes.PermissionDenied:
		w.Render("api:denied", nil)
	case codes.Unimplemented:
		w.Render("api:not-implemented", nil)
	case codes.Internal:
		log.Error("Internal Error.", e)
		w.Render("api:error", nil)
	case codes.Unavailable:
		w.Render("api:not-available", nil)
	default:
		log.Error("Unhandled err:", e)
	}
}

func (w *StdWriter) handleClientError(e *models.ClientErr) {
	if e.Title != "" {
		w.Render("validation:title", e.Title)
	}
	for i, m := range e.Msgs {
		if o := getDescOverride(m.Code); o != "" {
			e.Msgs[i].Desc = o
		}
	}
	if len(e.Msgs) > 1 {
		for _, m := range e.Msgs {
			w.Render("validation:multi-line", m)
		}

	} else if len(e.Msgs) == 1 {
		w.Render("validation:one-line", e.Msgs[0])
	} else {
		log.Error("Unhandled error", e)
		fmt.Println(e.Error())
	}
}

var overrides = map[models.Code]string{
	//models.Code_MultipleSnippetsFound: "Multiple snippets found with that name",
	models.Code_InvalidPassword:       "Password must have one upper, one lower and one numeric",
	models.Code_InvalidUsername:       "Username must be bl",
	models.Code_EmailTaken:            "That email has been taken",
}

func getDescOverride(code models.Code) string {
	return overrides[code]
}
