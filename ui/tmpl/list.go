package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"github.com/rjarmstrong/tablewriter"
	"github.com/rjarmstrong/go-humanize"
	"strings"
	"bytes"
	"time"
	"fmt"
	"io"
	"sort"
	"golang.org/x/text/unicode/norm"
)

var mainMarkers = map[string]string{
	"go": "func main() {",
}

type CodeLine struct {
	Margin string
	Body   string
}

func StatusText(s *models.Snippet) string {
	if s.Ext == "url" {
		return "bookmark"
	}
	if s.RunStatus == models.UseStatusSuccess {
		return "success"
	} else if s.RunStatus == models.UseStatusFail {
		return "error"
	}
	return "static"
}

func printStatus(s *models.Snippet, includeText bool) string {
	if s.RunStatus == models.UseStatusSuccess {
		if includeText {
			return style.Fmt(style.Green, "✔") + "  Success"
		}
		return style.Fmt(style.Green, "✔") //"⚡" //"✓"//
	} else if s.RunStatus == models.UseStatusFail {
		if includeText {
			return "🔥  Error"
		}
		return "🔥" //style.Fmt(style.Red, "●") //
	}
	return "-" //"🔸"
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

	fmtHeader(w, r.Username, "", nil)
	fmt.Fprint(w, strings.Repeat(" ", 65), style.Fmt(style.Subdued, "◉  " + models.Principal.Username, "    TLS12")) //TLS ⇨ᗜ 🔑
	fmt.Fprint(w, TWOLINES)
	w.Write(listHorizontal(all))
	fmt.Fprint(w, "\n\n", MARGIN, style.Fmt(style.Subdued, "Community"),  "\n")

	com := []interface{}{}
	com = append(com, &models.Pouch{
		Name: style.Fmt(style.Cyan, "/kwk/") + "unicode",
		Username: "kwk",
		SnipCount: 12,
	}, &models.Pouch{
		Name:      style.Fmt(style.Cyan, "/kwk/") +"news",
		Username:  "kwk",
		SnipCount: 10,
	},
	&models.Pouch{
		Name:      style.Fmt(style.Cyan, "/kwk/") +"devops",
		Username:  "kwk",
		SnipCount: 18,
	})
	w.Write(listHorizontal(com))
	w.WriteString("\n")

	if len(r.Snippets) > 0 {
		fmt.Fprint(w, listSnippets(r))
	}

	//w.WriteString(style.Fmt(style.Subdued, fmt.Sprintf("%d/50 Pouches", len(r.Pouches)-1)))
	if models.ClientIsNew(r.LastUpgrade) {
		w.WriteString(style.Fmt(style.Subdued, fmt.Sprintf("          kwk auto-updated to %s %s", models.Client.Version, humanTime(r.LastUpgrade))))
	} else {
		w.WriteString("\n")
	}
	w.WriteString("\n\n")
	return w.String()
}

func printPouchHeadAndFoot(w *bytes.Buffer, list *models.ListView) {
	fmtHeader(w, list.Username, list.Pouch.Name, nil)
	fmt.Fprint(w, MARGIN, MARGIN, fmtLocked(list.Pouch.MakePrivate, true))
	fmt.Fprint(w, " Pouch")
	fmt.Fprint(w, MARGIN, MARGIN, len(list.Snippets), " snippets")
	fmt.Fprint(w, "\n")
}

func listPouch(list *models.ListView) string {
	w := &bytes.Buffer{}
	if models.Prefs().Naked {
		fmt.Fprint(w, listNaked(list))
	} else {
		sort.Slice(list.Snippets, func(i, j int) bool {
			return list.Snippets[i].Created <= list.Snippets[j].Created
		})
		if list.Pouch != nil {
			printPouchHeadAndFoot(w, list)
		}
		fmt.Fprint(w, listSnippets(list))
		if list.Pouch != nil && len(list.Snippets) > 10 {
			printPouchHeadAndFoot(w, list)
		}
		fmt.Fprint(w, "\n")
	}
	return w.String()
}

const timeLayout = "2 Jan 15:04 06"
func listNaked(list *models.ListView) interface{} {
	w := &bytes.Buffer{}
	tbl := tablewriter.NewWriter(w)
	tbl.SetHeader([]string{"Name", "Private", "Run status", "Run count", "View count", "Updated"})
	tbl.SetAutoWrapText(false)
	tbl.SetBorder(false)
	tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator("")
	tbl.SetRowLine(false)
	tbl.SetAutoFormatHeaders(false)
	tbl.SetHeaderLine(false)
	tbl.SetColWidth(5)
	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	for _, v := range list.Snippets {
		t := time.Unix(v.Created, 0)
		var private string
		if v.Private {
			private = "private"
		} else {
			private = "public"
		}
		tbl.Append([]string{
			v.String(),
			private,
			StatusText(v),
			fmt.Sprintf("%d", v.RunCount),
			fmt.Sprintf("%d", v.ViewCount),
			t.Format(timeLayout),
		})
	}
	tbl.Render()
	return w.String()
}

func listSnippets(list *models.ListView) string {
	w := &bytes.Buffer{}

	if len(list.Snippets) == 0 {
		fmt.Fprint(w, "\n")
		fmt.Fprint(w, MARGIN)
		fmt.Fprint(w, style.Fmt(style.Subdued, "<empty pouch>"))
		fmt.Fprint(w, TWOLINES)
		fmt.Fprint(w, MARGIN)
		fmt.Fprint(w, style.Fmt(style.Cyan, "Add new snippets to this pouch: "))
		fmt.Fprintf(w, "`kwk new <snippet> %s/<name>.<ext>`", list.Pouch.Name)
		fmt.Fprint(w, "\n")
		return w.String()
	}

	tbl := tablewriter.NewWriter(w)
	tbl.SetHeader([]string{"", "", "", ""})
	tbl.SetAutoWrapText(false)
	tbl.SetBorder(false)
	tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator(" ")
	if models.Prefs().RowLines {
		tbl.SetRowLine(true)
	}

	tbl.SetAutoFormatHeaders(false)
	tbl.SetHeaderLine(true)
	tbl.SetColWidth(5)

	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for i, v := range list.Snippets {
		var executed string
		if v.RunStatusTime > 0 {
			executed = fmt.Sprintf("%s  %s", printStatus(v, false), style.Fmt(style.Subdued, humanize.Time(time.Unix(v.RunStatusTime, 0))))
		} else {
			executed = printStatus(v, false)
		}

		// col1
		name := &bytes.Buffer{}
		name.WriteString(printIcon(v))
		name.WriteString("  ")
		nt := style.Fmt(style.Cyan, v.SnipName.String())
		name.WriteString(nt)
		if models.Prefs().AlwaysExpandRows {
			name.WriteString("\n\n")
			name.WriteString(style.FmtBox(v.Description, 25, 3))
		}
		// Add special instructions:
		if v.Role == models.SnipRoleEnvironment {
			name.WriteString("\n\n")
			name.WriteString(style.Fmt(style.Subdued, "short-cut: kwk edit env"))
		}

		// col2
		var lines int
		if models.Prefs().AlwaysExpandRows {
			lines = models.Prefs().ExpandedRows
		} else {
			lines = models.Prefs().SlimRows
		}
		status := &bytes.Buffer{}
		if v.RunCount > 0 {
			status.WriteString(executed)
			status.WriteString(" ")
		}
		status.WriteString(fmtRunCount(v.RunCount))

		//col3
		snip := FmtSnippet(v, 60, lines, (i+1)%2==0)
		if models.Prefs().RowSpaces {
			snip = snip + "\n"
		}

		// //strings.Join(v.Tags, ", "),
		tbl.Append([]string{
			name.String(),
			status.String(),
			snip,
			style.FmtBox(v.Preview, 18, 1),
		})
	}
	tbl.Render()


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

	return w.String()
}

func printIcon(v *models.Snippet) string {
	if v.IsApp() {
		return "💫" // 📦
	} else if v.Ext == "url" {
		return "🌎"
	} else if v.RunCount > 0 {
		return  "🔸" //"⚡" //☰
	} else {
		return "📄"
	}
}

func fmtRunCount(count int64) string {
	return style.Fmt(style.Subdued, fmt.Sprintf("↻ %d", count))
}

func fmtHeader(w io.Writer, username string, pouch string, s *models.SnipName) {
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, MARGIN)
	fmt.Fprint(w, style.Start)
	fmt.Fprint(w, "7m")
	fmt.Fprint(w, " ❯ ")
	fmt.Fprint(w, KWK_HOME)
	fmt.Fprint(w, "/")
	if pouch == "" && s == nil {
		fmt.Fprint(w, style.Start)
		fmt.Fprint(w, "1m")
		fmt.Fprint(w, username)
		fmt.Fprint(w, " ")
		fmt.Fprint(w, style.End)
		return
	}
	fmt.Fprint(w, username)
	fmt.Fprint(w, "/")
	if s == nil {
		fmt.Fprint(w, style.Start)
		fmt.Fprint(w, "1m")
		fmt.Fprint(w,  pouch)
		fmt.Fprint(w, " ")
		fmt.Fprint(w, style.End)
		return
	}
	if pouch != "" {
		fmt.Fprint(w, pouch)
		fmt.Fprint(w,"/")
	}
	fmt.Fprint(w, style.Start)
	fmt.Fprint(w, "1m")
	fmt.Fprint(w, s.String())
	fmt.Fprint(w, " ")
	fmt.Fprint(w, style.End)
}

func pad(width int, in string) *bytes.Buffer {
	buff := &bytes.Buffer{}
	diff := width - len([]rune(in))
	if diff > 0 {
		buff.WriteString(in)
		buff.WriteString(strings.Repeat(" ", diff))
	} else {
		var ia norm.Iter
		ia.InitString(norm.NFKD, in)
		nc := 0
		for !ia.Done() && nc < width  {
			nc += 1
			buff.Write(ia.Next())
		}
	}
	return buff
}