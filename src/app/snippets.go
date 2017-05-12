package app

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
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
	headers  Headers
	s        types.SnippetsClient
	v        types.SnippetsServer
	runner   Runner
	settings Persister
	Dialog
	vwrite.Writer
}

func NewSnippet(s types.SnippetsClient, r Runner, d Dialog, w vwrite.Writer, t Persister) *snippets {
	return &snippets{s: s, runner: r, Dialog: d, Writer: w, settings: t}
}

func (sc *snippets) Search(args ...string) error {
	term := strings.Join(args, " ")
	res, err := sc.s.Alpha(sc.headers.Context(), &types.AlphaRequest{Term: term})
	if err != nil {
		return err
	}
	return sc.EWrite(out.AlphaSearchResult(res))
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

func (sc *snippets) Suggest(term string) error {
	res, err := sc.s.TypeAhead(sc.headers.Context(), &types.TypeAheadRequest{Term: term})
	if err != nil {
		return err
	}
	if res.Total > 0 {
		return sc.EWrite(out.AlphaTypeAhead(res))
	}
	return nil
}

func (sc *snippets) run(selected *types.Snippet, args []string) error {
	return sc.runner.Run(selected, args)
}

func (sc *snippets) Run(uri string, args []string) error {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return err
	}
	list, err := sc.s.Get(sc.headers.Context(), &types.GetRequest{Alias: a})
	if err != nil {
		sc.typeAhead(uri, sc.run)
	}
	alias := sc.handleMultiResponse(uri, list)
	if alias != nil {
		//TODO: If username is not the current user or 'kwk' then prompt before executing.
		return sc.runner.Run(alias, args)
	}
	return sc.typeAhead(uri, sc.run)
}

func GuessArgs(a string, b string) (*types.Alias, string, error) {
	aIsNotSUri := types.IsDefNotPouchedSnippetURI(a)
	bIsNotSUri := types.IsDefNotPouchedSnippetURI(b)

	Debug("ARG1:%v (snip=%t) ARG2:%v (snip=%t)", a, aIsNotSUri, b, bIsNotSUri)

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

func (sc *snippets) Create(args []string) error {
	a1 := &types.Alias{}
	var snippet string
	if len(args) == 1 {
		if types.IsDefNotPouchedSnippetURI(args[0]) {
			a, err := types.ParseAlias(args[0])
			if err != nil {
				return err
			}
			a1 = a
		} else {
			Debug("Assuming first arg is the snippet.")
			snippet = args[0]
			a1 = &types.Alias{}
		}
	} else if len(args) > 1 {
		a, s, err := GuessArgs(args[0], args[1])
		if err != nil {
			return err
		}
		a1 = a
		snippet = s
	}
	if snippet == "" {
		snippet = stdInAsString()
	}
	res, err := sc.s.Create(snippet, *a1, types.RoleStandard)
	if err != nil {
		return err
		//TODO: If snippet is similar to an existing one prompt for it here.
	}
	sc.List("", types.PouchRoot)
	return sc.EWrite(out.SnippetCreated(res.Snippet))
}

func (sc *snippets) Edit(distinctName string) error {
	innerEdit := func(s *types.Snippet) error {
		sc.Write(out.SnippetEditing(s))
		err := sc.runner.Edit(s)
		if err != nil {
			return err
		}
		return sc.EWrite(out.SnippetEdited(s))
	}
	if distinctName == "env" {
		distinctName = NewSetupAlias(distinctName, "yml", true).String()
	}
	list, _, err := sc.getSnippet(distinctName)
	if err != nil {
		if errs.HasCode(err, errs.CodeNotFound) {
			a, err := types.ParseAlias(distinctName)
			if err != nil {
				return err
			}
			r := sc.Dialog.Modal(out.SnippetEditNewPrompt(a.String()), false)
			if r.Ok {
				r, err := sc.s.Create("", *a, types.RoleStandard)
				if err != nil {
					return err
				}
				return innerEdit(r.Snippet)
			}
		}
		return err
	}
	snippet := sc.handleMultiResponse(distinctName, list)
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
				out.FreeText(fmt.Sprintf("Enter new description for %s: ", s.SnipName.String())), false)
			if err != nil {
				return err
			}
			Debug("Form result: %+v", res.Value)
			return sc.Describe(s.SnipName.String(), res.Value.(string))
		})
	}
	if err != nil {
		return err
	}
	alias, err := sc.s.Update(*a, description)
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetDescriptionUpdated(alias.String(), description))
}

func (sc *snippets) InspectListOrRun(distinctName string, forceInspect bool, args ...string) error {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return err
	}
	v, err := sc.s.GetRoot(a.Username, true)
	if err != nil {
		return err
	}
	if a.Ext == "" && v.IsPouch(a.Name) {
		p := v.GetPouch(a.Name)
		if p.Type == types.PouchTypeVirtual {
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
	list, err := sc.s.Get(*a)
	if err != nil {
		return sc.typeAhead(distinctName, sc.run)
	}
	s := sc.handleMultiResponse(distinctName, list)
	if forceInspect || Prefs().RequireRunKeyword {
		sc.Write(out.SnippetView(s))
	}
	return sc.runner.Run(s, args)
}

func (sc *snippets) Delete(args []string) error {
	r, err := sc.s.GetRoot("", true)
	if err != nil {
		return err
	}
	if r.IsPouch(args[0]) {
		return sc.deletePouch(args[0])
	}
	return sc.deleteSnippet(args)
}

func (sc *snippets) Lock(pouch string) error {
	_, err := sc.s.MakePrivate(pouch, true)
	if err != nil {
		return err
	}
	return sc.EWrite(out.PouchLocked(pouch))
}

func (sc *snippets) UnLock(pouch string) error {
	res := sc.Dialog.Modal(out.PouchCheckUnLock(pouch), false)
	if res.Ok {
		_, err := sc.s.MakePrivate(pouch, false)
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
	root, err := sc.s.GetRoot("", true)
	if err != nil {
		return err
	}
	last := args[len(args)-1]
	if root.IsPouch(args[0]) {
		p, err := sc.s.RenamePouch(args[0], last)
		if err != nil {
			return err
		}
		sc.List("", types.PouchRoot)
		return sc.EWrite(out.PouchRenamed(args[0], p))
	} else if !root.IsPouch(last) && len(args) == 2 {
		// rename single snippet
		snip, original, err := sc.rename(args[0], args[1])
		if err != nil {
			return err
		}
		sc.List("", snip.Pouch)
		return sc.EWrite(out.SnippetRenamed(original.String(), snip.String()))
		return nil
	}
	as, source, err := types.ParseMany(args[0 : len(args)-1])
	if err != nil {
		return err
	}
	// move snippets into a pouch
	p, err := sc.s.Move("", source, last, as)
	if err != nil {
		return err
	}
	sc.List("", last)
	return sc.EWrite(out.SnippetsMoved(as, p))
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
	list, err := sc.s.Get(*a)
	if err != nil {
		if errs.HasCode(err, errs.CodeNotFound) {
			return sc.typeAhead(distinctName, func(s *types.Snippet, args []string) error {
				return sc.EWrite(out.SnippetCat(s))
				return nil
			})
		}
		return nil
	}
	if len(list.Snippets) == 0 {
		return errs.NotFound
	}
	if len(list.Snippets) == 1 {
		// TODO: use echo instead so that we can do variable substitution
		return sc.EWrite(out.SnippetCat(list.Snippets[0]))
	}
	return sc.EWrite(out.SnippetAmbiguousCat(list.Snippets))
}

func (sc *snippets) Patch(distinctName string, target string, patch string) error {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return err
	}
	alias, err := sc.s.Patch(*a, target, patch)
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
	alias, err := sc.s.Clone(*a, *newA)
	if err != nil {
		return err
	}
	err = sc.List("", newA.Pouch)
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetClonedAs(alias.String()))
}

func (sc *snippets) Tag(distinctName string, tags ...string) error {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return err
	}
	alias, err := sc.s.Tag(*a, tags...)
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
	alias, err := sc.s.UnTag(*a, tags...)
	if err != nil {
		return err
	}
	return sc.EWrite(out.UnTag(alias.String(), tags))
}

func (sc *snippets) CreatePouch(name string) error {
	_, err := sc.s.CreatePouch(name)
	if err != nil {
		return err
	}
	return sc.EWrite(out.PouchCreated(name))
}

func (sc *snippets) Flatten(username string) error {
	p := &ListParams{Username: username, IgnorePouches: true, All: Prefs().ListAll}
	list, err := sc.s.List(p)
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetList(list))
}

// GetEra lists snippets by special filters: @today @week @month @old
func (sc *snippets) GetEra(virtualPouch string) error {
	p := &ListParams{Username: "", IgnorePouches: true, All: Prefs().ListAll}
	list, err := sc.s.List(p)
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
	for _, v := range list.Snippets {
		if v.RunStatusTime.Unix() > since && v.RunStatusTime.Unix() < latest {
			era = append(era, v)
		}
	}
	sort.Slice(era, func(i, j int) bool {
		return era[i].RunStatusTime.Unix() < era[j].RunStatusTime.Unix()
	})
	list.Snippets = era
	return sc.EWrite(out.SnippetList(list))
}

func (sc *snippets) List(username string, pouch string) error {
	if pouch == "" {
		r, err := sc.s.GetRoot(username, Prefs().ListAll)
		if err != nil {
			return err
		}
		return sc.EWrite(out.PrintRoot(r))
	}
	var size int64
	p := &ListParams{Username: username, Pouch: pouch, Size: size, All: Prefs().ListAll}
	list, err := sc.s.List(p)
	if err != nil {
		return err
	}
	return sc.listSnippets(list)
}

func (sc *snippets) listSnippets(l *types.ListResponse) error {
	a := types.NewAlias(l.Username, l.Pouch.Name, "", "")
	_, err := sc.s.LogUse(sc.headers.Context(),
		&types.UseContext{
			Alias:   a,
			Type:    types.UseType_View,
			Status:  types.UseStatus_Success,
			Preview: "",
		})
	if err != nil {
		return err
	}
	return sc.EWrite(out.SnippetList(l))
}

func (sc *snippets) typeAhead(term string, onSelect func(s *types.Snippet, args []string) error) error {
	res, err := sc.s.TypeAhead(sc.headers.Context(), &types.TypeAheadRequest{Term: term})
	if err != nil {
		return err
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

func (sc *snippets) handleMultiResponse(distinctName string, list *types.ListView) *types.Snippet {
	list.Version = cliInfo.String()
	if list.Total == 1 {
		return list.Snippets[0]
	} else if list.Total > 1 {
		s, _ := sc.MultiChoice(out.FreeText("Multiple matches. Choose a snippet to run:  "), list.Snippets)
		return s
	}
	return nil
}

func (sc *snippets) getSnippet(uri string) (*types.ListView, *types.Alias, error) {
	a, err := types.ParseAlias(uri)
	if err != nil {
		return nil, nil, err
	}
	list, err := sc.s.Get(sc.headers.Context(), &types.GetRequest{Alias: a})
	if err != nil {
		return nil, a, err
	}
	return list, a, nil
}

func (sc *snippets) rename(distinctName string, newSnipName string) (*types.Snippet, *types.SnipName, error) {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return nil, nil, err
	}
	sn, err := types.ParseSnipName(newSnipName)
	if err != nil {
		return nil, nil, err
	}
	return sc.s.Rename(*a, *sn)
}

func (sc *snippets) deleteSnippet(args []string) error {
	sn, pouch, err := types.ParseMany(args)
	if err != nil {
		return err
	}
	r := sc.Modal(out.SnippetCheckDelete(sn), Prefs().AutoYes)
	if !r.Ok {
		sc.EWrite(out.SnippetsNotDeleted(sn))
	}
	res, err := sc.s.Delete(sc.headers.Context(), &types.DeleteRequest{Pouch: pouch, Names: sn})
	if err != nil {
		return err
	}
	sc.listSnippets(res.List)
	return sc.EWrite(out.SnippetsDeleted(sn))
}

func (sc *snippets) deletePouch(pouch string) error {
	res := sc.Dialog.Modal(out.PouchCheckDelete(pouch), false)
	if res.Ok {
		_, err := sc.s.DeletePouch(pouch)
		if err != nil {
			return err
		}
		sc.List("", types.PouchRoot)
		return sc.EWrite(out.PouchDeleted(pouch))
	}
	return sc.EWrite(out.PouchNotDeleted(pouch))
}
