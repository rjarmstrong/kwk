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
	"strconv"
	"strings"
	"time"
	"fmt"
)

type SnippetCli struct {
	search   search.Term
	service  snippets.Service
	runner   cmd.Runner
	system   sys.Manager
	settings config.Settings
	dlg.Dialog
	tmpl.Writer
}

func NewSnippetCli(a snippets.Service, r cmd.Runner, s sys.Manager, d dlg.Dialog, w tmpl.Writer, t config.Settings, search search.Term) *SnippetCli {
	return &SnippetCli{service: a, runner: r, system: s, Dialog: d, Writer: w, settings: t, search: search}
}

func (a *SnippetCli) Share(fullKey string, destination string) {
	k := a.getAbsAlias(fullKey)
	if list, err := a.service.Get(k); err != nil {
		a.HandleErr(err)
	} else {
		if alias := a.handleMultiResponse(fullKey, list); alias != nil {
			gmail := &models.Snippet{Runtime: "url", Extension: "url"}
			gmail.Snip = "https://mail.google.com/mail/?ui=2&view=cm&fs=1&tf=1&su=&body=http%3A%2F%2Faus.kwk.co%2F" + alias.Username + "%2f" + alias.FullName
			a.runner.Run(gmail, []string{})
		} else {
			a.Render("snippet:notfound", map[string]interface{}{"fullKey": fullKey})
		}
	}
}

func (s *SnippetCli) SetEnv(fullName string) {
	a := s.getAbsAlias(fullName)
	if _, err := s.runner.SetEnv(a); err != nil {
		panic(err)
	} else {
		s.Render("env:changed", fullName)
	}
}

func (a *SnippetCli) Run(fullKey string, args []string) {
	k := a.getAbsAlias(fullKey)
	if list, err := a.service.Get(k); err != nil {
		a.HandleErr(err)
	} else {
		if alias := a.handleMultiResponse(fullKey, list); alias != nil {
			if err = a.runner.Run(alias, args); err != nil {
				fmt.Println(err)
			}
		} else {
			if res, err := a.search.Execute(fullKey); err != nil {
				a.HandleErr(err)
			} else if res.Total > 0 {
				a.Render("search:alphaSuggest", res)
				return
			}
			a.Render("snippet:notfound", map[string]interface{}{"FullKey": fullKey})
		}
	}
}

func (a *SnippetCli) New(uri string, fullKey string) {
	if createAlias, err := a.service.Create(uri, fullKey); err != nil {
		a.HandleErr(err)
	} else {
		if createAlias.Snippet != nil {
			if createAlias.Snippet.Private {
				a.Render("snippet:newprivate", createAlias.Snippet)
			} else {
				a.Render("snippet:new", createAlias.Snippet)
			}
			a.system.CopyToClipboard(createAlias.Snippet.FullName)
			//if createAlias.Snippet.Runtime != "url" {
			//	a.Edit(createAlias.Snippet.FullKey)
			//}
		} else {
			matches := createAlias.TypeMatch.Matches
			r := a.MultiChoice("snippet:chooseruntime", "Choose a type for this snippet:", matches)
			winner := r.Value.(models.Match)
			if winner.Score == -1 {
				ca, _ := a.service.Create("_", "_")
				matches = ca.TypeMatch.Matches
				winner = a.MultiChoice("snippet:chooseruntime", "Choose a type for this snippet:", matches).Value.(models.Match)
			}
			fk := fullKey + "." + winner.Extension
			a.New(uri, fk)
		}
	}
}

func (a *SnippetCli) Edit(fullKey string) {
	if list, err := a.service.Get(a.getAbsAlias(fullKey)); err != nil {
		a.HandleErr(err)
	} else {
		if alias := a.handleMultiResponse(fullKey, list); alias != nil {
			a.Render("snippet:editing", alias)
			if err := a.runner.Edit(alias); err != nil {
				a.HandleErr(err)
			} else {
				a.Render("snippet:edited", alias)
			}
		} else {
			a.Render("snippet:notfound", &models.Snippet{FullName: fullKey})
		}
	}
}

func (a *SnippetCli) Describe(fullKey string, description string) {
	if alias, err := a.service.Update(fullKey, description); err != nil {
		a.HandleErr(err)
	} else {
		a.Render("snippet:updated", alias)
	}
}

func (a *SnippetCli) Inspect(fullKey string) {
	if list, err := a.service.Get(a.getAbsAlias(fullKey)); err != nil {
		a.HandleErr(err)
	} else {
		a.Render("snippet:inspect", list)
	}
}

func (a *SnippetCli) Delete(fullKey string) {
	alias := &models.Snippet{FullName: fullKey}
	if r := a.Modal("snippet:delete", alias); r.Ok {
		if err := a.service.Delete(fullKey); err != nil {
			a.HandleErr(err)
		}
		a.Render("snippet:deleted", alias)
	} else {
		a.Render("snippet:notdeleted", alias)
	}
}

func (a *SnippetCli) Cat(fullKey string) {
	if list, err := a.service.Get(a.getAbsAlias(fullKey)); err != nil {
		a.HandleErr(err)
	} else {
		if len(list.Items) == 0 {
			a.Render("snippet:notfound", &models.Snippet{FullName: fullKey})
		} else if len(list.Items) == 1 {
			a.Render("snippet:cat", list.Items[0])
		} else {
			a.Render("snippet:ambiguouscat", list)
		}
	}
}

func (a *SnippetCli) Patch(fullKey string, target string, patch string) {
	if alias, err := a.service.Patch(fullKey, target, patch); err != nil {
		a.HandleErr(err)
	} else {
		a.Render("snippet:patched", alias)
	}
}

func (a *SnippetCli) Clone(fullKey string, newFullKey string) {
	if alias, err := a.service.Clone(a.getAbsAlias(fullKey), newFullKey); err != nil {
		a.HandleErr(err)
	} else {
		a.Render("snippet:cloned", alias)
	}
}

func (a *SnippetCli) Rename(fullKey string, newKey string) {
	if alias, originalFullKey, err := a.service.Rename(fullKey, newKey); err != nil {
		a.HandleErr(err)
	} else {
		if alias.Private {
			a.Render("snippet:madeprivate", &map[string]string{
				"fullKey": originalFullKey,
			})
		} else {
			a.Render("snippet:renamed", &map[string]string{
				"fullKey":    originalFullKey,
				"newFullKey": alias.FullName,
			})
		}
	}
}

func (a *SnippetCli) Tag(fullKey string, tags ...string) {
	if alias, err := a.service.Tag(fullKey, tags...); err != nil {
		a.HandleErr(err)
	} else {
		a.Render("snippet:tag", alias)
	}
}

func (a *SnippetCli) UnTag(fullKey string, tags ...string) {
	if alias, err := a.service.UnTag(fullKey, tags...); err != nil {
		a.HandleErr(err)
	} else {
		a.Render("snippet:untag", alias)
	}
}

func (a *SnippetCli) List(args ...string) {
	var size int64
	var tags = []string{}
	u := &models.User{}
	var username = ""
	for i, v := range args {
		if num, err := strconv.Atoi(v); err == nil {
			size = int64(num)
		} else {
			if i == 0 && v[0] == '/' {
				username = strings.Replace(v, "/", "", -1)
			} else {
				tags = append(tags, v)
			}
		}
	}
	if username == "" {
		if err := a.settings.Get(models.ProfileFullKey, u, 0); err != nil {
			a.Render("api:not-authenticated", nil)
			return
		} else {
			username = u.Username
		}
	}
	if list, err := a.service.List(username, size, int64(time.Now().Unix()*1000), tags...); err != nil {
		a.HandleErr(err)
	} else {
		list.Username = username
		a.Render("snippet:list", list)
	}
}

func (a *SnippetCli) handleMultiResponse(fullKey string, list *models.SnippetList) *models.Snippet {
	if list.Total == 1 {
		return &list.Items[0]
	} else if list.Total > 1 {
		r := a.MultiChoice("dialog:choose", "Multiple matches. Choose a snippet to run:", list.Items)
		s := r.Value.(models.Snippet)
		return &s
	} else {
		return nil
	}
}

func (a *SnippetCli) getAbsAlias(fullKey string) *models.Alias {
	u := &models.User{}
	if err := a.settings.Get(models.ProfileFullKey, u, 0); err != nil {
		a.Render("api:not-authenticated", nil)
		return nil
	}
	k := &models.Alias{}
	k.Username = u.Username
	// TODO: MOVE LOGIC SERVER SIDE
	if strings.Contains(fullKey, "/") {
		tokens := strings.Split(fullKey, "/")
		k.Username = tokens[0]
		k.FullKey = tokens[1]
	} else {
		k.FullKey = fullKey
	}
	return k
}
