package tmpl

import (
	"bitbucket.com/sharingmachine/kwkcli/ui/style"
	"fmt"
	"bytes"
	"text/tabwriter"
	"bitbucket.com/sharingmachine/kwkcli/models"
)

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
					item.WriteString(style.Fmt(style.DarkGrey, "ⓟ")) //"🔒")
				} else {
					if pch.SnipCount == 0 {
						item.WriteString(style.Fmt(style.DarkGrey, "▆") )
					}
					if pch.SnipCount > 0 && pch.SnipCount < 20 {
						item.WriteString(style.Fmt(style.White, "▆") )
					}
					if pch.SnipCount > 20 {
						item.WriteString(style.Fmt(style.LightRed, "▆") )
					}
					//item.WriteString(style.Fmt(style.LightRed, "▆") ) //▇") //👝 ▇")
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