package openers

import (
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"strings"
	"fmt"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
	"os"
	"time"
	"github.com/kwk-links/kwk-cli/libs/models"
	"github.com/kwk-links/kwk-cli/libs/services/aliases"
)

const (
	filecache = "filecache"
)

type Opener struct {
	aliases aliases.IAliases
	system  system.ISystem
}

func NewOpener(system system.ISystem, aliases aliases.IAliases) *Opener {
	return &Opener{aliases:aliases, system:system}
}

var iterationCount = 0

func printUri(uri string) {
	fmt.Printf(gui.Colour(gui.LightBlue, " %d - %s\n"), iterationCount, uri)
}

func (o *Opener) Edit(alias *models.Alias) error {
	filePath, err := o.system.WriteToFile(filecache, alias.FullKey, alias.Uri)
	if err != nil {
		return err
	}
	fi, _ := os.Stat(filePath);
	o.system.ExecSafe("open", filePath)
	fmt.Println(gui.Colour(gui.LightBlue, "Editing file in default editor. Please save and close to continue. Or Ctrl+C to abort."))
	edited := false
	for edited == false {
		if fi2, _ := os.Stat(filePath); fi2.ModTime().UnixNano() > fi.ModTime().UnixNano() {
			edited = true
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
	if text, err := o.system.ReadFromFile(filecache, alias.FullKey); err != nil {
		return err
	} else {
		alias = o.aliases.Patch(alias.FullKey, text)
		fmt.Println(gui.Colour(gui.LightBlue, "Successfully updated " + alias.FullKey))
		return nil
	}
}

func (o *Opener) Open(alias *models.Alias, args []string) {

	if args[0] == "covert" {
		o.OpenCovert(alias.Uri)
		return
	}

	uri := alias.Uri
	iterationCount += 1
	if iterationCount > 3 {
		fmt.Println("Max recursion reached.")
		return
	}
	//printUri(link.Uri)

	if alias.Media == "script" {
		if alias.Runtime == "app" {
			// Should we be string replacing??
			if len(args) == 1 {
				uri = strings.Replace(uri, "${1}", args[0], 1)
			}
			o.system.ExecSafe("/bin/bash", "-c", uri)
		}
		if alias.Runtime == "bash" {
			o.system.ExecSafe("/bin/bash", append([]string{"-c", alias.Uri}, args...)...)
			return
		}
		if alias.Runtime == "nodejs" {
			args = append([]string{uri}, args...)
			// -r (require flag)
			o.system.ExecSafe("node", append([]string{"-e"}, args...)...)
			return
		}
		if alias.Runtime == "python" {
			args = append([]string{uri}, args...)
			o.system.ExecSafe("python", append([]string{"-c"}, args...)...)
			return
		}
		if alias.Runtime == "php" {
			args = append([]string{uri}, args...)
			o.system.ExecSafe("php", append([]string{"-r"}, args...)...)
			return
		}
		if alias.Runtime == "csharp" {
			return
		}
		if alias.Runtime == "golang" {
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
		if alias.Runtime == "rust" {
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
		if alias.Runtime == "scala" {
			// scalac HelloWorld.scala
			// args
			// scala HelloWorld
		}
		if alias.Runtime == "java" {
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
	independants := strings.Split(alias.Uri, " && ")
	// This model is a bit odd but necessary to get around the locking issue, will have to redirect stdin and out to make piping work.
	for _, v := range independants {
		if len(v) > 3 && v[0:4] == "kwk " {
			args := strings.Split(v, " ")
			firstArg := args[1]
			if firstArg == "upgrade" {
				o.system.Upgrade()
				return
			}
			alias := o.aliases.Get(firstArg)
			if alias.Total == 1 {
				o.Open(&alias.Items[0], args)
			} else if alias.Total > 1 {
				fmt.Println("More than one option please give a file extension.")
			} else {
				fmt.Printf(gui.Colour(gui.Yellow, "Can't run sub-command: '%s' - has it been deleted?\n"), v)
			}

		} else {
			o.system.ExecSafe("/bin/bash", "-c", v)
		}
	}
}

func (o *Opener) OpenCovert(uri string) {
	o.system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	o.system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}