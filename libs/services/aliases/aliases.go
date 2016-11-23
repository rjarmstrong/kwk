package aliases

import (
	"time"
	"google.golang.org/grpc"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/models"
	"bitbucket.com/sharingmachine/kwkcli/libs/rpc"
	"bitbucket.com/sharingmachine/rpc/src/aliasesRpc"
)

const TimeLayout = time.RFC3339

type Aliases struct {
	Settings settings.ISettings
	client   aliasesRpc.AliasesRpcClient
	headers  *rpc.Headers
}

func New(conn *grpc.ClientConn, s settings.ISettings, h *rpc.Headers) IAliases {
	return &Aliases{Settings:s, client:aliasesRpc.NewAliasesRpcClient(conn), headers:h}
}

func (a *Aliases) Update(fullKey string, description string) (*models.Alias, error) {
	if r, err := a.client.Update(a.headers.GetContext(), &aliasesRpc.UpdateRequest{FullKey:fullKey, Description:description}); err != nil {
		return nil, err
	} else {
		m := &models.Alias{}
		mapAlias(r.Alias, m)
		return m, nil
	}
}

func (a *Aliases) List(username string, page int64, size int64, tags ...string) (*models.AliasList, error) {
	if res, err := a.client.List(a.headers.GetContext(), &aliasesRpc.ListRequest{Username:username, Page:page, Size:size, Tags:tags}); err != nil {
		return nil, err
	} else {
		list := &models.AliasList{}
		mapAliasList(res, list)
		return list, nil
	}
}

func (a *Aliases) Get(k *models.KwkKey) (*models.AliasList, error) {
	if res, err := a.client.Get(a.headers.GetContext(), &aliasesRpc.GetRequest{Username: k.Username, FullKey: k.FullKey}); err != nil {
		return nil, err
	} else {
		list := &models.AliasList{}
		mapAliasList(res, list)
		return list, nil
	}
}

func (a *Aliases) Delete(fullKey string) error {
	_, err := a.client.Delete(a.headers.GetContext(), &aliasesRpc.DeleteRequest{FullKey:fullKey})
	return err
}

func (a *Aliases) Create(uri string, path string) (*models.CreateAlias, error) {
	if res, err := a.client.Create(a.headers.GetContext(), &aliasesRpc.CreateRequest{Uri:uri, FullKey:path}); err != nil {
		return nil, err
	} else {
		createAlias := &models.CreateAlias{}
		if res.Alias != nil {
			alias := &models.Alias{}
			mapAlias(res.Alias, alias)
			createAlias.Alias = alias
		} else {
			createAlias.TypeMatch = &models.TypeMatch{
				Matches: []models.Match{},
			}
			for _, v := range res.TypeMatch.Matches {
				m := models.Match{
					Extension:v.Extension,
					Media:v.Media,
					Runtime:v.Runtime,
					Score:v.Score,
				}
				createAlias.TypeMatch.Matches = append(createAlias.TypeMatch.Matches, m)
			}
		}
		return createAlias, nil
	}
}

func (a *Aliases) Rename(fullKey string, newFullKey string) (*models.Alias, string, error) {
	if res, err := a.client.Rename(a.headers.GetContext(), &aliasesRpc.RenameRequest{FullKey:fullKey, NewFullKey:newFullKey}); err != nil {
		return nil, "", err
	} else {
		alias := &models.Alias{}
		mapAlias(res.Alias, alias)
		return alias, res.OriginalFullKey, nil
	}
}

func (a *Aliases) Patch(fullKey string, target string, patch string) (*models.Alias, error) {
	if res, err := a.client.Patch(a.headers.GetContext(), &aliasesRpc.PatchRequest{FullKey:fullKey, Target:target, Patch:patch}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res.Alias, alias)
		return alias, nil
	}
}

func (a *Aliases) Clone(k *models.KwkKey, newFullKey string) (*models.Alias, error) {
	if res, err := a.client.Clone(a.headers.GetContext(), &aliasesRpc.CloneRequest{Username: k.Username, FullKey:k.FullKey, NewFullKey:newFullKey}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res.Alias, alias)
		return alias, nil
	}
}

func (a *Aliases) Tag(fullKey string, tags ...string) (*models.Alias, error) {
	if res, err := a.client.Tag(a.headers.GetContext(), &aliasesRpc.TagRequest{FullKey:fullKey, Tags:tags}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res.Alias, alias)
		return alias, nil
	}
}

func (a *Aliases) UnTag(fullKey string, tags ...string) (*models.Alias, error) {
	if res, err := a.client.UnTag(a.headers.GetContext(), &aliasesRpc.UnTagRequest{FullKey:fullKey, Tags:tags}); err != nil {
		return nil, err
	} else {
		alias := &models.Alias{}
		mapAlias(res.Alias, alias)
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

	model.Description = rpc.Description
	model.ForkedFromFullKey = rpc.ForkedFromFullKey
	model.ForkedFromVersion  = rpc.ForkedFromVersion
	model.Private	 = rpc.Private
	model.RunCount  = rpc.RunCount
	model.ForkCount  = rpc.ForkCount
}

func mapAliasList(rpc *aliasesRpc.AliasListResponse, model *models.AliasList) {
	model.Total = rpc.Total
	model.Page = rpc.Page
	for _, v := range rpc.Items {
		item := &models.Alias{}
		mapAlias(v, item)
		model.Items = append(model.Items, *item)
	}
	model.Size = rpc.Size
}

