package app

import (
	"strconv"
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/gui"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
)

type AliasController struct {
	service aliases.IAliases
	openers openers.IOpen
	system  system.ISystem
	gui.IDialogues
	gui.ITemplateWriter
}

func NewAliasController(a aliases.IAliases, o openers.IOpen, s system.ISystem, d gui.IDialogues, w gui.ITemplateWriter) *AliasController {
	return &AliasController{service:a, openers:o, system:s, IDialogues:d, ITemplateWriter:w}
}

func (a *AliasController) Open(fullKey string, args []string) {
	if list, err := a.service.Get(fullKey); err != nil {
		a.Render("error", err)
	} else {
		if alias := a.handleMultiResponse(fullKey, list); alias != nil {
			a.openers.Open(alias, args)
		} else {
			a.Render("alias:notfound", map[string]interface{}{"fullKey":fullKey})
		}
	}
}

func (a *AliasController) New(uri string, fullKey string) {
	if createAlias, err := a.service.Create(uri, fullKey); err != nil {
		a.Render("error", err)
	} else {
		if createAlias.Alias != nil {
			a.Render("alias:new", createAlias.Alias)
			a.system.CopyToClipboard(createAlias.Alias.FullKey)
		} else {
			r := a.MultiChoice("alias:chooseruntime", "Choose the correct type for this link:", createAlias.TypeMatch.Matches)
			winner := r.Value.(models.Match)
			a.New(uri, fullKey + "." + winner.Extension)
		}
	}
}

func (a *AliasController) Edit(fullKey string) {
	if list, err := a.service.Get(fullKey); err != nil {
		a.Render("error", err)
	} else {
		if alias := a.handleMultiResponse(fullKey, list); alias != nil {
			if err := a.openers.Edit(alias); err != nil {
				a.Render("error", err)
			} else {
				a.Render("alias:edited", alias)
			}
		} else {
			a.Render("alias:notfound", map[string]interface{}{"fullKey":fullKey})
		}
	}
}

func (a *AliasController) Inspect(fullKey string) {
	if list, err := a.service.Get(fullKey); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:inspect", list)
	}
}

func (a *AliasController) Delete(fullKey string) {
	data := map[string]string{"fullKey" : fullKey}
	if r := a.Modal("alias:delete", data); r.Ok {
		if err := a.service.Delete(fullKey); err != nil {
			a.Render("error", err)
		}
		a.Render("alias:deleted", data)
	} else {
		a.Render("alias:notdeleted", map[string]string{"fullKey" : fullKey})
	}
}

func (a *AliasController) Cat(fullKey string) {
	if list, err := a.service.Get(fullKey); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:cat", list)
	}
}

func (a *AliasController) Patch(fullKey string, uri string) {
	if alias, err := a.service.Patch(fullKey, uri); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:patched", alias)
	}
}

func (a *AliasController) Clone(fullKey string, newFullKey string) {
	if alias, err := a.service.Clone(fullKey, newFullKey); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:cloned", alias)
	}
}

func (a *AliasController) Rename(fullKey string, newKey string) {
	if alias, err := a.service.Rename(fullKey, newKey); err != nil {
		a.Render("error", err)
	} else {
		a.Render("alias:renamed", alias)
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
	var page, size int32
	var tags = []string{}
	for _, v := range args {
		if num, err := strconv.Atoi(v); err == nil {
			if page == 0 {
				page = int32(num)
			} else {
				size = int32(num)
			}
		} else {
			tags = append(tags, v)
		}
	}

	if list, err := a.service.List("richard", page, size, tags...); err != nil {
		a.Render("error", err)
	} else {
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

