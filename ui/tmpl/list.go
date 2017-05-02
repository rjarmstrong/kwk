package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"github.com/rjarmstrong/tablewriter"
	"strings"
	"bytes"
	"fmt"
	"io"
	"golang.org/x/text/unicode/norm"
	"sort"
	"bitbucket.com/sharingmachine/types"
)

var mainMarkers = map[string]string{
	"go": "func main() {",
}

type CodeLine struct {
	Margin string
	Body   string
}

func StatusText(s *types.Snippet) string {
	if s.Ext == "url" {
		return "bookmark"
	}
	if s.RunStatus == types.UseStatusSuccess {
		return "success"
	} else if s.RunStatus == types.UseStatusFail {
		return "error"
	}
	return "static"
}

func FStatus(s *types.Snippet, includeText bool) string {
	if s.RunStatus == types.UseStatusSuccess {
		if includeText {
			return style.Fmt256(style.Color_YesGreen, style.Icon_Tick) + "  Success"
		}
		return style.Fmt256(style.Color_YesGreen, style.Icon_Tick)
	} else if s.RunStatus == types.UseStatusFail {
		if includeText {
			return style.Fmt256(style.Color_BrightRed, style.Icon_Broke) +  "  Error"
		}
		return style.Fmt256(style.Color_BrightRed, style.Icon_Broke)
	}
	return style.Fmt16(style.Subdued, "? ")
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
	fmt.Fprint(w, strings.Repeat(" ", 50), style.Fmt16(style.Subdued, "◉  "+models.Principal.Username + "    TLS12"))
	fmt.Fprint(w, style.TWOLINES)
	w.Write(listHorizontal(all, &r.UserStats))

	if len(r.Snippets) > 0 {
		fmt.Fprintf(w, "\n%sLast:", style.MARGIN)
		fmt.Fprint(w, listSnippets(r, true))
	}

	if models.ClientIsNew(r.LastUpgrade) {
		w.WriteString(style.Fmt16(style.Subdued, fmt.Sprintf("\n\n%skwk auto-updated to %s %s", style.MARGIN, models.Client.Version, humanTime(r.LastUpgrade))))
	} else {
		w.WriteString("\n")
	}
	w.WriteString("\n\n")
	return w.String()
}

//func printCommunity(w *bytes.Buffer) {
//	fmt.Fprint(w, "\n", style.MARGIN, style.Fmt(style.Subdued, "Community"), "\n")
//	com := []interface{}{}
//	com = append(com, &models.Pouch{
//		Name:       style.Fmt(style.Cyan, "/kwk/") + "unicode",
//		Username:   "kwk",
//		PouchStats: models.PouchStats{Runs: 12},
//	}, &models.Pouch{
//		Name:       style.Fmt(style.Cyan, "/kwk/") + "news",
//		Username:   "kwk",
//		PouchStats: models.PouchStats{Runs: 12},
//	},
//		&models.Pouch{
//			Name:       style.Fmt(style.Cyan, "/kwk/") + "devops",
//			Username:   "kwk",
//			PouchStats: models.PouchStats{Runs: 12},
//		})
//	w.Write(listHorizontal(com, nil))
//	w.WriteString("\n")
//}

func printPouchHeadAndFoot(w *bytes.Buffer, list *models.ListView) {
	fmtHeader(w, list.Username, list.Pouch.Name, nil)
	fmt.Fprint(w, style.MARGIN, style.MARGIN, pouchIcon(list.Pouch, false))
	fmt.Fprint(w, "  ")
	fmt.Fprint(w, locked(list.Pouch.MakePrivate))
	fmt.Fprint(w, " Pouch")
	fmt.Fprint(w, style.MARGIN, style.MARGIN, len(list.Snippets), " snippets")
	fmt.Fprint(w, "\n")
}

func locked(locked bool) string {
	if locked {
		return "Locked (Private)"
	}
	return "Public"
}

func listPouch(list *models.ListView) string {
	w := &bytes.Buffer{}
	if models.Prefs().Naked {
		fmt.Fprint(w, listNaked(list))
	} else {
		//sort.Slice(list.Snippets, func(i, j int) bool {
		//	return list.Snippets[i].Created <= list.Snippets[j].Created
		//})
		if list.Pouch != nil {
			printPouchHeadAndFoot(w, list)
		}
		fmt.Fprint(w, listSnippets(list, false))
		if list.Pouch != nil && len(list.Snippets) > 10 && !models.Prefs().HorizontalLists {
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
	tbl.SetHeader([]string{"Name", "Username", "Pouch", "Ext", "Private", "Status", "Runs", "Views", "Deps", "LastActive", "Updated"})
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
		var private string
		if v.Private {
			private = "private"
		} else {
			private = "public"
		}
		tbl.Append([]string{
			v.Name,
			v.Username,
			v.Pouch,
			v.Ext,
			private,
			StatusText(v),
			fmt.Sprintf("%d", v.Runs),
			fmt.Sprintf("%d", v.Views),
			fmt.Sprintf("%d", len(v.Dependencies)),
			v.RunStatusTime.Format(timeLayout),
			v.Created.Format(timeLayout),
		})
	}
	tbl.Render()
	return w.String()
}

func listSnippets(list *models.ListView, fullName bool) string {
	if models.Prefs() != nil && models.Prefs().HorizontalLists {
		sort.Slice(list.Snippets, func(i,j int) bool {
			return list.Snippets[i].Name < list.Snippets[j].Name
		})
		l := []interface{}{}
		for _, v := range list.Snippets {
			l = append(l, v)
		}
		return "\n\n" + string(listHorizontal(l, nil)) + "\n\n"
	}

	w := &bytes.Buffer{}

	if len(list.Snippets) == 0 {
		fmt.Fprint(w, "\n")
		fmt.Fprint(w, style.MARGIN)
		fmt.Fprint(w, style.Fmt16(style.Subdued, "<empty pouch>"))
		fmt.Fprint(w, style.TWOLINES)
		fmt.Fprint(w, style.MARGIN)
		fmt.Fprint(w, style.Fmt16(style.Cyan, "Add new snippets to this pouch: "))
		if list.Pouch != nil {
			fmt.Fprintf(w, "`kwk new <snippet> %s/<name>.<ext>`", list.Pouch.Name)
		}
		fmt.Fprint(w, "\n")
		return w.String()
	}

	tbl := tablewriter.NewWriter(w)
	tbl.SetHeader([]string{"", "", ""})
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
	tbl.SetColWidth(1)

	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for i, v := range list.Snippets {
		// col1
		name := &bytes.Buffer{}
		name.WriteString(snippetIcon(v))
		name.WriteString("  ")
		sn := v.SnipName.String()
		if fullName {
			sn = v.String()
		}
		nt := style.Fmt256(style.Color_BrighterWhite, sn)
		name.WriteString(nt)
		if v.Description != "" {
			name.WriteString("\n\n")
			name.WriteString(style.Fmt256(style.Color_MonthGrey, style.FBox(v.Description, 25, 3)))
		}
		if v.Role == types.RoleEnvironment {
			name.WriteString("\n\n")
			name.WriteString(style.Fmt16(style.Subdued, "short-cut: kwk edit env"))
		}
		// col2
		var lines int
		if models.Prefs().AlwaysExpandRows {
			lines = models.Prefs().ExpandedRows
		} else {
			lines = models.Prefs().SlimRows
		}
		status := &bytes.Buffer{}
		runCount := fmtRunCount(v.Runs)
		status.WriteString(PadRight(runCount, " ", 21))
		status.WriteString(" ")
		if !v.RunStatusTime.IsZero() {
			h := PadLeft(style.Time(v.RunStatusTime), " ", 4)
			t := fmt.Sprintf("%s", style.Fmt256(239, h))
			status.WriteString(t)
		}
		//col3
		snip := FCodeview(v, 60, lines, (i+1)%2 == 0)
		if models.Prefs().RowSpaces {
			snip = snip + "\n"
		}
		if len(v.Preview) >= 10 {
			snip = snip + "\n\n" + style.MARGIN + style.Fmt256(style.Color_MonthGrey, style.FPreview(v.Preview, 120, 1))
		}

		tbl.Append([]string{
			name.String(),
			status.String(),
			snip,
		})
	}
	tbl.Render()

	//fmt.Fprint(w, style.Start)
	//fmt.Fprintf(w, "%dm", style.Subdued)
	//fmt.Fprint(w, style.MARGIN)
	//fmt.Fprintf(w,"Expand list: `kwk expand %s`", list.Pouch)
	//fmt.Fprint(w, style.MARGIN)
	//fmt.Fprint(w, style.End)
	////fmt.Fprint(w, style.Start)
	//fmt.Fprintf(w, "%dm", style.Subdued)
	//fmt.Fprint(w, style.MARGIN)
	//fmt.Fprintf(w, "%d of max 32 snippets in pouch", len(list.Snippets))
	//fmt.Fprint(w, style.End)

	return w.String()
}

func PadRight(str, pad string, length int) string {
	if len(str) < length {
		return str + strings.Repeat(pad, length-len(str))
	}
	return str
}

func PadLeft(str, pad string, length int) string {
	if len(str) < length {
		return strings.Repeat(pad, length-len(str)) + str
	}
	return str
}

func snippetIcon(v *types.Snippet) string {
	icon := style.Icon_Snippet
	if v.IsApp() {
		icon = style.Icon_App
	} else if v.Ext == "url" {
		icon = style.Icon_Bookmark
	}
	if v.RunStatus == types.UseStatusSuccess {
		return style.Fmt256(122, icon)
	} else if v.RunStatus == types.UseStatusFail {
		return style.Fmt256(196, icon)
	}
	return style.Fmt256(style.Color_MonthGrey, icon)
}

func fmtRunCount(count int64) string {
	return fmt.Sprintf(style.Fmt256(247, "↻ %0d"), count)
}

func fmtHeader(w io.Writer, username string, pouch string, s *types.SnipName) {
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, style.MARGIN)
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
		fmt.Fprint(w, pouch)
		fmt.Fprint(w, " ")
		fmt.Fprint(w, style.End)
		return
	}
	if pouch != "" {
		fmt.Fprint(w, pouch)
		fmt.Fprint(w, "/")
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
		for !ia.Done() && nc < width {
			nc += 1
			buff.Write(ia.Next())
		}
	}
	return buff
}
