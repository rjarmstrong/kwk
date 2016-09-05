package openers

import (
	"github.com/kwk-links/kwk-cli/system"
	"strings"
	"github.com/kwk-links/kwk-cli/api"
	"fmt"
	"github.com/kwk-links/kwk-cli/gui"
)

type Opener struct {
	apiClient *api.ApiClient
}

func NewOpener(apiClient *api.ApiClient) *Opener {
	return &Opener{apiClient:apiClient}
}

var iterationCount = 0

func printUri(uri string){
	fmt.Printf(gui.Colour(gui.LightBlue, " %d - %s\n"), iterationCount, uri)
}

func (o *Opener) Open(uri string) {
	iterationCount += 1
	if iterationCount > 3 {
		fmt.Println("Max recursion reached.")
		return
	}
	printUri(uri)
	tokens := strings.Split(uri, " ")
	if tokens[0] == "sudo"{
		script := strings.Replace(uri, "sudo ", "", -1)
		independants := strings.Split(script, " && ")
		for _, v := range independants {
			if len(v) > 3 && v[0:4] == "kwk " {
				args :=strings.Split(v, " ")
				firstArg := args[1]
				if firstArg == "upgrade" {
					system.Upgrade()
					return
				}
				link := o.apiClient.Get(firstArg)
				if link.Uri != "" {
					o.Open(link.Uri)
				} else {
					fmt.Printf(gui.Colour(gui.Yellow, "Can't run sub-command: '%s' - has it been deleted?\n"), v)
				}

			} else {
				system.ExecSafe("/bin/bash", "-c", v)
			}
		}
		return
	}
	if tokens[0] == "node"{
		i := strings.Join(tokens[1:], " ")
		system.ExecSafe("node", "-e", i)
		return
	}
	system.ExecSafe("open", uri)
}

func (o *Opener) OpenCovert(uri string){
	system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}