package app

import (
	"bytes"
	"github.com/kwk-super-snippets/cli/src/models"
	"github.com/kwk-super-snippets/types"
	"github.com/lunixbochs/vtclean"
	"google.golang.org/grpc"
	"runtime"
	"strings"
	"time"
	"bitbucket.com/sharingmachine/rpc-deprecated/src/snipsRpc"
)

type SnippetsGrpc struct {
	pc     snipsRpc.PouchesRpcClient
	client snipsRpc.SnipsRpcClient
	h      Headers
}

func New(conn *grpc.ClientConn, client types.AppInfo) types.SnippetsClient {

		client: snipsRpc.NewSnipsRpcClient(conn),
		pc:     snipsRpc.NewPouchesRpcClient(conn),
		h:      Headers{version: client.Version},
	}
}


func (rs *SnippetsGrpc) AlphaSearch(term string) (*types.SearchTermResponse, error) {
	if res, err := rs.client.Alpha(rs.h.Context(), &snipsRpc.AlphaRequest{
		Term: term,
	}); err != nil {
		return nil, err
	} else {
		r := &types.SearchTermResponse{}
		r.Results = []*types.SearchResult{}
		r.Took = res.Took
		r.Total = res.Total
		for _, v := range res.Results {
			s := &types.Snippet{}
			rs.mapSnip(v.Snippet, s, false)
			sr := &types.SearchResult{
				Snippet:    s,
				Highlights: v.Highlights,
			}
			r.Results = append(r.Results, sr)
		}
		return r, nil
	}
}

func (rs *SnippetsGrpc) Update(a types.Alias, description string) (*types.Snippet, error) {
	if res, err := rs.client.Update(rs.h.Context(), &snipsRpc.UpdateRequest{Alias: mapAlias(a), Description: description}); err != nil {
		return nil, err
	} else {
		m := &types.Snippet{}
		rs.mapSnip(res.Snip, m, true)
		return m, nil
	}
}

// since unix time in milliseconds
func (rs *SnippetsGrpc) List(l *models.ListParams) (*models.ListView, error) {
	if res, err := rs.client.List(rs.h.Context(), &snipsRpc.ListRequest{
		Username:      l.Username,
		Pouch:         l.Pouch,
		Since:         l.Since,
		Size:          l.Size,
		Tags:          l.Tags,
		IgnorePouches: l.IgnorePouches,
		Category:      l.Category,
		All:           l.All,
	}); err != nil {
		return nil, err
	} else {
		list := &models.ListView{}
		rs.mapSnippetList(res, list, true)
		return list, nil
	}
}

func (rs *SnippetsGrpc) Get(k types.Alias) (*models.ListView, error) {
	if res, err := rs.client.Get(rs.h.Context(), &snipsRpc.GetRequest{Alias: mapAlias(k), Version: k.Version}); err != nil {
		return nil, err
	} else {
		list := &models.ListView{}
		rs.mapSnippetList(res, list, false)
		return list, nil
	}
}



func (rs *SnippetsGrpc) Move(username string, sourcePouch string, targetPouch string, names []*types.SnipName) (string, error) {
	ns := []*snipsRpc.SnipName{}
	for _, v := range names {
		ns = append(ns, &snipsRpc.SnipName{Name: v.Name, Extension: v.Ext})
	}
	mv := &snipsRpc.MoveRequest{Username: username, SourcePouch: sourcePouch, TargetPouch: targetPouch, SnipNames: ns}
	r, err := rs.client.Move(rs.h.Context(), mv)
	if err != nil {
		return "", err
	}
	return r.Pouch, nil
}

func (rs *SnippetsGrpc) Create(snip string, a types.Alias, role types.SnipRole) (*models.CreateSnippetResponse, error) {
	if res, err := rs.client.Create(rs.h.Context(), &snipsRpc.CreateRequest{Snip: snip, Alias: mapAlias(a), Role: snipsRpc.Role(role)}); err != nil {
		return nil, err
	} else {
		cs := &models.CreateSnippetResponse{}
		if res.Snip != nil {
			snip := &types.Snippet{}
			rs.mapSnip(res.Snip, snip, true)
			cs.Snippet = snip
		}
		return cs, nil
	}
}

func (rs *SnippetsGrpc) Rename(a types.Alias, new types.SnipName) (*types.Snippet, *types.SnipName, error) {
	if res, err := rs.client.Rename(rs.h.Context(), &snipsRpc.RenameRequest{Alias: mapAlias(a), NewName: &snipsRpc.SnipName{Name: new.Name, Extension: new.Ext}}); err != nil {
		return nil, nil, err
	} else {
		snip := &types.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, &types.SnipName{Name: res.Original.Name, Ext: res.Original.Extension}, nil
	}
}

func (rs *SnippetsGrpc) Patch(a types.Alias, target string, patch string) (*types.Snippet, error) {
	if res, err := rs.client.Patch(rs.h.Context(), &snipsRpc.PatchRequest{Alias: mapAlias(a), Target: target, Patch: patch}); err != nil {
		return nil, err
	} else {
		snip := &types.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (rs *SnippetsGrpc) Clone(a types.Alias, new types.Alias) (*types.Snippet, error) {
	if res, err := rs.client.Clone(rs.h.Context(), &snipsRpc.CloneRequest{Alias: mapAlias(a), New: mapAlias(new)}); err != nil {
		return nil, err
	} else {
		snip := &types.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

type UseContext struct {
	Preview     string
	Runner      string
	Level       int64
	CallerAlias string
}

func (rs *SnippetsGrpc) LogUse(a types.Alias, s types.UseStatus, u types.UseType, ctx *UseContext) {
	req := &snipsRpc.LogUseRequest{
		Alias: mapAlias(a), Status: snipsRpc.UseStatus(s),
		Type: snipsRpc.UseType(u),
		Time: time.Now().Unix(),
	}

	if ctx != nil {
		req.Preview = LimitPreview(ctx.Preview, 50)
		req.Runner = ctx.Runner
		req.Level = ctx.Level
		req.CallerAlias = ctx.CallerAlias
		req.Os = runtime.GOOS
	}
	_, err := rs.client.LogUse(rs.h.Context(), req)
	if err != nil {
		models.Debug("Error sending LogRun")
		models.LogErr(err)
	}
}

/*
 Limits a preview adding an ascii escape at the end and fixing the length.
*/
func LimitPreview(in string, length int) string {
	in = vtclean.Clean(in, true)
	return Limit(in, length-5) + "\033[0m"
}

func Limit(in string, length int) string {
	in = strings.Replace(in, "\n", "  ", -1)
	in = strings.TrimSpace(in)
	var numRunes = 0
	b := bytes.Buffer{}
	for _, r := range in {
		if numRunes == length {
			return strings.TrimSpace(b.String())
		}
		numRunes++
		if r == '\n' {
			b.WriteRune(' ')
			b.WriteRune(' ')
			continue
		}
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

func (rs *SnippetsGrpc) Tag(a types.Alias, tags ...string) (*types.Snippet, error) {
	if res, err := rs.client.Tag(rs.h.Context(), &snipsRpc.TagRequest{Alias: mapAlias(a), Tags: tags}); err != nil {
		return nil, err
	} else {
		snip := &types.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (rs *SnippetsGrpc) UnTag(a types.Alias, tags ...string) (*types.Snippet, error) {
	if res, err := rs.client.UnTag(rs.h.Context(), &snipsRpc.UnTagRequest{Alias: mapAlias(a), Tags: tags}); err != nil {
		return nil, err
	} else {
		snip := &types.Snippet{}
		rs.mapSnip(res.Snip, snip, true)
		return snip, nil
	}
}

func (rs *SnippetsGrpc) GetRoot(username string, all bool) (*types.ListView, error) {
	r, err := rs.pc.GetRoot(rs.h.Context(), &snipsRpc.RootRequest{Username: username, All: all})
	if err != nil {
		return nil, err
	}
	l := &types.ListView{}
	res := &snipsRpc.ListResponse{Items: r.Snips}
	rs.mapSnippetList(res, l, true)
	pl := rs.mapPouchList(r.Pouches)
	perL := rs.mapPouchList(r.Personal)
	//record := &update.Record{}
	root := &types.ListView{
		IsRoot:   true,
		Snippets: l.Snippets,
		Pouches:  pl,
		Personal: perL,
		Username: r.Username,
		UserStats: types.UserStats{
			LastPouch:        r.Stats.LastPouch,
			RecentPouches:    r.Stats.RecentPouches,
			MaxUsePerPouch:   r.Stats.MaxUsePerPouch,
			MaxSnipsPerPouch: r.Stats.MaxSnipsPerPouch,
		},
	}
	//if e := rs.persister.Get(update.RecordFile, record, 0); e == nil {
	//	root.LastUpgrade = record.LastUpdate
	//}
	return root, err
}

func (rs *SnippetsGrpc) CreatePouch(pouch string) (string, error) {
	_, err := rs.pc.Create(rs.h.Context(), &snipsRpc.CreatePouchRequest{Name: pouch})
	if err != nil {
		return "", err
	}
	return pouch, nil
}

func (rs *SnippetsGrpc) RenamePouch(pouch string, newPouch string) (string, error) {
	req := &snipsRpc.RenamePouchRequest{Name: pouch, NewName: newPouch}
	r, err := rs.pc.Rename(rs.h.Context(), req)
	if err != nil {
		return "", err
	}
	return r.OriginalName, nil
}

func (rs *SnippetsGrpc) MakePrivate(pouch string, private bool) (bool, error) {
	req := &snipsRpc.MakePrivateRequest{MakePrivate: private, Name: pouch}
	_, err := rs.pc.MakePrivate(rs.h.Context(), req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rs *SnippetsGrpc) DeletePouch(pouch string) (bool, error) {
	req := &snipsRpc.DeletePouchRequest{Name: pouch}
	_, err := rs.pc.Delete(rs.h.Context(), req)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rs *SnippetsGrpc) mapPouchList(in []*snipsRpc.Pouch) []*types.Pouch {
	out := []*types.Pouch{}
	for _, v := range in {
		p := mapPouch(v)
		out = append(out, p)
	}
	return out
}

/*
  cache will add the snippet to the local cache to deal with eventual consistency user experience.
*/
func (rs *SnippetsGrpc) mapSnip(rpc *snipsRpc.Snip, model *types.Snippet, cache bool) {
	model.Id = rpc.SnipId
	model.Alias = types.Alias{
		Username: rpc.Alias.Username,
		Pouch:    rpc.Alias.Pouch,
		SnipName: types.SnipName{Name: rpc.Alias.SnipName.Name, Ext: rpc.Alias.SnipName.Extension},
	}
	// if encrypted, decrypt
	// if checksum doesn't match then throw warning
	// check that checksum signature is valid with public key.
	model.Snip = rpc.Snip
	model.Version = rpc.SnipVersion
	model.Tags = rpc.Tags
	model.Created = time.Unix(rpc.Created, 0)
	model.Description = rpc.Description
	model.ClonedFromAlias = rpc.ClonedFromAlias
	model.ClonedFromVersion = rpc.ClonedFromVersion
	model.Private = rpc.Private
	model.Runs = rpc.Stats.Runs
	model.Views = rpc.Stats.Views
	model.Clones = rpc.Stats.Clones

	model.Role = types.SnipRole(rpc.Role)
	model.RunStatus = types.UseStatus(rpc.RunStatus)
	model.RunStatusTime = time.Unix(rpc.RunStatusTime, 0)
	model.Preview = rpc.Preview
	model.SnipChecksum = rpc.SnipChecksum
	model.Attribution = rpc.Attribution
	model.Dependencies = rpc.Dependencies
	model.Apps = rpc.Apps
	model.SupportedOs = rpc.SupportedOs
}

func (rs *SnippetsGrpc) mapSnippetList(rpc *snipsRpc.ListResponse, model *models.ListView, isList bool) {
	model.Username = rpc.Username
	model.Pouch = mapPouch(rpc.Pouch)
	model.Total = rpc.Total
	model.Since = time.Unix(rpc.Since/1000, 0)
	model.Size = rpc.Size
	for _, v := range rpc.Items {
		item := &types.Snippet{}
		rs.mapSnip(v, item, false)
		model.Snippets = append(model.Snippets, item)
	}
}

func mapAlias(a types.Alias) *snipsRpc.Alias {
	return &snipsRpc.Alias{
		Username: a.Username,
		Pouch:    a.Pouch,
		SnipName: &snipsRpc.SnipName{Name: a.Name, Extension: a.Ext},
	}
}

func mapPouch(p *snipsRpc.Pouch) *types.Pouch {
	if p == nil {
		return nil
	}
	return &types.Pouch{
		Encrypt:     p.Encrypt,
		MakePrivate: p.MakePrivate,
		Modified:    time.Unix(p.Modified, 0),
		Name:        p.Name,
		Id:          p.PouchId,
		SharedWith:  p.SharedWith,
		LastUse:     time.Unix(p.LastUse, 0),
		Type:        types.PouchType(p.Type),
		PouchCounts: types.PouchCounts{
			Runs:   p.Stats.Runs,
			Views:  p.Stats.Views,
			Clones: p.Stats.Clones,
			Green:  p.Stats.Green,
			Red:    p.Stats.Red,
			Snips:  p.Stats.Snips,
		},
		UnOpened: p.UnOpened,
		Username: p.Username,
	}
}
