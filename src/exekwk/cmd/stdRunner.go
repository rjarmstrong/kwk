package cmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/kwk-super-snippets/cli/src/exekwk/setup"
	"github.com/kwk-super-snippets/cli/src/gokwk"
	"github.com/kwk-super-snippets/cli/src/models"
	"github.com/kwk-super-snippets/cli/src/persist"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type StdRunner struct {
	snippets gokwk.Snippets
	file     persist.IO
}

func NewStdRunner(s persist.IO, ss gokwk.Snippets) *StdRunner {
	return &StdRunner{snippets: ss, file: s}
}

func (r *StdRunner) Edit(s *types.Snippet) error {
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

	filePath, err := r.file.Write(setup.SNIP_CACHE_PATH, s.String(), s.Snip, true)
	if err != nil {
		return err
	}
	replaceVariables(&cli, filePath, s)
	models.Debug("EDITING:%v %v", s.Alias, cli)
	if err != nil {
		return err
	}
	editor := cli[0]
	cliEditors := map[string]bool{
		"vi":   true,
		"nano": true,
	}
	if cliEditors[editor] {
		done := make(chan bool)
		go func() {
			models.Debug("EDIT asynchronously.")
			err = r.execEdit(s.Alias, editor, cli[1:]...)
			done <- true
			if err != nil {
				models.Debug("Error editing:")
				models.LogErr(err)
			}
		}()
		<-done
	} else {
		err = r.execEdit(s.Alias, editor, cli[1:]...)
		rdr := bufio.NewReader(os.Stdin)
		rdr.ReadLine()
	}

	text, err := r.file.Read(setup.SNIP_CACHE_PATH, s.String(), true, 0)
	if err != nil {
		return err
	}
	_, err = r.snippets.Patch(s.Alias, s.Snip, text)
	if err != nil {
		return err
	}
	return nil
}

func (r *StdRunner) Run(s *types.Snippet, args []string) error {
	if !s.VerifyChecksum() {
		return errs.SnippetNotVerified
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
		if filePath, err := r.file.Write(setup.SNIP_CACHE_PATH, s.String(), s.Snip, true); err != nil {
			return err
		} else {
			_, compile := getSubSection(&comp, "compile")
			if compile != nil {
				replaceVariables(&compile, filePath, s)
				models.Debug("COMPILE: %s", compile)
				err := r.exec(s.Alias, args, compile[0], compile[1:]...)
				if err != nil {
					return err
				}
			}
			_, run := getSubSection(&comp, "run")
			replaceVariables(&run, filePath, s)

			models.Debug("RUN: %s", run)
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

func (r *StdRunner) execEdit(a types.Alias, editor string, arg ...string) error {
	models.Debug("EXEC EDIT: %s %s %s", a.String(), editor, strings.Join(arg, " "))
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(models.Prefs().CommandTimeout)*time.Second)
	c := exec.CommandContext(ctx, editor, arg...)
	c.Stdin = os.Stdin
	var stderr bytes.Buffer
	c.Stdout = os.Stdout
	c.Stderr = &stderr
	err := c.Run()
	if err != nil {
		return err
	}
	//if stderr.Len() > 0 {
	//	desc := fmt.Sprintf("Error editing '%s' (editor: %s)\n\n%s", a.String(), editor, stderr.String())
	//	return nil, models.ErrOneLine(models.Code_RunnerExitError, desc)
	//}
	return nil
}

/*
exec realArgs are args that were passed to the snippet, and not the derived args which are passed to the runner.
*/
func (r *StdRunner) exec(a types.Alias, snipArgs []string, runner string, arg ...string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(models.Prefs().CommandTimeout)*time.Second)
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
	node, err := getCurrentNode(a, runner, snipArgs, c)
	if node.Level > types.MaxProcessLevel {
		return errs.New(errs.CodeErrTooDeep, "Maximum levels in an app is %d. Not executing:%s", types.MaxProcessLevel, node.AliasString)
	}
	if err != nil {
		return err
	}
	// CHECK FOR COMMON VULNERABILITIES AND ABORT (COULD BE USED WHEN EDITING?)
	err = scanVulnerabilities(strings.Join(arg, " "), a.Ext)
	if err != nil {
		e := err.(*errs.Error)
		r.snippets.LogUse(a, types.UseStatusFail, types.UseTypeRun,
			&gokwk.UseContext{
				Preview:     e.Message,
				CallerAlias: node.Caller.AliasString,
				Level:       node.Level,
				Runner:      node.Runner,
			},
		)
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
		models.Debug("INTERRUPTED: %s|Level:%d|Caller:%s|Message:%s", node.AliasString, node.Level, caller, res.String())
	}()
	err = c.Run()
	if c.ProcessState != nil {
		node.Complete(c.ProcessState.Pid())
	}
	if stderr.Len() > 0 {
		desc := fmt.Sprintf("Error running '%s' (runner: '%s' %s)\n\n%s", a.String(), runner, err.Error(), stderr.String())
		r.snippets.LogUse(a, types.UseStatusFail, types.UseTypeRun,
			&gokwk.UseContext{
				Preview:     stderr.String(),
				CallerAlias: node.Caller.AliasString,
				Level:       node.Level,
				Runner:      node.Runner,
			},
		)
		return errs.New(errs.CodeRunnerExit, desc)
	}
	if err != nil {
		exErr, ok := err.(*exec.ExitError)
		if !ok {
			// kwk error return and handle
			return err
		}
		// Was an interrupt
		models.Debug("Interupted:%+v", exErr)
	}
	r.snippets.LogUse(a, types.UseStatusSuccess, types.UseTypeRun,
		&gokwk.UseContext{
			Preview:     outBuff.String(),
			CallerAlias: node.Caller.AliasString,
			Level:       node.Level,
			Runner:      node.Runner,
		},
	)
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
func replaceVariables(cliArgs *[]string, filePath string, s *types.Snippet) {
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

//TODO: Optimise this and add flag to disable it
func scanVulnerabilities(snip string, ext string) error {
	if strings.Contains(snip, "rm -rf") || strings.Contains(snip, "rm ") {
		return errs.New(errs.CodeSnippetVulnerable, "kwk constraint: Shell scripts cannot contain 'rm '.")
	}
	if strings.Contains(snip, ":(){") || strings.Contains(snip, "./$0|./$0&") {
		return errs.New(errs.CodeSnippetVulnerable, "kwk constraint: Fork bomb detected.")
	}
	if strings.Contains(snip, "fork") {
		return errs.New(errs.CodeSnippetVulnerable, "kwk constraint: 'fork' not allowed in script.")
	}
	if strings.Contains(snip, "/dev/sd") {
		return errs.New(errs.CodeSnippetVulnerable, "kwk constraint: '/dev/sd' is not allowed in scripts.")
	}
	if strings.Contains(snip, "/dev/null") {
		return errs.New(errs.CodeSnippetVulnerable, "kwk constraint: '/dev/null' is not allowed in scripts.")
	}
	if strings.Contains(snip, "| sh") || strings.Contains(snip, "| bash") {
		return errs.New(errs.CodeSnippetVulnerable, "kwk constraint: piping directly into terminal not allowed in scripts.")
	}
	if strings.Contains(snip, "nohup") {
		return errs.New(errs.CodeSnippetVulnerable, "kwk constraint: 'nohup' command is not allowed.")
	}
	if (ext == "sh" || ext == "js") && strings.Contains(snip, "eval") {
		m := "kwk constraint: 'eval' command is not allowed."
		if ext == "sh" {
			m += "  Tip: try using '($VAR)' instead of 'eval $VAR' to execute commands.\n"
			m += "  See: /richard/cli/basheval.url\n"
		}
		return errs.New(errs.CodeSnippetVulnerable, m)
	}
	return nil
}
