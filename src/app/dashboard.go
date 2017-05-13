package app

import (
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/vwrite"
	"io"
)

type Dashboard struct {
	s types.SnippetsClient
	vwrite.Writer
	errs.Handler
}

func NewDashBoard(w vwrite.Writer, eh errs.Handler, s types.SnippetsClient) *Dashboard {
	return &Dashboard{Writer: w, s: s, Handler: eh}
}

func (d *Dashboard) GetWriter() func(w io.Writer, templ string, data interface{}) {
	return d.writer
}

func (d *Dashboard) writer(w io.Writer, templ string, data interface{}) {
	if principal == nil {
		d.Write(out.SignedOut())
		return
	}
	r, err := d.s.GetRoot(GetCtx(), &types.RootRequest{})
	if err != nil {
		d.Handle(err)
		return
	}
	d.Write(out.Dashboard(r))
}
