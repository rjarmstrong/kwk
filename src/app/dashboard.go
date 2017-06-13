package app

import (
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/runtime"
	"io"
)

type Dashboard struct {
	rg runtime.RootGetter
	vwrite.Writer
	errs.Handler
}

func NewDashBoard(w vwrite.Writer, eh errs.Handler, rg runtime.RootGetter) *Dashboard {
	return &Dashboard{Writer: w, rg: rg, Handler: eh}
}

func (d *Dashboard) GetWriter() func(w io.Writer, templ string, data interface{}) {
	return func(w io.Writer, templ string, data interface{}) {
		if !principal.HasAccessToken() {
			d.Write(out.SignedOut())
			return
		}
		r, err := d.rg(&types.RootRequest{PrivateView: prefs.PrivateView})
		if err != nil {
			d.Handle(err)
			return
		}
		d.Write(out.Dashboard(prefs, &info, r, &principal.User))
	}
}
