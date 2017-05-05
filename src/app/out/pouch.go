package out

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/kwkcli/src/style"
	"bitbucket.com/sharingmachine/types"
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"
)

const oneMin = int64(60)
const oneHour = oneMin * 60
const oneDay = oneHour * 24

func listHorizontal(l []interface{}, stats *types.UserStats) []byte {
	var buff bytes.Buffer
	w := tabwriter.NewWriter(&buff, 5, 1, 3, ' ', tabwriter.TabIndent)
	var item bytes.Buffer
	colWidths := map[int]int{}
	for i, v := range l {
		if i%5 == 0 {
			item.WriteString("  ")
		}
		if sn, ok := v.(*types.Snippet); ok {
			item.WriteString(snippetIcon(sn))
			item.WriteString("  ")
			item.WriteString(style.Fmt256(style.ColorBrighterWhite, sn.SnipName.String()))
			//item.WriteString(FStatus(sn, false))
		}
		if pch, ok := v.(*types.Pouch); ok {
			if models.Prefs().ListAll || !pch.MakePrivate {
				if colWidths[i%5] < len(pch.Name) {
					colWidths[i%5] = len(pch.Name)
				}
				isLast := stats.LastPouch == pch.Id
				item.WriteString(pouchIcon(pch, isLast))
				if isLast {
					item.WriteString(style.Fmt256(style.AnsiCode(style.ColorBrightestWhite), "  ❯ "+pch.Name))
				} else {
					item.WriteString("  ")
					item.WriteString(style.Fmt256(decayColor(pch.LastUse.Unix(), true), pch.Name))
				}
				if pch.Type == types.PouchTypeVirtual {
					item.WriteString(style.Pad11)
				} else {
					item.WriteString(style.Fmt256(style.ColorDimStat, fmt.Sprintf(" %d", pch.Snips)))
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
	b.WriteString(fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\n", style.Margin, style.Pad12, style.Pad12, style.Pad12, style.Pad12, style.Pad12))
}

func pouchIcon(pch *types.Pouch, isLast bool) string {
	if pch.MakePrivate {
		return colorPouch(isLast, pch.LastUse.Unix(), pch.Red, style.IconPrivatePouch)
	} else {
		return colorPouch(isLast, pch.LastUse.Unix(), pch.Red, style.IconPouch)
	}
}

func colorPouch(lastPouch bool, lastUsed int64, reddy int64, icon string) string {
	var color style.AnsiCode
	if lastPouch && reddy > 0 {
		color = style.ColorBrightRed
	} else if lastPouch {
		color = style.ColorPouchCyan
	} else if reddy > 0 {
		color = style.ColorDimRed
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
			return style.AnsiCode(style.ColorBrightWhite)
		}
		return style.AnsiCode(style.ColorPouchCyan)
	}
	if newerThan(unix, 7*oneDay) {
		return style.AnsiCode(style.ColorWeekGrey)
	}
	if newerThan(unix, 28*oneDay) {
		return style.AnsiCode(style.ColorMonthGrey)
	}
	return style.AnsiCode(style.ColorOldGrey)
}
