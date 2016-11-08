package app

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/search"
	"strings"
)

type SearchController struct {
	service search.ISearch
	gui.ITemplateWriter
	gui.IDialogues
}

func NewSearchController(search search.ISearch, w gui.ITemplateWriter, d gui.IDialogues) *SearchController {
	return &SearchController{service:search, ITemplateWriter: w, IDialogues: d}
}

func (c *SearchController) Search(args ...string){
	term := strings.Join(args, " ")
	if res, err := c.service.Search(term); err != nil {
		c.Render("error", err)
	} else {
		c.Render("search:alpha", res)
	}
}