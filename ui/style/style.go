package style

import (
	"fmt"
	"strings"
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
	CyanBg       AnsiCode = 46
	DarkGrey     AnsiCode = 90
	LightRed     AnsiCode = 91
	LightGreen   AnsiCode = 92
	LightYellow  AnsiCode = 93
	LightBlue    AnsiCode = 94
	LightMagenta AnsiCode = 95
	LightCyan    AnsiCode = 96
	White        AnsiCode = 97

	Bold AnsiCode = 21
	Dim AnsiCode = 22
	Underline AnsiCode = 23

	ClearLine = "\033[1K"
	MoveBack  = "\033[9D"
	Block     = "2588"
	Start     = "\033["
	End       = "\033[0m"
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
)

func Build(quant int, unicode string) string {
	var str string
	for i := 0; i <= quant; i++ {
		str = str + unicode
	}
	return str
}

// Fmt formats output for the CLI.
func Fmt(c AnsiCode, params ...interface{}) string {
	var text string
	for _, v := range params {
		text = text + fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("\033[%dm%v\033[0m", c, text)
}

func ColourSpan(colour AnsiCode, text string, openTag string, closeTag string, surroundingColor AnsiCode) string {
	text = strings.Replace(text, openTag, fmt.Sprintf("%s\033[%dm[", End, colour), -1)
	text = strings.Replace(text, closeTag, fmt.Sprintf("\033[0m%s%dm", Start, surroundingColor), -1)
	return text
}
