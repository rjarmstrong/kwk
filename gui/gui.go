package gui

import (
	"fmt"
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
)

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
