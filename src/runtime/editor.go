package runtime

import (
	"bufio"
	"bytes"
	"context"
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/kwk-super-snippets/cli/src/store"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"os"
	"os/exec"
	"strings"
	"time"
	"gopkg.in/yaml.v2"
)

type Editor interface {
	Edit(s *types.Snippet) error
}

func NewEditor(env *yaml.MapSlice, prefs *out.Prefs, p SnippetPatcher, f store.SnippetReadWriter) Editor {
	return &editor{env : env, prefs: prefs, patch: p, file : f }
}

type editor struct {
	env *yaml.MapSlice
	prefs *out.Prefs
	patch SnippetPatcher
	file  store.SnippetReadWriter
}

func (r *editor) Edit(s *types.Snippet) error {
	//TODO: if we pull out the env from getSection we can improve speed
	a, err := GetSection(r.env,"apps")
	eRoot, err := GetSection(r.env,"editors")

	if err != nil {
		return err
	}
	_, candidates := getSubSection(eRoot, s.Ext())
	if len(candidates) != 1 {
		return errs.New(0, "No editors have been specified for  %s. And default editor is not specified.", s.Ext())
	}
	_, cli := getSubSection(a, candidates[0])

	filePath, err := r.file.Write(s.Alias.URI(), s.Content)
	out.Debug("SAVED TO: %s", filePath)
	if err != nil {
		return err
	}
	replaceVariables(&cli, filePath, s)
	out.Debug("EDITING:%v %v", s.Alias.URI(), cli)
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
			out.Debug("EDIT asynchronously.")
			err = r.execEdit(s.Alias, editor, cli[1:]...)
			done <- true
			if err != nil {
				out.Debug("Error editing:")
				out.LogErr(err)
			}
		}()
		<-done
	} else {
		err = r.execEdit(s.Alias, editor, cli[1:]...)
		rdr := bufio.NewReader(os.Stdin)
		rdr.ReadLine()
	}

	text, err := r.file.Read(s.Alias.URI())
	if err != nil {
		return err
	}
	if s.Content == text {
		out.FreeText("File unchanged.")
		return nil
	}
	_, err = r.patch(&types.PatchRequest{Alias: s.Alias, Target: s.Content, Patch: text})
	if err != nil {
		return err
	}
	return nil
}

func (r *editor) execEdit(a *types.Alias, editor string, arg ...string) error {
	out.Debug("EXEC EDIT: %s %s %s", a.URI(), editor, strings.Join(arg, " "))
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(r.prefs.CommandTimeout)*time.Second)
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
	//	return nil, ErrOneLine(Code_RunnerExitError, desc)
	//}
	return nil
}
