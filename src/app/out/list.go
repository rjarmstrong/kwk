package out

import (
	"bytes"
	"fmt"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/kwk-super-snippets/types"
	"github.com/rjarmstrong/tablewriter"
	"golang.org/x/text/unicode/norm"
	"io"
	"sort"
	"strings"
	"time"
)

var mainMarkers = map[string]string{
	"go": "func main() {",
}

type CodeLine struct {
	Margin string
	Body   string
}

func StatusText(s *types.Snippet) string {
	if s.Ext() == "url" {
		return "bookmark"
	}
	if s.RunStatus == types.UseStatus_Success {
		return "success"
	} else if s.RunStatus == types.UseStatus_Fail {
		return "error"
	}
	return "static"
}

func FStatus(s *types.Snippet, includeText bool) string {
	if s.RunStatus == types.UseStatus_Success {
		if includeText {
			return style.Fmt256(style.ColorYesGreen, style.IconTick) + "  Success"
		}
		return style.Fmt256(style.ColorYesGreen, style.IconTick)
	} else if s.RunStatus == types.UseStatus_Fail {
		if includeText {
			return style.Fmt256(style.ColorBrightRed, style.IconBroke) + "  Error"
		}
		return style.Fmt256(style.ColorBrightRed, style.IconBroke)
	}
	return style.Fmt16(style.Subdued, "? ")
}

func printRoot(w io.Writer, prefs *Prefs, cli *types.AppInfo, r *types.RootResponse, p *types.User) {
	var all []*types.Pouch
	for _, v := range r.Pouches {
		if v.Name != "" {
			all = append(all, v)
		}
	}
	for _, v := range r.Personal {
		all = append(all, v)
	}

	fmtHeader(w, r.Username, "", "")
	status := fmt.Sprintf("%s  %s    TLS12", style.IconAccount, p.Username)
	fmt.Fprint(w, strings.Repeat(" ", 45), style.Fmt16(style.Subdued, status))
	if prefs.PrivateView {
		fmt.Fprint(w, style.Fmt256(style.Grey243," Pvt"))
	} else {
		fmt.Fprint(w, style.Fmt256(style.Grey243," Pub"))
	}
	fmt.Fprint(w, style.TwoLines)
	w.Write(horizontalPouches(prefs, all, r.Stats))

	if len(r.Snippets) > 0 {
		fmt.Fprintf(w, "\n%sLast:", style.Margin)
		printSnippets(w, prefs, "", r.Snippets, true)
	}

	if clientIsNew(cli.Time) {
		fmt.Fprint(w, style.Fmt16(style.Subdued, fmt.Sprintf("\n\n%skwk auto-updated to %s %s", style.Margin, cli.Version, style.Time(time.Unix(cli.Time, 0)))))
	} else {
		fmt.Fprintln(w, "")
	}
	fmt.Fprint(w, "\n\n")
}

func clientIsNew(t int64) bool {
	if t == 0 {
		return false
	}
	return t > (time.Now().Unix() - 60)
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

func printPouchHeadAndFoot(w io.Writer, pouch *types.Pouch, list []*types.Snippet) {
	fmtHeader(w, pouch.Username, pouch.Name, "")
	fmt.Fprint(w, style.Margin, style.Margin, pouchIcon(pouch, false))
	fmt.Fprint(w, "  ")
	fmt.Fprint(w, locked(pouch.MakePrivate))
	fmt.Fprint(w, " Pouch")
	fmt.Fprint(w, style.Margin, style.Margin, len(list), " snippets")
	fmt.Fprint(w, "\n")
}

func locked(locked bool) string {
	if locked {
		return "Locked (Private)"
	}
	return "Public"
}

func printPouchSnippets(w io.Writer, prefs *Prefs, list *types.ListResponse) {
	if prefs.Quiet {
		for _, v := range list.Items {
			fmt.Fprint(w, v.Alias.URI())
			fmt.Fprint(w, "\n")
		}
		return
	}
	if prefs.Naked {
		fmt.Fprint(w, listNaked(list))
		return
	}
	//sort.Slice(list.Snippets, func(i, j int) bool {
	//	return list.Snippets[i].Created <= list.Snippets[j].Created
	//})
	if list.Pouch != nil {
		printPouchHeadAndFoot(w, list.Pouch, list.Items)
	}
	printSnippets(w, prefs, list.Pouch.Name, list.Items, false)
	if list.Pouch != nil && len(list.Items) > 10 && !prefs.ListHorizontal {
		printPouchHeadAndFoot(w, list.Pouch, list.Items)
	}
	fmt.Fprint(w, "\n")
}

const timeLayout = "2 Jan 15:04 06"

func listNaked(list *types.ListResponse) interface{} {
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
	for _, v := range list.Items {
		var private string
		if !v.Public {
			private = "private"
		} else {
			private = "public"
		}
		tbl.Append([]string{
			v.Name(),
			v.Username(),
			v.Pouch(),
			v.Ext(),
			private,
			StatusText(v),
			fmt.Sprintf("%d", v.Stats.Runs),
			fmt.Sprintf("%d", v.Stats.Views),
			fmt.Sprintf("%d", len(v.Dependencies.Aliases)),
			time.Unix(v.RunStatusTime, 0).Format(timeLayout),
			time.Unix(v.Created, 0).Format(timeLayout),
		})
	}
	tbl.Render()
	return w.String()
}

// printSnippets is the standard way snippets are viewed in as a list.
func printSnippets(w io.Writer, prefs *Prefs, pouchName string, list []*types.Snippet, showFullName bool) {
	if prefs.ListHorizontal {
		sort.Slice(list, func(i, j int) bool {
			return list[i].Name() < list[j].Name()
		})
		fmt.Fprint(w, "\n\n"+string(horizontalSnippets(list))+"\n\n")
		return
	}
	if len(list) == 0 {
		fmt.Fprint(w, "\n")
		fmt.Fprint(w, style.Margin)
		fmt.Fprint(w, style.Fmt16(style.Subdued, "<empty pouch>"))
		fmt.Fprint(w, style.TwoLines)
		fmt.Fprint(w, style.Margin)
		fmt.Fprint(w, style.Fmt16(style.Cyan, "Add new snippets to this pouch: "))
		if pouchName != "" {
			fmt.Fprintf(w, "`kwk new <snippet> %s/<name>.<ext>`", pouchName)
		}
		fmt.Fprint(w, "\n")
		return
	}

	tbl := tablewriter.NewWriter(w)
	tbl.SetHeader([]string{"", "", ""})
	tbl.SetAutoWrapText(false)
	tbl.SetBorder(false)
	tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tbl.SetCenterSeparator("")
	tbl.SetColumnSeparator(" ")
	if prefs.RowLines {
		tbl.SetRowLine(true)
	}

	tbl.SetAutoFormatHeaders(false)
	tbl.SetHeaderLine(true)
	tbl.SetColWidth(1)

	tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for i, v := range list {
		// col1
		name := &bytes.Buffer{}
		name.WriteString(snippetIcon(v))
		name.WriteString("  ")
		sn := v.Alias.FileName()
		if showFullName {
			sn = v.Alias.URI()
		}
		nt := style.Fmt256(style.ColorBrighterWhite, sn)
		name.WriteString(nt)
		if v.Description != "" {
			name.WriteString("\n\n")
			name.WriteString(style.Fmt256(style.ColorMonthGrey, style.FBox(v.Description, 25, 3)))
		}
		if v.Role == types.Role_Environment {
			name.WriteString("\n\n")
			name.WriteString(style.Fmt16(style.Subdued, "short-cut: kwk edit env"))
		}
		// col2
		var lines int
		if prefs.ExpandedRows {
			lines = prefs.ExpandedThumbRows
		} else {
			lines = prefs.SnippetThumbRows
		}
		status := &bytes.Buffer{}
		runCount := fmtRunCount(v.Stats.Runs)
		status.WriteString(PadRight(runCount, " ", 21))
		status.WriteString(" ")
		if v.RunStatusTime > 0 {
			h := PadLeft(style.Time(time.Unix(v.RunStatusTime, 0)), " ", 4)
			t := fmt.Sprintf("%s", style.Fmt256(239, h))
			status.WriteString(t)
		}
		//col3
		snip := FCodeview(v, 60, lines, (i+1)%2 == 0, prefs.ExpandedRows)
		if prefs.RowSpaces {
			snip = snip + "\n"
		}
		if len(v.Preview) >= 10 {
			snip = snip + "\n\n" + style.Margin + style.Fmt256(style.ColorMonthGrey, Fpreview(v.Preview, prefs, 120, 1))
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

func snippetIcon(s *types.Snippet) string {
	icon := style.IconSnippet
	if s.IsApp() {
		icon = style.IconApp
	} else if s.Ext() == "url" {
		icon = style.IconBookmark
	}
	if s.RunStatus == types.UseStatus_Success {
		return style.Fmt256(122, icon)
	} else if s.RunStatus == types.UseStatus_Fail {
		return style.Fmt256(196, icon)
	}
	return style.Fmt256(style.ColorMonthGrey, icon)
}

func fmtRunCount(count int64) string {
	return fmt.Sprintf(style.Fmt256(247, "↻ %0d"), count)
}

func fmtHeader(w io.Writer, username string, pouch string, fileName string) {
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, style.Margin)
	fmt.Fprint(w, style.Esc)
	fmt.Fprint(w, "7m")
	fmt.Fprint(w, " ❯ ")
	fmt.Fprint(w, types.KwkHost)
	fmt.Fprint(w, "/")
	if pouch == "" && fileName == "" {
		fmt.Fprint(w, style.Esc)
		fmt.Fprint(w, "1m")
		fmt.Fprint(w, username)
		fmt.Fprint(w, " ")
		fmt.Fprint(w, style.End)
		return
	}
	fmt.Fprint(w, username)
	fmt.Fprint(w, "/")
	if fileName == "" {
		fmt.Fprint(w, style.Esc)
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
	fmt.Fprint(w, style.Esc)
	fmt.Fprint(w, "1m")
	fmt.Fprint(w, fileName)
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
