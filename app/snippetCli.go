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
	"strconv"
	"strings"
	"time"
	"fmt"
)

type SnippetCli struct {
	ps 	 snippets.PouchService
	search   search.Term
	service  snippets.Service
	runner   cmd.Runner
	system   sys.Manager
	settings config.Persister
	su       setup.Provider
	dlg.Dialog
	tmpl.Writer
}

func NewSnippetCli(a snippets.Service, r cmd.Runner, s sys.Manager, d dlg.Dialog, w tmpl.Writer, t config.Persister, search search.Term, su setup.Provider) *SnippetCli {
	return &SnippetCli{service: a, runner: r, system: s, Dialog: d, Writer: w, settings: t, search: search, su: su}
}

func (s *SnippetCli) Share(distinctName string, destination string) {
	if list, _, err := s.get(distinctName); err != nil {
		s.HandleErr(err)
	} else {
		if alias := s.handleMultiResponse(distinctName, list); alias != nil {
			snip := "https://mail.google.com/mail/?ui=2&view=cm&fs=1&tf=1&su=&body=http%3A%2F%2Faus.kwk.co%2F" + alias.Username + "%2f" + alias.FullName
			gmail := models.NewSnippet(snip)
			gmail.Ext = "url"
			s.runner.Run(gmail, []string{})
		} else {
			s.Render("snippet:notfound", map[string]interface{}{"fullKey": distinctName})
		}
	}
}

func (s *SnippetCli) Run(distinctName string, args []string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		s.Render("validation:one-line", err)
	}

	suggest := func(fn string) {
		if res, err := s.search.Execute(fn); err != nil {
			s.HandleErr(err)
		} else if res.Total > 0 {
			s.Render("search:alphaSuggest", res)
			return
		}
	}

	if list, err := s.service.Get(*a); err != nil {
		suggest(distinctName)
	} else {
		if alias := s.handleMultiResponse(distinctName, list); alias != nil {
			if err = s.runner.Run(alias, args); err != nil {
				//TODO: Move to template
				fmt.Println(err)
			}
		} else {
			suggest(distinctName)
		}
	}
}

func (s *SnippetCli) New(snippet string, distinctName string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		s.Render("validation:one-line", err)
	}
	if createAlias, err := s.service.Create(snippet, *a, models.RoleStandard); err != nil {
		s.HandleErr(err)
	} else {
		///if createAlias.Snippet != nil {
		if createAlias.Snippet.Private {
			s.Render("snippet:newprivate", createAlias.Snippet)
		} else {
			s.Render("snippet:new", createAlias.Snippet)
		}
		// TODO: Add similarity prompt here
		//} else {
		//	matches := createAlias.TypeMatch.Matches
		//	r := s.MultiChoice("snippet:chooseruntime", "Choose a type for this snippet:", matches)
		//	winner := r.Value.(models.Match)
		//	if winner.Score == -1 {
		//		ca, _ := s.service.Create("_", "_", models.RoleStandard)
		//		matches = ca.TypeMatch.Matches
		//		winner = s.MultiChoice("snippet:chooseruntime", "Choose a type for this snippet:", matches).Value.(models.Match)
		//	}
		//	fk := fullKey + "." + winner.Extension
		//	s.New(uri, fk)
		//}
	}
}

func (s *SnippetCli) Edit(distinctName string) {
	if distinctName == "env" || distinctName == "prefs" {
		distinctName = models.NewSetupAlias(distinctName, "yml").String()
	}
	if list, _, err := s.get(distinctName); err != nil {
		s.HandleErr(err)
	} else {
		if alias := s.handleMultiResponse(distinctName, list); alias != nil {
			s.Render("snippet:editing", alias)
			if err := s.runner.Edit(alias); err != nil {
				s.HandleErr(err)
			} else {
				s.Render("snippet:edited", alias)
			}
		} else {
			s.Render("snippet:notfound", &models.Snippet{FullName: distinctName})
		}
	}
}

func (s *SnippetCli) Describe(distinctName string, description string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		s.Render("validation:one-line", err)
	}
	if alias, err := s.service.Update(*a, description); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:updated", alias)
	}
}

func (s *SnippetCli) Inspect(distinctName string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		s.Render("validation:one-line", err)
	}
	if list, err := s.service.Get(*a); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:inspect", list)
	}
}

// Problem is that if many items are passed for deletion and any are ambiguous how do we handle this?
func (s *SnippetCli) Delete(args []string) {
	if len(args) == 1 {
		r, err := s.ps.GetRoot("", true)
		if err != nil {
			panic("implement")
		}
		if r.IsPouch(args[0]) {
			_, err := s.ps.Delete(args[0])
			if err != nil {
				panic("implement")
			} else {
				panic("implement")
			}
		}
	}
	as, pouch, err := models.ParseMany(args)
	if err != nil {
		s.Render("validation:one-line", err)
		return
	}
	if r := s.Modal("snippet:delete", as, s.su.Prefs().AutoYes); r.Ok {
		if err := s.service.Delete("", pouch, as); err != nil {
			s.HandleErr(err)
		}
		s.Render("pouch:deleted", pouch)
	} else {
		s.Render("pouch:notdeleted", pouch)
	}
}

// kwk mv regions.txt reference -- moves the reference pouch, if no reference pouch then move to reference.txt
// kwk mv examples/regions.txt reference
func (s *SnippetCli) Move(args []string) {
	// TODO: if its a pouch and is prefixed with a "." then make private
	root, err := s.ps.GetRoot("", true)
	if err != nil {
		panic(root)
	}
	last := args[len(args) - 1]
	if root.IsPouch(args[0]) {
		s.ps.Rename(args[0], last)
		return
	} else if !root.IsPouch(last) {
		s.rename(args[0], args[1])
		return
	}
	as, source, err := models.ParseMany(args[0:len(args)-1])
	if err != nil {
		s.Render("validation:one-line", err)
		return
	}
	_, err = s.service.Move("", source, last, as)
}

func (s *SnippetCli) Cat(distinctName string) {
	if list, a, err := s.get(distinctName); err != nil {
		s.HandleErr(err)
	} else {
		if len(list.Items) == 0 {
			s.Render("snippet:notfound", a)
		} else if len(list.Items) == 1 {
			s.Render("snippet:cat", list.Items[0])
		} else {
			s.Render("snippet:ambiguouscat", list)
		}
	}
}

func (s *SnippetCli) Patch(distinctName string, target string, patch string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		s.Render("validation:one-line", err)
		return
	}
	if alias, err := s.service.Patch(*a, target, patch); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:patched", alias)
	}
}

func (s *SnippetCli) Clone(distinctName string, newFullName string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		s.Render("validation:one-line", err)
		return
	}
	newA, err := models.ParseAlias(newFullName)
	if err != nil {
		s.Render("validation:one-line", err)
		return
	}

	if alias, err := s.service.Clone(*a, *newA); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:cloned", alias)
	}
}

func (s *SnippetCli) Tag(distinctName string, tags ...string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		s.Render("validation:one-line", err)
		return
	}
	if alias, err := s.service.Tag(*a, tags...); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:tag", alias)
	}
}

func (s *SnippetCli) UnTag(distinctName string, tags ...string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		s.Render("validation:one-line", err)
		return
	}
	if alias, err := s.service.UnTag(*a, tags...); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:untag", alias)
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
func (s *SnippetCli) List(args ...string) {
	a, err := models.ParseAlias(args[0])
	if err != nil {
		panic(err)
	}
	if len(args) == 0 || (a.Username != "" && a.Pouch == "") {
		r, err := s.ps.GetRoot("", s.su.Prefs().ListAll)
		if err != nil {

		}
		s.Render("pouch:list-root", r)
		return
	}

	var size int64
	var tags = []string{}
	u := &models.User{}
	var username = ""
	for i, v := range args {
		if num, err := strconv.Atoi(v); err == nil {
			size = int64(num)
		} else {
			if i == 0 && v[len(v) - 1:] == "/" {
				username = strings.Replace(v, "/", "", -1)
			} else {
				tags = append(tags, v)
			}
		}
	}
	if username == "" {
		if err := s.settings.Get(models.ProfileFullKey, u, 0); err != nil {
			s.Render("api:not-authenticated", nil)
			return
		} else {
			username = u.Username
		}
	}
	p := &models.ListParams{Username:username, Size:size, Since:int64(time.Now().Unix() * 1000), Tags:tags, All:s.su.Prefs().ListAll}
	if list, err := s.service.List(p); err != nil {
		s.HandleErr(err)
	} else {
		list.Username = username
		s.Render("snippet:list", list)
	}
}

func (s *SnippetCli) handleMultiResponse(distinctName string, list *models.SnippetList) *models.Snippet {
	if list.Total == 1 {
		return &list.Items[0]
	} else if list.Total > 1 {
		r := s.MultiChoice("dialog:choose", "Multiple matches. Choose a snippet to run:", list.Items)
		s := r.Value.(models.Snippet)
		return &s
	} else {
		return nil
	}
}

func (s *SnippetCli) get(distinctName string) (*models.SnippetList, *models.Alias, error) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		return nil, nil, err
	}
	if list, err := s.service.Get(*a); err != nil {
		return nil, a, err
	} else {
		return list, a, nil
	}
}


func (s *SnippetCli) rename(distinctName string, newSnipName string) {
	a, err := models.ParseAlias(distinctName)
	if err != nil {
		s.Render("validation:one-line", err)
		return
	}
	sn, err := models.ParseSnipName(newSnipName)
	if err != nil {
		s.Render("validation:one-line", err)
		return
	}
	if snip, original, err := s.service.Rename(*a, *sn); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:renamed", &map[string]string{
			"distinctName":   original.String(),
			"newDistinctName":    snip.SnipName.String(),
		})
	}
}