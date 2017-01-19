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

func (s *SnippetCli) Share(fullKey string, destination string) {
	k := s.getAbsAlias(fullKey)
	if list, err := s.service.Get(k); err != nil {
		s.HandleErr(err)
	} else {
		if alias := s.handleMultiResponse(fullKey, list); alias != nil {
			gmail := &models.Snippet{Runtime: "url", Extension: "url"}
			gmail.Snip = "https://mail.google.com/mail/?ui=2&view=cm&fs=1&tf=1&su=&body=http%3A%2F%2Faus.kwk.co%2F" + alias.Username + "%2f" + alias.FullName
			s.runner.Run(gmail, []string{})
		} else {
			s.Render("snippet:notfound", map[string]interface{}{"fullKey": fullKey})
		}
	}
}

func (s *SnippetCli) Run(fullKey string, args []string) {
	k := s.getAbsAlias(fullKey)
	if list, err := s.service.Get(k); err != nil {
		s.HandleErr(err)
	} else {
		if alias := s.handleMultiResponse(fullKey, list); alias != nil {
			if err = s.runner.Run(alias, args); err != nil {
				fmt.Println(err)
			}
		} else {
			if res, err := s.search.Execute(fullKey); err != nil {
				s.HandleErr(err)
			} else if res.Total > 0 {
				s.Render("search:alphaSuggest", res)
				return
			}
			s.Render("snippet:notfound", map[string]interface{}{"FullKey": fullKey})
		}
	}
}

func (s *SnippetCli) New(uri string, fullKey string) {
	if createAlias, err := s.service.Create(uri, fullKey); err != nil {
		s.HandleErr(err)
	} else {
		if createAlias.Snippet != nil {
			if createAlias.Snippet.Private {
				s.Render("snippet:newprivate", createAlias.Snippet)
			} else {
				s.Render("snippet:new", createAlias.Snippet)
			}
		} else {
			matches := createAlias.TypeMatch.Matches
			r := s.MultiChoice("snippet:chooseruntime", "Choose a type for this snippet:", matches)
			winner := r.Value.(models.Match)
			if winner.Score == -1 {
				ca, _ := s.service.Create("_", "_")
				matches = ca.TypeMatch.Matches
				winner = s.MultiChoice("snippet:chooseruntime", "Choose a type for this snippet:", matches).Value.(models.Match)
			}
			fk := fullKey + "." + winner.Extension
			s.New(uri, fk)
		}
	}
}

func (s *SnippetCli) Edit(fullKey string) {
	// TODO: ENCAPSULATE
	if fullKey == "env" || fullKey == "prefs" {
		fullKey += ".yml"
	}
	if fullKey == "env.yml" || fullKey == "prefs.yml" {
		fullKey = models.GetHostConfigName(fullKey)
	}

	if list, err := s.service.Get(s.getAbsAlias(fullKey)); err != nil {
		s.HandleErr(err)
	} else {
		if alias := s.handleMultiResponse(fullKey, list); alias != nil {
			s.Render("snippet:editing", alias)
			if err := s.runner.Edit(alias); err != nil {
				s.HandleErr(err)
			} else {
				s.Render("snippet:edited", alias)
			}
		} else {
			s.Render("snippet:notfound", &models.Snippet{FullName: fullKey})
		}
	}
}

func (s *SnippetCli) Describe(fullKey string, description string) {
	if alias, err := s.service.Update(fullKey, description); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:updated", alias)
	}
}

func (s *SnippetCli) Inspect(fullKey string) {
	if list, err := s.service.Get(s.getAbsAlias(fullKey)); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:inspect", list)
	}
}

func (s *SnippetCli) Delete(fullKey string) {
	alias := &models.Snippet{FullName: fullKey}
	if r := s.Modal("snippet:delete", alias, s.settings.GetPrefs().AutoYes); r.Ok {
		if err := s.service.Delete(fullKey); err != nil {
			s.HandleErr(err)
		}
		s.Render("snippet:deleted", alias)
	} else {
		s.Render("snippet:notdeleted", alias)
	}
}

func (s *SnippetCli) Cat(fullKey string) {
	if list, err := s.service.Get(s.getAbsAlias(fullKey)); err != nil {
		s.HandleErr(err)
	} else {
		if len(list.Items) == 0 {
			s.Render("snippet:notfound", &models.Snippet{FullName: fullKey})
		} else if len(list.Items) == 1 {
			s.Render("snippet:cat", list.Items[0])
		} else {
			s.Render("snippet:ambiguouscat", list)
		}
	}
}

func (s *SnippetCli) Patch(fullKey string, target string, patch string) {
	if alias, err := s.service.Patch(fullKey, target, patch); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:patched", alias)
	}
}

func (s *SnippetCli) Clone(fullKey string, newFullKey string) {
	if alias, err := s.service.Clone(s.getAbsAlias(fullKey), newFullKey); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:cloned", alias)
	}
}

func (s *SnippetCli) Rename(fullKey string, newKey string) {
	if alias, originalFullKey, err := s.service.Rename(fullKey, newKey); err != nil {
		s.HandleErr(err)
	} else {
		if alias.Private {
			s.Render("snippet:madeprivate", &map[string]string{
				"fullKey": originalFullKey,
			})
		} else {
			s.Render("snippet:renamed", &map[string]string{
				"fullKey":    originalFullKey,
				"newFullKey": alias.FullName,
			})
		}
	}
}

func (s *SnippetCli) Tag(fullKey string, tags ...string) {
	if alias, err := s.service.Tag(fullKey, tags...); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:tag", alias)
	}
}

func (s *SnippetCli) UnTag(fullKey string, tags ...string) {
	if alias, err := s.service.UnTag(fullKey, tags...); err != nil {
		s.HandleErr(err)
	} else {
		s.Render("snippet:untag", alias)
	}
}

func (s *SnippetCli) List(args ...string) {
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
		if err := s.settings.Get(models.ProfileFullKey, u, 0); err != nil {
			s.Render("api:not-authenticated", nil)
			return
		} else {
			username = u.Username
		}
	}
	if list, err := s.service.List(username, size, int64(time.Now().Unix()*1000), tags...); err != nil {
		s.HandleErr(err)
	} else {
		list.Username = username
		s.Render("snippet:list", list)
	}
}

func (s *SnippetCli) handleMultiResponse(fullKey string, list *models.SnippetList) *models.Snippet {
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

func (s *SnippetCli) getAbsAlias(fullKey string) *models.Alias {
	u := &models.User{}
	if err := s.settings.Get(models.ProfileFullKey, u, 0); err != nil {
		s.Render("api:not-authenticated", nil)
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
