package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"fmt"
	"bytes"
	"text/tabwriter"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"math"
)

func listHorizontal(l []interface{}, stats models.UserStats) []byte {
	//fmt.Printf("%+v\n\n", stats)
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
					item.WriteString(colorPouch(pch.Use, stats.MaxUsePerPouch, pch.Green, pch.Red,"ⓟ")) //"🔒")
				} else {
					//item.WriteString(fmt.Sprintf("[%d]", pch.Use))
					item.WriteString(colorPouch(pch.Use, stats.MaxUsePerPouch, pch.Green, pch.Red, "▆"))
					//if pch.PouchStats.Snips == 0 {
					//	item.WriteString(style.Fmt(style.DarkGrey, "▆") )
					//}
					//if pch.PouchStats.Snips > 0 && pch.PouchStats.Snips < 20 {
					//	item.WriteString(style.Fmt(style.White, "▆") )
					//}
					//if pch.PouchStats.Snips > 20 {
					//	item.WriteString(style.Fmt(style.LightRed, "▆") )
					//}
					//item.WriteString(style.Fmt(style.LightRed, "▆") ) //▇") //👝 ▇")
				}

				item.WriteString("  ")
				item.WriteString(pch.Name)
				item.WriteString(style.Fmt(style.Subdued, fmt.Sprintf(" (%d)", pch.PouchStats.Snips)))
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

var matrix = [][]int{
	{23, 29, 35 },
	{147, 186, 226},
	{52, 124, 196},
}

func colorPouch(use int64, maxUse int64, greeny int64, reddy int64, icon string) string {
	//237-255 = 18
	usage := float64(use) / float64(maxUse)
	brightness := int(math.Ceil(usage * 16))
	if greeny == 0 && reddy == 0 {
		return style.FmtFgBg(icon, style.AnsiCode(239+brightness), style.Black0)
	}
	var y int
	x := round(usage * 2)
	if reddy == 0 && greeny > 0 {
		y = 0
	} else if greeny == 0 && reddy > 0 {
		y = 2
	} else {
		redness := float64(reddy) / float64(greeny)
		y = round(redness)
	}
	color := matrix[y][x]
	return style.FmtFgBg(icon, style.AnsiCode(color), style.Black0)
	//return style.FmtFgBg(fmt.Sprintf("g:%d r:%d %d-%d %d %s", greeny, reddy, x, y, color, icon), style.AnsiCode(color), style.Black0)
}

func round(f float64) int {
	return int(f + math.Copysign(0.5, f))
}
