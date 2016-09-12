package openers

import (
	"github.com/kwk-links/kwk-cli/system"
	"strings"
	"github.com/kwk-links/kwk-cli/api"
	"fmt"
	"github.com/kwk-links/kwk-cli/gui"
	"github.com/kennygrant/sanitize"
	"io/ioutil"
	"os"
	"time"
)

type Opener struct {
	apiClient *api.ApiClient
}

func NewOpener(apiClient *api.ApiClient) *Opener {
	return &Opener{apiClient:apiClient}
}

var iterationCount = 0

func printUri(uri string) {
	fmt.Printf(gui.Colour(gui.LightBlue, " %d - %s\n"), iterationCount, uri)
}

func (o *Opener) Edit(key string) {
	kwklink := o.apiClient.Get(key)
	path := system.GetCachePath() + "/" + sanitize.Name(kwklink.Key) + ".js"
	if err := ioutil.WriteFile(path, []byte(kwklink.Uri), 0666); err != nil {
		panic(err)
	}
	fi, _ := os.Stat(path);

	system.ExecSafe("open", path)
	fmt.Println(gui.Colour(gui.LightBlue, "Editing file in default editor. Please save and close to continue. Or Ctrl+C to abort."))

	edited := false
	for edited == false {
		if fi2, _ := os.Stat(path); fi2.ModTime().UnixNano() > fi.ModTime().UnixNano() {
			edited = true
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}

	if bts, err := ioutil.ReadFile(path); err != nil {
		panic(err)
	} else {
		kwklink = o.apiClient.Patch(kwklink.Key, string(bts))
		fmt.Println(gui.Colour(gui.LightBlue, "Successfully updated " + kwklink.Key))
	}
}

func (o *Opener) Open(link *api.KwkLink, args string) {
	uri := link.Uri
	iterationCount += 1
	if iterationCount > 3 {
		fmt.Println("Max recursion reached.")
		return
	}
	printUri(link.Uri)
	if link.Media == "script" {
		if link.Type == "bash" {
			system.ExecSafe("/bin/bash", "-c", uri)
			return
		}
		if link.Type == "nodejs" {
			// -r (require flag)
			// node -e "script" args
			system.ExecSafe("node", "-e", uri)
			return
		}
		if link.Type == "python" {
			// -c "" args
			system.ExecSafe("python", "-c", uri)
			return
		}
		if link.Type == "php" {
			// -r "" -- args
			system.ExecSafe("php", "-r", uri)
			return
		}
		if link.Type == "csharp" {
			return
		}
		if link.Type == "golang" {
			// check if file exists
			// if not
			// write file to disk in cache
			// compile it
			//system.ExecSafe("go", "build", key .go)
			// java file name
			// run it
			// args
			//system.ExecSafe(key)
		}
		if link.Type == "rust" {
			// check if file exists
			// if not
			// write file to disk in cache
			// compile it
			//system.ExecSafe("rustc", key .rs)
			// java file name
			// run it
			// args
			//system.ExecSafe(file)
			return
		}
		if link.Type == "scala" {
			// scalac HelloWorld.scala
			// args
			// scala HelloWorld
		}
		if link.Type == "java" {
			// check if file exists
			// if not
			// write file to disk in cache
			// compile it
			//system.ExecSafe("javac", key .java)
			// java file name
			// run it
			// args
			//system.ExecSafe("java", key .class)
			return
		}

	}
	independants := strings.Split(link.Uri, " && ")
	// This model is a bit odd but necessary to get around the locking issue, will have to redirect stdin and out to make piping work.
	for _, v := range independants {
		if len(v) > 3 && v[0:4] == "kwk " {
			args := strings.Split(v, " ")
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
			if args != "" {
				v = strings.Replace(v, "[param1]", args, -1)
			}
			system.ExecSafe("/bin/bash", "-c", v)
		}
	}
	return
	system.ExecSafe("open", uri)
}

func (o *Opener) OpenCovert(uri string) {
	system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}