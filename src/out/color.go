package out

import (
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/rjarmstrong/kwk-types"
)

func SetColors(c Colors) {
	colors = c
}

var colors Colors

type Colors struct {
	Subdued     types.AnsiCode
	RecentPouch types.AnsiCode
}

func ColorsDefault() Colors {
	return Colors{
		Subdued:     2,
		RecentPouch: style.ColorPouchCyan,
	}
}

func ColorsAngry() Colors {
	return Colors{
		Subdued:     style.ColorBrighterWhite,
		RecentPouch: style.ColorYesGreen,
	}
}
