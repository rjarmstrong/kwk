package style

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	BrightWhite        = 1
	Subdued            = 2
	Black              = 30
	Red                = 31
	Green              = 32
	Yellow             = 33
	DarkBlue           = 34
	Pink               = 35
	LightBlue          = 36
	LightBlueHighlight = 46
	ClearLine          = "\033[1K"
	MoveBack           = "\033[9D"
	Block              = "2588"
	Start              = "\033["
	End                = "\033[0m"
	Space              = 20
	UBlock             = "\u2588"
	Lock               = "\xF0\x9F\x94\x92"
	OpenLock           = "\xF0\x9F\x94\x93"
	Warning            = "\xE2\x9A\xA0"
	Fire		   = "\xF0\x9F\x94\xA5"
	Tick		   = "\xE2\x9C\x93"
	Ambulance	   = "\xF0\x9F\x9A\x91"
)

func Build(quant int, unicode string) string {
	var str string
	for i := 0; i <= quant; i++ {
		str = str + unicode
	}
	return str
}

func Italic(str string) string {
	return Start + "3m" + str + End
}

func Underline(str string) string {
	return "\033[4m" + str + "\033[0m"
}

func Bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

// Colour creates colored output for the CLI.
func Colour(colour int, params ...interface{}) string {
	var text string
	for _, v := range params {
		text = text + fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("\033[%dm%v\033[0m", colour, text)
}

func ColourSpan(colour int, text string, openTag string, closeTag string, surroundingColor int) string {
	text = strings.Replace(text, openTag, End+"\033["+strconv.Itoa(colour)+"m", -1)
	text = strings.Replace(text, closeTag, "\033[0m"+Start+strconv.Itoa(surroundingColor)+"m", -1)
	return text
}
