package charting

import (
	tm "github.com/buger/goterm"
	"math"
	"fmt"
	"github.com/kwk-links/kwk-cli/libs/api"
	//"github.com/kwk-links/kwk-cli/system"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"sort"
)

func PrintSine() {
	chart := tm.NewLineChart(70, 10)
	d := &tm.DataTable{}
	d.AddColumn("")
	d.AddColumn("Sin(x)")
	d.AddColumn("Cos(x+1)")

	for i := 0.0; i < 20.0; i += 1.0 {
		d.AddRow(i, math.Sin(i), math.Cos(i + 1))
	}
	fmt.Println(chart.Draw(d))
}


type GraphItem struct {
	Key string
	Value int
}

type GraphItemSorter struct {
	Items []GraphItem
}

func (s GraphItemSorter) Len() int {
	return len(s.Items)
}
func (s GraphItemSorter) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}
func (s GraphItemSorter) Less(i, j int) bool {
	return s.Items[i].Value > s.Items[j].Value
}

func PrintTags(list *api.AliasList){
	tags := map[string]int{}
	for _, v := range list.Items {
		for _, t := range v.Tags {
			tags[t] += 1
		}
	}
	fmt.Print(gui.Colour(gui.LightBlue, "\n     kwklinks by tag\n\n"))
	sorter := GraphItemSorter{Items:[]GraphItem{}}
	for k, v := range tags {
		sorter.Items = append(sorter.Items, GraphItem{k, v})
	}
	sort.Sort(sorter)
	index := 0
	for index <= 10 {
		item := sorter.Items[index]
		fmt.Println(gui.Build(5, " ") + fmt.Sprintf("%-12s %2d   ", item.Key, item.Value) + gui.Build(item.Value, gui.UBlock))
		index += 1
	}
	fmt.Println()
}
