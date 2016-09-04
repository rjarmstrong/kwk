package openers

import (
	"github.com/kwk-links/kwk-cli/system"
	"strings"
)

func Open(uri string) {
	tokens := strings.Split(uri, " ")
	if tokens[0] == "sudo"{
		system.ExecSafe("/bin/bash", "-c", uri)
		return
	}
	if tokens[0] == "node"{
		i := strings.Join(tokens[1:], " ")
		system.ExecSafe("node", "-e", i)
		return
	}
	system.ExecSafe("open", uri)
}

func OpenCovert(uri string){
	system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}