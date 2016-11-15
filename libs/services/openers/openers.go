package openers

import (
	"os"
	"time"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
)

const (
	filecache = "filecache"
)

type Opener struct {
	aliases aliases.IAliases
	system  system.ISystem
	writer gui.ITemplateWriter
}

func New(system system.ISystem, aliases aliases.IAliases, w gui.ITemplateWriter) *Opener {
	return &Opener{aliases:aliases, system:system, writer:w}
}

func (o *Opener) Edit(alias *models.Alias) error {
	filePath, err := o.system.WriteToFile(filecache, alias.FullKey, alias.Uri)
	if err != nil {
		return err
	}
	fi, _ := os.Stat(filePath);
	o.system.ExecSafe("open", filePath)
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
		if alias, err = o.aliases.Patch(alias.FullKey, alias.Uri, text); err != nil {
			return err
		}
		return nil
	}
}

func (o *Opener) Open(alias *models.Alias, args []string) error {
	if len(args) > 0 && args[0] == "covert" {
		o.OpenCovert(alias.Uri)
		return nil
	}
	uri := alias.Uri
	if alias.Runtime == "url" {
		o.system.ExecSafe("open", uri)
		return nil
	}
	if alias.Runtime == "bash" {
		o.system.ExecSafe("/bin/bash", append([]string{"-c", uri}, args...)...)
		return nil
	}
	if alias.Runtime == "nodejs" {
		args = append([]string{uri}, args...)
		// -r (require flag)
		o.system.ExecSafe("node", append([]string{"-e"}, args...)...)
		return nil
	}
	if alias.Runtime == "python" {
		args = append([]string{uri}, args...)
		o.system.ExecSafe("python", append([]string{"-c"}, args...)...)
		return nil
	}
	if alias.Runtime == "php" {
		args = append([]string{uri}, args...)
		o.system.ExecSafe("php", append([]string{"-r"}, args...)...)
		return nil
	}
	if alias.Runtime == "csharp" {
		return nil
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
		return nil
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
		return nil
	}
	return nil
}

func (o *Opener) OpenCovert(uri string) {
	o.system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", uri)
	o.system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}