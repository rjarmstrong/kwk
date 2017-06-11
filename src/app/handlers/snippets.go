package handlers

import (
	"github.com/rjarmstrong/kwk/src/cli"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/runtime"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/age"
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"sort"
	"strings"
	"time"
)

type Snippets struct {
	prefs  *out.Prefs
	client types.SnippetsClient
	runner runtime.Runner
	editor runtime.Editor
	out.Dialog
	vwrite.Writer
	cxf         cli.ContextFunc
	rootPrinter cli.RootPrinter
}

func NewSnippets(p *out.Prefs, s types.SnippetsClient, r runtime.Runner,
	e runtime.Editor, w vwrite.Writer, c cli.ContextFunc, rp cli.RootPrinter, d out.Dialog) *Snippets {

	return &Snippets{
		rootPrinter: rp,
		prefs:       p,
		client:      s,
		runner:      r,
		Dialog:      d,
		Writer:      w,
		editor:      e,
		cxf:         c,
	}
}

// Create a new Snippet
func (sc *Snippets) Create(args []string, pipe bool) error {
	content, alias, err := resolveCreateArgs(args)
	if err != nil {
		return err
	}
	if pipe {
		content = stdInAsString()
	}
	if !alias.IsSnippet() {
		return errs.New(errs.CodeInvalidArgument, "You must provide a pouch when creating a snippet")
	}
	res, err := sc.client.Create(sc.cxf(), &types.CreateRequest{Content: content, Alias: alias, Role: types.Role_Standard})
	if err != nil {
		return err
	}
	err = sc.EWrite(out.SnippetList(sc.prefs, res.List))
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetCreated(res.Snippet))
}

// Search for snippets
func (sc *Snippets) Search(username string, args ...string) error {
	term := strings.Join(args, " ")
	req := &types.AlphaRequest{Term: term, PrivateView: sc.prefs.PrivateView, Username: username}
	res, err := sc.client.Alpha(sc.cxf(), req)
	if err != nil {
		return err
	}
	return sc.EWrite(out.AlphaSearchResult(sc.prefs, res))
}

func (sc *Snippets) ViewListOrRun(uri string, forceView bool, args ...string) error {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	// TASK: Heavy handed, cache preferable
	rr, err := sc.client.GetRoot(sc.cxf(), &types.RootRequest{Username: a.Username, PrivateView: sc.prefs.PrivateView})
	if err != nil {
		return err
	}
	if a.IsUsername() {
		return sc.rootPrinter(rr)
	}
	if a.Ext == "" && rr.IsPouch(a.Name) {
		sc.List(a.Username, a.Name)
		return nil
	}
	s, err := sc.getSnippet(uri)
	if forceView || sc.prefs.RequireRunKeyword {
		return sc.EWrite(out.SnippetView(sc.prefs, s))
	}
	return sc.runner.Run(s, args)
}

// Cat
func (sc *Snippets) Cat(uri string) error {
	snippet, err := sc.getSnippet(uri)
	if err != nil {
		return err
	}
	if snippet != nil {
		return sc.EWrite(out.SnippetCat(snippet))
	}
	return errs.NotFound
}

// Run a snippet
func (sc *Snippets) Run(uri string, args []string) error {
	snippet, err := sc.getSnippet(uri)
	if err != nil {
		return err
	}
	return sc.runner.Run(snippet, args)
}

// RunNode is a call to run a snippet from within an app (i.e. as a new process)
func (sc *Snippets) RunNode(pr cli.UserWithToken, prefs *out.Prefs, node *runtime.ProcessNode, uri string, args []string) error {
	out.Debug("RUN:%s %s", uri, args)
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	list, err := sc.client.Get(sc.cxf(), &types.GetRequest{Alias: a, Suggest: false})
	if err != nil {
		return err
	}
	if len(list.Items) > 1 {
		return sc.EWrite(out.SnippetAmbiguous(node.URI, uri))
	}
	if !prefs.RunAllSnippets && list.Items[0].Username() != pr.User.Username {
		return sc.EWrite(out.RunAllSnippetsNotTrue(node.URI, uri))
	}
	return sc.runner.Run(list.Items[0], args)
}

// Edit Snippet
func (sc *Snippets) Edit(uri string) error {
	// TASK: test
	edit := func(s *types.Snippet) error {
		sc.Write(out.SnippetEditing(s))
		err := sc.editor.Invoke(s, func(s types.Snippet) {})
		if err != nil {
			return err
		}
		changes, err := sc.editor.Close(s)
		if err != nil {
			return err
		}
		if changes > 0 {
			return sc.EWrite(out.SnippetEdited(s))
		}
		return sc.EWrite(out.SnippetNoChanges(s))
	}
	snippet, err := sc.getSnippet(uri)
	if err != nil {
		if errs.HasCode(err, errs.CodeNotFound) {
			a, err := types.ParseAlias(uri)
			if err != nil {
				return err
			}
			r := sc.Dialog.Modal(out.SnippetEditNewPrompt(a.String()), false)
			if r.Ok {
				r, err := sc.client.Create(sc.cxf(), &types.CreateRequest{Alias: a, Role: types.Role_Standard})
				if err != nil {
					return err
				}
				return edit(r.Snippet)
			}
		}
		return err
	}
	return edit(snippet)
}

// Describe
func (sc *Snippets) Describe(uri string, description string) error {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	alias, err := sc.client.Update(sc.cxf(), &types.UpdateRequest{Alias: a, Description: description})
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetDescriptionUpdated(alias.String(), description))
}

// Delete
func (sc *Snippets) Delete(args []string) error {
	// TASK: Use a lighter-weight method to get all pouches
	r, err := sc.client.GetRoot(sc.cxf(), &types.RootRequest{PrivateView: sc.prefs.PrivateView})
	if err != nil {
		return err
	}
	if r.IsPouch(args[0]) {
		return sc.deletePouch(args[0])
	}
	return sc.deleteSnippet(args)
}

// Move
func (sc *Snippets) Move(args []string) error {
	// Task: Add tests refactor into parts
	if len(args) < 2 {
		return errs.TwoArgumentsReqForMove
	}
	// Task: Make lighter weight
	root, err := sc.client.GetRoot(sc.cxf(), &types.RootRequest{PrivateView: sc.prefs.PrivateView})
	if err != nil {
		return err
	}
	last := args[len(args)-1]
	// If first argument is pouch is a pouch rename
	if root.IsPouch(args[0]) {
		res, err := sc.client.RenamePouch(sc.cxf(), &types.RenamePouchRequest{Name: args[0], NewName: args[1]})
		if err != nil {
			return err
		}
		return sc.rootPrinter(res.Root)
	} else if !root.IsPouch(last) && len(args) == 2 {
		// rename single snippet
		snip, original, err := sc.renameSnippet(args[0], args[1])
		if err != nil {
			return err
		}
		sc.List("", snip.Pouch())
		return sc.EWrite(out.SnippetRenamed(original.String(), snip.String()))
	}
	as, source, err := types.ParseMany(args[0 : len(args)-1])
	if err != nil {
		return err
	}
	// move snippets into a pouch
	p, err := sc.client.Move(sc.cxf(), &types.MoveRequest{SourcePouch: source, TargetPouch: last, SnipNames: as})
	if err != nil {
		return err
	}
	sc.List("", last)
	return sc.EWrite(out.SnippetsMoved(as, p.Pouch))
}

// Patch
func (sc *Snippets) Patch(uri string, target string, patch string) error {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	alias, err := sc.client.Patch(sc.cxf(), &types.PatchRequest{Alias: a, Target: target, Patch: patch})
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetPatched(alias.String()))
}

// Clone
func (sc *Snippets) Clone(uri string, newFullName string) error {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	newA, err := types.ParseAlias(newFullName)
	if err != nil {
		return err
	}
	res, err := sc.client.Clone(sc.cxf(), &types.CloneRequest{Alias: a, New: newA})
	if err != nil {
		return err
	}
	err = sc.EWrite(out.SnippetList(sc.prefs, res.List))
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetClonedAs(res.Snippet.Alias.URI()))
}

// Tag
func (sc *Snippets) Tag(uri string, tags ...string) error {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	alias, err := sc.client.Tag(sc.cxf(), &types.TagRequest{Alias: a, Tags: tags})
	if err != nil {
		return err
	}
	return sc.EWrite(out.Tagged(alias.String(), tags))
}

// Untag
func (sc *Snippets) UnTag(uri string, tags ...string) error {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	alias, err := sc.client.UnTag(sc.cxf(), &types.UnTagRequest{Alias: a, Tags: tags})
	if err != nil {
		return err
	}
	return sc.EWrite(out.UnTag(alias.String(), tags))
}

// Dump writes out all snippets as one long list
func (sc *Snippets) Dump(username string) error {
	list, err := sc.client.List(sc.cxf(), &types.ListRequest{Username: username, PrivateView: sc.prefs.PrivateView})
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetList(sc.prefs, list))
}

// GetEra lists snippets by special filters: @today @week @month @old
func (sc *Snippets) GetEra(virtualPouch string) error {
	// Task: Not working!!
	list, err := sc.client.List(sc.cxf(), &types.ListRequest{PrivateView: sc.prefs.PrivateView})
	if err != nil {
		return err
	}
	era := []*types.Snippet{}
	var since, latest int64
	sod := age.StartOfDay(time.Now()).Unix()
	isoYear, isoWeek := time.Now().ISOWeek()
	fdw := age.FirstDayOfISOWeek(isoYear, isoWeek, time.Local).Unix()
	som := age.StartOfMonth(time.Now()).Unix()
	if virtualPouch == "@today" {
		since = sod
		latest = time.Now().Unix()
	} else if virtualPouch == "@week" {
		since = fdw
		latest = sod
	} else if virtualPouch == "@month" {
		since = som
		latest = fdw
	} else if virtualPouch == "@old" {
		since = 0
		latest = som
	}
	for _, v := range list.Items {
		if v.RunStatusTime > since && v.RunStatusTime < latest {
			era = append(era, v)
		}
	}
	sort.Slice(era, func(i, j int) bool {
		return era[i].RunStatusTime < era[j].RunStatusTime
	})
	list.Items = era
	return sc.EWrite(out.SnippetList(sc.prefs, list))
}

func (sc *Snippets) List(username string, pouch string) error {
	if pouch == "" {
		r, err := sc.client.GetRoot(sc.cxf(),
			&types.RootRequest{Username: username, PrivateView: sc.prefs.PrivateView})
		if err != nil {
			return err
		}
		return sc.rootPrinter(r)
	}
	var size int64
	list, err := sc.client.List(sc.cxf(),
		&types.ListRequest{Username: username, Pouch: pouch, Limit: size, PrivateView: sc.prefs.PrivateView})
	if err != nil {
		return err
	}
	a := types.NewAlias(list.Username, list.Pouch.Name, "", "")
	_, err = sc.client.LogUse(sc.cxf(),
		&types.UseContext{
			Alias:  a,
			Type:   types.UseType_View,
			Status: types.UseStatus_Success,
			Time:   types.KwkTime(time.Now()),
		})
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetList(sc.prefs, list))
}

func (sc *Snippets) getSnippet(uri string) (*types.Snippet, error) {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return nil, err
	}
	res, err := sc.client.Get(sc.cxf(), &types.GetRequest{Alias: a, Suggest: true})
	if res.Suggested && len(res.Items) == 1 {
		mres := sc.Modal(out.DidYouMean(res.Items[0].Alias.URI()), false)
		if mres.Ok {
			return res.Items[0], nil
		}
		return nil, nil
	}
	return sc.ChooseSnippet(res.Items), nil
}

func (sc *Snippets) renameSnippet(uri string, newSnipName string) (*types.Snippet, *types.SnipName, error) {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return nil, nil, err
	}
	sn, err := types.ParseSnipName(newSnipName)
	if err != nil {
		return nil, nil, err
	}
	res, err := sc.client.Rename(sc.cxf(), &types.RenameRequest{Alias: a, NewName: sn})
	return res.Snippet, res.Original, nil
}

func (sc *Snippets) deleteSnippet(args []string) error {
	// TASK: Awkward interface as requires pouched snippet format
	sn, pouch, err := types.ParseMany(args)
	if err != nil {
		return err
	}
	r := sc.Modal(out.SnippetCheckDelete(sn), sc.prefs.AutoYes)
	if !r.Ok {
		sc.EWrite(out.SnippetsNotDeleted(sn))
	}
	res, err := sc.client.Delete(sc.cxf(), &types.DeleteRequest{Pouch: pouch, Names: sn})
	if err != nil {
		return err
	}
	err = sc.EWrite(out.SnippetList(sc.prefs, res.List))
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetsDeleted(sn))
}
