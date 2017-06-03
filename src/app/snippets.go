package app

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/cli/src/app/runtime"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/age"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/vwrite"
	"os"
	"sort"
	"strings"
	"time"
)

type snippets struct {
	s      types.SnippetsClient
	runner runtime.Runner
	editor runtime.Editor
	Dialog
	vwrite.Writer
}

func NewSnippets(s types.SnippetsClient, r runtime.Runner, e runtime.Editor, d Dialog, w vwrite.Writer) *snippets {
	return &snippets{s: s, runner: r, Dialog: d, Writer: w, editor: e}
}

func (sc *snippets) Search(args ...string) error {
	term := strings.Join(args, " ")
	req := &types.AlphaRequest{Term: term, PrivateView: prefs.PrivateView}
	if !prefs.PrivateView {
		req.Username = principal.User.Username
	}
	res, err := sc.s.Alpha(Ctx(), req)
	if err != nil {
		return err
	}
	return sc.EWrite(out.AlphaSearchResult(prefs, res))
}

//func (sc *SnippetCli) ListCategory(category string, args ...string) {
//	username, pouch, err := types.ParsePouch(args[0])
//	if err != nil {
//		sc.HandleErr(err)
//		return
//	}
//	var size int64
//var tags = []string{}
//for i, v := range args {
//	if num, err := strconv.Atoi(v); err == nil {
//		size = int64(num)
//	} else {
//		if i == 0 && v[len(v)-1:] == "/" {
//			username = strings.Replace(v, "/", "", -1)
//		} else {
//			tags = append(tags, v)
//		}
//	}
//}
//	p := &models.ListParams{Category: category, Username: username, Pouch: pouch, Size: size, All: models.Prefs().ListAll}
//	if list, err := sc.s.List(p); err != nil {
//		sc.HandleErr(err)
//	} else {
//		sc.Render("snippet:list", list)
//	}
//}

//func (sc *snippets) Share(distinctName string, destination string) error {
//	list, _, err := sc.getSnippet(distinctName)
//	if err != nil {
//		return err
//	}
//	alias := sc.handleMultiResponse(distinctName, list)
//	if alias != nil {
//		snip := "https://mail.google.com/mail/?ui=2&view=cm&fs=1&tf=1&su=&body=http%3A%2F%2Faus.kwk.co%2F" + alias.String()
//		gmail := gokwk.NewSnippet(snip)
//		gmail.Ext = "url"
//		sc.runner.Run(gmail, []string{})
//	}
//	return sc.EWrite(out.) sc.Render("snippet:notfound", map[string]interface{}{"fullKey": distinctName})
//}

func (sc *snippets) run(selected *types.Snippet, args []string) error {
	return sc.runner.Run(selected, args)
}

func (sc *snippets) NodeRun(uri string, args []string) error {
	out.Debug("RUN:%s %s", uri, args)
	a, err := types.ParseAlias(uri)
	// TASK: Add config to allow this check?
	//if a.Version < 1 {
	//	return sc.EWrite(out.VersionMissing(a.URI()))
	//}
	if err != nil {
		return err
	}
	list, err := sc.s.Get(Ctx(), &types.GetRequest{Alias: a})
	if err != nil {
		if errs.HasCode(err, errs.CodeNoSnipNamesFound) {
			return sc.EWrite(out.NotFoundInApp(node.URI, uri))
		}
		return err
	}
	s := sc.handleMultiResponse(uri, list.Items)
	if s == nil {
		return sc.EWrite(out.NotFoundInApp(node.URI, uri))
	}
	//TODO: If username is not the current user or 'kwk' then prompt before executing.
	return sc.runner.Run(s, args)
}

func (sc *snippets) Run(uri string, args []string) error {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	list, err := sc.s.Get(Ctx(), &types.GetRequest{Alias: a})
	if err != nil {
		sc.typeAhead(uri, sc.run)
	}
	alias := sc.handleMultiResponse(uri, list.Items)
	if alias != nil {
		return sc.runner.Run(alias, args)
	}
	return sc.typeAhead(uri, sc.run)
}

func GuessArgs(a string, b string) (*types.Alias, string, error) {
	aIsNotSUri := types.IsDefNotPouchedSnippetURI(a)
	bIsNotSUri := types.IsDefNotPouchedSnippetURI(b)

	out.Debug("ARG1:%v (snip=%t) ARG2:%v (snip=%t)", a, aIsNotSUri, b, bIsNotSUri)

	if aIsNotSUri && bIsNotSUri {
		return nil, "", errs.New(errs.CodeInvalidArgument,
			"Please specify a pouch to create a snippet."+
				"\n   e.g."+
				"\n  `kwk new <pouch_name>/<snip_name>[.<ext>]  <snippet>`"+
				"\n  `kwk new <snippet> <pouch_name>/<snip_name>[.<ext>]`"+
				"\n  `<cmd> | kwk new <pouch_name>/<snip_name>[.<ext>]`",
		)
	}
	if !aIsNotSUri && !bIsNotSUri {
		return nil, "", errs.New(errs.CodeInvalidArgument, "It looks like both arguments could be either be a path or kwk URIs, please add an extension or fully quality the kwk URI. e.g. kwk.co/richard/dill/name.path")
	}
	if !aIsNotSUri {
		alias, err := types.ParseAlias(a)
		if err != nil {
			return nil, "", err
		}
		return alias, b, nil

	}
	alias, err := types.ParseAlias(b)
	if err != nil {
		return nil, "", err
	}
	return alias, a, nil
}

func (sc *snippets) Create(args []string, pipe bool) error {
	a1 := &types.Alias{}
	var snippet string
	if len(args) == 1 {
		if types.IsDefNotPouchedSnippetURI(args[0]) {
			out.Debug("Assuming the only arg is the snippet.")
			snippet = args[0]
			a1 = &types.Alias{}
		} else {
			a, err := types.ParseAlias(args[0])
			if err != nil {
				return err
			}
			a1 = a
		}
	} else if len(args) > 1 {
		a, s, err := GuessArgs(args[0], args[1])
		if err != nil {
			return err
		}
		a1 = a
		snippet = s
	}
	if pipe {
		snippet = stdInAsString()
	}
	res, err := sc.s.Create(Ctx(), &types.CreateRequest{Content: snippet, Alias: a1, Role: types.Role_Standard})
	if err != nil {
		return err
		//TODO: If snippet is similar to an existing one prompt for it here.
	}
	sc.List("", types.PouchRoot)
	return sc.EWrite(out.SnippetCreated(res.Snippet))
}

func (sc *snippets) Edit(uri string) error {
	innerEdit := func(s *types.Snippet) error {
		sc.Write(out.SnippetEditing(s))
		err := sc.editor.Edit(s)
		if err != nil {
			return err
		}
		return sc.EWrite(out.SnippetEdited(s))
	}
	if uri == "env" {
		uri = runtime.GetEnvURI()
	}
	list, _, err := sc.getSnippet(uri)
	if err != nil {
		if errs.HasCode(err, errs.CodeNotFound) {
			a, err := types.ParseAlias(uri)
			if err != nil {
				return err
			}
			r := sc.Dialog.Modal(out.SnippetEditNewPrompt(a.String()), false)
			if r.Ok {
				r, err := sc.s.Create(Ctx(), &types.CreateRequest{Alias: a, Role: types.Role_Standard})
				if err != nil {
					return err
				}
				return innerEdit(r.Snippet)
			}
		}
		return err
	}
	snippet := sc.handleMultiResponse(uri, list.Items)
	if snippet != nil {
		return innerEdit(snippet)
	}
	return errs.NotFound
}

func (sc *snippets) Describe(distinctName string, description string) error {
	a, err := types.ParseAlias(distinctName)
	if description == "" {
		return sc.typeAhead(distinctName, func(s *types.Snippet, args []string) error {
			res, err := sc.FormField(
				out.FreeText(fmt.Sprintf("Enter new description for %s: ", s.Alias.FileName())), false)
			if err != nil {
				return err
			}
			out.Debug("Form result: %+v", res.Value)
			return sc.Describe(s.Alias.FileName(), res.Value.(string))
		})
	}
	if err != nil {
		return err
	}
	alias, err := sc.s.Update(Ctx(), &types.UpdateRequest{Alias: a, Description: description})
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetDescriptionUpdated(alias.String(), description))
}

func (sc *snippets) InspectListOrRun(uri string, forceView bool, args ...string) error {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	// TASK: Heavy handed, cache preferable
	rr, err := sc.s.GetRoot(Ctx(), &types.RootRequest{Username: a.Username, PrivateView: prefs.PrivateView})
	if err != nil {
		return err
	}
	if a.Ext == "" && rr.IsPouch(a.Name) {
		p := rr.GetPouch(a.Name)
		if p.Type == types.PouchType_Virtual {
			fmt.Println("List virtual.")
		} else {
			sc.List(a.Username, a.Name)
		}
		return nil
	}
	if err != nil {
		return err
	}

	// GET SNIPPET
	list, err := sc.s.Get(Ctx(), &types.GetRequest{Alias: a})
	if err != nil {
		return sc.typeAhead(uri, sc.run)
	}
	s := sc.handleMultiResponse(uri, list.Items)
	if s == nil {
		return out.LogErrM("Snippet not found", errs.NotFound)
	}
	if forceView || prefs.RequireRunKeyword {
		out.Debug("RUN KEYWORD REQUIRED, VIEWING.")
		sc.Write(out.SnippetView(prefs, s))
		return nil
	}
	return sc.runner.Run(s, args)
}

func (sc *snippets) Delete(args []string) error {
	r, err := sc.s.GetRoot(Ctx(), &types.RootRequest{PrivateView: prefs.PrivateView})
	if err != nil {
		return err
	}
	if r.IsPouch(args[0]) {
		return sc.deletePouch(args[0])
	}
	return sc.deleteSnippet(args)
}

func (sc *snippets) Lock(pouch string) error {
	_, err := sc.s.MakePouchPrivate(Ctx(), &types.MakePrivateRequest{Name: pouch, MakePrivate: true})
	if err != nil {
		return err
	}
	return sc.EWrite(out.PouchLocked(pouch))
}

func (sc *snippets) UnLock(pouch string) error {
	res := sc.Dialog.Modal(out.PouchCheckUnLock(pouch), false)
	if res.Ok {
		_, err := sc.s.MakePouchPrivate(Ctx(), &types.MakePrivateRequest{Name: pouch, MakePrivate: false})
		if err != nil {
			return err
		}
		return sc.EWrite(out.PouchUnLocked(pouch))
	}
	return sc.EWrite(out.PouchNotUnLocked(pouch))
}

// kwk mv regions.txt reference -- moves the reference pouch, if no reference pouch then move to reference.txt
// kwk mv examples/regions.txt reference
func (sc *snippets) Move(args []string) error {
	if len(args) < 2 {
		return errs.TwoArgumentsReqForMove
	}
	root, err := sc.s.GetRoot(Ctx(), &types.RootRequest{PrivateView: prefs.PrivateView})
	if err != nil {
		return err
	}
	last := args[len(args)-1]
	// If first argument is pouch is a pouch rename
	if root.IsPouch(args[0]) {
		res, err := sc.s.RenamePouch(Ctx(), &types.RenamePouchRequest{Name: args[0], NewName: args[1]})
		if err != nil {
			return err
		}
		return sc.EWrite(out.PrintRoot(prefs, &cliInfo, res.Root, &principal.User))
	} else if !root.IsPouch(last) && len(args) == 2 {
		// rename single snippet
		snip, original, err := sc.renameSnippet(args[0], args[1])
		if err != nil {
			return err
		}
		sc.List("", snip.Pouch())
		return sc.EWrite(out.SnippetRenamed(original.String(), snip.String()))
		return nil
	}
	as, source, err := types.ParseMany(args[0 : len(args)-1])
	if err != nil {
		return err
	}
	// move snippets into a pouch
	p, err := sc.s.Move(Ctx(), &types.MoveRequest{SourcePouch: source, TargetPouch: last, SnipNames: as})
	if err != nil {
		return err
	}
	sc.List("", last)
	return sc.EWrite(out.SnippetsMoved(as, p.Pouch))
}

type MoveResult struct {
	Pouch string
	Quant int
}

func (sc *snippets) Cat(distinctName string) error {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return err
	}
	if !a.IsSnippet() {
		return errs.SnippetNameRequired
	}
	list, err := sc.s.Get(Ctx(), &types.GetRequest{Alias: a})
	if err != nil {
		if errs.HasCode(err, errs.CodeNotFound) {
			return sc.typeAhead(distinctName, func(s *types.Snippet, args []string) error {
				return sc.EWrite(out.SnippetCat(s))
				return nil
			})
		}
		return nil
	}
	if len(list.Items) == 0 {
		return errs.NotFound
	}
	if len(list.Items) == 1 {
		// TODO: use echo instead so that we can do variable substitution
		return sc.EWrite(out.SnippetCat(list.Items[0]))
	}
	return sc.EWrite(out.SnippetAmbiguousCat(list.Items))
}

func (sc *snippets) Patch(distinctName string, target string, patch string) error {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return err
	}
	alias, err := sc.s.Patch(Ctx(), &types.PatchRequest{Alias: a, Target: target, Patch: patch})
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetPatched(alias.String()))
}

func (sc *snippets) Clone(distinctName string, newFullName string) error {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return err
	}
	newA, err := types.ParseAlias(newFullName)
	if err != nil {
		return err
	}
	res, err := sc.s.Clone(Ctx(), &types.CloneRequest{Alias: a, New: newA})
	if err != nil {
		return err
	}
	sc.EWrite(out.SnippetList(prefs, res.List))
	return sc.EWrite(out.SnippetClonedAs(res.Snippet.Alias.URI()))
}

func (sc *snippets) Tag(distinctName string, tags ...string) error {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return err
	}
	alias, err := sc.s.Tag(Ctx(), &types.TagRequest{Alias: a, Tags: tags})
	if err != nil {
		return err
	}
	return sc.EWrite(out.Tagged(alias.String(), tags))
}

func (sc *snippets) UnTag(distinctName string, tags ...string) error {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return err
	}
	alias, err := sc.s.UnTag(Ctx(), &types.UnTagRequest{Alias: a, Tags: tags})
	if err != nil {
		return err
	}
	return sc.EWrite(out.UnTag(alias.String(), tags))
}

func (sc *snippets) CreatePouch(name string) error {
	res, err := sc.s.CreatePouch(Ctx(), &types.CreatePouchRequest{Name: name})
	if err != nil {
		return err
	}
	sc.EWrite(out.PrintRoot(prefs, &cliInfo, res.Root, &principal.User))
	return sc.EWrite(out.PouchCreated(name))
}

func (sc *snippets) Flatten(username string) error {
	list, err := sc.s.List(Ctx(), &types.ListRequest{Username: username, PrivateView: prefs.PrivateView})
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetList(prefs, list))
}

// GetEra lists snippets by special filters: @today @week @month @old
func (sc *snippets) GetEra(virtualPouch string) error {
	list, err := sc.s.List(Ctx(), &types.ListRequest{PrivateView: prefs.PrivateView})
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
	return sc.EWrite(out.SnippetList(prefs, list))
}

func (sc *snippets) List(username string, pouch string) error {
	if pouch == "" {
		r, err := sc.s.GetRoot(Ctx(), &types.RootRequest{Username: username, PrivateView: prefs.PrivateView})
		if err != nil {
			return err
		}
		return sc.EWrite(out.PrintRoot(prefs, &cliInfo, r, &principal.User))
	}
	var size int64
	list, err := sc.s.List(Ctx(), &types.ListRequest{Username: username, Pouch: pouch, Limit: size, PrivateView: prefs.PrivateView})
	if err != nil {
		return err
	}
	return sc.listSnippets(list)
}

func (sc *snippets) listSnippets(l *types.ListResponse) error {
	a := types.NewAlias(l.Username, l.Pouch.Name, "", "")
	_, err := sc.s.LogUse(Ctx(),
		&types.UseContext{
			Alias:  a,
			Type:   types.UseType_View,
			Status: types.UseStatus_Success,
			Time: types.KwkTime(time.Now()),
		})
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetList(prefs, l))
}

func (sc *snippets) typeAhead(term string, onSelect func(s *types.Snippet, args []string) error) error {
	res, err := sc.s.TypeAhead(Ctx(), &types.TypeAheadRequest{Term: term})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return errs.NotFound
	}
	if len(res.Results) == 1 {
		// TASK: Add 'Did you mean <snippet>?
	}
	snips := []*types.Snippet{}
	for _, v := range res.Results {
		snips = append(snips, v.Snippet)
	}
	res2, err := sc.Dialog.MultiChoice(out.FreeText("Please choose a snippet:  "), snips)
	if err != nil {
		return err
	}
	if res2 != nil {
		return onSelect(res2, nil)
	}
	return nil
}

func stdInAsString() string {
	scanner := bufio.NewScanner(os.Stdin)
	in := bytes.Buffer{}
	for scanner.Scan() {
		in.WriteString(scanner.Text() + "\n")
	}
	return in.String()
}

func (sc *snippets) handleMultiResponse(distinctName string, list []*types.Snippet) *types.Snippet {
	if len(list) == 1 {
		return list[0]
	} else if len(list) > 1 {
		s, _ := sc.MultiChoice(out.FreeText("Multiple matches. Choose a snippet to run:  "), list)
		return s
	}
	return nil
}

func (sc *snippets) getSnippet(uri string) (*types.ListResponse, *types.Alias, error) {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return nil, nil, err
	}
	list, err := sc.s.Get(Ctx(), &types.GetRequest{Alias: a})
	if err != nil {
		return nil, a, err
	}
	return list, a, nil
}

func (sc *snippets) renameSnippet(uri string, newSnipName string) (*types.Snippet, *types.SnipName, error) {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return nil, nil, err
	}
	sn, err := types.ParseSnipName(newSnipName)
	if err != nil {
		return nil, nil, err
	}
	res, err := sc.s.Rename(Ctx(), &types.RenameRequest{Alias: a, NewName: sn})
	return res.Snippet, res.Original, nil
}

func (sc *snippets) deleteSnippet(args []string) error {
	sn, pouch, err := types.ParseMany(args)
	if err != nil {
		return err
	}
	r := sc.Modal(out.SnippetCheckDelete(sn), prefs.AutoYes)
	if !r.Ok {
		sc.EWrite(out.SnippetsNotDeleted(sn))
	}
	res, err := sc.s.Delete(Ctx(), &types.DeleteRequest{Pouch: pouch, Names: sn})
	if err != nil {
		return err
	}
	sc.EWrite(out.SnippetList(prefs, res.List))
	return sc.EWrite(out.SnippetsDeleted(sn))
}

func (sc *snippets) deletePouch(pouch string) error {
	res := sc.Dialog.Modal(out.PouchCheckDelete(pouch), false)
	if !res.Ok {
		return sc.EWrite(out.PouchNotDeleted(pouch))
	}
	dres, err := sc.s.DeletePouch(Ctx(), &types.DeletePouchRequest{Name: pouch})
	if err != nil {
		return err
	}
	sc.EWrite(out.PrintRoot(prefs, &cliInfo, dres.Root, &principal.User))
	return sc.EWrite(out.PouchDeleted(pouch))
}
