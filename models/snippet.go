package models

import (
	"bitbucket.com/sharingmachine/types"
	"bytes"
	"strings"
)

const (
	ProfileFullKey = "profile.json"
)

func NewSnippet(snippet string) *types.Snippet {
	return &types.Snippet{Snip: snippet, Alias: types.Alias{SnipName: types.SnipName{}}}
}


type CreateSnippetResponse struct {
	Snippet   *types.Snippet
}

type ListParams struct {
	All           bool
	Pouch         string
	Username      string
	Size          int64
	Since         int64
	Tags          []string
	IgnorePouches bool
	Category      string
}

func Limit(in string, length int) string {
	in = strings.Replace(in, "\n", "  ", -1)
	in = strings.TrimSpace(in)
	var numRunes = 0
	b := bytes.Buffer{}
	for _, r := range in {
		if numRunes == length {
			return strings.TrimSpace(b.String())
		}
		numRunes++
		if r == '\n' {
			b.WriteRune(' ')
			b.WriteRune(' ')
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

//TODO: Optimise this
func ScanVulnerabilities(snip string, ext string) error {
	if strings.Contains(snip, "rm -rf") || strings.Contains(snip, "rm ") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: Shell scripts cannot contain 'rm '.")
	}
	if strings.Contains(snip, ":(){") || strings.Contains(snip, "./$0|./$0&") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: Fork bomb detected.")
	}
	if strings.Contains(snip, "fork") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: 'fork' not allowed in script.")
	}
	if strings.Contains(snip, "/dev/sd") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: '/dev/sd' is not allowed in scripts.")
	}
	if strings.Contains(snip, "/dev/null") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: '/dev/null' is not allowed in scripts.")
	}
	if strings.Contains(snip, "| sh") || strings.Contains(snip, "| bash") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: piping directly into terminal not allowed in scripts.")
	}
	if strings.Contains(snip, "nohup") {
		return ErrOneLine(Code_SnippetVulnerable, "kwk constraint: 'nohup' command is not allowed.")
	}
	if (ext == "sh" || ext == "js") && strings.Contains(snip, "eval") {
		m := "kwk constraint: 'eval' command is not allowed."
		if ext == "sh" {
			m += "  Tip: try using '($VAR)' instead of 'eval $VAR' to execute commands.\n"
			m += "  See: /richard/cli/basheval.url\n"
		}
		return ErrOneLine(Code_SnippetVulnerable, m)
	}
	return nil
}
