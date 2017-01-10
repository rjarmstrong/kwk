package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"gopkg.in/yaml.v2"
	"strings"
)

type StdRunner struct {
	snippets snippets.Service
	system   sys.Manager
	writer   tmpl.Writer
}

func NewStdRunner(s sys.Manager, ss snippets.Service, w tmpl.Writer) *StdRunner {
	return &StdRunner{snippets: ss, system: s, writer: w}
}

func (r *StdRunner) cacheFile(fullName string, snip string) (string, error) {
	return r.system.WriteToFile(filecache, fullName, snip)
}

func (r *StdRunner) Edit(alias *models.Snippet) error {
	filePath, err := r.cacheFile(alias.FullName, alias.Snip)
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

const (
	filecache = "filecache"
)

func getSection(yml *yaml.MapSlice, name string) (yaml.MapSlice, []string) {
	for _, v := range *yml {
		if v.Key == name {
			if slice, ok := v.Value.(yaml.MapSlice); ok {
				return slice, nil
			} else {
				items := []string{}
				for _, v2 := range v.Value.([]interface{}) {
					items = append(items, v2.(string))
				}
				return nil, items
			}
		}
	}
	return nil, nil
}

func getRunnersSection() (*yaml.MapSlice, error) {
	f, err := ioutil.ReadFile("./cmd/config.yml")
	if err != nil {
		fmt.Print(err)
	}
	c := &yaml.MapSlice{}
	if err := yaml.Unmarshal(f, c); err != nil {
		return nil, err
	}
	rs, _ := getSection(c, "runners")
	if rs == nil {
		panic("no runners section")
	}
	return &rs, nil
}

func (r *StdRunner) Run(s *models.Snippet, args []string) error {
	rs, err := getRunnersSection()
	if err != nil {
		return err
	}
	comp, interp := getSection(rs, s.Extension)
	if interp == nil && comp == nil {
		_, interp = getSection(rs, "default")
	}
	if interp == nil && comp == nil {
		panic("no default")
	}
	if comp != nil {
		if filePath, err := r.cacheFile(s.FullName, s.Snip); err != nil {
			return err
		} else {
			_, compile := getSection(&comp, "compile")
			if compile != nil {
				for i := range compile {
					compile[i] = strings.Replace(compile[i], "$FULL_NAME", filePath, -1)
					compile[i] = strings.Replace(compile[i], "$CACHE_DIR", strings.Replace(filePath, s.FullName, "", -1), -1)
					compile[i] = strings.Replace(compile[i], "$NAME", strings.Replace(filePath, "."+s.Extension, "", -1), -1)
				}
				fmt.Println("compile", compile)
				r.system.ExecSafe(compile[0], compile[1:]...).Close()
			}
			_, run := getSection(&comp, "run")
			for i := range run {
				run[i] = strings.Replace(run[i], "$FULL_NAME", filePath, -1)
				run[i] = strings.Replace(run[i], "$NAME", strings.Replace(filePath, "." + s.Extension, "", -1), -1)
			}

			//fmt.Println("run", run)
			run = append(run, args...)
			r.system.ExecSafe(run[0], run[1:]...).Close()
		}
	} else {
		//fmt.Println(runner)
		for i, v := range interp {
			interp[i] = strings.Replace(v, "$SNIP", s.Snip, -1)
		}
		interp = append(interp, args...)
		//fmt.Println(runner)
		r.system.ExecSafe(interp[0], interp[1:]...).Close()
	}
	return nil
}

func (r *StdRunner) OpenWeb(alias *models.Snippet) {
	r.system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", fmt.Sprintf("https://www.kwk.co/%s/%s", alias.Username, alias.FullName))
	r.system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}

func (r *StdRunner) OpenCovert(snippet string) {
	r.system.ExecSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", snippet)
	r.system.ExecSafe("osascript", "-e", "activate application \"Google Chrome\"")
}
