package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"github.com/siddontang/go/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc"
	"io"
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

func (w *StdWriter) HandleErr(err error) {
	switch e := err.(type) {
	case *models.ClientErr:
		w.handleClientError(e)
	case *rpc.RpcError:
		log.Debug("Handling error as remote. Consider replacing with a ClientErr.")
		w.handleRemoteError(e)
	default:
	}
}

/*
 handleRemoteError should only come into effect where we have inadvertently not
 transposed it to a ClientErr in the respective integration service.
 */
func (w *StdWriter) handleRemoteError(e *rpc.RpcError) {
	code := grpc.Code(e)
	switch code {
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
		log.Debug("Internal Error.", e)
		w.Render("api:error", nil)
	case codes.Unavailable:
		w.Render("api:not-available", nil)
	default:
		log.Debug("Unhandled err:", e)
		panic(e)
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
		panic(e)
	}
}

var overrides = map[models.Code]string{
	models.Code_MultipleSnippetsFound: "Multiple snippets found with that name",
	models.Code_InvalidPassword:       "Password must have one upper, one lower and one numeric",
	models.Code_InvalidUsername:       "Username must be bl",
	models.Code_EmailTaken:            "That email has been taken",
}

func getDescOverride(code models.Code) string {
	return overrides[code]
}
