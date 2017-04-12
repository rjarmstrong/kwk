package app

import (
	"io"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"fmt"
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

func renderSignoutDash() {
	fmt.Print("<Signed out dash>\n\nkwk signin | kwk signup\n")
}

func (d *Dashboard) writer(out io.Writer, templ string, data interface{}) {
	if len(models.Principal.Token) == 0 {
		renderSignoutDash()
		return
	}
	r, err := d.s.GetRoot("", true)
	ce, ok := err.(*models.ClientErr)
	if ok && ce.Contains(models.Code_SnippetVulnerable) {
		renderSignoutDash()
		return
	}
	if err != nil {
		d.HandleErr(err)
		return
	}
	d.Render("dashboard", r)
}
