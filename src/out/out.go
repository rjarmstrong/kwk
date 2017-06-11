// Package out provides the UX for kwk super snippets CLI
package out

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/lunixbochs/vtclean"
	"io"
	"text/tabwriter"
	"time"
)

func FreeText(text string) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprint(w, text)
	})
}

func formatTime(millis int64) string {
	t := time.Unix(millis/1000, 0)
	return style.Time(t)
}

func Dashboard(prefs *Prefs, cli *types.AppInfo, rr *types.RootResponse, u *types.User) vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		PrintRoot(prefs, cli, rr, u).Write(w)
	})
}

func SignedOut() vwrite.Handler {
	return vwrite.HandlerFunc(func(w io.Writer) {
		fmt.Fprint(w, "<Signed out dash goes here>  \n\nkwk signin | kwk signup\n\n")
	})
}

func printSnipNames(w io.Writer, snipNames []*types.SnipName) {
	for i, v := range snipNames {
		fmt.Fprintf(w, "%s", v.String())
		if i-1 < len(snipNames) {
			fmt.Fprint(w, ", ")
		}
	}
}

func multiChoice(w io.Writer, in interface{}) {
	list := in.([]*types.Snippet)
	fmt.Fprint(w, "\n")
	if len(list) == 1 {
		fmt.Fprintf(w, "%sDid you mean: %s? y/n\n\n", style.Margin, style.Fmt256(style.ColorPouchCyan, list[0].String()))
		return
	}
	t := tabwriter.NewWriter(w, 5, 1, 3, ' ', tabwriter.TabIndent)
	for i, v := range list {
		if i%3 == 0 {
			t.Write([]byte(style.Margin))
		}
		fmt256 := style.Fmt16(style.Cyan, i+1)
		t.Write([]byte(fmt.Sprintf("%s %s", fmt256, v.Alias.FileName())))
		x := i + 1
		if x%3 == 0 {
			t.Write([]byte("\n"))
		} else {
			t.Write([]byte("\t"))
		}
	}
	t.Write([]byte("\n\n"))
	t.Flush()
	fmt.Fprint(w, style.Margin+style.Fmt256(style.ColorPouchCyan, "Please select a snippet: "))
}

func Fpreview(in string, prefs *Prefs, wrapAt int, lines int) string {
	if prefs.DisablePreview {
		return ""
	}
	in = vtclean.Clean(in, false)
	return style.FBox(in, wrapAt, lines) + style.End
}
