package cmd

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/setup"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"gopkg.in/yaml.v2"
	"strings"
	"os/exec"
	"errors"
	"bytes"
	"time"
	"fmt"
	"io"
	"os"
	"context"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"bufio"
	"os/signal"
	"syscall"
)

type StdRunner struct {
	snippets snippets.Service
	system   sys.Manager
	settings config.Persister
	setup    setup.Provider
}

func NewStdRunner(s sys.Manager, ss snippets.Service, setup setup.Provider) *StdRunner {
	return &StdRunner{snippets: ss, system: s, setup: setup}
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

	log.Debug("EDITING:%v %v", s.Alias, cli)
	err = r.execEdit(s.Alias, cli[0], cli[1:]...)
	if err != nil {
		return err
	}

	rdr := bufio.NewReader(os.Stdin)
	rdr.ReadLine()

	if text, err := r.system.ReadFromFile(setup.SNIP_CACHE_PATH, s.String(), true, 0); err != nil {
		return err
	} else if _, err = r.snippets.Patch(s.Alias, s.Snip, text); err != nil {
		return err
	}
	return nil
}

func (r *StdRunner) Run(s *models.Snippet, args []string) error {
	if !s.VerifyChecksum() {
		return models.ErrOneLine(models.Code_SnippetNotVerified, "The checksum doesn't match the snippet.")
	}
	rs, err := r.getEnvSection("runners")
	if err != nil {
		return err
	}
	yamlKey := s.Ext
	if err != nil {
		return err
	}
	if models.Prefs().Covert {
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
				log.Debug("COMPILE: %s", compile)
				err := r.exec(s.Alias, args, compile[0], compile[1:]...)
				if err != nil {
					return err
				}
			}
			_, run := getSubSection(&comp, "run")
			replaceVariables(&run, filePath, s)

			log.Debug("RUN: %s", run)
			run = append(run, args...)
			err := r.exec(s.Alias, args, run[0], run[1:]...)
			if err != nil {
				return err
			}
		}
	} else {
		if len(interp) > 1 && interp[0] == "echo" && interp[1] == "$SNIP" {
			fmt.Println("Not executable. Printing to screen.")
			fmt.Println(s.Snip)
			return nil
		}
		if s.Ext == "sh" || s.Ext == "bash" {
			// Set unofficial safe-mode
			s.Snip = "set -euo pipefail;\n\n" + s.Snip
		}
		for i, v := range interp {
			interp[i] = strings.Replace(v, "$SNIP", s.Snip, -1)
		}
		interp = append(interp, args...)
		//fmt.Println(runner)

		err := r.exec(s.Alias, args, interp[0], interp[1:]...)
		if err != nil {
			return err
		}
	}
	return nil
}

const PROCESS_NODE = "PROCESS_NODE"

func (r *StdRunner) execEdit(a models.Alias, editor string, arg ...string) error {
	ctx, stop := context.WithTimeout(context.Background(), time.Duration(models.Prefs().CommandTimeout)*time.Second)
	c := exec.CommandContext(ctx, editor, arg...)
	c.Stdin = os.Stdin
	out, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	defer out.Close()
	err = c.Run()
	if err != nil {
		stop()
		_, ok := err.(*exec.ExitError)
		if ok {
			//fmt.Println("Interupted", exErr.UserTime())
		} else {
			desc := fmt.Sprintf("Error editing '%s' (editor: %s)\n\n%s", a.String(), editor, err.Error())
			return models.ErrOneLine(models.Code_RunnerExitError, desc)
		}
		// In this case the error is a snippet runtime error and is treated already above from app errors.
		return nil
	}
	return nil
}

/*
exec realArgs are args that were passed to the snippet, and not the derived args which are passed to the runner.
 */
func (r *StdRunner) exec(a models.Alias, snipArgs []string, runner string, arg ...string) error {
	ctx, stop := context.WithTimeout(context.Background(), time.Duration(models.Prefs().CommandTimeout)*time.Second)
	c := exec.CommandContext(ctx, runner, arg...)
	c.Stdin = os.Stdin
	out, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	defer out.Close()

	var stderr bytes.Buffer
	outBuff := &bytes.Buffer{}
	mw := io.MultiWriter(os.Stdout, outBuff)
	c.Stdout = mw
	c.Stderr = &stderr

	// KEEP TRACK OF PROCESS GRAPH
	node, err := getCurrentNode(a, snipArgs, c)
	if node.Level > models.MaxProcessLevel {
		return models.ErrOneLine(models.Code_ProcessTooDeep,
			"Maximum levels in an app is %d. Not executing:%s", models.MaxProcessLevel, node.AliasString)
	}
	if err != nil {
		return err
	}
	// CHECK FOR COMMON VULNERABILITIES AND ABORT (COULD BE USED WHEN EDITING?)
	err = models.ScanVulnerabilities(strings.Join(arg, " "))
	if err != nil {
		e := err.(*models.ClientErr)
		r.snippets.LogUse(a, models.UseStatusFail, models.UseTypeRun, e.Msgs[0].Desc)
		return err
	}
	// CAPTURE INTERRUPT SO WE CAN LOG PART OF THE EXECUTION IF IS ONGOING e.g. real-time analytics.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)
	go func() {
		res := <-sigChan
		var caller string
		if node.Caller != nil {
			caller = node.Caller.AliasString
		}
		log.Debug("INTERRUPTED: %s|Level:%d|Caller:%s|Message:%s", node.AliasString, node.Level, caller, res.String())
	}()

	err = c.Run()
	if err != nil {
		stop()
		_, ok := err.(*exec.ExitError)
		if ok {
			//fmt.Println("Interupted", exErr.UserTime())
		} else {
			desc := fmt.Sprintf("Error running '%s' (runner: %s %s)\n\n%s", a.String(), runner, err.Error(), stderr.String())
			return models.ErrOneLine(models.Code_RunnerExitError, desc)
		}
		r.snippets.LogUse(a, models.UseStatusFail, models.UseTypeRun, stderr.String())
		// In this case the error is a snippet runtime error and is treated already above from app errors.
		return nil
	}
	r.snippets.LogUse(a, models.UseStatusSuccess, models.UseTypeRun, outBuff.String())
	return nil
}

func (r *StdRunner) getEnvSection(name string) (*yaml.MapSlice, error) {
	rs, _ := getSubSection(models.Env(), name)
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
		(*cliArgs)[i] = strings.Replace((*cliArgs)[i], "$DIR", strings.Replace(filePath, s.String(), "", -1), -1)
		(*cliArgs)[i] = strings.Replace((*cliArgs)[i], "$NAME", strings.Replace(filePath, "."+s.Ext, "", -1), -1)
	}
}

func getSubSection(yml *yaml.MapSlice, name string) (yaml.MapSlice, []string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("The yml config section '%s' is not valid please check it.", name)
		}
	}()
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
