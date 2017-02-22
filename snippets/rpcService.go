package snippets

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/rpc/src/snipsRpc"
	"google.golang.org/grpc"
	"time"
	"fmt"
	"bitbucket.com/sharingmachine/kwkcli/update"
	"bitbucket.com/sharingmachine/kwkcli/log"
)

type RpcService struct {
	persister config.Persister
	pc        snipsRpc.PouchesRpcClient
	client    snipsRpc.SnipsRpcClient
	h         *rpc.Headers
}

func New(conn *grpc.ClientConn, s config.Persister, h *rpc.Headers) Service {
	return &RpcService{persister: s,
		client:               snipsRpc.NewSnipsRpcClient(conn),
		pc:                   snipsRpc.NewPouchesRpcClient(conn),
		h:              h,
	}
}

func (rs *RpcService) AlphaSearch(term string) (*models.SearchTermResponse, error) {
	if res, err := rs.client.Alpha(rs.h.Context(), &snipsRpc.AlphaRequest{
		Term: term,
	}); err != nil {
		return nil, err
	} else {
		r := &models.SearchTermResponse{}
		r.Results = []*models.SearchResult{}
		r.Took = res.Took
		r.Total = res.Total
		for _, v := range res.Results {
			s := &models.Snippet{}
			rs.mapSnip(v.Snippet, s, false)
			sr := &models.SearchResult{
				Snippet: s,
				Highlights: v.Highlights,
			}
			r.Results = append(r.Results, sr)
		}
		return r, nil
	}
}

func (rs *RpcService) Update(a models.Alias, description string) (*models.Snippet, error) {
	if res, err := rs.client.Update(rs.h.Context(), &snipsRpc.UpdateRequest{Alias: mapAlias(a), Description: description}); err != nil {
		return nil, err
	} else {
		m := &models.Snippet{}
		rs.mapSnip(res.Snip, m, true)
		return m, nil
	}
}

// since unix time in milliseconds
func (rs *RpcService) List(l *models.ListParams) (*models.SnippetList, error) {
	if res, err := rs.client.List(rs.h.Context(), &snipsRpc.ListRequest{Username: l.Username, Pouch: l.Pouch, Since: l.Since, Size: l.Size, Tags: l.Tags, All: l.All}); err != nil {
		return nil, err
	} else {
		list := &models.SnippetList{}
		rs.mapSnippetList(res, list, true)
		return list, nil
	}
}

func (rs *RpcService) Get(k models.Alias) (*models.SnippetList, error) {
	if res, err := rs.client.Get(rs.h.Context(), &snipsRpc.GetRequest{Alias: mapAlias(k)}); err != nil {
		return nil, err
	} else {
		list := &models.SnippetList{}
		rs.mapSnippetList(res, list, false)
		return list, nil
	}
}

func (rs *RpcService) Delete(username string, pouch string, names []*models.SnipName) error {
	sn := []*snipsRpc.SnipName{}
	for _, v := range names {
		sn = append(sn, &snipsRpc.SnipName{Name: v.Name, Extension: v.Ext })
	}
	_, err := rs.client.Delete(rs.h.Context(), &snipsRpc.DeleteRequest{Username: username, Pouch: pouch, Names: sn})
	//if err != nil {
	//	rs.Settings.Upsert(DELETED_SNIPPET, )
	//}
	return err
}

func (rs *RpcService) Move(username string, sourcePouch string, targetPouch string, names []*models.SnipName) (string, error) {
	ns := []*snipsRpc.SnipName{}
	for _, v := range names {
	   ns = append(ns, &snipsRpc.SnipName{Name:v.Name, Extension:v.Ext})
	}
	mv := &snipsRpc.MoveRequest{Username: username, SourcePouch: sourcePouch, TargetPouch: targetPouch,SnipNames:ns}
	r, err := rs.client.Move(rs.h.Context(), mv)
	if err != nil {
		return "", err
	}
	return r.Pouch, nil
}

func (rs *RpcService) Create(snip string, a models.Alias, role models.SnipRole) (*models.CreateSnippetResponse, error) {
	if res, err := rs.client.Create(rs.h.Context(), &snipsRpc.CreateRequest{Snip: snip, Alias: mapAlias(a), Role: snipsRpc.Role(role)}); err != nil {
		return nil, err
	} else {
		cs := &models.CreateSnippetResponse{}
		if res.Snip != nil {
			snip := &models.Snippet{}
			rs.mapSnip(res.Snip, snip, true)
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

func (rs *RpcService) Rename(a models.Alias, new models.SnipName) (*models.Snippet, *models.SnipName, error) {
	if res, err := rs.client.Rename(rs.h.Context(), &snipsRpc.RenameRequest{Alias: mapAlias(a), NewName: &snipsRpc.SnipName{Name: new.Name, Extension: new.Ext}}); err != nil {
		return nil, nil, err
	} else {
		snip := &models.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, &models.SnipName{Name: res.Original.Name, Ext: res.Original.Extension}, nil
	}
}

func (rs *RpcService) Patch(a models.Alias, target string, patch string) (*models.Snippet, error) {
	if res, err := rs.client.Patch(rs.h.Context(), &snipsRpc.PatchRequest{Alias: mapAlias(a), Target: target, Patch: patch}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (rs *RpcService) Clone(a models.Alias, new models.Alias) (*models.Snippet, error) {
	if res, err := rs.client.Clone(rs.h.Context(), &snipsRpc.CloneRequest{Alias: mapAlias(a), New: mapAlias(new)}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (rs *RpcService) LogRun(a models.Alias, s models.RunStatus) {
	_, err := rs.client.LogRun(rs.h.Context(), &snipsRpc.LogRunRequest{
		Alias: mapAlias(a), Status: snipsRpc.RunStatus(s), Time: time.Now().Unix() })
	if err != nil {
		log.Error("Error sending LogRun", err)
	}
}

func (rs *RpcService) Tag(a models.Alias, tags ...string) (*models.Snippet, error) {
	if res, err := rs.client.Tag(rs.h.Context(), &snipsRpc.TagRequest{Alias: mapAlias(a), Tags: tags}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (rs *RpcService) UnTag(a models.Alias, tags ...string) (*models.Snippet, error) {
	if res, err := rs.client.UnTag(rs.h.Context(), &snipsRpc.UnTagRequest{Alias: mapAlias(a), Tags: tags}); err != nil {
		return nil, err
	} else {
		snip := &models.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (rs *RpcService) GetRoot(username string, all bool) (*models.Root, error) {
	r, err := rs.pc.GetRoot(rs.h.Context(), &snipsRpc.RootRequest{Username: username, All: all})
	if err != nil {
		return nil, err
	}
	l := &models.SnippetList{}
	res := &snipsRpc.ListResponse{Items: r.Snips}
	rs.mapSnippetList(res, l, true)
	pl := rs.mapPouchList(r.Pouches)
	perL := rs.mapPouchList(r.Personal)
	record := &update.Record{}
	root := &models.Root{Snippets: l.Items, Pouches: pl, Personal:perL, Username:r.Username}
	if e := rs.persister.Get(update.RecordFile, record, 0); e == nil {
		root.LastUpdate = record.LastUpdate
	}
	return root, err
}

func (rs *RpcService) CreatePouch(pouch string) (string, error) {
	_, err := rs.pc.Create(rs.h.Context(), &snipsRpc.CreatePouchRequest{Name: pouch})
	if err != nil {
		return "", err
	}
	return pouch, nil
}

func (rs *RpcService) RenamePouch(pouch string, newPouch string) (string, error) {
	req := &snipsRpc.RenamePouchRequest{ Name:pouch, NewName:newPouch}
	r, err := rs.pc.Rename(rs.h.Context(), req)
	if err != nil {
		return "", err
	}
	return r.OriginalName, nil
}

func (rs *RpcService) MakePrivate(pouch string, private bool) (bool, error) {
	req := &snipsRpc.MakePrivateRequest{ MakePrivate:private, Name:pouch }
	_, err := rs.pc.MakePrivate(rs.h.Context(), req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rs *RpcService) DeletePouch(pouch string) (bool, error) {
	req := &snipsRpc.DeletePouchRequest{ Name:pouch }
	_, err := rs.pc.Delete(rs.h.Context(), req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func MillisToTime(in int64) time.Time {
	return time.Unix(in/1000, 0)
}

func (rs *RpcService) mapPouchList(in []*snipsRpc.Pouch) []*models.Pouch {
	out := []*models.Pouch{}
	for _, v := range in {
		p := &models.Pouch{
			Name:        v.Name,
			Username:    v.Username,
			Encrypt:     v.Encrypt,
			MakePrivate: v.MakePrivate,
			SharedWith:  v.SharedWith,
			SnipCount:   v.SnipCount,
			Modified:    MillisToTime(v.Modified),
			PouchId:     v.PouchId,
		}
		out = append(out, p)
	}
	return out
}

/*
  cache will add the snippet to the local cache to deal with eventual consistency user experience.
 */
func (rs *RpcService) mapSnip(rpc *snipsRpc.Snip, model *models.Snippet, cache bool) {
	model.Id = rpc.SnipId
	model.Alias = models.Alias{
		Username: rpc.Alias.Username,
		Pouch:    rpc.Alias.Pouch,
		SnipName: models.SnipName{Name: rpc.Alias.SnipName.Name, Ext: rpc.Alias.SnipName.Extension},
	}
	// if encrypted, decrypt
	// if checksum doesn't match then throw warning
	// check that checksum signature is valid with public key.
	model.Snip = rpc.Snip
	model.Version = rpc.SnipVersion
	model.Tags = rpc.Tags
	model.Created = MillisToTime(rpc.Created)
	model.Description = rpc.Description
	model.ClonedFromAlias = rpc.ClonedFromAlias
	model.ClonedFromVersion = rpc.ClonedFromVersion
	model.Private = rpc.Private
	model.RunCount = rpc.RunCount
	model.CloneCount = rpc.CloneCount
	model.Role = models.SnipRole(rpc.Role)
	model.RunStatus = models.RunStatus(rpc.RunStatus)
	model.RunStatusTime = rpc.RunStatusTime
	if cache {
		rs.persister.Upsert(LATEST_SNIPPET, model)
	}
}

const LATEST_SNIPPET = "latest-snippet.json"
const DELETED_SNIPPET = "deleted-snippet.json"

func (rs *RpcService) mapSnippetList(rpc *snipsRpc.ListResponse, model *models.SnippetList, isList bool) {
	model.Username = rpc.Username
	model.Pouch = rpc.Pouch
	model.Total = rpc.Total
	model.Since = time.Unix(rpc.Since/1000, 0)
	model.Size = rpc.Size
	newSnip := &models.Snippet{}
	// TODO: Monitor eventual consistency and tweak cache duration.
	// Test with: go build;./kwkcli new "dong1" zing.sh;./kwkcli ls;sleep 11;./kwkcli ls;
	rs.persister.Get(LATEST_SNIPPET, newSnip, time.Now().Unix()-10)
	isInList := false
	for _, v := range rpc.Items {
		item := &models.Snippet{}
		rs.mapSnip(v, item, false)
		model.Items = append(model.Items, item)
		if item.Id == newSnip.Id {
			isInList = true
		}
	}
	if isList && !isInList && newSnip.Alias.Name != "" {
		// TODO: add to logger
		fmt.Println("Adding from cache")
		model.Items = append([]*models.Snippet{newSnip}, model.Items...)
	}
}

func mapAlias(a models.Alias) *snipsRpc.Alias {
	return &snipsRpc.Alias{
		Username: a.Username,
		Pouch:    a.Pouch,
		SnipName: &snipsRpc.SnipName{Name: a.Name, Extension: a.Ext},
	}
}
