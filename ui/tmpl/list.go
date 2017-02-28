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

func statusString(s *models.Snippet, includeText bool) string {
	if s.Ext == "url" {
		if includeText {
			return "🌎  Bookmark"
		}
		return "🌎"
	}
	if s.RunStatus == models.UseStatusSuccess {
		if includeText {
			return "⚡  Success"
		}
		return "⚡" //"✓"//
	} else if s.RunStatus == models.UseStatusFail {
		if includeText {
			return "🔥  Error"
		}
		return "🔥" //style.Fmt(style.Red, "●") //
	}
	if includeText {
		return "📄  Not run/runnable"
	}
	return "📄" //"🔸"
}

func fmtLocked(locked bool, includeText bool) string {
	if locked {
		if includeText {
			return "🔒  Private"
		}
		return "🔒"
	}
	if includeText {
		return "🔓  Public"
	}
	return "🔓"
}

func listRoot(r *models.ListView) string {
	w := &bytes.Buffer{}
	var all []interface{}
	for _, v := range r.Pouches {
		if v.Name != "" {
			all = append(all, v)
		}
	}
	for _, v := range r.Personal {
		all = append(all, v)
	}

	fmtHeader(w, r.HidePrivate, r.Username, r.Pouch, nil)
	w.Write(listHorizontal(all))

	if len(r.Snippets) > 0 {
		fmt.Fprint(w, listPouchSnippets(r))
	}

	//w.WriteString(style.Fmt(style.Subdued, fmt.Sprintf("%d/50 Pouches", len(r.Pouches)-1)))
	if models.ClientIsNew(r.LastUpdate) {
		w.WriteString(style.Fmt(style.Subdued, fmt.Sprintf("          kwk auto-updated to %s %s", models.Client.Version, humanTime(r.LastUpdate))))
	} else {
		w.WriteString("\n")
	}
	w.WriteString("\n\n")
	return w.String()
}

func listPouch(list *models.ListView) string {
	w := &bytes.Buffer{}
	fmtHeader(w, list.HidePrivate, list.Username, list.Pouch, nil)
	fmt.Fprint(w, listPouchSnippets(list))
	return w.String()
}

func listPouchSnippets(list *models.ListView) string {
	w := &bytes.Buffer{}

	//fmt.Fprintf(w,"%s", fmtLocked(list.HidePrivate, false))
	//if list.Pouch != "" {
	//	fmt.Fprintf(w, "%s👝  \n", MARGIN)
	//}

	tbl := tablewriter.NewWriter(w)
	tbl.SetHeader([]string{"", "", "", ""})
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
			executed = fmt.Sprintf("%s  %s", statusString(v, false), style.Fmt(style.Subdued, humanize.Time(time.Unix(v.RunStatusTime, 0))))
		} else {
			executed = statusString(v, false)
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
			name.WriteString(style.Fmt(style.Subdued, "short-cut: kwk edit env"))
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
		fmt.Fprint(w, MARGIN)
		fmt.Fprint(w, style.Fmt(style.Subdued, "<empty pouch>\n"))
	}

	fmt.Fprint(w, "\n")

	//fmt.Fprint(w, style.Start)
	//fmt.Fprintf(w, "%dm", style.Subdued)
	//fmt.Fprint(w, MARGIN)
	//fmt.Fprintf(w,"Expand list: `kwk expand %s`", list.Pouch)
	//fmt.Fprint(w, MARGIN)
	//fmt.Fprint(w, style.End)
	////fmt.Fprint(w, style.Start)
	//fmt.Fprintf(w, "%dm", style.Subdued)
	//fmt.Fprint(w, MARGIN)
	//fmt.Fprintf(w, "%d of max 32 snippets in pouch", len(list.Snippets))
	//fmt.Fprint(w, style.End)


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

func fmtHeader(w io.Writer, locked bool, username string, pouch string, s *models.SnipName) {
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, MARGIN)
	fmt.Fprint(w, style.Start)
	fmt.Fprintf(w, "%dm", style.Cyan)
	fmt.Fprint(w, KWK_HOME)
	fmt.Fprint(w, "/")
	if pouch == "" && s == nil {
		fmt.Fprint(w, style.End)
		fmt.Fprint(w, username)
		fmt.Fprint(w, FOOTER)
		return
	}
	fmt.Fprint(w, username)
	fmt.Fprint(w, "/")
	if s == nil {
		fmt.Fprint(w, style.End)
		fmt.Fprint(w,  pouch)
		fmt.Fprint(w, FOOTER)
		return
	}
	if pouch != "" {
		fmt.Fprint(w, pouch)
		fmt.Fprint(w,"/")
	}
	fmt.Fprint(w, style.End)
	fmt.Fprint(w, s.String())
	fmt.Fprint(w, FOOTER)
}

func FmtOutPreview(s *models.Snippet) string {
	chunks := strings.Split(s.Preview, "\n")
	lines := []string{}
	for i := 0; i < len(chunks) && i < 3; i ++ {
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
			item.WriteString("  ")
		}
		if sn, ok := v.(*models.Snippet); ok {
			item.WriteString(statusString(sn, false))
			item.WriteString("  ")
			item.WriteString(style.Fmt(style.Cyan, sn.SnipName.Name))
			item.WriteString(style.Fmt(style.Subdued, "."+sn.SnipName.Ext))
			item.WriteString(" ")
		}
		if pch, ok := v.(*models.Pouch); ok {
			if models.Prefs().ListAll || !pch.MakePrivate {
				if pch.Name == "inbox" {
					if pch.UnOpened > 0 {
						item.WriteString(fmt.Sprintf("📬%d", pch.UnOpened))
					} else {
						item.WriteString("📪")
					}
				} else if pch.Name == "settings" {
					item.WriteString("⚙")
				} else if pch.MakePrivate {
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
