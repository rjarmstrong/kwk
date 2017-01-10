package tmpl

import (
	"io"
	"google.golang.org/grpc/codes"
	"bitbucket.com/sharingmachine/kwkcli/models"
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
	if t := Templates[templateName]; t != nil {
		Templates[templateName].Execute(w.Writer, data)
	} else {
		panic("Template not found: " + templateName)
	}
}

/*
 HandleErr requires a *models.ClientErr param.
 Make sure that any other errors are converted to a *models.ClientErr prior to calling this method.
 */
func (w *StdWriter) HandleErr(err error) {
	e, ok := err.(*models.ClientErr)
	if !ok {
		panic(err)
	}
	switch e.TransportCode {
	case codes.InvalidArgument:
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
			panic(e)
		}
	case codes.Unauthenticated:
		w.Render("api:not-authenticated", nil)
	case codes.NotFound:
		w.Render("api:not-found", nil)
	case codes.AlreadyExists:
		w.Render("api:exists", nil)
	case codes.PermissionDenied:
		w.Render("api:not-found", nil)
	case codes.Unimplemented:
		panic("not implemented")
	case codes.Internal:
		w.Render("api:error", nil)
	case codes.Unavailable:
		w.Render("api:not-available", nil)
	default:
	}
}

var overrides = map[models.Code]string{
	models.Code_MultipleSnippetsFound: "Multiple snippets found with that name",
	models.Code_InvalidPassword: "Password must have one upper, one lower and one numeric",
	models.Code_InvalidUsername: "Username must be bl",
	models.Code_EmailTaken:      "That email has been taken",
}

func getDescOverride(code models.Code) string {
	return overrides[code]
}
