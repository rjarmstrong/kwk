package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"github.com/rjarmstrong/tablewriter"
	"github.com/rjarmstrong/go-humanize"
	"bytes"
	"fmt"
	"time"
	"strings"
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
)

func inspect(s *models.Snippet) string {

	w := &bytes.Buffer{}
	w.WriteString("\n")
	w.WriteString(MARGIN)
	fmtHeader(w,  s.Username, s.Pouch, &s.SnipName)
	w.WriteString(strings.Repeat(" ", 4))
	w.WriteString(printIcon(s))
	if s.IsApp() {
		w.WriteString(style.Fmt(style.Subdued,"  App"))
	} else if s.Ext == "url" {
		w.WriteString(style.Fmt(style.Subdued, "  Link"))
	} else {
		w.WriteString(style.Fmt(style.Subdued, "  Snippet"))
	}
	fmt.Fprint(w,"\n")
	fmt.Fprint(w, TWOLINES)
	fmt.Fprint(w, FmtSnippet(s, 100, 0, false))
	fmt.Fprint(w,"\n\n")

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

	if s.IsApp() {
		tbl.Append([]string{style.Fmt(style.Cyan, "App Details:"), "", "", ""})
	} else if s.Ext == "url" {
		tbl.Append([]string{style.Fmt(style.Cyan, "Link Details:"), "", "", ""})
	} else {
		tbl.Append([]string{style.Fmt(style.Cyan, "Snippet Details:"), "", "", ""})
	}

	var lastRun string
	if s.RunCount < 1 {
		lastRun = "never"
	} else {
		lastRun = pad(20, humanize.Time(time.Unix(s.RunStatusTime, 0))).String()
	}
	tbl.Append([]string{
		style.Fmt(style.Subdued,"Run Status:"), pad(20, printStatus(s, true)).String(),
		style.Fmt(style.Subdued,"Last Run:"), lastRun,
	})
	tbl.Append([]string{
		style.Fmt(style.Subdued,"Run Count: "), fmt.Sprintf("↻ %2d", s.RunCount),
		style.Fmt(style.Subdued,"View count:") , fmt.Sprintf("🔦  %2d", s.ViewCount )}) //👁 👀
	if s.IsApp() {
		tbl.Append([]string{
			style.Fmt(style.Subdued,"App Dependencies:"), strings.Join(s.Dependencies, ", "), "", ""})
	}
	tbl.Append([]string{
		style.Fmt(style.Subdued,"Description:"), style.FmtBox(fmtEmpty(s.Description), 25, 3), "", ""})

	tbl.Append([]string{
		style.Fmt(style.Subdued,"Preview:"), style.FmtPreview(s.Preview, 25, 2), "", ""})

	tbl.Append([]string{
		style.Fmt(style.Subdued,"Tags:"), fmtTags(s.Tags), "", ""})
	tbl.Append([]string{
		style.Fmt(style.Subdued,"sha256:"), fmtVerified(s) })
	tbl.Append([]string{
		style.Fmt(style.Subdued,"Updated:"), humanTime(s.Created) })

	tbl.Render()



	//fmt.Fprint(w, style.Start)
	//fmt.Fprintf(w, "%dm", style.Subdued)
	//fmt.Fprint(w, MARGIN)
	//fmt.Fprint(w,"Snippet details: `kwk <name>`")
	//fmt.Fprint(w, MARGIN)
	//fmt.Fprint(w,"Run snippet: `kwk run <name>`")
	//fmt.Fprint(w, MARGIN)
	//fmt.Fprint(w, style.End)
	fmt.Fprint(w,"\n")

	return w.String()
}
func fmtVerified(s *models.Snippet) string {
	var buff bytes.Buffer
	if s.VerifyChecksum() {
		buff.WriteString(style.Fmt(style.Green, "✓  "))
		buff.WriteString(pad(12, s.CheckSum).String())
		buff.WriteString("...")
	} else {
		buff.WriteString(" ☠  Invalid Checksum: ")
		buff.WriteString(fmtEmpty(s.CheckSum))
	}
	return buff.String()
}

func FmtSnippet(s *models.Snippet, width int, lines int, odd bool) string {
	if s.Snip == "" {
		s.Snip = "<empty>"
	}
	chunks := strings.Split(s.Snip, "\n")
	//if s.Ext == "url" {
	//	return uri(s.Snip)
	//}
	code := []CodeLine{}
	// Add line numbers and pad
	for i, v := range chunks {
		code = append(code, CodeLine{
			Margin: style.FmtStart(style.Subdued, fmt.Sprintf("%3d ", i+1)),
			Body:   fmt.Sprintf("    %s", strings.Replace(v, "\t", "  ", -1)),
		})
	}

	lastLine := code[len(chunks)-1]

	// Add to preview starting from most important line
	marker := mainMarkers[s.Ext]
	if marker != "" {
		var clipped []CodeLine
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
	var preview []CodeLine
	if crop {
		preview = code[0:lines]
	} else {
		preview = code
	}

	rightTrim := style.FmtStart(style.Subdued, "|")
	if width > 0 {
		for i, v := range preview {
			preview[i].Body = pad(width, v.Body).String() + rightTrim
		}
	}

	// Add page tear and last line
	if models.Prefs().AlwaysExpandRows && crop && lines < len(code) {
		preview = append(preview, CodeLine{
			style.FmtStart(style.Subdued, "----"),
			style.FmtStart(style.Subdued, strings.Repeat("-", width)+"|"),
		})
		lastLine.Body = pad(width, lastLine.Body).String() + rightTrim
		preview = append(preview, lastLine)
	}

	buff := bytes.Buffer{}
	for i, v := range preview {
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
		buff.WriteString(" ")
		if i < len(preview) - 1 {
			buff.WriteString("\n")
		}
	}
	return buff.String() // fmt.Sprintf("%q", buff.String())
}

func fmtTags(tags []string) string {
	if len(tags) == 0 {
		return fmtEmpty("")
	}
	return strings.Join(tags, ", ")
}

func fmtEmpty(in string) string {
	if in == "" {
		return "<none>"
	}
	return in
}
