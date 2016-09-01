package openers

import (
	"github.com/kwk-links/kwk-cli/system"
	"strings"
	"fmt"
	"github.com/kwk-links/kwk-cli/gui"
)

func Open(uri string) {
	tokens := strings.Split(uri, " ")
	if tokens[0] == "sudo"{
		fmt.Println(gui.Colour(gui.LightBlue, "kwk > ", uri))
		system.ExecSafe(tokens[1], tokens[2:]...)
		return
	}
	system.ExecSafe("open", uri)
}

func OpenCovert(uri string){
	system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}