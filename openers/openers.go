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

func (o *Opener) Open(link *api.KwkLink, param string) {
	iterationCount += 1
	if iterationCount > 3 {
		fmt.Println("Max recursion reached.")
		return
	}
	printUri(link.Uri)
	if link.Media == "script" {
		uri := link.Uri
		if strings.Contains(uri, "[param1]") && param == "" {
			fmt.Println(gui.Colour(gui.Yellow, "A [param1] is required."))
			return
		}
		if param != "" {
			uri = strings.Replace(uri, "[param1]", param, -1)
		}
		if link.Type == "bash" {
			system.ExecSafe("/bin/bash", "-c", uri)
			return
		}
		if link.Type == "nodejs" {
			system.ExecSafe("node", "-e", uri)
			return
		}
		if link.Type == "python" {
			system.ExecSafe("python", "-c", uri)
			return
		}

	}
	tokens := strings.Split(link.Uri, " ")
	if tokens[0] == "sudo"{
		script := strings.Replace(link.Uri, "sudo ", "", -1)
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
					o.Open(link, args[2])
				} else {
					fmt.Printf(gui.Colour(gui.Yellow, "Can't run sub-command: '%s' - has it been deleted?\n"), v)
				}

			} else {
				if param != "" {
					v = strings.Replace(v, "[param1]", param, -1)
				}
				system.ExecSafe("/bin/bash", "-c", v)
			}
		}
		return
	}
	system.ExecSafe("open", link.Uri)
}

func (o *Opener) OpenCovert(uri string){
	system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}