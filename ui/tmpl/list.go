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
	fmt.Fprint(w, strings.Repeat(" ", 90), "👤  ", models.Principal.Username)
	fmt.Fprint(w, FOOTER)
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
		printPouchHeadAndFoot(w, list)
		fmt.Fprint(w, listSnippets(list))
		if len(list.Snippets) > 10 {
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
		if v.RunCount > 0 {
			status.WriteString(executed)
			status.WriteString("\n")
		}
		status.WriteString(fmtRunCount(v.RunCount))

		//col3

		// //strings.Join(v.Tags, ", "),
		tbl.Append([]string{
			name.String(),
			status.String(),
			FmtSnippet(v, 60, lines),
			FmtOutPreview(v.Preview),
		})
	}
	tbl.Render()

	if len(list.Snippets) == 0 {
		fmt.Fprint(w, MARGIN)
		fmt.Fprint(w, style.Fmt(style.Subdued, "<empty pouch>\n"))
	}

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
	return style.Fmt(style.Subdued, fmt.Sprintf("↻ %2d", count))
}

func fmtDescription(w io.Writer, in string, width int) {
	t := strings.Split(style.WrapString(in, width), "\n")
	for i, v := range t {
		t[i] = style.Fmt(style.Subdued, v)
	}
	join := strings.Join(t, "\n")
	fmt.Fprint(w, join)
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

func FmtOutPreview(in string) string {
	in = strings.Replace(in, "\n\n", "\n", -1)
	in = strings.TrimSpace(in)
	return style.WrapString(pad(30, in).String(), 30)
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

func listHorizontal(l []interface{}) []byte {
	var buff bytes.Buffer
	w := tabwriter.NewWriter(&buff, 20, 3, 2, ' ', tabwriter.DiscardEmptyColumns)
	var item bytes.Buffer
	colWidths := map[int]int{}
	for i, v := range l {
		if i%5 == 0 {
			item.WriteString("  ")
		}
		if sn, ok := v.(*models.Snippet); ok {
			item.WriteString(printStatus(sn, false))
			item.WriteString("  ")
			item.WriteString(style.Fmt(style.Cyan, sn.SnipName.Name))
			item.WriteString(style.Fmt(style.Subdued, "."+sn.SnipName.Ext))
			item.WriteString(" ")
		}
		if pch, ok := v.(*models.Pouch); ok {
			if models.Prefs().ListAll || !pch.MakePrivate {
				if colWidths[i%5] < len(pch.Name) {
					colWidths[i%5] = len(pch.Name)
				}
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

		x := i + 1
		if x%20 == 0 {
			item.WriteString(MARGIN)
			item.WriteString("\n\t\t\t\t")
			item.WriteString("\n")
		} else if x%5 == 0 {
			item.WriteString("\n")
		} else {
			item.WriteString("\t")
		}

		fmt.Fprint(w, fmt.Sprintf("%s", item.String()))
		item.Reset()
	}
	w.Flush()
	return buff.Bytes()
}
