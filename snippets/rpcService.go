package snippets

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/rpc/src/snipsRpc"
	"google.golang.org/grpc"
	"time"
	"fmt"
)

type RpcService struct {
	Settings config.Settings
	client   snipsRpc.SnipsRpcClient
	headers  *rpc.Headers
}

func New(conn *grpc.ClientConn, s config.Settings, h *rpc.Headers) Service {
	return &RpcService{Settings: s, client: snipsRpc.NewSnipsRpcClient(conn), headers: h}
}

func (a *RpcService) Update(fullKey string, description string) (*models.Snippet, error) {
	if r, err := a.client.Update(a.headers.GetContext(), &snipsRpc.UpdateRequest{FullName: fullKey, Description: description}); err != nil {
		return nil, err
	} else {
		m := &models.Snippet{}
		mapSnip(r.Snip, m)
		return m, nil
	}
}

// since unix time in milliseconds
func (a *RpcService) List(username string, size int64, since int64, tags ...string) (*models.SnippetList, error) {
	if res, err := a.client.List(a.headers.GetContext(), &snipsRpc.ListRequest{Username: username, Since: since, Size: size, Tags: tags}); err != nil {
		return nil, err
	} else {
		list := &models.SnippetList{}
		mapSnippetList(res, list)
		return list, nil
	}
}

func (a *RpcService) Get(k *models.Alias) (*models.SnippetList, error) {
	if res, err := a.client.Get(a.headers.GetContext(), &snipsRpc.GetRequest{Username: k.Username, FullName: k.FullKey}); err != nil {
		return nil, err
	} else {
		list := &models.SnippetList{}
		mapSnippetList(res, list)
		return list, nil
	}
}

func (a *RpcService) Delete(fullKey string) error {
	_, err := a.client.Delete(a.headers.GetContext(), &snipsRpc.DeleteRequest{FullName: fullKey})
	return err
}

func (a *RpcService) Create(snip string, path string) (*models.CreateSnippetRequest, error) {
	// encrypt if requested
	fmt.Println(snip, path)
	if res, err := a.client.Create(a.headers.GetContext(), &snipsRpc.CreateRequest{Snip: snip, FullName: path}); err != nil {
		return nil, err
	} else {
		cs := &models.CreateSnippetRequest{}
		if res.Snip != nil {
			snip := &models.Snippet{}
			mapSnip(res.Snip, snip)
			cs.Snippet = snip
		} else {
			cs.TypeMatch = &models.TypeMatch{
				Matches: []models.Match{},
			}
			for _, v := range res.TypeMatch.Matches {
				m := models.Match{
					Extension: v.Extension,
					Score:     v.Score,
				}
				cs.TypeMatch.Matches = append(cs.TypeMatch.Matches, m)
			}
		}
		return cs, nil
	}
}

func (a *RpcService) Rename(fullKey string, newFullName string) (*models.Snippet, string, error) {
	if res, err := a.client.Rename(a.headers.GetContext(), &snipsRpc.RenameRequest{FullName: fullKey, NewFullName: newFullName}); err != nil {
		return nil, "", err
	} else {
		snip := &models.Snippet{}
		mapSnip(res.Snip, snip)
		return snip, res.OriginalFullName, nil
	}
}

func (a *RpcService) Patch(fullKey string, target string, patch string) (*models.Snippet, error) {
	if res, err := a.client.Patch(a.headers.GetContext(), &snipsRpc.PatchRequest{FullName: fullKey, Target: target, Patch: patch}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		mapSnip(res.Snip, snip)
		return snip, nil
	}
}

func (a *RpcService) Clone(k *models.Alias, newFullName string) (*models.Snippet, error) {
	if res, err := a.client.Clone(a.headers.GetContext(), &snipsRpc.CloneRequest{Username: k.Username, FullName: k.FullKey, NewFullName: newFullName}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		mapSnip(res.Snip, snip)
		return snip, nil
	}
}

func (a *RpcService) Tag(fullKey string, tags ...string) (*models.Snippet, error) {
	if res, err := a.client.Tag(a.headers.GetContext(), &snipsRpc.TagRequest{FullName: fullKey, Tags: tags}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		mapSnip(res.Snip, snip)
		return snip, nil
	}
}

func (a *RpcService) UnTag(fullKey string, tags ...string) (*models.Snippet, error) {
	if res, err := a.client.UnTag(a.headers.GetContext(), &snipsRpc.UnTagRequest{FullName: fullKey, Tags: tags}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		mapSnip(res.Snip, snip)
		return snip, nil
	}
}

func mapSnip(rpc *snipsRpc.Snip, model *models.Snippet) {
	model.Id = rpc.SnipId
	model.FullName = rpc.FullName
	model.Username = rpc.Username
	model.Name = rpc.Name
	model.Extension = rpc.Extension
	// if encrypted, decrypt
	// if checksum doesn't match then throw warning
	// check that checksum signature is valid with public key.
	model.Snip = rpc.Snip
	model.Version = rpc.SnipVersion
	model.Tags = rpc.Tags
	model.Created = time.Unix(rpc.Created/1000, 0)
	model.Description = rpc.Description
	model.ClonedFromFullName = rpc.ClonedFromFullName
	model.ClonedFromVersion = rpc.ClonedFromVersion
	model.Private = rpc.Private
	model.RunCount = rpc.RunCount
	model.CloneCount = rpc.CloneCount
}

func mapSnippetList(rpc *snipsRpc.ListResponse, model *models.SnippetList) {
	model.Total = rpc.Total
	model.Since = time.Unix(rpc.Since/1000, 0)
	model.Size = rpc.Size
	for _, v := range rpc.Items {
		item := &models.Snippet{}
		mapSnip(v, item)
		model.Items = append(model.Items, *item)
	}
}
