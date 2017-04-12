package style

import (
	"fmt"
	"strings"
	"bytes"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"github.com/lunixbochs/vtclean"
)

type AnsiCode int

const (
	BrightWhite  AnsiCode = 1
	Subdued      AnsiCode = 2

	// 16 COLORS
	Black        AnsiCode = 30
	Red          AnsiCode = 31
	Green        AnsiCode = 32
	Yellow       AnsiCode = 33
	Blue         AnsiCode = 34
	Magenta      AnsiCode = 35
	Cyan         AnsiCode = 36
	LightGrey    AnsiCode = 37
	White        AnsiCode = 38

	CyanBg       AnsiCode = 46
	DarkGrey     AnsiCode = 90
	LightRed     AnsiCode = 91
	LightGreen   AnsiCode = 92
	LightYellow  AnsiCode = 93
	LightBlue    AnsiCode = 94
	LightMagenta AnsiCode = 95
	LightCyan    AnsiCode = 96
	White97      AnsiCode = 97

	// 256 COLORS
	LightBlue104 AnsiCode = 104
	Black0     AnsiCode = 0
	Black231     AnsiCode = 231
	Black232     AnsiCode = 232
	Black233     AnsiCode = 233
	Black234     AnsiCode = 234
	Grey234      AnsiCode = 234
	Grey236      AnsiCode = 236
	Grey238      AnsiCode = 238
	Grey240      AnsiCode = 240
	Grey241      AnsiCode = 241
	Grey243      AnsiCode = 243
	White15      AnsiCode = 15
	OffWhite248  AnsiCode = 248
	OffWhite249  AnsiCode = 249
	OffWhite250  AnsiCode = 250
	OffWhite253  AnsiCode = 253
	OffWhite254  AnsiCode = 254
	OffWhite255  AnsiCode = 255

	Bold AnsiCode = 1
	Dim AnsiCode = 22
	Regular AnsiCode = 5
	Underline AnsiCode = 4

	ClearLine = "\033[1K"
	MoveBack  = "\033[9D"
	Block     = "2588"

	Start     = "\033["
	End       = "\033[0m"
	Start255     = "\033[48;5;"
	End255       = "\033[0;00m"
	HideCursor   = "\033[?25l"
	ShowCursor   = "\033[?25h"

	Warning        = "\xE2\x9A\xA0"
	Fire           = "\xF0\x9F\x94\xA5"
	Ambulance      = "\xF0\x9F\x9A\x91"
	InfoDeskPerson = "\xF0\x9F\x92\x81"

	Icon_App     = "✿" //✱ ▚ ❖  ꌳ ⧓ ⧗ 〓 ⁘ ꌳ ⁑⁘ ⁙ ѧꊞ ▚ 囙"
	Icon_Snippet = "◆"
	Icon_View = "❍"  // 274d
	Icon_Tick = "✓"  // 2713
	Icon_Cross = "✘" // 2718
	Icon_PrivatePouch = "◤"
	Icon_Broke = "▦"

	Color_BrightRed = 196
	Color_PouchCyan = 122
	Color_DimRed = 124
	Color_BrightWhite = 250
	Color_BrighterWhite = 252
	Color_BrightestWhite = 254
	Color_WeekGrey = 247
	Color_MonthGrey = 245
	Color_OldGrey = 242
	Color_DimStat = 238
	Color_YesGreen = 119

	MARGIN   = "  "
	TWOLINES = "\n\n"

	// 🔰 👝 🔒 🔸 ⚡ ✓ ⇨ ᗜ 🔑 ● 🌎 ◯ ⚡ ☰ 💫 📦 ▻ ▸ ► ▷ ◦ ▲ ⚙ ⿳ ▣ ⬤ ⬜ 👁 👀
)

var(
	Pad_1_2 = strings.Repeat(Fmt256(100, " "), 2)
	Pad_2_2 = strings.Repeat(Fmt256(100, "  "), 2)
	Pad_1_1 = strings.Repeat(Fmt256(100, " "), 1)
	Pad_0_0 = strings.Repeat(Fmt256(30, ""), 1)

	Pad16_0_0 = strings.Repeat(Fmt16(0, ""), 1)
	Pad16_1_1 = strings.Repeat(Fmt16(0, " "), 1)
	Pad16_2_1 = strings.Repeat(Fmt16(0, "  "), 1)
	Pad16_3_1 = strings.Repeat(Fmt16(0, "   "), 1)
	Pad16_4_1 = strings.Repeat(Fmt16(0, "    "), 1)
	Pad16_1_2 = strings.Repeat(Fmt16(0, " "), 2)
	Pad16_2_2 = strings.Repeat(Fmt16(0, "  "), 2)
	Pad16_3_2 = strings.Repeat(Fmt16(0, "   "), 2)
	Pad16_3_4 = strings.Repeat(Fmt16(0, " "), 4)
)

func FStart(c AnsiCode, in interface{}) string {
	return fmt.Sprintf("\033[%dm%v", c, in)
}

func fmtColor (c AnsiCode, in interface{}, ansiPattern string) string {
	a := strings.Split(fmt.Sprintf("%v", in), "\n")
	for i, v := range a {
		ansi := fmt.Sprintf("%s%s%dm%s%s", Start, ansiPattern, c, v, End)
		if models.Prefs() != nil && models.Prefs().PrintAnsi {
			a[i] = fmt.Sprintf("%q", ansi)
		} else {
			a[i] = ansi
		}
	}
	return strings.Join(a, "\n")
}

func Fmt16(c AnsiCode, in interface{}) string {
	return fmtColor(c, in,"")
}

func Fmt256(c AnsiCode, in interface{}) string {
	return fmtColor(c, in,"38;5;")
}

func FmtFgBg(in string, fg AnsiCode, bg AnsiCode) string {
	r := fmt.Sprintf("%s38;5;%dm%s48;5;%dm%s%s", Start, fg, Start, bg, in, End)
	return r
}

func FPreview(in string, wrapAt int, lines int) string {
	if models.Prefs().DisablePreview {
		return ""
	}
	in = vtclean.Clean(in, false)
	return FBox(in, wrapAt, lines) + End
}

/*
 Creates a text box constrained by width (number of runes) and number of lines.
 */
func FBox(in string, wrapAt int, lines int) string {
	in = strings.Replace(in, "\n", "  ", -1)
	in = strings.TrimSpace(in)
	var numRunes = 0
	b := bytes.Buffer{}
	var trim bool
	lineCount := 1
	for _, r := range in {
		numRunes++
		if trim && r == ' ' {
			continue
		}
		trim = false
		b.WriteRune(r)
		if numRunes%wrapAt == 0 {
			if lineCount >= lines {
				return strings.TrimSpace(b.String())
			}
			b.WriteString("\n")
			lineCount++
			trim = true
		}
	}
	return strings.TrimSpace(b.String())
}

/*
 Squeeze pushes text into a given width truncating the middle.
 */
func Squeeze(text string) string {
	if len(text) >= 40 {
		text = text[0:10] + "..." + text[len(text)-30:]
	}
	return " " + text
}
