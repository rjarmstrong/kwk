package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"github.com/rjarmstrong/tablewriter"
	"bytes"
	"fmt"
	"time"
	"strings"
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
)

func view(s *models.Snippet) string {
	w := &bytes.Buffer{}
	w.WriteString("\n")
	w.WriteString(style.MARGIN)
	fmtHeader(w,  s.Username, s.Pouch, &s.SnipName)
	w.WriteString(strings.Repeat(" ", 4))
	w.WriteString(snippetIcon(s))
	w.WriteString("  ")
	w.WriteString(FSnippetType(s))
	fmt.Fprint(w,"\n")
	fmt.Fprint(w, style.TWOLINES)
	fmt.Fprint(w, FCodeview(s, 100, 0, false))
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
	tbl.Append([]string{style.Fmt256(style.Color_PouchCyan, FSnippetType(s) + " Details:"), "", "", ""})

	var lastRun string
	if s.RunCount < 1 {
		lastRun = "never"
	} else {
		lastRun = pad(20, style.Time(time.Unix(s.RunStatusTime, 0))).String()
	}
	tbl.Append([]string{
		style.Fmt16(style.Subdued,"Run Status:"), FStatus(s, true),
		style.Fmt16(style.Subdued,"Last Run:"), lastRun,
	})
	tbl.Append([]string{
		style.Fmt16(style.Subdued,"Run Count: "), fmt.Sprintf("↻ %2d", s.RunCount),
		style.Fmt16(style.Subdued,"View count:") , fmt.Sprintf("%s  %2d", style.Icon_View, s.ViewCount )})
	if s.IsApp() {
		tbl.Append([]string{
			style.Fmt16(style.Subdued,"App Deps:"),
			style.FBox(strings.Join(s.Dependencies, ", "), 50, 5)})
	}
	tbl.Append([]string{
		style.Fmt16(style.Subdued,"Description:"), style.FBox(FEmpty(s.Description), 50, 3), "", ""})

	tbl.Append([]string{
		style.Fmt16(style.Subdued,"Preview:"), style.FPreview(s.Preview, 50, 1), "", ""})

	tbl.Append([]string{
		style.Fmt16(style.Subdued,"Tags:"), FTags(s.Tags), "", ""})
	tbl.Append([]string{
		style.Fmt16(style.Subdued,"sha256:"), FVerified(s) })
	tbl.Append([]string{
		style.Fmt16(style.Subdued,"Updated:"), fmt.Sprintf("%s - %s", humanTime(s.Created), fmt.Sprintf("v%d", s.Version )),
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
	fmt.Fprint(w,"\n")

	return w.String()
}
func FSnippetType(s *models.Snippet) string {
	if s.IsApp() {
		return "App"
	} else if s.Ext == "url" {
		return "Bookmark"
	} else {
		return "Snippet"
	}
}
func FVerified(s *models.Snippet) string {
	var buff bytes.Buffer
	if s.VerifyChecksum() {
		buff.WriteString(style.Fmt256(style.Color_YesGreen,  style.Icon_Tick + " "))
		buff.WriteString(pad(12, s.CheckSum).String())
		buff.WriteString("...")
	} else {
		buff.WriteString(" " + style.Icon_Cross + "  Invalid Checksum: ")
		buff.WriteString(FEmpty(s.CheckSum))
	}
	return buff.String()
}

func FCodeview(s *models.Snippet, width int, lines int, odd bool) string {
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
			Margin: style.FStart(style.Subdued, fmt.Sprintf("%3d ", i+1)),
			Body:   fmt.Sprintf("    %s", strings.Replace(v, "\t", "  ", -1)),
		})
	}

	lastLine := code[len(chunks)-1]

	// Add to  starting from most important line
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
	var codeview []CodeLine
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
	if models.Prefs().AlwaysExpandRows && crop && lines < len(code) {
		codeview = append(codeview, CodeLine{
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
		if i < len(codeview) - 1 {
			buff.WriteString("\n")
		}
	}
	return buff.String()
}

func FTags(tags []string) string {
	if len(tags) == 0 {
		return FEmpty("")
	}
	return strings.Join(tags, ", ")
}

func FEmpty(in string) string {
	if in == "" {
		return style.Fmt16(style.Subdued, "<none>")
	}
	return in
}
