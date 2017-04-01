package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"fmt"
	"bytes"
	"text/tabwriter"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"math"
	"strings"
	"time"
)

var Pad_1_2 = strings.Repeat(style.Fmt256(100, " "), 2)
var Pad_2_2 = strings.Repeat(style.Fmt256(100, "  "), 2)
var Pad_1_1 = strings.Repeat(style.Fmt256(100, " "), 1)
var Pad_0_0 = strings.Repeat(style.Fmt256(30, ""), 1)

var Pad16_0_0 = strings.Repeat(style.Fmt(0, ""), 1)
var Pad16_1_1 = strings.Repeat(style.Fmt(0, " "), 1)
var Pad16_2_1 = strings.Repeat(style.Fmt(0, "  "), 1)
var Pad16_3_1 = strings.Repeat(style.Fmt(0, "   "), 1)
var Pad16_4_1 = strings.Repeat(style.Fmt(0, "    "), 1)
var Pad16_1_2 = strings.Repeat(style.Fmt(0, " "), 2)
var Pad16_2_2 = strings.Repeat(style.Fmt(0, "  "), 2)
var Pad16_3_2 = strings.Repeat(style.Fmt(0, "   "), 2)
var Pad16_3_4 = strings.Repeat(style.Fmt(0, " "), 4)

const oneMin = int64(60)
const oneHour = oneMin * 60
const oneDay = oneHour * 24

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
				isLast := stats.LastPouch == pch.PouchId
				item.WriteString(pouchIcon(pch, isLast))
				if isLast {
					item.WriteString(style.Fmt256(style.AnsiCode(254), "  ❯ "+pch.Name))
				} else {
					item.WriteString("  ")
					item.WriteString(style.Fmt256(decayColor(pch.LastUse, true), pch.Name))
				}
				if pch.Type == models.PouchType_Virtual {
					item.WriteString(Pad_1_1)
				} else {
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
	b.WriteString(fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\n", MARGIN, Pad_1_2, Pad_1_2, Pad_1_2, Pad_1_2, Pad_1_2))
}

func pouchIcon(pch *models.Pouch, isLast bool) string {
	if pch.Type == models.PouchType_Virtual {
		return style.Fmt256(242, style.Icon_Pouch)
		//} else if pch.Name == "inbox" {
		//	if pch.UnOpened > 0 {
		//		item.WriteString(fmt.Sprintf("📬%d ", pch.UnOpened))
		//	} else {
		//		item.WriteString(style.Fmt256(242, "▉ "))
		//	}
	} else if pch.MakePrivate {
		return colorPouch(isLast, pch.LastUse, pch.Red, "◤")
	} else {
		return colorPouch(isLast, pch.LastUse, pch.Red, style.Icon_Pouch)
	}
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

func colorPouch(lastPouch bool, lastUsed int64, reddy int64, icon string) string {
	var color style.AnsiCode
	if lastPouch && reddy > 0 {
		color = 196
	} else if lastPouch {
		color = 122
	} else if reddy > 0 {
		color = 124
	} else {
		color = decayColor(lastUsed, false)
	}
	return style.Fmt256(color, icon)
}

func newerThan(unix int64, seconds int64) bool {
	return time.Now().Unix()-unix < 5*seconds
}

func decayColor(unix int64, whiteToday bool) style.AnsiCode {
	local := time.Now()
	pouchT := time.Unix(unix, 0)
	today := local.YearDay() == pouchT.YearDay() && local.Year() == pouchT.Year()
	if today {
		if whiteToday {
			return style.AnsiCode(250)
		}
		return style.AnsiCode(122)
	}
	if newerThan(unix, 7*oneDay) {
		return style.AnsiCode(247)
	}
	if newerThan(unix, 28*oneDay) {
		return style.AnsiCode(245)
	}
	return style.AnsiCode(242)
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
	} else {
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
