package app

import (
	"strconv"
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"strings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/search"
	"time"
)

type AliasController struct {
	search   search.ISearch
	service  aliases.IAliases
	openers  openers.IOpen
	system   system.ISystem
	settings settings.ISettings
	gui.IDialogues
	gui.ITemplateWriter
}

func NewAliasController(a aliases.IAliases, o openers.IOpen, s system.ISystem, d gui.IDialogues, w gui.ITemplateWriter, t settings.ISettings, search search.ISearch) *AliasController {
	return &AliasController{service:a, openers:o, system:s, IDialogues:d, ITemplateWriter:w, settings: t, search:search}
}

func (a *AliasController) Share(fullKey string, destination string) {
	k := a.getKwkKey(fullKey);
	if list, err := a.service.Get(k); err != nil {
		a.Render("error", err)
	} else {
		if alias := a.handleMultiResponse(fullKey, list); alias != nil {
			gmail := &models.Alias{Runtime:"url", Extension:"url"}
			gmail.Snip = "https://mail.google.com/mail/?ui=2&view=cm&fs=1&tf=1&su=&body=http%3A%2F%2Faus.kwk.co%2F" + alias.Username + "%2f" + alias.FullKey
			a.openers.Open(gmail, []string{})
		} else {
			a.Render("alias:notfound", map[string]interface{}{"fullKey":fullKey})
		}
	}
}

func (a *AliasController) Open(fullKey string, args []string) {
	k := a.getKwkKey(fullKey);
	if list, err := a.service.Get(k); err != nil {
		a.Render("error", err)
	} else {
		if alias := a.handleMultiResponse(fullKey, list); alias != nil {
			a.openers.Open(alias, args)
		} else {
			if res, err := a.search.Search(fullKey); err != nil {
				a.Render("error", err)
			} else if res.Total > 0 {
				a.Render("search:beta", res)
				return
			}
			a.Render("alias:notfound", map[string]interface{}{"FullKey":fullKey})
		}
	}
}

func (a *AliasController) New(uri string, fullKey string) {
	if createAlias, err := a.service.Create(uri, fullKey); err != nil {
		a.Render("error", err)
	} else {
		if createAlias.Alias != nil {
			if createAlias.Alias.Private {
				a.Render("alias:newprivate", createAlias.Alias)
			} else {
				a.Render("alias:new", createAlias.Alias)
			}
				a.system.CopyToClipboard(createAlias.Alias.FullKey)
			//if createAlias.Alias.Runtime != "url" {
			//	a.Edit(createAlias.Alias.FullKey)
			//}
		} else {
			matches := createAlias.TypeMatch.Matches
			r := a.MultiChoice("alias:chooseruntime", "Choose a type for this alias:", matches)
			winner := r.Value.(models.Match)
			if winner.Score == -1 {
				ca, _ := a.service.Create("_", "_")
				matches = ca.TypeMatch.Matches
				winner = a.MultiChoice("alias:chooseruntime", "Choose a type for this alias:", matches).Value.(models.Match)
			}
			fk := fullKey + "." + winner.Extension
			a.New(uri, fk)
		}
	}
}

func (a *AliasController) Edit(fullKey string) {
	if list, err := a.service.Get(a.getKwkKey(fullKey)); err != nil {
		a.Render("error", err)
	} else {
		if alias := a.handleMultiResponse(fullKey, list); alias != nil {
			a.Render("alias:editing", alias)
			if err := a.openers.Edit(alias); err != nil {
				a.Render("error", err)
			} else {
				a.Render("alias:edited", alias)
			}
		} else {
			a.Render("alias:notfound", &models.Alias{FullKey:fullKey})
		}
	}
}

func (a *AliasController) Describe(fullKey string, description string) {
	if alias, err := a.service.Update(fullKey, description); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:updated", alias)
	}
}

func (a *AliasController) Inspect(fullKey string) {
	if list, err := a.service.Get(a.getKwkKey(fullKey)); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:inspect", list)
	}
}

func (a *AliasController) Delete(fullKey string) {
	alias := &models.Alias{FullKey:fullKey}
	if r := a.Modal("alias:delete", alias); r.Ok {
		if err := a.service.Delete(fullKey); err != nil {
			a.Render("error", err)
		}
		a.Render("alias:deleted", alias)
	} else {
		a.Render("alias:notdeleted", alias)
	}
}

func (a *AliasController) Cat(fullKey string) {
	if list, err := a.service.Get(a.getKwkKey(fullKey)); err != nil {
		a.Render("error", err)
	} else {
		if len(list.Items) == 0 {
			a.Render("alias:notfound", &models.Alias{FullKey:fullKey})
		} else if (len(list.Items) == 1) {
			a.Render("alias:cat", list.Items[0])
		} else {
			a.Render("alias:ambiguouscat", list)
		}
	}
}

func (a *AliasController) Patch(fullKey string, target string, patch string) {
	if alias, err := a.service.Patch(fullKey, target, patch); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:patched", alias)
	}
}

func (a *AliasController) Clone(fullKey string, newFullKey string) {
	if alias, err := a.service.Clone(a.getKwkKey(fullKey), newFullKey); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:cloned", alias)
	}
}

func (a *AliasController) Rename(fullKey string, newKey string) {
	if alias, originalFullKey, err := a.service.Rename(fullKey, newKey); err != nil {
		a.Render("error", err)
	} else {
		if alias.Private {
			a.Render("alias:madeprivate", &map[string]string{
				"fullKey":originalFullKey,
			})
		} else {
			a.Render("alias:renamed", &map[string]string{
				"fullKey":originalFullKey,
				"newFullKey":alias.FullKey,
			})
		}
	}
}

func (a *AliasController) Tag(fullKey string, tags ...string) {
	if alias, err := a.service.Tag(fullKey, tags...); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:tag", alias)
	}
}

func (a *AliasController) UnTag(fullKey string, tags ...string) {
	if alias, err := a.service.UnTag(fullKey, tags...); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:untag", alias)
	}
}

func (a *AliasController) List(args ...string) {
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
		if err := a.settings.Get(models.ProfileFullKey, u); err != nil {
			a.Render("account:notloggedin", nil)
			return
		} else {
			username = u.Username
		}
	}
	if list, err := a.service.List(username, size, int64(time.Now().Unix()*1000), tags...); err != nil {
		a.Render("error", err)
	} else {
		list.Username = username
		a.Render("alias:list", list)
	}
}

func (a *AliasController) handleMultiResponse(fullKey string, list *models.AliasList) *models.Alias {
	if list.Total == 1 {
		return &list.Items[0]
	} else if list.Total > 1 {
		r := a.MultiChoice("alias:choose", nil, list.Items)
		return r.Value.(*models.Alias)
	} else {
		return nil
	}
}

func (a *AliasController) getKwkKey(fullKey string) *models.KwkKey {
	u := &models.User{}
	if err := a.settings.Get(models.ProfileFullKey, u); err != nil {
		a.Render("account:notloggedin", nil)
		return nil
	}
	k := &models.KwkKey{}
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

