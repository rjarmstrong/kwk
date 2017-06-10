package runtime

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/kwk-super-snippets/cli/src/store"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Editor interface {
	Invoke(s *types.Snippet, onchange func(s types.Snippet)) error
	Close(s *types.Snippet) (uint, error)
}

type EditOptions struct {
	CommandTimeout int64
}

type AppInvoker func(a *types.Alias, app string, args []string, opts EditOptions) error

func NewEditor(env *yaml.MapSlice, prefs *out.Prefs, p SnippetPatcher, f store.SnippetReadWriter) Editor {
	return &editor{env: env, prefs: prefs, patch: p, file: f, inline: inlineInvoker, gui: guiInvoker}
}

type editor struct {
	env    *yaml.MapSlice
	prefs  *out.Prefs
	patch  SnippetPatcher
	file   store.SnippetReadWriter
	inline AppInvoker
	gui    AppInvoker
}

func (ed *editor) Invoke(s *types.Snippet, onchange func(a types.Snippet)) error {
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
	out.Debug("EDIT:%v %v", s.Alias.URI(), edArgs)
	app := edArgs[0]
	opts := EditOptions{CommandTimeout: ed.prefs.CommandTimeout}

	var invoke AppInvoker
	if editInline(app) {
		invoke = ed.inline
	} else {
		invoke = ed.gui
	}
	printEditor(edArgs)
	err = invoke(s.Alias, app, edArgs[1:], opts)
	if err != nil {
		return err
	}
	return nil
}

func printEditor(edArgs []string) {
	// TASK: Move this presentation out of editor
	editor := "default editor"
	for _, v := range edArgs {
		if v == "open" ||
			strings.HasPrefix(v, "-") ||
			strings.HasPrefix(v, "/") ||
			strings.HasPrefix(v, "$") {
			continue
		}
		editor = v
		break
	}
	fmt.Printf("\n%s%s\n\n", style.Margin, style.Fmt16(style.Subdued, "("+editor+")"))
}

func (ed *editor) Close(s *types.Snippet) (uint, error) {
	defer func() {
		if s.Role == types.Role_Preferences || s.Role == types.Role_Environment {
			// We want to keep these to save remote requests.
			return
		}
		err := ed.file.RmDir(s.Alias.URI())
		if err != nil {
			out.LogErr(err)
		}
	}()
	text, err := ed.file.Read(s.Alias.URI())
	if err != nil {
		return 0, err
	}
	out.Debug("EDIT: new:%s, orig: %s, checksum: %s", text, s.Content, s.Checksum)
	if text == s.Content && s.VerifyChecksum() {
		out.Debug("EDIT: Content not changed, not patching")
		return 0, nil
	}
	out.Debug("EDIT: Content changed, patching...")
	res, err := ed.patch(&types.PatchRequest{Alias: s.Alias, Target: s.Content, Patch: text})
	if err != nil {
		return 1, err
	}
	s.Alias.Version = res.Snippet.Alias.Version
	s.Content = res.Snippet.Content
	s.Updated = res.Snippet.Updated
	return 1, nil
}

var guiInvoker = func(a *types.Alias, app string, args []string, opts EditOptions) error {
	err := execEdit(a, app, args, opts)
	if err != nil {
		return err
	}
	rdr := bufio.NewReader(os.Stdin)
	rdr.ReadLine()
	return nil
}

var inlineInvoker = func(a *types.Alias, app string, args []string, opts EditOptions) error {
	done := make(chan bool)
	go func() {
		out.Debug("Editing inline.")
		err := execEdit(a, app, args, opts)
		if err != nil {
			out.LogErrM("EDIT ERR:", err)
		}
		done <- true
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
	a, err := GetSection(env, "apps")
	eRoot, err := GetSection(env, "editors")
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.CommandTimeout)*time.Second)
	defer cancel()
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
