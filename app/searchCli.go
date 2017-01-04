package app

import (
	"bitbucket.com/sharingmachine/kwkcli/search"
	"strings"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
)

type SearchCli struct {
	service search.Term
	tmpl.Writer
	dlg.Dialog
}

func NewSearchCli(search search.Term, w tmpl.Writer, d dlg.Dialog) *SearchCli {
	return &SearchCli{service: search, Writer: w, Dialog: d}
}

func (c *SearchCli) Search(args ...string) {
	term := strings.Join(args, " ")
	if res, err := c.service.Execute(term); err != nil {
		c.Render("error", err)
	} else {
		res.Term = term
		c.Render("search:alpha", res)
	}
}
