package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"github.com/rjarmstrong/tablewriter"
	"github.com/rjarmstrong/go-humanize"
	"text/tabwriter"
	"strings"
	"bytes"
	"time"
	"fmt"
	"io"
)

var mainMarkers = map[string]string{
	"go": "func main() {",
}

type CodeLine struct {
	Margin string
	Body   string
}

func statusString(s *models.Snippet) string {
	if s.Ext == "url" {
		return "🌎"
	}
	if s.RunStatus == models.RunStatusSuccess {
		return "⚡"  //"✓"//
	} else if s.RunStatus == models.RunStatusFail {
		return "🔥" //style.Fmt(style.Red, "●") //
	}
	return "📄" //"🔸"
}

func listRoot(r *models.ListView) string {
	buff := &bytes.Buffer{}

	fmtHeader(buff, r)

	var all []interface{}
	for _, v := range r.Pouches {
		if v.Name != "" {
			all = append(all, v)
		}
	}

	buff.WriteString(listPouch(r))
	buff.Write(listHorizontal(all))

	buff.WriteString(style.Fmt(style.Subdued, fmt.Sprintf("\n\n   %d/50 Pouches", len(r.Pouches)-1)))
	if models.ClientIsNew(r.LastUpdate) {
		buff.WriteString(style.Fmt(style.Subdued, fmt.Sprintf("          kwk auto-updated to %s %s", models.Client.Version, humanTime(r.LastUpdate))))
	} else {
		buff.WriteString("\n")
	}
	buff.WriteString("\n\n")
	for _, v := range r.Personal {
		if v.Name == "inbox" {
			if v.UnOpened > 0 {
				buff.WriteString(fmt.Sprintf(" 📬 Inbox %d", v.UnOpened))

			} else {
				buff.WriteString(" 📪  inbox")
			}
		} else if v.Name == "settings" {
			buff.WriteString("   ⚙  settings")
		}
	}
	buff.WriteString("\n\n")
	return buff.String()
}

func listPouch(list *models.ListView) string {
	w := &bytes.Buffer{}
	fmtHeader(w, list)

	tbl := tablewriter.NewWriter(w)
	tbl.SetHeader([]string{"Name", "Status", "Snippet", "Preview"})
	tbl.SetAutoWrapText(false)
	tbl.SetBorder(false)
	tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator(" ")
	tbl.SetRowLine(true)
	tbl.SetAutoFormatHeaders(false)
	tbl.SetHeaderLine(true)
	tbl.SetColWidth(5)

	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for _, v := range list.Snippets {
		var executed string
		if v.RunStatusTime > 0 {
			executed = fmt.Sprintf("%s  %s", statusString(v), style.Fmt(style.Subdued, humanize.Time(time.Unix(v.RunStatusTime, 0))))
		} else {
			executed = statusString(v)
		}

		// col1
		name := &bytes.Buffer{}
		name.WriteString(style.Fmt(style.Cyan, v.SnipName.String()))
		name = pad(25, name.String())
		if models.Prefs().AlwaysExpandLists {
			name.WriteString("\n\n")
			fmtDescription(name, v.Description, 20)
		}
		// Add special instructions:
		if v.Role == models.SnipRoleEnvironment {
			name.WriteString("\n\n")
			name.WriteString(style.Fmt(style.Subdued,"short-cut: kwk edit env"))
		}

		// col2
		var lines int
		if models.Prefs().AlwaysExpandLists {
			lines = models.Prefs().ExpandedLines
		} else {
			lines = models.Prefs().SlimLines
		}
		status := &bytes.Buffer{}
		status.WriteString(executed)
		status.WriteString("\n")
		status.WriteString(style.Fmt(style.Subdued, fmt.Sprintf("↻ %2d", v.RunCount)))

		//col3

		// //strings.Join(v.Tags, ", "),
		tbl.Append([]string{
			name.String(),
			status.String(),
			FmtSnippet(v, 60, lines),
			FmtOutPreview(v),
		})
	}
	tbl.Render()
	if len(list.Snippets) == 0 {
		fmt.Println(style.Fmt(style.Yellow, "Create some snippets to fill this view!\n"))
	}
	fmt.Fprintf(w, "\n%d of %d records\n\n", len(list.Snippets), list.Total)
	fmt.Fprint(w, "\n\n")

	return w.String()
}

func fmtDescription(w io.Writer, in string, width int) {
	t := strings.Split(style.WrapString(in, width), "\n")
	for i, v := range t {
		t[i] = style.Fmt(style.Subdued, v)
	}
	join := strings.Join(t, "\n")
	fmt.Fprint(w, join)
}

func fmtHeader(w io.Writer, list *models.ListView) {
	fmt.Print(w, "\n")
	fmt.Print(MARGIN)
	fmt.Print(w, style.Fmt(style.Cyan, "kwk.co/"+list.Username+"/"))
	if !list.IsRoot {
		fmt.Print(w, list.Pouch)
		fmt.Print(w, "/")
	}
	fmt.Print(w, "\n\n")
}

func FmtOutPreview(s *models.Snippet) string{
	chunks := strings.Split(s.Preview, "\n")
	lines := []string{}
	for i :=0; i < len(chunks) && i < 3; i ++ {
		if chunks[i] == "" {
			continue
		}
		l := style.Fmt(style.Subdued, pad(20, chunks[i]).String())
		lines = append(lines, l)
	}
	return strings.Join(lines, "\n")
}

func FmtSnippet(s *models.Snippet, width int, lines int) string {
	if s.Snip == "" {
		s.Snip = "<empty>"
	}
	chunks := strings.Split(s.Snip, "\n")

	//// Return any non code previews
	//if s.Role == models.SnipRolePreferences {
	//	return `(Global prefs) 'kwk edit prefs'`
	//} else if s.Role == models.SnipRoleEnvironment {
	//	return `(Local environment) 'kwk edit env'`
	//} else

	if s.Ext == "url" {
		return uri(s.Snip)
	}

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
	if models.Prefs().AlwaysExpandLists && crop && lines < len(code) {
		preview = append(preview, CodeLine{
			style.FmtStart(style.Subdued, "----"),
			style.FmtStart(style.Subdued, strings.Repeat("-", width)+"|"),
		})
		lastLine.Body = pad(width, lastLine.Body).String() + rightTrim
		preview = append(preview, lastLine)
	}

	buff := bytes.Buffer{}
	for _, v := range preview {
		// Style
		m := style.Fmt256(style.GreyBg238, true, v.Margin)
		b := style.Fmt256(style.GreyBg236, true, v.Body)
		buff.WriteString(m)
		buff.WriteString(b)
		buff.WriteString("  ")
		buff.WriteString("\n")
	}
	return buff.String()
}

func pad(width int, in string) *bytes.Buffer {
	buff := &bytes.Buffer{}
	diff := width - len(in)
	if diff > 0 {
		buff.WriteString(in)
		buff.WriteString(strings.Repeat(" ", diff))
	} else {
		buff.WriteString(in[0:width])
	}
	return buff
}

func listHorizontal(l []interface{}) []byte {
	var buff bytes.Buffer
	w := tabwriter.NewWriter(&buff, 20, 3, 2, ' ', tabwriter.DiscardEmptyColumns)
	var item bytes.Buffer
	for i, v := range l {
		if i%5 == 0 {
			item.WriteString("   ")
		}
		if sn, ok := v.(*models.Snippet); ok {
			item.WriteString(statusString(sn))
			item.WriteString("  ")
			item.WriteString(style.Fmt(style.Cyan, sn.SnipName.Name))
			item.WriteString(style.Fmt(style.Subdued, "."+sn.SnipName.Ext))
			item.WriteString(" ")
		}
		if pch, ok := v.(*models.Pouch); ok {
			if models.Prefs().ListAll || !pch.MakePrivate {
				if pch.MakePrivate {
					item.WriteString("🔒")
				} else {
					item.WriteString("👝")
				}
				item.WriteString("  ")
				item.WriteString(pch.Name)
				item.WriteString(style.Fmt(style.Subdued, fmt.Sprintf(" (%d)", pch.SnipCount)))
			}
		}

		item.WriteString(" \t")
		x := i + 1
		if x%5 == 0 {
			item.WriteString("\n")
		}
		if x%20 == 0 {
			item.WriteString("\n")
		}
		fmt.Fprint(w, fmt.Sprintf("%s", item.String()))
		item.Reset()
	}
	w.Flush()
	return buff.Bytes()
}
