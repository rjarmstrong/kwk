package out

import (
	"fmt"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/rjarmstrong/kwk/src/style"
	"io"
	"strings"
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

func AlphaTypeAhead(res *types.TypeAheadResponse) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		//printSnippets(w, view, true)
	})
}

func AlphaSearchResult(prefs *Prefs, res *types.AlphaResponse) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {

		if res.FallbackTitle != "" {
			fmt.Fprintln(w)
			fmt.Fprintln(w, style.Margin, res.FallbackTitle)
		}

		fmt.Fprintf(w, "\n%s\033[7m  \"%s\" found in %d results      %d ms  \033[0m",
			style.Margin, res.Term, res.Total, res.Took)
		fmt.Fprint(w, "\n")

		hltStart := fmt.Sprintf("%s%dm", style.Start255Fg, style.ColorPouchCyan)
		subdued := fmt.Sprintf("%s%dm", style.Start255Fg, style.Grey243)
		for _, v := range res.Results {
			fmt.Fprintf(w, "\n%s%s %s\n", style.Margin, style.Fmt256(style.ColorPouchCyan, style.IconSnippet), v.Snippet.Alias.URI())
			if v.Highlights == nil {
				v.Highlights = map[string]string{}
			}
			if v.Highlights["content"] == "" {
				if len(v.Snippet.Content) > 50 {
					v.Highlights["content"] = v.Snippet.Content[:50] + "\n..."
				} else {
					v.Highlights["content"] = v.Snippet.Content
				}
			}
			lines := highlightsToLines(v.Highlights)

			for _, line := range lines {
				fmt.Fprint(w, style.Margin)
				fmt.Fprint(w, subdued)
				fmt.Fprint(w, "|  ")
				fmt.Fprint(w, strings.Replace(strings.Replace(line.Line,
					"<em>", hltStart, -1), "</em>", style.End+subdued, -1))
				fmt.Fprint(w, "\n")
			}
			fmt.Fprint(w, style.End)
			if len(v.Related) > 0 {
				fmt.Fprintf(w, "%s%dm", style.Start255Fg, style.Grey241)
				fmt.Fprint(w, style.Margin)
				fmt.Fprint(w, style.IconRelated)
				for _, v2 := range v.Related {
					fmt.Fprintf(w, " %s ", v2.Name)
				}
				fmt.Fprint(w, style.End)
			}
			fmt.Fprint(w, "\n")
		}
		fmt.Fprint(w, "\n")

		if len(res.Pouches) > 0 {
			fmt.Fprint(w, "Matching pouches:\n")
			for _, v := range res.Pouches {
				fmt.Fprintf(w, "%s\n", v)
			}
		}
	})
}

//fmt.Fprint(w, "\n\n")
// {{ .Username }}{{ \"/\" }}{{ .Name | blue }}.{{ .Extension | subdued }}\n{{ . | result}}\n

//view := &models.ListView{Snippets: []*types.Snippet{}}
//for _, v := range res.Results {
//	view.Snippets = append(view.Snippets, v.Snippet)
//}
//
//printSnippets(w, view, true)
