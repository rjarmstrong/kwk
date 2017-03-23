package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"fmt"
	"bytes"
	"text/tabwriter"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"math"
	"strings"
)

var gs2 = strings.Repeat(style.Fmt256(100," "), 2)
var gs = strings.Repeat(style.Fmt256(100," "), 1)
var gs16 = strings.Repeat(style.Fmt(0,"  "), 1)


func listHorizontal(l []interface{}, stats *models.UserStats) []byte {
	//fmt.Printf("%+v\n\n", stats)
	var buff bytes.Buffer
	w := tabwriter.NewWriter(&buff, 5, 1, 3, ' ', tabwriter.TabIndent)
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
				if pch.Type == models.PouchType_Virtual {
					item.WriteString(style.Fmt256(242, "▆ "))
				//} else if pch.Name == "inbox" {
				//	if pch.UnOpened > 0 {
				//		item.WriteString(fmt.Sprintf("📬%d ", pch.UnOpened))
				//	} else {
				//		item.WriteString(style.Fmt256(242, "▆ "))
				//	}
				} else if pch.Name == "settings" {
					item.WriteString("⚙ ")
				} else if stats == nil {
					item.WriteString(colorPouch8(false,
						0, 0, 0, 0, "▆ "))
				} else if pch.MakePrivate {
					item.WriteString(colorPouch8(
						contains(stats.RecentPouches, pch.PouchId),
						pch.Use, stats.MaxUsePerPouch, pch.Green, pch.Red, "Ⓟ ")) //"PV"))
				} else {
					//item.WriteString(fmt.Sprintf("[%d]", pch.Use))
					item.WriteString(colorPouch8(
						contains(stats.RecentPouches, pch.PouchId),
						pch.Use, stats.MaxUsePerPouch, pch.Green, pch.Red, "▆ "))
				}

				item.WriteString(" ")
				if pch.Type == models.PouchType_Virtual {
					item.WriteString(pch.Name)
					item.WriteString(gs)
				} else {
					item.WriteString(pch.Name)
					item.WriteString(style.Fmt256(238, fmt.Sprintf(" %d", pch.PouchStats.Snips)))
				}
			}
		}
		x := i + 1
		if x%20 == 0 {
			insertGridLine(&item)
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

func insertGridLine(b *bytes.Buffer) {
	b.WriteString(fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\n", MARGIN,  gs2, gs2, gs2, gs2, gs2))
}

//var matrix = [][]int{
//	{23, 29, 35},
//	{52, 124, 196},
//}


var matrix = [][]int{
		{15, 15, 15},
		{15, 15, 15},
}

func contains(in []string, val string) bool {
	for _, x := range in {
		if x == val {
			return true
		}
	}
	return false
}

func colorPouch8(recent bool, use int64, maxUse int64, greeny int64, reddy int64, icon string) string {
	var color style.AnsiCode
	if reddy > 0 && recent {
		color = 196
	} else if reddy > 0  {
		color =  160
	} else if recent {
		color = 122
	} else {
		color = usageColor(maxUse, use)
	}
	return style.Fmt256(color, icon)
}

func usageColor(maxUse int64, use int64) style.AnsiCode {
	if maxUse == 0 {
		maxUse = 1
	}
	usage := float64(use) / float64(maxUse)
	if usage > 0.5 {
		return 253
	} else {
		return 244
	}
}

func colorPouch256(recent bool, use int64, maxUse int64, greeny int64, reddy int64, icon string) string {
	if recent {
		return style.Fmt256(style.AnsiCode(15), icon)
	}
	if maxUse == 0 {
		maxUse = 1
	}
	usage := float64(use) / float64(maxUse)
	//239-255 = 16
	brightness := int(math.Ceil(usage * 16))
	if greeny == 0 && reddy == 0 {
		return style.Fmt256(style.AnsiCode(239+brightness), icon)
	}

	var y int
	if reddy > 0 {
		y = 1
	} else  {
		y = 0
	}
	x := round(usage * 2)
	color := matrix[y][x]
	return style.Fmt256(style.AnsiCode(color), icon)
	//return style.FmtFgBg(fmt.Sprintf("g:%d r:%d x:%d y:%d %d/%d %d %s", greeny, reddy, x, y, use, maxUse, color, icon), style.AnsiCode(color), style.Black0)
}

func round(f float64) int {
	return int(f + math.Copysign(0.5, f))
}
