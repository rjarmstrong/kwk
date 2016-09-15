package openers

import (
	"github.com/kwk-links/kwk-cli/system"
	"strings"
	"github.com/kwk-links/kwk-cli/api"
	"fmt"
	"github.com/kwk-links/kwk-cli/gui"
	"os"
	"time"
)

const (
	filecache = "filecache"
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

func (o *Opener) Edit(key string) error {
	kwklink := o.apiClient.Get(key)
	filePath, err := system.WriteToFile(filecache, kwklink.FullKey, kwklink.Uri)
	if err != nil {
		return err
	}
	fi, _ := os.Stat(filePath);
	system.ExecSafe("open", filePath)
	fmt.Println(gui.Colour(gui.LightBlue, "Editing file in default editor. Please save and close to continue. Or Ctrl+C to abort."))
	edited := false
	for edited == false {
		if fi2, _ := os.Stat(filePath); fi2.ModTime().UnixNano() > fi.ModTime().UnixNano() {
			edited = true
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
	if text, err := system.ReadFromFile(filecache, kwklink.FullKey); err != nil {
		return err
	} else {
		kwklink = o.apiClient.Patch(kwklink.FullKey, text)
		fmt.Println(gui.Colour(gui.LightBlue, "Successfully updated " + kwklink.FullKey))
		return nil
	}
}

func (o *Opener) Open(link *api.KwkLink, args []string) {
	uri := link.Uri
	iterationCount += 1
	if iterationCount > 3 {
		fmt.Println("Max recursion reached.")
		return
	}
	//printUri(link.Uri)

	if link.Media == "script" {
		if link.Runtime == "app" {
			// Should we be string replacing??
			if len(args) == 1 {
				uri = strings.Replace(uri, "${1}", args[0], 1)
			}
			system.ExecSafe("/bin/bash", "-c", uri)
		}
		if link.Runtime == "bash" {
			system.ExecSafe("/bin/bash", append([]string{"-c", link.Uri}, args...)...)
			return
		}
		if link.Runtime == "nodejs" {
			args = append([]string{uri}, args...)
			// -r (require flag)
			system.ExecSafe("node", append([]string{"-e"}, args...)...)
			return
		}
		if link.Runtime == "python" {
			args = append([]string{uri}, args...)
			system.ExecSafe("python", append([]string{"-c"}, args...)...)
			return
		}
		if link.Runtime == "php" {
			args = append([]string{uri}, args...)
			system.ExecSafe("php", append([]string{"-r"}, args...)...)
			return
		}
		if link.Runtime == "csharp" {
			return
		}
		if link.Runtime == "golang" {
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
		if link.Runtime == "rust" {
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
		if link.Runtime == "scala" {
			// scalac HelloWorld.scala
			// args
			// scala HelloWorld
		}
		if link.Runtime == "java" {
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
				o.Open(link, args)
			} else {
				fmt.Printf(gui.Colour(gui.Yellow, "Can't run sub-command: '%s' - has it been deleted?\n"), v)
			}

		} else {
			system.ExecSafe("/bin/bash", "-c", v)
		}
	}
}

func (o *Opener) OpenCovert(uri string) {
	system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}