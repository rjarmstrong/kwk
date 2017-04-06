package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"fmt"
	"bytes"
	"text/tabwriter"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"time"
)

const oneMin = int64(60)
const oneHour = oneMin * 60
const oneDay = oneHour * 24

func listHorizontal(l []interface{}, stats *models.UserStats) []byte {
	var buff bytes.Buffer
	w := tabwriter.NewWriter(&buff, 5, 1, 3, ' ', tabwriter.TabIndent)
	var item bytes.Buffer
	colWidths := map[int]int{}
	for i, v := range l {
		if i%5 == 0 {
			item.WriteString("  ")
		}
		if sn, ok := v.(*models.Snippet); ok {
			item.WriteString(snippetIcon(sn))
			item.WriteString("  ")
			item.WriteString(style.Fmt256(style.Color_BrighterWhite, sn.SnipName.String()))
			//item.WriteString(FStatus(sn, false))
		}
		if pch, ok := v.(*models.Pouch); ok {
			if models.Prefs().ListAll || !pch.MakePrivate {
				if colWidths[i%5] < len(pch.Name) {
					colWidths[i%5] = len(pch.Name)
				}
				isLast := stats.LastPouch == pch.PouchId
				item.WriteString(pouchIcon(pch, isLast))
				if isLast {
					item.WriteString(style.Fmt256(style.AnsiCode(style.Color_BrightestWhite), "  ❯ "+pch.Name))
				} else {
					item.WriteString("  ")
					item.WriteString(style.Fmt256(decayColor(pch.LastUse, true), pch.Name))
				}
				if pch.Type == models.PouchType_Virtual {
					item.WriteString(style.Pad_1_1)
				} else {
					item.WriteString(style.Fmt256(style.Color_DimStat, fmt.Sprintf(" %d", pch.PouchStats.Snips)))
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
	b.WriteString(fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\n", MARGIN, style.Pad_1_2, style.Pad_1_2, style.Pad_1_2, style.Pad_1_2, style.Pad_1_2))
}

func pouchIcon(pch *models.Pouch, isLast bool) string {
	if pch.MakePrivate {
		return colorPouch(isLast, pch.LastUse, pch.Red, style.Icon_PrivatePouch)
	} else {
		return colorPouch(isLast, pch.LastUse, pch.Red, style.Icon_Pouch)
	}
}

func colorPouch(lastPouch bool, lastUsed int64, reddy int64, icon string) string {
	var color style.AnsiCode
	if lastPouch && reddy > 0 {
		color = style.Color_BrightRed
	} else if lastPouch {
		color = style.Color_PouchCyan
	} else if reddy > 0 {
		color = style.Color_DimRed
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
			return style.AnsiCode(style.Color_BrightWhite)
		}
		return style.AnsiCode(style.Color_PouchCyan)
	}
	if newerThan(unix, 7*oneDay) {
		return style.AnsiCode(style.Color_WeekGrey)
	}
	if newerThan(unix, 28*oneDay) {
		return style.AnsiCode(style.Color_MonthGrey)
	}
	return style.AnsiCode(style.Color_OldGrey)
}