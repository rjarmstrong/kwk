package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"strings"
	"fmt"
	"io"
	"bitbucket.com/sharingmachine/types"
)

	// SEARCH
	//add("search:alpha", "{{ . | alphaSearchResult }}", template.FuncMap{"result": alphaSearchResult })
	//
	//add("search:alphaSuggest", "\n\033[7m Suggestions: \033[0m\n\n{{range .Results}}{{ .Username }}{{ \"/\" }}{{ .Name | blue }}.{{ .Extension | subdued }}\n{{end}}\n", template.FuncMap{"blue": blue, "subdued": subdued})
	//
	//add("search:typeahead", "{{range .Results}}{{ .String }}\n{{end}}", nil)

type SearchResultLine struct {
	Key  string
	Line string
}

func alphaSearchResult(w io.Writer, in interface{}) {
	res := in.(*models.SearchTermResponse)

	fmt.Fprintf(w, "\n\033[7m  \"%s\" found in %d results in %d ms  \033[0m", res.Term, res.Total, res.Took)
	//fmt.Fprint(w, "\n\n")
	// {{ .Username }}{{ \"/\" }}{{ .Name | blue }}.{{ .Extension | subdued }}\n{{ . | result}}\n

	view := &models.ListView{Snippets: []*types.Snippet{}}
	for _, v := range res.Results {
		view.Snippets = append(view.Snippets, v.Snippet)
	}

	fmt.Fprint(w, listSnippets(view, true))

	//for _, v := range res.Results {
	//	fmt.Fprintf(w, "%s%s\n", MARGIN, v.String())
		//if v.Highlights == nil {
		//	v.Highlights = map[string]string{}
		//}
		//if v.Highlights["snip"] == "" {
		//	v.Highlights["snip"] = v.Snip
		//}
		//lines := highlightsToLines(v.Highlights)
		//f := ""
		//for _, line := range lines {
		//	f = f + line.Key[:4] + "\u2847  " + line.Line + "\n"
		//}
		//f = style.Fmt(style.Subdued, f)
		//f = style.ColourSpan(style.Black, f, "<em>", "</em>", style.Subdued)
		//fmt.Fprint(w, f)
	//}
	//fmt.Fprint(w, "\n")
}

func highlightsToLines(highlights map[string]string) []SearchResultLine {
	allLines := []SearchResultLine{}
	for k, v := range highlights {
		lines := strings.Split(v, "\n")
		for _, line := range lines {
			allLines = append(allLines, SearchResultLine{Key: k, Line: line})
		}
	}
	return allLines
}
