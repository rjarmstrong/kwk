package cmd

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/setup"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"gopkg.in/yaml.v2"
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

type StdRunner struct {
	snippets snippets.Service
	system   sys.Manager
	settings config.Persister
	setup setup.Provider
}

func NewStdRunner(s sys.Manager, ss snippets.Service, setup setup.Provider) *StdRunner {
	return &StdRunner{snippets: ss, system: s,  setup: setup}
}

func (r *StdRunner) Edit(s *models.Snippet) error {
	//TODO: if we pull out the env from getSection we can improve speed
	a, err := r.getEnvSection("apps")
	eRoot, err := r.getEnvSection("editors")

	if err != nil {
		return err
	}
	_, candidates := getSubSection(eRoot, s.Ext)
	if len(candidates) != 1 {
		return errors.New("No editors have been specified for " + s.Ext + " . And default editor is not specified.")
	}
	_, cli := getSubSection(a, candidates[0])

	filePath, err := r.system.WriteToFile(setup.SNIP_CACHE_PATH, s.String(), s.Snip, true)
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

	if text, err := r.system.ReadFromFile(setup.SNIP_CACHE_PATH, s.String(), true, 0); err != nil {
		// if there was an error close the application anyway
		closer()
		return err
	} else {
		// else save and close the app
		if s, err = r.snippets.Patch(s.Alias, s.Snip, text); err != nil {
			closer()
			return err
		}
		closer()
		return nil
	}
}


func (r *StdRunner) Run(s *models.Snippet, args []string) error {
	rs, err := r.getEnvSection("runners")
	if err != nil {
		return err
	}
	yamlKey := s.Ext
	if r.setup.Prefs().Covert {
		yamlKey += "-covert"
	}
	comp, interp := getSubSection(rs, yamlKey)
	if comp != nil {
		if filePath, err := r.system.WriteToFile(setup.SNIP_CACHE_PATH, s.String(), s.Snip, true); err != nil {
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
			replaceVariables(&run, filePath, s)

			// TODO: ADD TO LOG
			fmt.Println("run", run)
			run = append(run, args...)
			execSafe(run[0], run[1:]...).Close()
		}
	} else {
		//fmt.Println(runner)
		// TODO: MULTI-STEP
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
	out, err := c.StdoutPipe()
	if err != nil {
		panic(err)
	}
	var stderr bytes.Buffer
	c.Stdout = os.Stdout
	c.Stderr = &stderr
	err = c.Run()
	if err != nil {
		panic(err)
	}
	return out
}

func (r *StdRunner) getEnvSection(name string) (*yaml.MapSlice, error) {
	rs, _ := getSubSection(r.setup.Env(), name)
	if rs == nil {
		return nil, errors.New(fmt.Sprintf("No %s section in env.yml", name))
	}
	return &rs, nil
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
		(*cliArgs)[i] = strings.Replace((*cliArgs)[i], "$NAME", strings.Replace(filePath, "."+s.Ext, "", -1), -1)
	}
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