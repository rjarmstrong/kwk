package out

import (
	"bytes"
	"fmt"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/kwk-super-snippets/types"
	"strings"
	"text/tabwriter"
	"time"
)

const oneMin = int64(60)
const oneHour = oneMin * 60
const oneDay = oneHour * 24

func listHorizontal(l []interface{}, listAll bool, stats *types.UserStats) []byte {
	var buff bytes.Buffer
	w := tabwriter.NewWriter(&buff, 5, 1, 3, ' ', tabwriter.TabIndent)
	var item bytes.Buffer
	colWidths := map[int]int{}
	pad11 := strings.Repeat(style.Fmt256(100, " "), 1)
	for i, v := range l {
		if i%5 == 0 {
			item.WriteString("  ")
		}
		if sn, ok := v.(*types.Snippet); ok {
			item.WriteString(snippetIcon(sn))
			item.WriteString("  ")
			item.WriteString(style.Fmt256(style.ColorBrighterWhite, sn.Alias.FileName()))
			//item.WriteString(FStatus(sn, false))
		}
		if pch, ok := v.(*types.Pouch); ok {
			if listAll || !pch.MakePrivate {
				if colWidths[i%5] < len(pch.Name) {
					colWidths[i%5] = len(pch.Name)
				}
				isLast := stats.LastPouch == pch.Id
				item.WriteString(pouchIcon(pch, isLast))
				if isLast {
					item.WriteString(style.Fmt256(style.ColorBrightestWhite, "  ❯ "+pch.Name))
				} else {
					item.WriteString("  ")
					item.WriteString(style.Fmt256(decayColor(pch.LastUse, true), pch.Name))
				}
				if pch.Type == types.PouchType_Virtual {
					item.WriteString(pad11)
				} else {
					item.WriteString(style.Fmt256(style.ColorDimStat, fmt.Sprintf(" %d", pch.Stats.Snips)))
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
	pad := strings.Repeat(style.Fmt256(100, " "), 2)
	b.WriteString(fmt.Sprintf("\n%s\t%s\t%s\t%s\t%s\t%s\n", style.Margin, pad, pad, pad, pad, pad))
}

func pouchIcon(pch *types.Pouch, isLast bool) string {
	if pch.MakePrivate {
		return colorPouch(isLast, pch.LastUse, pch.Stats.Red, style.IconPrivatePouch)
	} else {
		return colorPouch(isLast, pch.LastUse, pch.Stats.Red, style.IconPouch)
	}
}

func colorPouch(lastPouch bool, lastUsed int64, reddy int64, icon string) string {
	var col types.AnsiCode
	if lastPouch && reddy > 0 {
		col = style.ColorBrightRed
	} else if lastPouch {
		col = colors.RecentPouch
	} else if reddy > 0 {
		col = style.ColorDimRed
	} else {
		col = decayColor(lastUsed, false)
	}
	return style.Fmt256(col, icon)
}

func newerThan(unix int64, seconds int64) bool {
	return time.Now().Unix()-unix < 5*seconds
}

func decayColor(unix int64, whiteToday bool) types.AnsiCode {
	local := time.Now()
	pouchT := time.Unix(unix, 0)
	today := local.YearDay() == pouchT.YearDay() && local.Year() == pouchT.Year()
	if today {
		if whiteToday {
			return types.AnsiCode(style.ColorBrightWhite)
		}
		return colors.RecentPouch
	}
	if newerThan(unix, 7*oneDay) {
		return types.AnsiCode(style.ColorWeekGrey)
	}
	if newerThan(unix, 28*oneDay) {
		return types.AnsiCode(style.ColorMonthGrey)
	}
	return types.AnsiCode(style.ColorOldGrey)
}
