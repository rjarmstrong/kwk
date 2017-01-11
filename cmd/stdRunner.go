package cmd

import (
	"fmt"
	"os"
	"time"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"gopkg.in/yaml.v2"
	"strings"
	"runtime"
	"io"
	"os/exec"
	"bytes"
	"io/ioutil"

)

const (
	CONF_PATH       = "conf"
	ENV_PATH       = "env"
	FILE_CACHE_PATH = "filecache"
	ENV_USERNAME    = "env"
)

type StdRunner struct {
	snippets snippets.Service
	system   sys.Manager
	writer   tmpl.Writer
}

func NewStdRunner(s sys.Manager, ss snippets.Service, w tmpl.Writer) *StdRunner {
	return &StdRunner{snippets: ss, system: s, writer: w}
}

func (r *StdRunner) Edit(s *models.Snippet) error {
	filePath, err := r.system.WriteToFile(FILE_CACHE_PATH, s.FullName, s.Snip, true)
	if err != nil {
		return err
	}
	fi, _ := os.Stat(filePath)
	// replace with configured editor
	execSafe("open", filePath)

	// Will this work in headless mode? I think so, but need to check thread locking.
	edited := false
	for edited == false {
		if fi2, _ := os.Stat(filePath); fi2.ModTime().UnixNano() > fi.ModTime().UnixNano() {
			edited = true
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}

	closer := func() {
		if runtime.GOOS == sys.OS_DARWIN {
			execSafe("osascript", "-e",
				fmt.Sprintf("tell application %q to close (every window whose name is \"%s.%s\")", "XCode", s.Name, s.Extension))
			// This assumes we are using iTerm, we'll have to get the active shell (echo $TERM_PROGRAM on Mac)
			execSafe("osascript", "-e", "tell application \"iTerm2\" to activate")
		} else if runtime.GOOS == sys.OS_WINDOWS {
			//	// How will this work on:
			//	- windows https://technet.microsoft.com/en-us/library/ee176882.aspx
		}
	}

	if text, err := r.system.ReadFromFile(FILE_CACHE_PATH, s.FullName, true); err != nil {
		// if there was an error close the application anyway
		closer()
		return err
	} else {
		// else save and close the app
		if s, err = r.snippets.Patch(s.FullName, s.Snip, text); err != nil {
			closer()
			return err
		}
		closer()
		return nil
	}
}

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

func (r *StdRunner) getEnv() (string, error) {
	envFullName := fmt.Sprintf("%s-%s.yml", runtime.GOOS, runtime.GOARCH)
	if ok, _ := r.system.FileExists(ENV_PATH, envFullName, false); !ok {
		fmt.Println("getting remote")
		if l, err := r.snippets.Get(&models.Alias{FullKey:envFullName, Username:ENV_USERNAME}); err != nil {
			return "", err
		} else {
			env := l.Items[0].Snip
			if _, err := r.system.WriteToFile(ENV_PATH, envFullName,  env,false); err != nil {
				return "", err
			}
			return env, nil
		}
	} else {
		if e, err := r.system.ReadFromFile(ENV_PATH, envFullName, false); err != nil {
			return "", err
		} else {
			return e, nil
		}
	}
}

func (r *StdRunner) getRunnerEnv() (*yaml.MapSlice, error) {
	//env, err := r.getEnv()
	b, err := ioutil.ReadFile("./cmd/testEnv.yml")
	env := string(b)
	if err != nil {
		return nil, err
	}

	c := &yaml.MapSlice{}
	if err := yaml.Unmarshal([]byte(env), c); err != nil {
		return nil, err
	}
	rs, _ := getSection(c, "runners")
	if rs == nil {
		panic("no runners section")
	}
	return &rs, nil
}

func (r *StdRunner) Run(s *models.Snippet, args []string) error {
	rs, err := r.getRunnerEnv()
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
		if filePath, err := r.system.WriteToFile(FILE_CACHE_PATH, s.FullName, s.Snip, true); err != nil {
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
				execSafe(compile[0], compile[1:]...).Close()
			}
			_, run := getSection(&comp, "run")
			for i := range run {
				run[i] = strings.Replace(run[i], "$FULL_NAME", filePath, -1)
				run[i] = strings.Replace(run[i], "$NAME", strings.Replace(filePath, "."+s.Extension, "", -1), -1)
			}

			//fmt.Println("run", run)
			run = append(run, args...)
			execSafe(run[0], run[1:]...).Close()
		}
	} else {
		//fmt.Println(runner)
		for i, v := range interp {
			interp[i] = strings.Replace(v, "$SNIP", s.Snip, -1)
		}
		interp = append(interp, args...)
		//fmt.Println(runner)
		execSafe(interp[0], interp[1:]...).Close()
	}
	return nil
}

func execSafe(name string, arg ...string) io.ReadCloser {
	c := exec.Command(name, arg...)
	c.Stdin = os.Stdin
	out, _ := c.StdoutPipe()
	var stderr bytes.Buffer
	c.Stdout = os.Stdout
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
	return out
}

func (r *StdRunner) OpenWeb(alias *models.Snippet) {
	execSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", fmt.Sprintf("https://www.kwk.co/%s/%s", alias.Username, alias.FullName))
	execSafe("osascript", "-e", "activate application \"Google Chrome\"")
}

func (r *StdRunner) OpenCovert(snippet string) {
	execSafe("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", snippet)
	execSafe("osascript", "-e", "activate application \"Google Chrome\"")
}
