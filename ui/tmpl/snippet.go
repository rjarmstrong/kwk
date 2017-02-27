package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"github.com/rjarmstrong/tablewriter"
	"github.com/rjarmstrong/go-humanize"
	"bytes"
	"fmt"
	"time"
	"strings"
)

func inspect(s *models.Snippet) string {

	w := &bytes.Buffer{}
	fmtHeader(w, s.Username, s.Pouch, &s.SnipName)

	tbl := tablewriter.NewWriter(w)
	tbl.SetHeader([]string{"", ""})
	tbl.SetAutoWrapText(false)
	tbl.SetBorder(false)
	tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator("|")
	tbl.SetRowLine(true)
	tbl.SetAutoFormatHeaders(false)
	tbl.SetHeaderLine(false)
	tbl.SetColWidth(5)
	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	tbl.Append([]string{"Run Status", statusString(s)})
	tbl.Append([]string{"Last Run", humanize.Time(time.Unix(s.RunStatusTime, 0))})
	tbl.Append([]string{"Run Count ", fmt.Sprintf("↻ %2d", s.RunCount)})
	tbl.Append([]string{"Description", fmtEmpty(s.Description)})
	tbl.Append([]string{"Tags", strings.Join(s.Tags, ", ")})
	tbl.Append([]string{"Snippet", FmtSnippet(s, 100, 0)})

	tbl.Render()
	fmt.Fprint(w, )
	return w.String()
}

func fmtEmpty(in string) string {
	if in == "" {
		return "<none>"
	}
	return in
}
