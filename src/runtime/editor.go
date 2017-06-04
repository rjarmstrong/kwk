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

type EditOptions struct {
	CommandTimeout int64
}

type EditFunc func(a *types.Alias, app string, args []string, opts EditOptions) error

func NewEditor(env *yaml.MapSlice, prefs *out.Prefs, p SnippetPatcher, f store.SnippetReadWriter) Editor {
	return &editor{env : env, prefs: prefs, patch: p, file : f , inlineFunc: inlineRunner, guiFunc: guiRunner}
}

type editor struct {
	env        *yaml.MapSlice
	prefs      *out.Prefs
	patch      SnippetPatcher
	file       store.SnippetReadWriter
	inlineFunc EditFunc
	guiFunc    EditFunc
}

func (ed *editor) Edit(s *types.Snippet) error {
	filePath, err := ed.file.Write(s.Alias.URI(), s.Content)
	out.Debug("SAVED TO: %s", filePath)
	if err != nil {
		return err
	}
	edArgs, err := getEditArgs(ed.env, s.Ext())
	if err != nil {
		return err
	}
	replaceVariables(edArgs, filePath, s)
	out.Debug("EDITING:%v %v", s.Alias.URI(), edArgs)
	app := edArgs[0]
	opts := EditOptions{CommandTimeout: ed.prefs.CommandTimeout}

	if editInline(app) {
		ed.inlineFunc(s.Alias, app, edArgs[1:], opts)
	} else {
		ed.guiFunc(s.Alias, app, edArgs[1:], opts)
	}

	text, err := ed.file.Read(s.Alias.URI())
	if err != nil {
		return err
	}
	if s.Content == text {
		out.FreeText("File unchanged.")
		return nil
	}
	_, err = ed.patch(&types.PatchRequest{Alias: s.Alias, Target: s.Content, Patch: text})
	if err != nil {
		return err
	}
	return nil
}

var guiRunner = func (a *types.Alias, app string, args []string, opts EditOptions) error {
	err := execEdit(a, app, args, opts)
	if err != nil {
		return err
	}
	rdr := bufio.NewReader(os.Stdin)
	rdr.ReadLine()
	return nil
}

var inlineRunner = func (a *types.Alias, app string, args []string, opts EditOptions) error {
	done := make(chan bool)
	go func() {
		out.Debug("Editing inline.")
		err := execEdit(a, app, args, opts)
		done <- true
		if err != nil {
			out.Debug("Error editing:")
			out.LogErr(err)
		}
	}()
	<-done
	// TASK: Write golang func error back to current proc
	return nil
}

func editInline(editor string) bool {
	return map[string]bool{
		"vi":   true,
		"nano": true,
	}[editor]
}

func getEditArgs(env *yaml.MapSlice, ext string) ([]string, error) {
	a, err := GetSection(env,"apps")
	eRoot, err := GetSection(env,"editors")
	if err != nil {
		return nil, err
	}
	_, candidates := getSubSection(eRoot, ext)
	if len(candidates) != 1 {
		return nil,
			errs.New(0, "No editors have been specified for  %s. And default editor is not specified.", ext)
	}
	_, ed := getSubSection(a, candidates[0])
	return ed, nil
}

func execEdit(a *types.Alias, app string, arg []string, opts EditOptions) error {
	out.Debug("EXEC EDIT: %s %s %s", a.URI(), app, strings.Join(arg, " "))
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(opts.CommandTimeout)*time.Second)
	c := exec.CommandContext(ctx, app, arg...)
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
