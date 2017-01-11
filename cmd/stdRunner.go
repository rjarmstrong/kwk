package cmd

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"runtime"
	"os/exec"
	"errors"
	"bytes"
	"time"
	"fmt"
	"io"
	"os"
)

const (
	CONF_PATH       = "conf"
	ENV_PATH        = "env"
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
	//TODO: if we pull out the env from getSection we can improve speed
	a, err := r.getSection("apps")
	eRoot, err := r.getSection("editors")

	if err != nil {
		return err
	}
	_, candidates := getSubSection(eRoot, s.Extension)
	if len(candidates) != 1 {
		return errors.New("No editors have been specified for " + s.Extension + " . And default editor is not specified.")
	}
	_, cli := getSubSection(a, candidates[0])

	filePath, err := r.system.WriteToFile(FILE_CACHE_PATH, s.FullName, s.Snip, true)
	if err != nil {
		return err
	}
	replaceVariables(&cli, filePath, s)
	fi, _ := os.Stat(filePath)
	openTime := fi.ModTime().UnixNano()

	execSafe(cli[0], cli[1:]...)

	edited := false
	for edited == false {
		if fi2, _ := os.Stat(filePath); fi2.ModTime().UnixNano() > openTime {
			edited = true
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}

	closer := func() {
		if runtime.GOOS == sys.OS_DARWIN {
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

func (r *StdRunner) getEnv() (string, error) {
	if os.Getenv(sys.KWK_TESTMODE) != "" {
		testEnv := "./cmd/testEnv.yml"
		// TODO: use log
		fmt.Println(">> Running with:", testEnv, " <<")
		b, err := ioutil.ReadFile(testEnv)
		return string(b), nil
		if err != nil {
			return "", err
		}
	}
	envFullName := fmt.Sprintf("%s-%s.yml", runtime.GOOS, runtime.GOARCH)
	// TODO: check yml version is compatible with this build else force upgrade.
	if ok, _ := r.system.FileExists(ENV_PATH, envFullName, false); !ok {
		// TODO: use log
		fmt.Println(">> No local env.yml getting remote. <<")
		if l, err := r.snippets.Get(&models.Alias{FullKey:envFullName, Username:ENV_USERNAME}); err != nil {
			return "", err
		} else {
			env := l.Items[0].Snip
			if _, err := r.system.WriteToFile(ENV_PATH, envFullName, env, false); err != nil {
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

func (r *StdRunner) getSection(name string) (*yaml.MapSlice, error) {
	env, err := r.getEnv()
	if err != nil {
		return nil, err
	}
	c := &yaml.MapSlice{}
	if err := yaml.Unmarshal([]byte(env), c); err != nil {
		return nil, err
	}
	rs, _ := getSubSection(c, name)
	if rs == nil {
		return nil, errors.New(fmt.Sprintf("No %s section in env.yml", name))
	}
	return &rs, nil
}

func getSubSection(yml *yaml.MapSlice, name string) (yaml.MapSlice, []string) {
	f := func(yml *yaml.MapSlice, name string) (yaml.MapSlice, []string) {
		for _, v := range *yml {
			if v.Key == name {
				if slice, ok := v.Value.(yaml.MapSlice); ok {
					return slice, nil
				}
				if _, ok := v.Value.([]interface{}); ok {
					items := []string{}
					for _, v2 := range v.Value.([]interface{}) {
						items = append(items, v2.(string))
					}
					return nil, items
				}
				return nil, []string{v.Value.(string)}
			}
		}
		return nil, nil
	}
	if sub, bottom := f(yml, name); sub == nil && bottom == nil {
		return f(yml, "default")
	} else {
		return sub, bottom
	}
}

/*
 	$FULL_NAME = full name of the snippet e.g. `hello.js`
 	$NAME = name excluding extension e.g. `hello`
 	$DIR = directory of the snippet on disk. Useful when editing a file in a directory structure or when compilation needs it.
 	$CLASS_NAME = for java and scala these will be the class name in the snippet. Used when attempting to run the compiled file.
 */
func replaceVariables(cliArgs *[]string, filePath string, s *models.Snippet) {
	for i := range *cliArgs {
		(*cliArgs)[i] = strings.Replace((*cliArgs)[i], "$FULL_NAME", filePath, -1)
		(*cliArgs)[i] = strings.Replace((*cliArgs)[i], "$DIR", strings.Replace(filePath, s.FullName, "", -1), -1)
		(*cliArgs)[i] = strings.Replace((*cliArgs)[i], "$NAME", strings.Replace(filePath, "."+s.Extension, "", -1), -1)
	}
}

func (r *StdRunner) Run(s *models.Snippet, args []string) error {
	rs, err := r.getSection("runners")
	if err != nil {
		return err
	}
	comp, interp := getSubSection(rs, s.Extension)
	if comp != nil {
		if filePath, err := r.system.WriteToFile(FILE_CACHE_PATH, s.FullName, s.Snip, true); err != nil {
			return err
		} else {
			_, compile := getSubSection(&comp, "compile")
			if compile != nil {
				replaceVariables(&compile, filePath, s)
				// TODO: ADD TO LOG
				fmt.Println("compile", compile)
				execSafe(compile[0], compile[1:]...).Close()
			}
			_, run := getSubSection(&comp, "run")
			replaceVariables(&compile, filePath, s)

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
