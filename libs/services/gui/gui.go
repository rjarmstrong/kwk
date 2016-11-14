package gui

import (
	"fmt"
	"strings"
	"strconv"
)


const (
	BrightWhite = 1
	Subdued = 2
	Black = 30
	Red = 31
	Green = 32
	Yellow = 33
	DarkBlue = 34
	Pink = 35
	LightBlue = 36
	LightBlueHighlight = 46
	ClearLine = "\033[1K"
	MoveBack = "\033[9D"
	Block = "2588"
	Start = "\033["
	End = "\033[0m"
	Space = " "
	UBlock = "\u2588"
)

func Build(quant int, unicode string) string {
	var str string
	for i := 0; i<=quant; i++ {
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

func Colour(colour int, params ...interface{}) string {
	var text string
	for _, v := range params {
		text = text + fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("\033[%dm%v\033[0m", colour, text)
}

func ColourSpan(colour int, text string, start string, end string) string {
	text = strings.Replace(text, start, "\033[" + strconv.Itoa(colour) + "m", -1)
	text = strings.Replace(text, end, "\033[0m", -1)
	return text
}


