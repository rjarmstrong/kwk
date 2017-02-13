package app

import (
	"io"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
)

type Dashboard struct {
	s  snippets.Service
	tmpl.Writer
}

func NewDashBoard(w tmpl.Writer, s snippets.Service) *Dashboard {
	return &Dashboard{Writer: w, s:s}
}

func (d *Dashboard) GetWriter() func(w io.Writer, templ string, data interface{}) {
   return d.writer
}

func (d *Dashboard) writer(out io.Writer, templ string, data interface{}) {
	r, err := d.s.GetRoot("", true)
	if err != nil {
		d.HandleErr(err)
	}
	d.Render("dashboard", r)
}
