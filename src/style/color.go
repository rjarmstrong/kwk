package style

import (
	"github.com/kwk-super-snippets/types"
)

const (
	BrightWhite types.AnsiCode = 1
	Subdued     types.AnsiCode = 2

	// 16 COLORS
	Black     types.AnsiCode = 30
	Red       types.AnsiCode = 31
	Green     types.AnsiCode = 32
	Yellow    types.AnsiCode = 33
	Blue      types.AnsiCode = 34
	Magenta   types.AnsiCode = 35
	Cyan      types.AnsiCode = 36
	LightGrey types.AnsiCode = 37
	White     types.AnsiCode = 38

	CyanBg       types.AnsiCode = 46
	DarkGrey     types.AnsiCode = 90
	LightRed     types.AnsiCode = 91
	LightGreen   types.AnsiCode = 92
	LightYellow  types.AnsiCode = 93
	LightBlue    types.AnsiCode = 94
	LightMagenta types.AnsiCode = 95
	LightCyan    types.AnsiCode = 96
	White97      types.AnsiCode = 97

	// 256 COLORS
	LightBlue104 types.AnsiCode = 104
	Black0       types.AnsiCode = 0
	Black231     types.AnsiCode = 231
	Black232     types.AnsiCode = 232
	Black233     types.AnsiCode = 233
	Black234     types.AnsiCode = 234
	Grey234      types.AnsiCode = 234
	Grey236      types.AnsiCode = 236
	Grey238      types.AnsiCode = 238
	Grey240      types.AnsiCode = 240
	Grey241      types.AnsiCode = 241
	Grey243      types.AnsiCode = 243
	White15      types.AnsiCode = 15
	OffWhite248  types.AnsiCode = 248
	OffWhite249  types.AnsiCode = 249
	OffWhite250  types.AnsiCode = 250
	OffWhite253  types.AnsiCode = 253
	OffWhite254  types.AnsiCode = 254
	OffWhite255  types.AnsiCode = 255

	Bold      types.AnsiCode = 1
	Dim       types.AnsiCode = 22
	Regular   types.AnsiCode = 5
	Underline types.AnsiCode = 4

	ClearLine = "\033[1K"
	MoveBack  = "\033[9D"
	Block     = "2588"

	Start      = "\033["
	End        = "\033[0m"
	Start255   = "\033[48;5;"
	End255     = "\033[0;00m"
	HideCursor = "\033[?25l"
	ShowCursor = "\033[?25h"
	Margin     = "  "
	TwoLines   = "\n\n"

	Warning          = "\xE2\x9A\xA0"
	Fire             = "\xF0\x9F\x94\xA5"
	IconApp          = "✿" //✱ ▚ ❖  ꌳ ⧓ ⧗ 〓 ⁘ ꌳ ⁑⁘ ⁙ ѧꊞ ▚ 囙"
	IconSnippet      = "✦" //◆"
	IconView         = "❍" // 274d
	IconTick         = "✓" // 2713
	IconCross        = "✘" // 2718
	IconPrivatePouch = "◤"
	IconBroke        = "▦"
	InfoDeskPerson   = "\xF0\x9F\x92\x81"

	Ambulance = "\xF0\x9F\x9A\x91"

	ColorBrightRed      types.AnsiCode = 196
	ColorPouchCyan      types.AnsiCode = 122
	ColorDimRed         types.AnsiCode = 124
	ColorBrightWhite    types.AnsiCode = 250
	ColorBrighterWhite  types.AnsiCode = 252
	ColorBrightestWhite types.AnsiCode = 254
	ColorWeekGrey       types.AnsiCode = 247
	ColorMonthGrey      types.AnsiCode = 245
	ColorOldGrey        types.AnsiCode = 242
	ColorDimStat        types.AnsiCode = 238
	ColorYesGreen       types.AnsiCode = 119

	// 🔰 👝 🔒 🔸 ⚡ ✓ ⇨ ᗜ 🔑 ● 🌎 ◯ ⚡ ☰ 💫 📦 ▻ ▸ ► ▷ ◦ ▲ ⚙ ⿳ ▣ ⬤ ⬜ 👁 👀
)
