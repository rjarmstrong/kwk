package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/tmpl"
	"fmt"
	"io"
	"bitbucket.com/sharingmachine/kwkcli/src/gokwk"
)

type Dashboard struct {
	s gokwk.Snippets
	tmpl.Writer
}

func NewDashBoard(w tmpl.Writer, s gokwk.Snippets) *Dashboard {
	return &Dashboard{Writer: w, s: s}
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
	// TODO: Review as this is confusing?
	//ce, ok := err.(*models.ClientErr)
	//if ok && ce.Contains(models.Code_SnippetVulnerable) {
	//	renderSignoutDash()
	//	return
	//}
	if err != nil {
		d.HandleErr(err)
		return
	}
	d.Render("dashboard", r)
}
