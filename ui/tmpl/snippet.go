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
	fmtHeader(w, s.Private, s.Username, s.Pouch, &s.SnipName)

	p := tablewriter.NewWriter(w)
	p.SetAutoWrapText(false)
	p.SetBorder(false)
	p.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	p.SetCenterSeparator("")
	p.SetColumnSeparator("")
	p.SetRowLine(false)
	p.SetAutoFormatHeaders(false)
	p.SetHeaderLine(false)
	p.SetColWidth(5)
	p.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	p.Append([]string{ FmtSnippet(s, 100, 0)})
	p.Render()

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

	tbl.Append([]string{style.Fmt(style.Cyan,"Snippet Details:"), "", "", ""})
	tbl.Append([]string{
		style.Fmt(style.Subdued,"Run Status:"), pad(20, statusString(s, true)).String(),
		style.Fmt(style.Subdued,"Last Run:"), pad(20, humanize.Time(time.Unix(s.RunStatusTime, 0))).String(),
	})
	tbl.Append([]string{
		style.Fmt(style.Subdued,"Run Count: "), fmt.Sprintf("↻ %2d", s.RunCount),
		style.Fmt(style.Subdued,"View count:") , fmt.Sprintf("🔦  %2d", s.ViewCount )}) //👁 👀
	tbl.Append([]string{
		style.Fmt(style.Subdued,"Description:"), fmtEmpty(s.Description), "", ""})
	tbl.Append([]string{
		style.Fmt(style.Subdued,"Tags:"), fmtTags(s.Tags), "", ""})
	tbl.Append([]string{
		style.Fmt(style.Subdued,"sha256:"), fmtVerified(s) })

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
	if s.VerifySnippet() {
		buff.WriteString(style.Fmt(style.Green, "✓  "))
		buff.WriteString(pad(12, s.CheckSum).String())
		buff.WriteString("...")
	} else {
		buff.WriteString(" ☠  Invalid Checksum: ")
		buff.WriteString(fmtEmpty(s.CheckSum))
	}
	return buff.String()
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
