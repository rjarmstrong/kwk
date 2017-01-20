package cmd

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/models"
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
	conf models.ConfigStore
}

func NewStdRunner(s sys.Manager, ss snippets.Service, conf models.ConfigStore) *StdRunner {
	return &StdRunner{snippets: ss, system: s,  conf: conf}
}

func (r *StdRunner) Edit(s *models.Snippet) error {
	//TODO: if we pull out the env from getSection we can improve speed
	a, err := r.getEnvSection("apps")
	eRoot, err := r.getEnvSection("editors")

	if err != nil {
		return err
	}
	_, candidates := r.conf.GetSubSection(eRoot, s.Extension)
	if len(candidates) != 1 {
		return errors.New("No editors have been specified for " + s.Extension + " . And default editor is not specified.")
	}
	_, cli := r.conf.GetSubSection(a, candidates[0])

	filePath, err := r.system.WriteToFile(models.SNIP_CACHE_PATH, s.FullName, s.Snip, true)
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

	if text, err := r.system.ReadFromFile(models.SNIP_CACHE_PATH, s.FullName, true, 0); err != nil {
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


func (r *StdRunner) Run(s *models.Snippet, args []string) error {
	rs, err := r.getEnvSection("runners")
	if err != nil {
		return err
	}
	yamlKey := s.Extension
	if r.conf.Prefs().Covert {
		yamlKey += "-covert"
	}
	comp, interp := r.conf.GetSubSection(rs, yamlKey)
	if comp != nil {
		if filePath, err := r.system.WriteToFile(models.SNIP_CACHE_PATH, s.FullName, s.Snip, true); err != nil {
			return err
		} else {
			_, compile := r.conf.GetSubSection(&comp, "compile")
			if compile != nil {
				replaceVariables(&compile, filePath, s)
				// TODO: ADD TO LOG
				fmt.Println("compile", compile)
				execSafe(compile[0], compile[1:]...).Close()
			}
			_, run := r.conf.GetSubSection(&comp, "run")
			replaceVariables(&run, filePath, s)

			// TODO: ADD TO LOG
			fmt.Println("run", run)
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

func (r *StdRunner) getEnvSection(name string) (*yaml.MapSlice, error) {
	rs, _ := r.conf.GetSubSection(r.conf.Env(), name)
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
		(*cliArgs)[i] = strings.Replace((*cliArgs)[i], "$NAME", strings.Replace(filePath, "."+s.Extension, "", -1), -1)
	}
}