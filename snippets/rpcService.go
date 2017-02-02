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
	persister config.Persister
	client    snipsRpc.SnipsRpcClient
	headers   *rpc.Headers
}

func New(conn *grpc.ClientConn, s config.Persister, h *rpc.Headers) Service {
	return &RpcService{persister: s, client: snipsRpc.NewSnipsRpcClient(conn), headers: h}
}

func (r *RpcService) Update(a models.Alias, description string) (*models.Snippet, error) {
	if res, err := r.client.Update(r.headers.GetContext(), &snipsRpc.UpdateRequest{Alias: mapAlias(a), Description: description}); err != nil {
		return nil, err
	} else {
		m := &models.Snippet{}
		r.mapSnip(res.Snip, m, true)
		return m, nil
	}
}

// since unix time in milliseconds
func (r *RpcService) List(l *models.ListParams) (*models.SnippetList, error) {
	if res, err := r.client.List(r.headers.GetContext(), &snipsRpc.ListRequest{Username: l.Username, Since: l.Since, Size: l.Size, Tags: l.Tags, All:l.All}); err != nil {
		return nil, err
	} else {
		list := &models.SnippetList{}
		r.mapSnippetList(res, list, true)
		return list, nil
	}
}

func (r *RpcService) Get(k models.Alias) (*models.SnippetList, error) {
	if res, err := r.client.Get(r.headers.GetContext(), &snipsRpc.GetRequest{Alias:mapAlias(k)}); err != nil {
		return nil, err
	} else {
		list := &models.SnippetList{}
		r.mapSnippetList(res, list, false)
		return list, nil
	}
}

func (r *RpcService) Delete(username string, pouch string, names []*models.SnipName) error {
	sn := []*snipsRpc.SnipName{}
	for _, v := range names {
		sn = append(sn, &snipsRpc.SnipName{Name: v.Name, Extension: v.Ext })
	}
	_, err := r.client.Delete(r.headers.GetContext(), &snipsRpc.DeleteRequest{Username:username, Pouch:pouch, Names:sn})
	//if err != nil {
	//	r.Settings.Upsert(DELETED_SNIPPET, )
	//}
	return err
}

func (sm *RpcService) Move(username string, sourcePouch string, targetPouch string, names []*models.SnipName) (string, error) {
	panic("not imp")
}

func (r *RpcService) Create(snip string, a models.Alias, role models.SnipRole) (*models.CreateSnippetResponse, error) {
	if res, err := r.client.Create(r.headers.GetContext(), &snipsRpc.CreateRequest{Snip: snip, Alias: mapAlias(a), Role: snipsRpc.Role(role)}); err != nil {
		return nil, err
	} else {
		cs := &models.CreateSnippetResponse{}
		if res.Snip != nil {
			snip := &models.Snippet{}
			r.mapSnip(res.Snip, snip, true)
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

func (r *RpcService) Rename(a models.Alias, new models.SnipName) (*models.Snippet, *models.SnipName, error) {
	if res, err := r.client.Rename(r.headers.GetContext(), &snipsRpc.RenameRequest{Alias:mapAlias(a), NewName: &snipsRpc.SnipName{Name:new.Name, Extension:new.Ext}}); err != nil {
		return nil, nil, err
	} else {
		snip := &models.Snippet{}
		r.mapSnip(res.Snip, snip, true)
		return snip, &models.SnipName{Name:res.Original.Name, Ext:res.Original.Extension}, nil
	}
}

func (r *RpcService) Patch(a models.Alias, target string, patch string) (*models.Snippet, error) {
	if res, err := r.client.Patch(r.headers.GetContext(), &snipsRpc.PatchRequest{Alias:mapAlias(a), Target: target, Patch: patch}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		r.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (r *RpcService) Clone(a models.Alias, new models.Alias) (*models.Snippet, error) {
	if res, err := r.client.Clone(r.headers.GetContext(), &snipsRpc.CloneRequest{Alias:mapAlias(a), New:mapAlias(new)}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		r.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (r *RpcService) Tag(a models.Alias, tags ...string) (*models.Snippet, error) {
	if res, err := r.client.Tag(r.headers.GetContext(), &snipsRpc.TagRequest{Alias: mapAlias(a), Tags: tags}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		r.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (r *RpcService) UnTag(a models.Alias, tags ...string) (*models.Snippet, error) {
	if res, err := r.client.UnTag(r.headers.GetContext(), &snipsRpc.UnTagRequest{Alias: mapAlias(a), Tags: tags}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		r.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

/*
  cache will add the snippet to the local cache to deal with eventual consistency user experience.
 */
func (r *RpcService) mapSnip(rpc *snipsRpc.Snip, model *models.Snippet, cache bool) {
	model.Id = rpc.SnipId
	model.Alias = models.Alias{
		Username:rpc.Alias.Username,
		Pouch:rpc.Alias.Pouch,
		SnipName:models.SnipName{Name: rpc.Alias.SnipName.Name, Ext: rpc.Alias.SnipName.Extension},
	}
	// if encrypted, decrypt
	// if checksum doesn't match then throw warning
	// check that checksum signature is valid with public key.
	model.Snip = rpc.Snip
	model.Version = rpc.SnipVersion
	model.Tags = rpc.Tags
	model.Created = time.Unix(rpc.Created/1000, 0)
	model.Description = rpc.Description
	model.ClonedFromAlias = rpc.ClonedFromAlias
	model.ClonedFromVersion = rpc.ClonedFromVersion
	model.Private = rpc.Private
	model.RunCount = rpc.RunCount
	model.CloneCount = rpc.CloneCount
	model.Role = models.SnipRole(rpc.Role)
	if cache {
		r.persister.Upsert(LATEST_SNIPPET, model)
	}
}

const LATEST_SNIPPET = "latest-snippet.json"
const DELETED_SNIPPET = "deleted-snippet.json"

func (r *RpcService) mapSnippetList(rpc *snipsRpc.ListResponse, model *models.SnippetList, isList bool) {
	model.Total = rpc.Total
	model.Since = time.Unix(rpc.Since/1000, 0)
	model.Size = rpc.Size
	newSnip := &models.Snippet{}
	// TODO: Monitor eventual consistency and tweak cache duration.
	// Test with: go build;./kwkcli new "dong1" zing.sh;./kwkcli ls;sleep 11;./kwkcli ls;
	r.persister.Get(LATEST_SNIPPET, newSnip, time.Now().Unix()-10)
	isInList := false
	for _, v := range rpc.Items {
		item := &models.Snippet{}
		r.mapSnip(v, item, false)
		model.Items = append(model.Items, *item)
		if item.Id == newSnip.Id {
			isInList = true
		}
	}
	if isList && !isInList && newSnip.Alias.Name != "" {
		// TODO: add to logger
		fmt.Println("Adding from cache")
		model.Items = append([]models.Snippet{*newSnip}, model.Items...)
	}
}

func mapAlias(a models.Alias) *snipsRpc.Alias {
	return &snipsRpc.Alias{
		Username:a.Username,
		Pouch:   a.Pouch,
		SnipName:&snipsRpc.SnipName{Name:a.Name, Extension:a.Ext},
	}
}
