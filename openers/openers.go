package openers

import (
	"github.com/kwk-links/kwk-cli/system"
)

func Open(uri string) {
	system.ExecSafe("open", uri)
}

func OpenCovert(uri string){
	system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}