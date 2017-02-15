package app

import (
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/search"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"bitbucket.com/sharingmachine/kwkcli/setup"
	"fmt"
)

type SnippetCli struct {
	search   search.Term
	s        snippets.Service
	runner   cmd.Runner
	system   sys.Manager
	settings config.Persister
	su       setup.Provider
	dlg.Dialog
	tmpl.Writer
}

func NewSnippetCli(a snippets.Service, r cmd.Runner, s sys.Manager, d dlg.Dialog, w tmpl.Writer, t config.Persister, search search.Term, su setup.Provider) *SnippetCli {
	return &SnippetCli{s: a, runner: r, system: s, Dialog: d, Writer: w, settings: t, search: search, su: su}
}

func (sc *SnippetCli) Share(distinctName string, destination string) {
	if list, _, err := sc.get(distinctName); err != nil {
		sc.HandleErr(err)
	} else {
		if alias := sc.handleMultiResponse(distinctName, list); alias != nil {
			snip := "https://mail.google.com/mail/?ui=2&view=cm&fs=1&tf=1&su=&body=http%3A%2F%2Faus.kwk.co%2F" + alias.Username + "%2f" + alias.FullName
			gmail := models.NewSnippet(snip)
			gmail.Ext = "url"
			sc.runner.Run(gmail, []string{})
		} else {
			sc.Render("snippet:notfound", map[string]interface{}{"fullKey": distinctName})
		}
	}
}

func (sc *SnippetCli) Run(distinctName string, args []string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		sc.Render("validation:one-line", err)
	}

	suggest := func(fn string) {
		if res, err := sc.search.Execute(fn); err != nil {
			sc.HandleErr(err)
		} else if res.Total > 0 {
			sc.Render("search:alphaSuggest", res)
			return
		}
	}

	if list, err := sc.s.Get(*a); err != nil {
		suggest(distinctName)
	} else {
		if alias := sc.handleMultiResponse(distinctName, list); alias != nil {
			if err = sc.runner.Run(alias, args); err != nil {
				//TODO: Move to template
				fmt.Println(err)
			}
		} else {
			suggest(distinctName)
		}
	}
}

func (sc *SnippetCli) Create(args []string) {
	var alias *models.Alias
	var snippet string
	if len(args) == 0 {
		alias = &models.Alias{}
	} else if len(args) == 1 {
		a, err := models.ParseAlias(args[0])
		if err != nil {
			sc.HandleErr(err)
			return
		}
		if a.Ext != "" {
			alias = a
		} else {
			snippet = args[0]
		}
	} else  {
		a, err := models.ParseAlias(args[1])
		if err != nil {
			sc.HandleErr(err)
			return
		}
		alias = a
		snippet = args[0]
	}

	if createAlias, err := sc.s.Create(snippet, *alias, models.RoleStandard); err != nil {
		sc.HandleErr(err)
	} else {

		sc.List()
		sc.Render("snippet:new", createAlias.Snippet.String())
		// TODO: Add similarity prompt here
		//} else {
		//	matches := createAlias.TypeMatch.Matches
		//	r := s.MultiChoice("snippet:chooseruntime", "Choose a type for this snippet:", matches)
		//	winner := r.Value.(models.Match)
		//	if winner.Score == -1 {
		//		ca, _ := sc.s.Create("_", "_", models.RoleStandard)
		//		matches = ca.TypeMatch.Matches
		//		winner = s.MultiChoice("snippet:chooseruntime", "Choose a type for this snippet:", matches).Value.(models.Match)
		//	}
		//	fk := fullKey + "." + winner.Extension
		//	s.New(uri, fk)
		//}
	}
}

func (sc *SnippetCli) Edit(distinctName string) {
	if distinctName == "env"  {
		distinctName = models.NewSetupAlias(distinctName, "yml", true).String()
	}
	if list, _, err := sc.get(distinctName); err != nil {
		sc.HandleErr(err)
	} else {
		if alias := sc.handleMultiResponse(distinctName, list); alias != nil {
			sc.Render("snippet:editing", alias)
			if err := sc.runner.Edit(alias); err != nil {
				sc.HandleErr(err)
			} else {
				sc.Render("snippet:edited", alias)
			}
		} else {
			sc.Render("snippet:notfound", &models.Snippet{FullName: distinctName})
		}
	}
}

func (sc *SnippetCli) Describe(distinctName string, description string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		sc.Render("validation:one-line", err)
	}
	if alias, err := sc.s.Update(*a, description); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:updated", alias)
	}
}

func (sc *SnippetCli) Inspect(distinctName string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		sc.Render("validation:one-line", err)
	}
	if list, err := sc.s.Get(*a); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:inspect", list)
	}
}

// Problem is that if many items are passed for deletion and any are ambiguous how do we handle this?
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

func (sc *SnippetCli) deleteSnippet(args []string) {
	as, pouch, err := models.ParseMany(args)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if r := sc.Modal("snippet:check-delete", as, sc.su.Prefs().AutoYes); r.Ok {
		if err := sc.s.Delete("", pouch, as); err != nil {
			sc.HandleErr(err)
			return
		}
		sc.List()
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
		sc.List()
		sc.Render("pouch:deleted", pouch)
	} else {
		sc.Render("pouch:not-deleted", pouch)
	}
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
		sc.HandleErr(models.ErrOneLine(models.Code_TwoArgumentsRequiredForMove, "Two arguments are required for the move command."))
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
		p, err := sc.s.RenamePouch(args[0], last)
		if err != nil {
			sc.HandleErr(err)
			return
		}
		sc.List()
		sc.Render("pouch:renamed", p)
		return
	} else if !root.IsPouch(last) && len(args) == 2 {
		sc.List()
		sc.rename(args[0], args[1])
		return
	}
	as, source, err := models.ParseMany(args[0:len(args)-1])
	if err != nil {
		sc.HandleErr(err)
		return
	}
	p, err := sc.s.Move("", source, last, as)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if last == "" {
		sc.List()
		sc.Render("snippet:moved", "root")
	} else {
		sc.List()
		sc.Render("snippet:moved", p)
	}

}

func (sc *SnippetCli) Cat(distinctName string) {
	if list, a, err := sc.get(distinctName); err != nil {
		sc.HandleErr(err)
	} else {
		if len(list.Items) == 0 {
			sc.Render("snippet:notfound", a)
		} else if len(list.Items) == 1 {
			sc.Render("snippet:cat", list.Items[0])
		} else {
			sc.Render("snippet:ambiguouscat", list)
		}
	}
}

func (sc *SnippetCli) Patch(distinctName string, target string, patch string) {
	a, err := models.ParseAlias(distinctName)
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
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	newA, err := models.ParseAlias(newFullName)
	if err != nil {
		sc.HandleErr(err)
		return
	}

	if alias, err := sc.s.Clone(*a, *newA); err != nil {
		sc.HandleErr(err)
	} else {
		sc.List()
		sc.Render("snippet:cloned", alias)
	}
}

func (sc *SnippetCli) Tag(distinctName string, tags ...string) {
	a, err := models.ParseAlias(distinctName)
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
	a, err := models.ParseAlias(distinctName)
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

// List
// Use root list:
// kwk ls /richard
// kwk ls
//
// Use snippet list:
// kwk ls richard (this is a pouch in this case)
// kwk ls /richard/examples
func (sc *SnippetCli) List(args ...string) {
	if len(args) == 0 {
		r, err := sc.s.GetRoot("", sc.su.Prefs().ListAll)
		if err != nil {
			sc.HandleErr(err)
			return
		}
		sc.Render("pouch:list-root", r)
		return
	}
	username, pouch, err := models.ParsePouch(args[0])
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if pouch == "" {
		r, err := sc.s.GetRoot(username, sc.su.Prefs().ListAll)
		if err != nil {
			sc.HandleErr(err)
		}
		sc.Render("pouch:list-root", r)
		return
	}

	var size int64
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
	p := &models.ListParams{Username: username, Pouch: pouch, Size: size, All: sc.su.Prefs().ListAll}
	if list, err := sc.s.List(p); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:list", list)
	}
}

func (sc *SnippetCli) handleMultiResponse(distinctName string, list *models.SnippetList) *models.Snippet {
	if list.Total == 1 {
		return list.Items[0]
	} else if list.Total > 1 {
		r := sc.MultiChoice("dialog:choose", "Multiple matches. Choose a snippet to run:", list.Items)
		s := r.Value.(models.Snippet)
		return &s
	} else {
		return nil
	}
}

func (sc *SnippetCli) get(distinctName string) (*models.SnippetList, *models.Alias, error) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		return nil, nil, err
	}
	if list, err := sc.s.Get(*a); err != nil {
		return nil, a, err
	} else {
		return list, a, nil
	}
}

func (sc *SnippetCli) rename(distinctName string, newSnipName string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	sn, err := models.ParseSnipName(newSnipName)
	if err != nil {
		sc.HandleErr(err)
		return
	}
	if snip, original, err := sc.s.Rename(*a, *sn); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:renamed", &map[string]string{
			"originalName": original.String(),
			"newName":      snip.SnipName.String(),
		})
	}
}
