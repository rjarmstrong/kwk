package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/cmd"
	"bitbucket.com/sharingmachine/kwkcli/src/gokwk"
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/kwkcli/src/persist"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/tmpl"
	"bitbucket.com/sharingmachine/types"
	"bitbucket.com/sharingmachine/types/errs"
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type SnippetCli struct {
	s        gokwk.Snippets
	runner   cmd.Runner
	settings persist.Persister
	dlg.Dialog
	tmpl.Writer
}

func NewSnippetCli(a gokwk.Snippets, r cmd.Runner, f persist.IO, d dlg.Dialog, w tmpl.Writer, t persist.Persister) *SnippetCli {
	return &SnippetCli{s: a, runner: r, Dialog: d, Writer: w, settings: t}
}

func (sc *SnippetCli) Search(args ...string) {
	term := strings.Join(args, " ")
	if res, err := sc.s.AlphaSearch(term); err != nil {
		sc.HandleErr(err)
	} else {
		res.Term = term
		sc.Out("search:alpha", res)
	}
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

func (sc *SnippetCli) Share(distinctName string, destination string) {
	if list, _, err := sc.getSnippet(distinctName); err != nil {
		sc.HandleErr(err)
	} else {
		if alias := sc.handleMultiResponse(distinctName, list); alias != nil {
			snip := "https://mail.google.com/mail/?ui=2&view=cm&fs=1&tf=1&su=&body=http%3A%2F%2Faus.kwk.co%2F" + alias.String()
			gmail := gokwk.NewSnippet(snip)
			gmail.Ext = "url"
			sc.runner.Run(gmail, []string{})
		} else {
			sc.Render("snippet:notfound", map[string]interface{}{"fullKey": distinctName})
		}
	}
}

func (sc *SnippetCli) Suggest(term string) {
	if res, err := sc.s.AlphaSearch(term); err != nil {
		sc.HandleErr(err)
	} else if res.Total > 0 {
		//sc.Render("search:alphaSuggest", res)
		sc.Render("search:typeahead", res)
		return
	}
}

func (sc *SnippetCli) run(selected *types.Snippet, args []string) {
	sc.runner.Run(selected, args)
}

func (sc *SnippetCli) Run(distinctName string, args []string) {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		sc.Render("validation:one-line", err)
	}

	if list, err := sc.s.Get(*a); err != nil {
		sc.typeAhead(distinctName, sc.run)
	} else {
		if alias := sc.handleMultiResponse(distinctName, list); alias != nil {
			//TODO: If username is not the current user or 'kwk' then prompt before executing.
			if err = sc.runner.Run(alias, args); err != nil {
				sc.HandleErr(err)
			}
		} else {
			sc.typeAhead(distinctName, sc.run)
		}
	}
}

func GuessArgs(a string, b string) (*types.Alias, string, error) {
	fstIsAlias := types.IsAlias(a)
	sndIsAlias := types.IsAlias(b)

	models.Debug("ARG1:%v (alias=%t) ARG2:%v (alias=%t)", a, fstIsAlias, b, sndIsAlias)

	if !fstIsAlias && !sndIsAlias {
		return nil, "", errs.New(errs.CodeInvalidArgument,
			"Please specify a pouch to create a snippet."+
				"\n   e.g."+
				"\n  `kwk new <pouch_name>/<snip_name>[.<ext>]  <snippet>`"+
				"\n  `kwk new <snippet> <pouch_name>/<snip_name>[.<ext>]`"+
				"\n  `<cmd> | kwk new <pouch_name>/<snip_name>[.<ext>]`",
		)
	}
	if fstIsAlias && sndIsAlias {
		return nil, "", errs.New(errs.CodeInvalidArgument,
			"Ambiguous",
		)
	}
	if fstIsAlias {
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

func (sc *SnippetCli) Create(args []string) {
	a1 := &types.Alias{}
	var snippet string
	if len(args) == 1 {
		if types.IsAlias(args[0]) {
			a, err := types.ParseAlias(args[0])
			if err != nil {
				sc.HandleErr(err)
				return
			}
			a1 = a
		} else {
			models.Debug("Assuming first arg is the snippet.")
			snippet = args[0]
			a1 = &types.Alias{}
		}
	} else if len(args) > 1 {
		a, s, err := GuessArgs(args[0], args[1])
		if err != nil {
			sc.HandleErr(err)
		}
		a1 = a
		snippet = s
	}
	if snippet == "" {
		snippet = stdInAsString()
	}
	if createAlias, err := sc.s.Create(snippet, *a1, types.RoleStandard); err != nil {
		sc.HandleErr(err)
		//TODO: If snippet is similar to an existing one prompt for it here.
	} else {
		sc.List("", types.PouchRoot)
		sc.Render("snippet:new", createAlias.Snippet.String())
	}
}

func (sc *SnippetCli) Edit(distinctName string) {
	innerEdit := func(s *types.Snippet) {
		sc.Render("snippet:editing", s)
		if err := sc.runner.Edit(s); err != nil {
			sc.HandleErr(err)
		} else {
			sc.Render("snippet:edited", &s)
		}
	}
	if distinctName == "env" {
		distinctName = models.NewSetupAlias(distinctName, "yml", true).String()
	}
	if list, _, err := sc.getSnippet(distinctName); err != nil {
		if errs.HasCode(err, errs.CodeNotFound) {
			a, err := types.ParseAlias(distinctName)
			if err != nil {
				sc.HandleErr(err)
			}
			r := sc.Dialog.Modal("snippet:edit-prompt", a, false)
			if r.Ok {
				r, err := sc.s.Create("", *a, types.RoleStandard)
				if err != nil {
					sc.HandleErr(err)
					return
				}
				innerEdit(r.Snippet)
			}
			return
		}
		sc.HandleErr(err)
	} else {
		if snippet := sc.handleMultiResponse(distinctName, list); snippet != nil {
			innerEdit(snippet)
		} else {
			sc.Render("snippet:notfound", &types.Snippet{})
		}
	}
}

func (sc *SnippetCli) Describe(distinctName string, description string) {
	a, err := types.ParseAlias(distinctName)
	if description == "" {
		sc.typeAhead(distinctName, func(s *types.Snippet, args []string) {
			cm := fmt.Sprintf("Enter new description for %s: ", s.SnipName.String())
			res, err := sc.FormField(cm)
			if err != nil {
				sc.HandleErr(err)
			}
			models.Debug("Form result: %+v", res.Value)
			sc.Describe(s.SnipName.String(), res.Value.(string))
		})
		return
	}
	if err != nil {
		sc.Render("validation:one-line", err)
	}
	if alias, err := sc.s.Update(*a, description); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:updated", alias)
	}
}

func (sc *SnippetCli) InspectListOrRun(distinctName string, forceInspect bool, args ...string) {
	a, err := types.ParseAlias(distinctName)
	v, err := sc.s.GetRoot(a.Username, true)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if a.Ext == "" && v.IsPouch(a.Name) {
		p := v.GetPouch(a.Name)
		if p.Type == types.PouchTypeVirtual {
			fmt.Println("List virtual.")
		} else {
			sc.List(a.Username, a.Name)
		}
		return
	}
	if err != nil {
		sc.Render("validation:one-line", err)
	}

	// GET SNIPPET
	if list, err := sc.s.Get(*a); err != nil {
		sc.typeAhead(distinctName, sc.run)
	} else {
		s := sc.handleMultiResponse(distinctName, list)
		if forceInspect || models.Prefs().RequireRunKeyword {
			sc.Render("snippet:inspect", s)
		} else {
			if err = sc.runner.Run(s, args); err != nil {
				sc.HandleErr(err)
			}
		}
	}
}

func (sc *SnippetCli) Delete(args []string) {
	r, err := sc.s.GetRoot("", true)
	if err != nil {
		sc.HandleErr(err)
	}
	if r.IsPouch(args[0]) {
		sc.deletePouch(args[0])
		return
	}
	sc.deleteSnippet(args)
}

func (sc *SnippetCli) Lock(pouch string) {
	_, err := sc.s.MakePrivate(pouch, true)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	sc.Render("pouch:locked", pouch)
}

func (sc *SnippetCli) UnLock(pouch string) {
	res := sc.Dialog.Modal("pouch:check-unlock", pouch, false)
	if res.Ok {
		_, err := sc.s.MakePrivate(pouch, false)
		if err != nil {
			sc.HandleErr(err)
			return
		}
		sc.Render("pouch:unlocked", pouch)
	} else {
		sc.Render("pouch:not-unlocked", pouch)
	}
}

// kwk mv regions.txt reference -- moves the reference pouch, if no reference pouch then move to reference.txt
// kwk mv examples/regions.txt reference
func (sc *SnippetCli) Move(args []string) {
	if len(args) < 2 {
		sc.HandleErr(errs.TwoArgumentsReqForMove)
		return
	}
	root, err := sc.s.GetRoot("", true)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	last := args[len(args)-1]
	// . (current folder is an alias for root directory)
	if last == "." {
		last = ""
	}
	if root.IsPouch(args[0]) {
		// rename pouch
		p, err := sc.s.RenamePouch(args[0], last)
		if err != nil {
			sc.HandleErr(err)
			return
		}
		sc.List("", types.PouchRoot)
		sc.Render("pouch:renamed", p)
		return
	} else if !root.IsPouch(last) && len(args) == 2 {
		// rename single snippet
		snip, original, err := sc.rename(args[0], args[1])
		if err != nil {
			sc.HandleErr(err)
		}
		sc.List("", snip.Pouch)
		sc.Render("snippet:renamed", &map[string]string{
			"originalName": original.String(),
			"newName":      snip.SnipName.String(),
		})
		return
	}
	as, source, err := types.ParseMany(args[0 : len(args)-1])
	if err != nil {
		sc.HandleErr(err)
		return
	}
	// move snippets into a pouch
	p, err := sc.s.Move("", source, last, as)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if last == "" {
		sc.List("", types.PouchRoot)
		sc.Render("snippet:moved-root", MoveResult{Quant: len(as)})
	} else {
		sc.List("", last)
		sc.Render("snippet:moved-pouch", MoveResult{Pouch: p, Quant: len(as)})
	}

}

type MoveResult struct {
	Pouch string
	Quant int
}

func (sc *SnippetCli) Cat(distinctName string) {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if !a.IsSnippet() {
		sc.HandleErr(errs.SnippetNameRequired)
		return
	}
	list, err := sc.s.Get(*a)
	if err != nil {
		if errs.HasCode(err, errs.CodeNotFound) {
			sc.typeAhead(distinctName, func(s *types.Snippet, args []string) {
				sc.Render("snippet:cat", s)
			})
			return
		}
		sc.HandleErr(err)
		return
	}

	if len(list.Snippets) == 0 {
		sc.Render("snippet:notfound", a)
		return
	}
	if len(list.Snippets) == 1 {
		// TODO: use echo instead so that we can do variable substitution
		sc.Render("snippet:cat", list.Snippets[0])
		return
	}
	sc.Render("snippet:ambiguouscat", list)
}

func (sc *SnippetCli) Patch(distinctName string, target string, patch string) {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if alias, err := sc.s.Patch(*a, target, patch); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:patched", alias)
	}
}

func (sc *SnippetCli) Clone(distinctName string, newFullName string) {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	newA, err := types.ParseAlias(newFullName)
	if err != nil {
		sc.HandleErr(err)
		return
	}

	if alias, err := sc.s.Clone(*a, *newA); err != nil {
		sc.HandleErr(err)
	} else {
		sc.List("", newA.Pouch)
		sc.Render("snippet:cloned", alias)
	}
}

func (sc *SnippetCli) Tag(distinctName string, tags ...string) {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if alias, err := sc.s.Tag(*a, tags...); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:tag", alias)
	}
}

func (sc *SnippetCli) UnTag(distinctName string, tags ...string) {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if alias, err := sc.s.UnTag(*a, tags...); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:untag", alias)
	}
}

func (sc *SnippetCli) CreatePouch(name string) {
	if _, err := sc.s.CreatePouch(name); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("pouch:created", name)
	}
}

func (sc *SnippetCli) Flatten(username string) {
	p := &models.ListParams{Username: username, IgnorePouches: true, All: models.Prefs().ListAll}
	if list, err := sc.s.List(p); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:list", list)
	}
}

func Bom(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func firstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}

// GetEra lists snippets by special filters: @today @week @month @old
func (sc *SnippetCli) GetEra(virtualPouch string) {
	p := &models.ListParams{Username: "", IgnorePouches: true, All: models.Prefs().ListAll}
	if list, err := sc.s.List(p); err != nil {
		sc.HandleErr(err)
	} else {
		era := []*types.Snippet{}
		var since, latest int64
		bod := Bod(time.Now()).Unix()
		isoYear, isoWeek := time.Now().ISOWeek()
		fod := firstDayOfISOWeek(isoYear, isoWeek, time.Local).Unix()
		bom := Bom(time.Now()).Unix()
		if virtualPouch == "@today" {
			since = bod
			latest = time.Now().Unix()
		} else if virtualPouch == "@week" {
			since = fod
			latest = bod
		} else if virtualPouch == "@month" {
			since = bom
			latest = fod
		} else if virtualPouch == "@old" {
			since = 0
			latest = bom
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
		sc.Render("snippet:list", list)
	}
}

// List
// Use root list:
// kwk /richard
// kwk
//
// Use snippet list:
// kwk richard (this is a pouch in this case)
// kwk /richard/examples
func (sc *SnippetCli) List(username string, pouch string) {
	if pouch == "" {
		r, err := sc.s.GetRoot(username, models.Prefs().ListAll)
		if err != nil {
			sc.HandleErr(err)
			return
		}
		sc.Render("pouch:list-root", r)
		return
	}
	var size int64
	p := &models.ListParams{Username: username, Pouch: pouch, Size: size, All: models.Prefs().ListAll}
	if list, err := sc.s.List(p); err != nil {
		sc.HandleErr(err)
	} else {
		a := types.NewAlias(username, pouch, "", "")
		sc.Render("snippet:list", list)
		sc.s.LogUse(
			*a,
			types.UseStatusSuccess,
			types.UseTypeView,
			nil,
		)
	}
}

func (sc *SnippetCli) typeAhead(term string, onSelect func(s *types.Snippet, args []string)) {
	if res, err := sc.s.AlphaSearch(term); err != nil {
		sc.HandleErr(err)
	} else {
		res.Term = term
		snips := []*types.Snippet{}
		for _, v := range res.Results {
			snips = append(snips, v.Snippet)
		}
		res := sc.Dialog.MultiChoice("dialog:choose", "Please choose one:", snips)
		if res != nil {
			onSelect(res, nil)
		}
	}
}

func stdInAsString() string {
	scanner := bufio.NewScanner(os.Stdin)
	in := bytes.Buffer{}
	for scanner.Scan() {
		in.WriteString(scanner.Text() + "\n")
	}
	return in.String()
}

func (sc *SnippetCli) handleMultiResponse(distinctName string, list *models.ListView) *types.Snippet {
	list.Version = CLIInfo.String()
	if list.Total == 1 {
		return list.Snippets[0]
	} else if list.Total > 1 {
		return sc.MultiChoice("dialog:choose", "Multiple matches. Choose a snippet to run:", list.Snippets)
	} else {
		return nil
	}
}

func (sc *SnippetCli) getSnippet(distinctName string) (*models.ListView, *types.Alias, error) {
	a, err := types.ParseAlias(distinctName)
	if err != nil {
		return nil, nil, err
	}
	list, err := sc.s.Get(*a)
	if err != nil {
		return nil, a, err
	}
	return list, a, nil
}

func (sc *SnippetCli) rename(distinctName string, newSnipName string) (*types.Snippet, *types.SnipName, error) {
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

func (sc *SnippetCli) deleteSnippet(args []string) {
	as, pouch, err := types.ParseMany(args)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if r := sc.Modal("snippet:check-delete", as, models.Prefs().AutoYes); r.Ok {
		if err := sc.s.Delete("", pouch, as); err != nil {
			sc.HandleErr(err)
			return
		}
		sc.List("", pouch)
		sc.Render("snippet:deleted", pouch)
	} else {
		sc.Render("snippet:not-deleted", pouch)
	}
}

func (sc *SnippetCli) deletePouch(pouch string) {
	res := sc.Dialog.Modal("pouch:check-delete", pouch, false)
	if res.Ok {
		_, err := sc.s.DeletePouch(pouch)
		if err != nil {
			sc.HandleErr(err)
			return
		}
		sc.List("", types.PouchRoot)
		sc.Render("pouch:deleted", pouch)
	} else {
		sc.Render("pouch:not-deleted", pouch)
	}
}
