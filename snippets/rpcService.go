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

func (r *RpcService) Update(fullKey string, description string) (*models.Snippet, error) {
	if res, err := r.client.Update(r.headers.GetContext(), &snipsRpc.UpdateRequest{FullName: fullKey, Description: description}); err != nil {
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
		r.mapSnippetList(res, list)
		return list, nil
	}
}

func (r *RpcService) Get(k *models.Alias) (*models.SnippetList, error) {
	if res, err := r.client.Get(r.headers.GetContext(), &snipsRpc.GetRequest{Username: k.Username, FullName: k.FullKey}); err != nil {
		return nil, err
	} else {
		list := &models.SnippetList{}
		r.mapSnippetList(res, list)
		return list, nil
	}
}

func (r *RpcService) Delete(fullName string) error {
	_, err := r.client.Delete(r.headers.GetContext(), &snipsRpc.DeleteRequest{FullName: fullName})
	//if err != nil {
	//	r.Settings.Upsert(DELETED_SNIPPET, )
	//}
	return err
}

func (r *RpcService) Create(snip string, path string, role models.SnipRole) (*models.CreateSnippetResponse, error) {
	// encrypt if requested
	fmt.Println(snip, path)
	if res, err := r.client.Create(r.headers.GetContext(), &snipsRpc.CreateRequest{Snip: snip, FullName: path, Role: snipsRpc.Role(role)}); err != nil {
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

func (r *RpcService) Rename(fullKey string, newFullName string) (*models.Snippet, string, error) {
	if res, err := r.client.Rename(r.headers.GetContext(), &snipsRpc.RenameRequest{FullName: fullKey, NewFullName: newFullName}); err != nil {
		return nil, "", err
	} else {
		snip := &models.Snippet{}
		r.mapSnip(res.Snip, snip, true)
		return snip, res.OriginalFullName, nil
	}
}

func (r *RpcService) Patch(fullKey string, target string, patch string) (*models.Snippet, error) {
	if res, err := r.client.Patch(r.headers.GetContext(), &snipsRpc.PatchRequest{FullName: fullKey, Target: target, Patch: patch}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		r.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (r *RpcService) Clone(k *models.Alias, newFullName string) (*models.Snippet, error) {
	if res, err := r.client.Clone(r.headers.GetContext(), &snipsRpc.CloneRequest{Username: k.Username, FullName: k.FullKey, NewFullName: newFullName}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		r.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (r *RpcService) Tag(fullKey string, tags ...string) (*models.Snippet, error) {
	if res, err := r.client.Tag(r.headers.GetContext(), &snipsRpc.TagRequest{FullName: fullKey, Tags: tags}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		r.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (r *RpcService) UnTag(fullKey string, tags ...string) (*models.Snippet, error) {
	if res, err := r.client.UnTag(r.headers.GetContext(), &snipsRpc.UnTagRequest{FullName: fullKey, Tags: tags}); err != nil {
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
	model.Role = models.SnipRole(rpc.Role)
	if cache {
		r.Settings.Upsert(LATEST_SNIPPET, model)
	}
}

const LATEST_SNIPPET = "latest-snippet.json"
const DELETED_SNIPPET = "deleted-snippet.json"

func (r *RpcService) mapSnippetList(rpc *snipsRpc.ListResponse, model *models.SnippetList) {
	model.Total = rpc.Total
	model.Since = time.Unix(rpc.Since/1000, 0)
	model.Size = rpc.Size
	newSnip := &models.Snippet{}
	// TODO: Monitor eventual consistency and tweak cache duration.
	// Test with: go build;./kwkcli new "dong1" zing.sh;./kwkcli ls;sleep 11;./kwkcli ls;
	r.Settings.Get(LATEST_SNIPPET, newSnip, time.Now().Unix()-10)
	isInList := false
	for _, v := range rpc.Items {
		item := &models.Snippet{}
		r.mapSnip(v, item, false)
		model.Items = append(model.Items, *item)
		if item.Id == newSnip.Id {
			isInList = true
		}
	}
	if !isInList && newSnip.Name != "" {
		// TODO: add to logger
		fmt.Println("Adding from cache")
		model.Items = append([]models.Snippet{*newSnip}, model.Items...)
	}
}
