package cmd

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"os"
	"time"
	"fmt"
)

const (
	filecache = "filecache"
)

type StdRunner struct {
	snippets snippets.Service
	system   sys.Manager
	writer   tmpl.Writer
}

func NewStdRunner(s sys.Manager, ss snippets.Service, w tmpl.Writer) *StdRunner {
	return &StdRunner{snippets: ss, system: s, writer: w}
}

func (r *StdRunner) Edit(alias *models.Snippet) error {
	filePath, err := r.system.WriteToFile(filecache, alias.FullName, alias.Snip)
	if err != nil {
		return err
	}
	fi, _ := os.Stat(filePath)
	r.system.ExecSafe("open", filePath)
	edited := false
	for edited == false {
		if fi2, _ := os.Stat(filePath); fi2.ModTime().UnixNano() > fi.ModTime().UnixNano() {
			edited = true
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}

	closer := func() {
		r.system.ExecSafe("osascript", "-e",
			fmt.Sprintf("tell application %q to close (every window whose name is \"%s.%s\")", "XCode", alias.Name, alias.Extension))
		r.system.ExecSafe("osascript", "-e", "tell application \"iTerm2\" to activate")
	}

	if text, err := r.system.ReadFromFile(filecache, alias.FullName); err != nil {
		closer()
		return err
	} else {
		if alias, err = r.snippets.Patch(alias.FullName, alias.Snip, text); err != nil {
			closer()
			return err
		}
		closer()
		return nil
	}
}

func (r *StdRunner) Run(alias *models.Snippet, args []string) error {
	if len(args) > 0 {
		if args[0] == "covert" {
			r.OpenCovert(alias.Snip)
			return nil
		}
		if args[0] == "web" {
			r.OpenWeb(alias)
			return nil
		}
	}

	snippet := alias.Snip
	if alias.Runtime == "url" {
		r.system.ExecSafe("open", snippet)
		return nil
	}
	if alias.Runtime == "bash" {
		args = append([]string{snippet}, args...)
		r.system.ExecSafe("/bin/bash", append([]string{"-c"}, args...)...)
		return nil
	}
	if alias.Runtime == "nodejs" {
		args = append([]string{snippet}, args...)
		// -r (require flag)
		r.system.ExecSafe("node", append([]string{"-e"}, args...)...)
		return nil
	}
	if alias.Runtime == "python" {
		args = append([]string{snippet}, args...)
		r.system.ExecSafe("python", append([]string{"-c"}, args...)...)
		return nil
	}
	if alias.Runtime == "php" {
		args = append([]string{snippet}, args...)
		r.system.ExecSafe("php", append([]string{"-r"}, args...)...)
		return nil
	}
	if alias.Runtime == "csharp" {
		return nil
	}
	if alias.Runtime == "applescript" {
		r.system.ExecSafe("osascript", append([]string{"-e"}, args...)...)
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

func (r *StdRunner) OpenWeb(alias *models.Snippet) {
	r.system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", fmt.Sprintf("http://aus.kwk.co/%s/%s", alias.Username, alias.FullName))
	r.system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}

func (r *StdRunner) OpenCovert(snippet string) {
	r.system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", snippet)
	r.system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}
