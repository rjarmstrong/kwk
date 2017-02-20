package fzf

import (
	"math"
	"sort"
	"unicode"

	"github.com/rjarmstrong/fzf/src/tui"
	"github.com/rjarmstrong/fzf/src/util"
)

// Offset holds two 32-bit integers denoting the offsets of a matched substring
type Offset [2]int32

type colorOffset struct {
	offset [2]int32
	color  tui.ColorPair
	attr   tui.Attr
	index  int32
}

type rank struct {
	points [4]uint16
	index  int32
}

type Result struct {
	item *Item
	rank rank
}

func buildResult(item *Item, offsets []Offset, score int, trimLen int) *Result {
	if len(offsets) > 1 {
		sort.Sort(ByOrder(offsets))
	}

	result := Result{item: item, rank: rank{index: item.index}}
	numChars := item.text.Length()
	minBegin := math.MaxUint16
	minEnd := math.MaxUint16
	maxEnd := 0
	validOffsetFound := false
	for _, offset := range offsets {
		b, e := int(offset[0]), int(offset[1])
		if b < e {
			minBegin = util.Min(b, minBegin)
			minEnd = util.Min(e, minEnd)
			maxEnd = util.Max(e, maxEnd)
			validOffsetFound = true
		}
	}

	for idx, criterion := range sortCriteria {
		val := uint16(math.MaxUint16)
		switch criterion {
		case byScore:
			// Higher is better
			val = math.MaxUint16 - util.AsUint16(score)
		case byLength:
			// If offsets is empty, trimLen will be 0, but we don't care
			val = util.AsUint16(trimLen)
		case byBegin, byEnd:
			if validOffsetFound {
				whitePrefixLen := 0
				for idx := 0; idx < numChars; idx++ {
					r := item.text.Get(idx)
					whitePrefixLen = idx
					if idx == minBegin || !unicode.IsSpace(r) {
						break
					}
				}
				if criterion == byBegin {
					val = util.AsUint16(minEnd - whitePrefixLen)
				} else {
					val = util.AsUint16(math.MaxUint16 - math.MaxUint16*(maxEnd-whitePrefixLen)/trimLen)
				}
			}
		}
		result.rank.points[idx] = val
	}

	return &result
}

// Sort criteria to use. Never changes once fzf is started.
var sortCriteria []criterion

// Index returns ordinal index of the Item
func (result *Result) Index() int32 {
	return result.item.index
}

func minRank() rank {
	return rank{index: 0, points: [4]uint16{math.MaxUint16, 0, 0, 0}}
}

func (result *Result) colorOffsets(matchOffsets []Offset, theme *tui.ColorTheme, color tui.ColorPair, attr tui.Attr, current bool) []colorOffset {
	itemColors := result.item.Colors()

	// No ANSI code, or --color=no
	if len(itemColors) == 0 {
		var offsets []colorOffset
		for _, off := range matchOffsets {
			offsets = append(offsets, colorOffset{offset: [2]int32{off[0], off[1]}, color: color, attr: attr})
		}
		return offsets
	}

	// Find max column
	var maxCol int32
	for _, off := range matchOffsets {
		if off[1] > maxCol {
			maxCol = off[1]
		}
	}
	for _, ansi := range itemColors {
		if ansi.offset[1] > maxCol {
			maxCol = ansi.offset[1]
		}
	}
	cols := make([]int, maxCol)

	for colorIndex, ansi := range itemColors {
		for i := ansi.offset[0]; i < ansi.offset[1]; i++ {
			cols[i] = colorIndex + 1 // XXX
		}
	}

	for _, off := range matchOffsets {
		for i := off[0]; i < off[1]; i++ {
			cols[i] = -1
		}
	}

	// sort.Sort(ByOrder(offsets))

	// Merge offsets
	// ------------  ----  --  ----
	//   ++++++++      ++++++++++
	// --++++++++--  --++++++++++---
	curr := 0
	start := 0
	var colors []colorOffset
	add := func(idx int) {
		if curr != 0 && idx > start {
			if curr == -1 {
				colors = append(colors, colorOffset{
					offset: [2]int32{int32(start), int32(idx)}, color: color, attr: attr})
			} else {
				ansi := itemColors[curr-1]
				fg := ansi.color.fg
				bg := ansi.color.bg
				if theme != nil {
					if fg == -1 {
						if current {
							fg = theme.Current
						} else {
							fg = theme.Fg
						}
					}
					if bg == -1 {
						if current {
							bg = theme.DarkBg
						} else {
							bg = theme.Bg
						}
					}
				}
				colors = append(colors, colorOffset{
					offset: [2]int32{int32(start), int32(idx)},
					color:  tui.NewColorPair(fg, bg),
					attr:   ansi.color.attr.Merge(attr)})
			}
		}
	}
	for idx, col := range cols {
		if col != curr {
			add(idx)
			start = idx
			curr = col
		}
	}
	add(int(maxCol))
	return colors
}

// ByOrder is for sorting substring offsets
type ByOrder []Offset

func (a ByOrder) Len() int {
	return len(a)
}

func (a ByOrder) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByOrder) Less(i, j int) bool {
	ioff := a[i]
	joff := a[j]
	return (ioff[0] < joff[0]) || (ioff[0] == joff[0]) && (ioff[1] <= joff[1])
}

// ByRelevance is for sorting Items
type ByRelevance []*Result

func (a ByRelevance) Len() int {
	return len(a)
}

func (a ByRelevance) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByRelevance) Less(i, j int) bool {
	return compareRanks((*a[i]).rank, (*a[j]).rank, false)
}

// ByRelevanceTac is for sorting Items
type ByRelevanceTac []*Result

func (a ByRelevanceTac) Len() int {
	return len(a)
}

func (a ByRelevanceTac) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByRelevanceTac) Less(i, j int) bool {
	return compareRanks((*a[i]).rank, (*a[j]).rank, true)
}

func compareRanks(irank rank, jrank rank, tac bool) bool {
	for idx := 0; idx < 4; idx++ {
		left := irank.points[idx]
		right := jrank.points[idx]
		if left < right {
			return true
		} else if left > right {
			return false
		}
	}
	return (irank.index <= jrank.index) != tac
}
