package aliases

import (
	"strconv"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkweb/rpc/aliasesRpc"
	"github.com/kwk-links/kwk-cli/libs/rpc"
	"github.com/kwk-links/kwk-cli/libs/models"
	"time"
)

const TimeLayout = time.RFC3339

type Aliases struct {
	Settings settings.ISettings
	client   aliasesRpc.AliasesRpcClient
	rpc.Headers
}

func New(s settings.ISettings) IAliases {
	return &Aliases{Settings:s}
}

func (a *Aliases) List(args []string) (*models.AliasList, error) {
	var page, size int
	var tags = []string{}
	for _, v := range args {
		if num, err := strconv.Atoi(v); err == nil {
			if page == 0 {
				page = num
			} else {
				size = num
			}
		} else {
			tags = append(tags, v)
		}
	}
	if res, err := a.client.List(a.GetContext(), &aliasesRpc.ListRequest{Username:"richard", Page:page, Size:size, Tags:tags}); err != nil {
		return nil, err
	} else {
		list := &models.AliasList{}
		mapAliasList(res, list)
		return list, nil
	}
}

func (a *Aliases) Get(fullKey string) (*models.AliasList, error) {

	if res, err := a.client.Get(a.GetContext(), &aliasesRpc.GetRequest{FullKey:fullKey}); err != nil {
		return nil, err
	} else {
		list := &models.AliasList{}
		mapAliasList(res, list)
		return list
	}
}

func (a *Aliases) Delete(fullKey string) error {
	_, err := a.client.Delete(a.GetContext(), &aliasesRpc.DeleteRequest{FullKey:fullKey})
	return err
}

func (a *Aliases) Create(uri string, path string) (*models.Alias, error) {
	if res, err := a.client.Create(a.GetContext(), &aliasesRpc.CreateRequest{Uri:uri, FullKey:path}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res, alias)
		return alias, nil
	}
}

func (a *Aliases) Rename(fullKey string, newFullKey string) (*models.Alias, error) {
	if res, err := a.client.Rename(a.GetContext(), &aliasesRpc.RenameRequest{FullKey:fullKey, NewFullKey:newFullKey}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res, alias)
		return alias, nil
	}
}

func (a *Aliases) Patch(fullKey string, uri string) (*models.Alias, error) {
	if res, err := a.client.Patch(a.GetContext(), &aliasesRpc.PatchRequest{FullKey:fullKey, Uri:uri}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res, alias)
		return alias, nil
	}
}

func (a *Aliases) Clone(fullKey string, newFullKey string) *models.Alias {
	if res, err := a.client.Clone(a.GetContext(), &aliasesRpc.CloneRequest{FullKey:fullKey, NewFullKey:newFullKey}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res, alias)
		return alias, nil
	}
}

func (a *Aliases) Tag(fullKey string, tags ...string) *models.Alias {
	if res, err := a.client.Tag(a.GetContext(), &aliasesRpc.TagRequest{FullKey:fullKey, Tags:tags}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res, alias)
		return alias, nil
	}
}

func (a *Aliases) UnTag(fullKey string, tags ...string) *models.Alias {
	if res, err := a.client.UnTag(a.GetContext(), &aliasesRpc.UnTagRequest{FullKey:fullKey, Tags:tags}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res, alias)
		return alias, nil
	}
}

func mapAlias(rpc *aliasesRpc.AliasResponse, model *models.Alias) {
	model.Id = rpc.Id
	model.FullKey = rpc.FullKey
	model.Username = rpc.Username
	model.Key = rpc.Key
	model.Extension = rpc.Extension
	model.Uri = rpc.Uri
	model.Version = rpc.Version
	model.Media = rpc.Media
	model.Runtime = rpc.Runtime
	model.Tags = rpc.Tags
	created, _ := time.Parse(TimeLayout, rpc.CreatedUTC)
	model.Created = created
	updated, _ := time.Parse(TimeLayout, rpc.UpdatedUTC)
	model.Updated = updated
}

func mapAliasList(rpc *aliasesRpc.AliasListResponse, model *models.AliasList) {
	model.Total = rpc.Total
	model.Page = rpc.Page
	for _, v := range rpc.Items {
		item := &models.Alias{}
		mapAlias(&v, item)
		model.Items = append(model.Items, item)
	}
	model.Size = rpc.Size
}

