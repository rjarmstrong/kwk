package out

import (
	"bytes"
	"fmt"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk/src/style"
	"github.com/rjarmstrong/tablewriter"
	"io"
	"strings"
)

func printSnippetView(w io.Writer, prefs *Prefs, s *types.Snippet) {
	fmt.Fprintln(w, "")
	fmt.Fprint(w, style.Margin)
	fmtHeader(w, s.Username(), s.Pouch(), s.Alias.FileName())
	fmt.Fprint(w, strings.Repeat(" ", 4))
	fmt.Fprint(w, snippetIcon(s))
	fmt.Fprint(w, "  ")
	fmt.Fprint(w, fSnippetType(s))
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, style.TwoLines)
	fmt.Fprint(w, fCodeview(s, 100, 0, false, prefs.ExpandedRows))
	fmt.Fprint(w, "\n\n")

	tbl := tablewriter.NewWriter(w)
	tbl.SetAutoWrapText(false)
	tbl.SetBorder(false)
	tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator("")
	tbl.SetRowLine(true)
	tbl.SetAutoFormatHeaders(false)
	tbl.SetHeaderLine(false)
	tbl.SetColWidth(20)
	tbl.Append([]string{style.Fmt256(colors.RecentPouch, fSnippetType(s)+" Details:"), "", "", ""})

	var lastRun string
	if s.Stats.Runs < 1 {
		lastRun = "never"
	} else {
		lastRun = pad(20, formatTime(s.RunStatusTime)).String()
	}
	tbl.Append([]string{
		style.Fmt16(colors.Subdued, "Run Status:"), FStatus(s, true),
		style.Fmt16(style.Subdued, "Last Run:"), lastRun,
	})
	tbl.Append([]string{
		style.Fmt16(style.Subdued, "Run Count: "), fmt.Sprintf("↻ %2d", s.Stats.Runs),
		style.Fmt16(style.Subdued, "View count:"), fmt.Sprintf("%s  %2d", style.IconView, s.Stats.Views)})
	if s.IsApp() {
		var uris []string
		for _, v := range s.Dependencies.Aliases {
			uris = append(uris, v.URI())
		}
		tbl.Append([]string{
			style.Fmt16(style.Subdued, "App Deps:"),
			style.FBox(strings.Join(uris, ", "), 50, 5)})
	}
	var apps []string
	for _, v := range s.Apps.Aliases {
		apps = append(apps, v.URI())
	}
	tbl.Append([]string{
		style.Fmt16(style.Subdued, "Used by:"),
		style.FBox(strings.Join(apps, ", "), 50, 5)})

	var oss []string
	for k := range s.SupportedOn.Oss {
		oss = append(oss, k)
	}
	tbl.Append([]string{
		style.Fmt16(style.Subdued, "Supported OS:"),
		style.FBox(strings.Join(oss, ", "), 50, 5)})
	tbl.Append([]string{
		style.Fmt16(style.Subdued, "Description:"), style.FBox(fEmpty(s.Description), 50, 3), "", ""})

	tbl.Append([]string{
		style.Fmt16(style.Subdued, "Preview:"), fPreview(s.Preview, prefs, 50, 1), "", ""})

	var tags []string
	for k := range s.Tags.Names {
		tags = append(tags, k)
	}
	tbl.Append([]string{
		style.Fmt16(style.Subdued, "Tags:"), fTags(tags), "", ""})
	tbl.Append([]string{
		style.Fmt16(style.Subdued, "sha256:"), fVerified(s)})
	tbl.Append([]string{
		style.Fmt16(style.Subdued, "Updated:"), fmt.Sprintf("%s - %s", formatTime(s.Updated), fmt.Sprintf("v%d", s.Version())),
	})
	tbl.Render()

	//fmt.Fprint(w, style.Start)
	//fmt.Fprintf(w, "%dm", style.Subdued)
	//fmt.Fprint(w, MARGIN)
	//fmt.Fprint(w,"Snippet details: `kwk <name>`")
	//fmt.Fprint(w, MARGIN)
	//fmt.Fprint(w,"Run snippet: `kwk run <name>`")
	//fmt.Fprint(w, MARGIN)
	//fmt.Fprint(w, style.End)
	fmt.Fprint(w, "\n")
}

func fSnippetType(s *types.Snippet) string {
	if s.IsApp() {
		return "App"
	} else if s.Ext() == "url" {
		return "Bookmark"
	} else {
		return "Snippet"
	}
}
func fVerified(s *types.Snippet) string {
	var buff bytes.Buffer
	if s.VerifyChecksum() {
		buff.WriteString(style.Fmt256(style.ColorYesGreen, style.IconTick+" "))
		buff.WriteString(pad(12, s.Checksum).String())
		buff.WriteString("...")
	} else {
		buff.WriteString(" " + style.IconCross + "  Invalid Checksum: ")
		buff.WriteString(fEmpty(s.Checksum))
	}
	return buff.String()
}

func fCodeview(s *types.Snippet, width int, lines int, odd bool, alwaysExpandRows bool) string {
	if s.Content == "" {
		s.Content = "<empty>"
	}
	chunks := strings.Split(s.Content, "\n")
	//if s.Ext == "url" {
	//	return uri(s.Snip)
	//}
	code := []codeLine{}
	// Add line numbers and pad
	for i, v := range chunks {
		code = append(code, codeLine{
			Margin: style.FStart(style.Subdued, fmt.Sprintf("%3d ", i+1)),
			Body:   fmt.Sprintf("    %s", strings.Replace(v, "\t", "  ", -1)),
		})
	}

	lastLine := code[len(chunks)-1]

	// Add to  starting from most important line
	marker := mainMarkers[s.Ext()]
	if marker != "" {
		var clipped []codeLine
		var startPreview int
		for i, v := range code {
			if strings.Contains(v.Body, marker) {
				startPreview = i
			}
		}
		for i, v := range code {
			if startPreview <= i {
				clipped = append(clipped, v)
			}
		}
		code = clipped
	}

	crop := len(code) >= lines && lines != 0

	// crop width
	var codeview []codeLine
	if crop {
		codeview = code[0:lines]
	} else {
		codeview = code
	}

	rightTrim := style.FStart(style.Subdued, "|")
	if width > 0 {
		for i, v := range codeview {
			codeview[i].Body = pad(width, v.Body).String() + rightTrim
		}
	}

	// Add page tear and last line
	if alwaysExpandRows && crop && lines < len(code) {
		codeview = append(codeview, codeLine{
			style.FStart(style.Subdued, "----"),
			style.FStart(style.Subdued, strings.Repeat("-", width)+"|"),
		})
		lastLine.Body = pad(width, lastLine.Body).String() + rightTrim
		codeview = append(codeview, lastLine)
	}

	buff := bytes.Buffer{}
	for i, v := range codeview {
		// Style
		var m, b string
		if odd {
			m = style.FmtFgBg(v.Margin, style.OffWhite248, style.Grey240)
			b = style.FmtFgBg(v.Body, style.OffWhite250, style.Grey238)
		} else {
			m = style.FmtFgBg(v.Margin, style.OffWhite248, style.Grey238)
			b = style.FmtFgBg(v.Body, style.OffWhite250, style.Grey236)
		}
		buff.WriteString(m)
		buff.WriteString(b)
		if i < len(codeview)-1 {
			buff.WriteString("\n")
		}
	}
	return buff.String()
}

func fTags(tags []string) string {
	if len(tags) == 0 {
		return fEmpty("")
	}
	return strings.Join(tags, ", ")
}

func fEmpty(in string) string {
	if in == "" {
		return style.Fmt16(style.Subdued, "<none>")
	}
	return in
}
