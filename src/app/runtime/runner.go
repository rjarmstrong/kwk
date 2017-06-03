package runtime

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/cli/src/store"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/vwrite"
	"github.com/lunixbochs/vtclean"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
	rt "runtime"
)

type Runner interface {
	Run(s *types.Snippet, args []string) error
}

type runner struct {
	prefs *out.Prefs
	env   *yaml.MapSlice
	file  store.SnippetReadWriter
	w     vwrite.Writer
	ul    UseLogger
}

func NewRunner(prefs *out.Prefs, env *yaml.MapSlice, w vwrite.Writer, f store.SnippetReadWriter, ul UseLogger) Runner {
	return &runner{prefs:prefs, env: env, file: f, w: w, ul: ul}
}

func (r *runner) Run(s *types.Snippet, args []string) error {
	if !s.VerifyChecksum() {
		return errs.SnippetNotVerified
	}
	rs, err := GetSection(r.env, "runners")
	if err != nil {
		return err
	}
	yamlKey := s.Ext()
	if err != nil {
		return err
	}
	if r.prefs.Covert {
		yamlKey += "-covert"
	}
	comp, interp := getSubSection(rs, yamlKey)
	if comp != nil {
		if filePath, err := r.file.Write(s.Alias.URI(), s.Content); err != nil {
			return err
		} else {
			_, compile := getSubSection(&comp, "compile")
			if compile != nil {
				replaceVariables(&compile, filePath, s)
				out.Debug("COMPILE: %s", compile)
				err := r.exec(s.Alias, args, compile[0], compile[1:]...)
				if err != nil {
					return err
				}
			}
			_, run := getSubSection(&comp, "run")
			replaceVariables(&run, filePath, s)

			out.Debug("RUN: %s", run)
			run = append(run, args...)
			err := r.exec(s.Alias, args, run[0], run[1:]...)
			if err != nil {
				return err
			}
		}
	} else {
		if len(interp) > 1 && interp[0] == "echo" && interp[1] == "$SNIP" {
			err := r.w.EWrite(out.NotExecutable(s))
			r.logUse(s.Alias, "", NewProcessNode(*s.Alias, "", args, nil), types.UseStatus_Success)
			return err
		}
		if s.Ext() == "sh" || s.Ext() == "bash" {
			// Set unofficial safe-mode
			s.Content = "set -euo pipefail;\n\n" + s.Content
		}
		for i, v := range interp {
			interp[i] = strings.Replace(v, "$SNIP", s.Content, -1)
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

/*
exec realArgs are args that were passed to the snippet, and not the derived args which are passed to the runner.
*/
func (r *runner) exec(a *types.Alias, snipArgs []string, runner string, arg ...string) error {
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.prefs.CommandTimeout)*time.Second)
	c := exec.CommandContext(ctx, runner, arg...)
	c.Stdin = os.Stdin
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	defer stdout.Close()

	var stderr bytes.Buffer
	outBuff := &bytes.Buffer{}
	mw := io.MultiWriter(os.Stdout, outBuff)
	c.Stdout = mw
	c.Stderr = &stderr

	// KEEP TRACK OF PROCESS GRAPH
	node, err := getCurrentNode(*a, runner, snipArgs, c)
	if node.Level > types.MaxProcessLevel {
		return errs.New(errs.CodeErrTooDeep, "Maximum levels in an app is %d. Not executing: %s", types.MaxProcessLevel, node.URI)
	}
	if err != nil {
		return err
	}
	// CHECK FOR COMMON VULNERABILITIES AND ABORT (COULD BE USED WHEN EDITING?)
	err = scanVulnerabilities(strings.Join(arg, " "), a.Ext)
	if err != nil {
		e := err.(*errs.Error)
		r.logUse(a, e.Message, node, types.UseStatus_Fail)
		return err
	}
	// CAPTURE INTERRUPT SO WE CAN LOG PART OF THE EXECUTION IF IS ONGOING e.g. real-time analytics.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT)
	go func() {
		res := <-sigChan
		var caller string
		if node.Caller != nil {
			caller = node.Caller.URI
		}
		out.Debug("INTERRUPTED: %s|Level:%d|Caller:%s|Message:%s", node.URI, node.Level, caller, res.String())
	}()
	err = c.Run()
	if c.ProcessState != nil {
		node.Complete(c.ProcessState.Pid())
	}
	if stderr.Len() > 0 {
		var desc string
		if err != nil {
			desc = fmt.Sprintf("Error running '%s' (runner: '%s' %s)\n\n%s", a.URI(), runner, err.Error(), stderr.String())
		} else {
			desc = fmt.Sprintf("Error running '%s' (runner: '%s')\n\n%s", a.URI(), runner, stderr.String())
		}
		out.Debug("RUNNER:%s", desc)
		r.logUse(a, stderr.String(), node, types.UseStatus_Fail)
		return errs.New(errs.CodeRunnerExit, desc)
	}
	if err != nil {
		exErr, ok := err.(*exec.ExitError)
		if !ok {
			// kwk error return and handle
			return err
		}
		// Was an interrupt
		out.Debug("Interupted:%+v", exErr)
	}
	r.logUse(a, outBuff.String(), node, types.UseStatus_Success)
	return nil
}

/*
 Limits a preview adding an ascii escape at the end and fixing the length.
*/
func limitPreview(in string, length int) string {
	in = vtclean.Clean(in, true)
	return types.Limit(in, length-5) + "\033[0m"
}

func (r *runner) logUse(a *types.Alias, output string, node *ProcessNode, s types.UseStatus) {
	r.ul(&types.UseContext{
		Alias:       a,
		Type:        types.UseType_Run,
		Status:      s,
		Preview:     limitPreview(output, types.PreviewMaxRuneLength),
		CallerAlias: node.Caller.URI,
		Level:       node.Level,
		Runner:      node.Runner,
		Os:          rt.GOOS,
		Time:        types.KwkTime(time.Now()),
	})
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
		(*cliArgs)[i] = strings.Replace((*cliArgs)[i], "$DIR", strings.Replace(filePath, s.Alias.URI(), "", -1), -1)
		(*cliArgs)[i] = strings.Replace((*cliArgs)[i], "$NAME", strings.Replace(filePath, "."+s.Ext(), "", -1), -1)
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
