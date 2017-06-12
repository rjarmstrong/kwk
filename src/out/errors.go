package out

import (
	"fmt"
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"io"
	"os"
)

func NewErrHandler(w io.Writer) errs.Handler {
	return &handlerWrapper{Writer: vwrite.New(w)}
}

type handlerWrapper struct {
	vwrite.Writer
}

func (e *handlerWrapper) Handle(err error) {
	e.Write(errHandler(err))
}

func errHandler(e error) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		if e == nil {
			return
		}
		LogErr(e)
		ce, ok := e.(*errs.Error)
		h := internalError(e)
		if ok {
			switch ce.Code {
			case errs.CodeInvalidArgument:
				h = invalidArgument(ce)
			case errs.CodeNotAuthenticated:
				h = NotAuthenticated()
			case errs.CodeNotFound:
				h = notFound("")
			case errs.CodeAlreadyExists:
				h = ItemExists()
			case errs.CodePermissionDenied:
				h = notPermitted()
			case errs.CodeNotImplemented:
				h = notImplemented()
			case errs.CodeNotAvailable:
				h = notAvailable()
			default:
				h = internalError(ce)
			}
		}
		h.Write(w)
	})
}

func NotAuthenticated() vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "Please login to continue: kwk login")
	})
}

func ItemExists() vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "An item with that identifier already exists.")
	})
}

func notFound(name string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		if name == "" {
			name = "Resource"
		}
		Warn(w, "%s couldn't be found.\n", name)
	})
}

func internalError(err error) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		ce, ok := err.(*errs.Error)
		if ok {
			fmt.Fprintln(w, ce.Message, ce.Code)
			return
		}
		LogErr(err)
		Fatal(w, "Something broke. \n- To report type: kwk upload-errors \n"+
			"- You can also try to upgrade: npm update kwkcli -g\n")
		os.Exit(1)
	})
}

func notImplemented() vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "Your CLI may be out of date. Please run: kwk update")
	})
}

func notPermitted() vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "Permission denied.")
	})
}

func notAvailable() vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "Kwk is DOWN! Please try again in a bit.\n\n\n")
	})
}

func invalidArgument(err *errs.Error) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		Warn(w, "%s.\n", err.Message)
	})
}
