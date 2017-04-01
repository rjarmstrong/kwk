package style

import (
	"fmt"
	"strings"
	"bytes"
	"bitbucket.com/sharingmachine/kwkcli/models"
)

type AnsiCode int

const (
	BrightWhite  AnsiCode = 1
	Subdued      AnsiCode = 2

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
	Start255     = "\033[48;5;"
	End       = "\033[0m"
	End255       = "\033[0;00m"


	Space     = 20
	UBlock    = "\u2588"

	// UTF-8 Hex
	Lock           = "\xF0\x9F\x94\x92"
	OpenLock       = "\xF0\x9F\x94\x93"
	Warning        = "\xE2\x9A\xA0"
	Fire           = "\xF0\x9F\x94\xA5"
	Tick           = "\xE2\x9C\x93"
	Ambulance      = "\xF0\x9F\x9A\x91"
	Glasses        = "\xF0\x9F\x91\x93"
	InfoDeskPerson = "\xF0\x9F\x92\x81"
	Folder         = "\xF0\x9F\x93\x81"
	Task           = "\xE2\x98\xB0"
	BriefCase      = "\xF0\x9F\x92\xBC"
	Pouch          = "\xF0\x9F\x91\x9D"
	BlueDiamond    = "\xF0\x9F\x94\xB9"
	YelloDiamond   = "\xF0\x9F\x94\xB8"

	Icon_App     = "▚" //❖  ꌳ ⧓ ⧗ 〓 ⁘ ꌳ ⁑⁘ ⁙ ѧꊞ ▚ 囙"
	Icon_Snippet = "◆"

	// 🔰 👝 🔒 🔸 ⚡ ✓ ⇨ ᗜ 🔑 ● 🌎 ◯ ⚡ ☰ 💫 📦 ▻ ▸ ► ▷ ◦ ▲ ⚙ ⿳ ▣ ⬤ ⬜
)

func Build(quant int, unicode string) string {
	var str string
	for i := 0; i <= quant; i++ {
		str = str + unicode
	}
	return str
}

func FmtStart(c AnsiCode, in interface{}) string {
	return fmt.Sprintf("\033[%dm%v", c, in)
}

// Fmt formats output for the CLI.
func Fmt(c AnsiCode, params ...interface{}) string {
	var text string
	for _, v := range params {
		text = text + fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("\033[%0dm%v\033[0m", c, text)
}

func FmtFgBg(in string, fg AnsiCode, bg AnsiCode) string {
	r := fmt.Sprintf("%s38;5;%dm%s48;5;%dm%s%s", Start, fg, Start, bg, in, End)
	return r
}

func Fmt256(c AnsiCode, in interface{}) string {
	return fmt.Sprintf("%s38;5;%00dm%s%s", Start, c, in, End)
}

func ColourSpan(colour AnsiCode, text string, openTag string, closeTag string, surroundingColor AnsiCode) string {
	text = strings.Replace(text, openTag, fmt.Sprintf("%s\033[%dm[", End, colour), -1)
	text = strings.Replace(text, closeTag, fmt.Sprintf("\033[0m%s%dm", Start, surroundingColor), -1)
	return text
}

//func ColourSpan256(colour AnsiCode, text string, openTag string, closeTag string, surroundingColor AnsiCode) string {
//	text = Fmt256(surroundingColor, false, text)
//	text = strings.Replace(text, openTag, fmt.Sprintf("%s%s%dm", End255, Start255, colour), -1)
//	text = strings.Replace(text, closeTag, fmt.Sprintf("%s%s%dm", End255, Start255, surroundingColor), -1)
//	return text
//}

/*
 StyleLines formats each line with a foreground and background color.
 */
func AnsiLinesFgBg(in string, fg AnsiCode, bg AnsiCode) string {
	t := strings.Split(in, "\n")
	for i, v := range t {
		t[i] = FmtFgBg(v, fg, bg)
	}
	join := strings.Join(t, "\n")
	return join
}

/*
 StyleLines formats each line with a foreground  color.
 */
func AnsiLines(in string, fg AnsiCode) string {
	t := strings.Split(in, "\n")
	for i, v := range t {
		t[i] = Fmt(fg, v)
	}
	join := strings.Join(t, "\n")
	return join
}

func FmtPreview (in string, wrapAt int, lines int) string {
	if models.Prefs().DisablePreview {
		return ""
	}
	return FmtBox(in, wrapAt, lines)
}

/*
 Creates a text box constrained by width (number of runes) and number of lines.
 */
func FmtBox(in string, wrapAt int, lines int) string {
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
