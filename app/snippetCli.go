package app

import (
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"github.com/rjarmstrong/fzf/src"
	"strings"
	"fmt"
	"os"
	"bufio"
	"bytes"
)

type SnippetCli struct {
	s        snippets.Service
	runner   cmd.Runner
	system   sys.Manager
	settings config.Persister
	dlg.Dialog
	tmpl.Writer
}

func NewSnippetCli(a snippets.Service, r cmd.Runner, s sys.Manager, d dlg.Dialog, w tmpl.Writer, t config.Persister) *SnippetCli {
	return &SnippetCli{s: a, runner: r, system: s, Dialog: d, Writer: w, settings: t}
}

func (sc *SnippetCli) Search(args ...string) {
	term := strings.Join(args, " ")
	if res, err := sc.s.AlphaSearch(term); err != nil {
		sc.HandleErr(err)
	} else {
		res.Term = term
		sc.Render("search:alpha", res)
	}
}

func (sc *SnippetCli) ListCategory(category string, args ...string) {
	username, pouch, err := models.ParsePouch(args[0])
	if err != nil {
		sc.HandleErr(err)
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
	p := &models.ListParams{Category: category, Username: username, Pouch: pouch, Size: size, All: models.Prefs().ListAll}
	if list, err := sc.s.List(p); err != nil {
		sc.HandleErr(err)
	} else {
		sc.Render("snippet:list", list)
	}
}

func (sc *SnippetCli) Share(distinctName string, destination string) {
	if list, _, err := sc.getSnippet(distinctName); err != nil {
		sc.HandleErr(err)
	} else {
		if alias := sc.handleMultiResponse(distinctName, list); alias != nil {
			snip := "https://mail.google.com/mail/?ui=2&view=cm&fs=1&tf=1&su=&body=http%3A%2F%2Faus.kwk.co%2F" + alias.String()
			gmail := models.NewSnippet(snip)
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

func (sc *SnippetCli) Run(distinctName string, args []string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		sc.Render("validation:one-line", err)
	}
	rerun := func(selected string) {
		r := sc.Dialog.FormField(fmt.Sprintf("kwk %s ", selected))
		argstr := r.Value.(string)
		sc.Run(selected, strings.Split(argstr, " "))
	}
	if list, err := sc.s.Get(*a); err != nil {
		sc.typeAhead(distinctName, rerun)
	} else {
		if alias := sc.handleMultiResponse(distinctName, list); alias != nil {
			//TODO: If username is not the current user or 'kwk' then prompt before executing.
			if err = sc.runner.Run(alias, args); err != nil {
				sc.HandleErr(err)
			}
		} else {
			sc.typeAhead(distinctName, rerun)
		}
	}
}

func (sc *SnippetCli) Create(args []string) {
	alias := &models.Alias{}
	var snippet string
	// TODO: 'http://www.fileformat.info/info/unicode/char/1f526/index.htm'
	// When there is only one argument a path or url will incorrectly be guessed as being an alias
	// suggest that if has more than 3 segments then we treat as the snippet instead.
	if len(args) == 1 {
		a, err := models.ParseAlias(args[0])
		if err != nil {
			sc.HandleErr(err)
			return
		}
		if a.Ext == "" {
			log.Debug("Assuming first arg is the snippet.")
			snippet = args[0]
			a = &models.Alias{}
		}
		alias = a
	} else if len(args) > 1 {
		log.Debug("Assuming first arg is the snippet and second is alias.")
		a, err := models.ParseAlias(args[1])
		if err != nil {
			sc.HandleErr(err)
			return
		}
		alias = a
		snippet = args[0]
	}
	if snippet == "" {
		snippet = stdInAsString()
	}

	if createAlias, err := sc.s.Create(snippet, *alias, models.SnipRoleStandard); err != nil {
		sc.HandleErr(err)
	} else {
		sc.List("", models.ROOT_POUCH)
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
	innerEdit := func(s *models.Snippet) {
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
		if models.HasErrCode(err, models.Code_NotFound) {
			a, err := models.ParseAlias(distinctName)
			if err != nil {
				sc.HandleErr(err)
			}
			r := sc.Dialog.Modal("snippet:edit-prompt", a, false)
			if r.Ok {
				r, err := sc.s.Create("", *a, models.SnipRoleStandard)
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
			sc.Render("snippet:notfound", &models.Snippet{})
		}
	}
}

func (sc *SnippetCli) Describe(distinctName string, description string) {
	a, err := models.ParseAlias(distinctName)
	if description == "" {
		sc.typeAhead(distinctName, func(input string) {
			cm := fmt.Sprintf("Enter new description for %s: ", input)
			if res := sc.FormField(cm); res.Ok {
				log.Debug("Form result: %+v", res.Value)
				sc.Describe(input, res.Value.(string))
			} else {
				log.Debug("not ok")
			}
			return
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
	a, err := models.ParseAlias(distinctName)
	v, err := sc.s.GetRoot(a.Username, true)
	if err != nil {
		sc.HandleErr(err)
		log.Error("Error getting root, but not critical to 'run'", err)
		return
	} else if a.Ext == "" && v.IsPouch(a.Name) {
		p := v.GetPouch(a.Name)
		if p.Type == models.PouchType_Virtual {
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
		sc.HandleErr(err)
	} else {
		s := sc.handleMultiResponse(distinctName, list)
		if forceInspect || models.Prefs().RequireRunKeyword {
			sc.Render("snippet:inspect", s)
		} else {
			sc.Run(distinctName, args)
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
		// rename pouch
		p, err := sc.s.RenamePouch(args[0], last)
		if err != nil {
			sc.HandleErr(err)
			return
		}
		sc.List("", models.ROOT_POUCH)
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
	as, source, err := models.ParseMany(args[0:len(args)-1])
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
		sc.List("", models.ROOT_POUCH)
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
	if list, a, err := sc.getSnippet(distinctName); err != nil {
		if models.HasErrCode(err, models.Code_NotFound) {
			sc.typeAhead(distinctName, func(str string) {
				_ = sc.Dialog.FormField(fmt.Sprintf("kwk cat %s ", str))
				sc.Cat(str)
			})
		} else {
			sc.HandleErr(err)
		}
	} else {
		if len(list.Snippets) == 0 {
			//sc.suggest(distinctName)
			sc.Render("snippet:notfound", a)
		} else if len(list.Snippets) == 1 {
			// TODO: use echo instead so that we can do variable substitution
			sc.Render("snippet:cat", list.Snippets[0])
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
		sc.List("", newA.Pouch)
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

func (sc *SnippetCli) Flatten(username string) {
	p := &models.ListParams{Username: username, IgnorePouches: true, All: models.Prefs().ListAll}
	if list, err := sc.s.List(p); err != nil {
		sc.HandleErr(err)
	} else {
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
		a := models.NewAlias(username, pouch, "", "")
		sc.Render("snippet:list", list)
		sc.s.LogUse(
			*a,
			models.UseStatusSuccess,
			models.UseTypeView,
			"",
		)
	}
}

func (sc *SnippetCli) typeAhead(distinctName string, onSelect func(name string)) {
	exe, _ := os.Executable()
	opt := fzf.ParseOptionsAs(fmt.Sprintf("--preview=%s cat %s", exe, "{}"), "-1", "--preview-window=right:70%", "--header=   Suggestions:   ", "--query="+distinctName, "--reverse", "--margin=2,6,2,2", "--height=40%", "--no-mouse", "--color=prompt:008,header:0,headerbg:008,fg:255,hl:006,pointer:014,hl+:014,fg+:006,bg+:000")
	opt.Printer = onSelect
	fzf.Run(fmt.Sprintf("%s suggest %s", exe, distinctName), opt)
}

func stdInAsString() string {
	scanner := bufio.NewScanner(os.Stdin)
	in := bytes.Buffer{}
	for scanner.Scan() {
		in.WriteString(scanner.Text() + "\n")
	}
	return in.String()
}

func (sc *SnippetCli) handleMultiResponse(distinctName string, list *models.ListView) *models.Snippet {
	if list.Total == 1 {
		return list.Snippets[0]
	} else if list.Total > 1 {
		r := sc.MultiChoice("dialog:choose", "Multiple matches. Choose a snippet to run:", list.Snippets)
		s := r.Value.(*models.Snippet)
		return s
	} else {
		return nil
	}
}

func (sc *SnippetCli) getSnippet(distinctName string) (*models.ListView, *models.Alias, error) {
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

func (sc *SnippetCli) rename(distinctName string, newSnipName string) (*models.Snippet, *models.SnipName, error) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		return nil, nil, err
	}
	sn, err := models.ParseSnipName(newSnipName)
	if err != nil {
		return nil, nil, err
	}
	return sc.s.Rename(*a, *sn)
}

func (sc *SnippetCli) deleteSnippet(args []string) {
	as, pouch, err := models.ParseMany(args)
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
		sc.List("", models.ROOT_POUCH)
		sc.Render("pouch:deleted", pouch)
	} else {
		sc.Render("pouch:not-deleted", pouch)
	}
}
