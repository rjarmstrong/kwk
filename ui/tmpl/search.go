package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"strings"
)

type SearchResultLine struct {
	Key  string
	Line string
}

func alphaSearchResult(result models.SearchResult) string {
	if result.Highlights == nil {
		result.Highlights = map[string]string{}
	}
	if result.Highlights["snip"] == "" {
		result.Highlights["snip"] = result.Snip
	}
	lines := highlightsToLines(result.Highlights)
	f := ""
	for _, line := range lines {
		f = f + line.Key[:4] + "\u2847  " + line.Line + "\n"
	}
	f = style.Fmt(style.Subdued, f)
	f = style.ColourSpan(style.Black, f, "<em>", "</em>", style.Subdued)
	return f
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
