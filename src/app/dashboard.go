package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/app/out"
	"bitbucket.com/sharingmachine/kwkcli/src/gokwk"
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/types/errs"
	"bitbucket.com/sharingmachine/types/vwrite"
	"io"
)

type Dashboard struct {
	s gokwk.Snippets
	vwrite.Writer
	errs.Handler
}

func NewDashBoard(w vwrite.Writer, eh errs.Handler, s gokwk.Snippets) *Dashboard {
	return &Dashboard{Writer: w, s: s, Handler: eh}
}

func (d *Dashboard) GetWriter() func(w io.Writer, templ string, data interface{}) {
	return d.writer
}

func (d *Dashboard) writer(w io.Writer, templ string, data interface{}) {
	if len(models.Principal.Token) == 0 {
		d.Write(out.SignedOut())
		return
	}
	r, err := d.s.GetRoot("", true)
	if err != nil {
		d.Handle(err)
		return
	}
	d.Write(out.Dashboard(r))
}
