package app

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/search"
	"strings"
	"bitbucket.com/sharingmachine/kwkcli/libs/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/libs/ui/dlg"
)

type SearchController struct {
	service search.ISearch
	tmpl.Writer
	dlg.Dialogue
}

func NewSearchController(search search.ISearch, w tmpl.Writer, d dlg.Dialogue) *SearchController {
	return &SearchController{service: search, Writer: w, Dialogue: d}
}

func (c *SearchController) Search(args ...string) {
	term := strings.Join(args, " ")
	if res, err := c.service.Search(term); err != nil {
		c.Render("error", err)
	} else {
		res.Term = term
		c.Render("search:alpha", res)
	}
}
