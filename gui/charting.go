package gui

import (
	tm "github.com/buger/goterm"
	"math"
	"fmt"
)

func PrintSine(){
	chart := tm.NewLineChart(70, 10)
	d := &tm.DataTable{}
	d.AddColumn("")
	d.AddColumn("Sin(x)")
	d.AddColumn("Cos(x+1)")

	for i := 0.0; i < 20.0; i += 1.0 {
		d.AddRow(i, math.Sin(i), math.Cos(i+1))
	}
	fmt.Println(chart.Draw(d))
}
