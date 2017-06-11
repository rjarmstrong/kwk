package handlers

import (
	"github.com/rjarmstrong/kwk/src/cli"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"reflect"
	rt "runtime"
	"strings"
	"testing"
)

var (
	snippetClient *fakeSnipClient
	runner        *fakeRunner
	ed            *fakeEditor
	dlg           *fakeDialogue
	writer        *fakeWWriter
	cxf           cli.ContextFunc
	prefs         *out.Prefs
	snippets      *Snippets
	rootCalled    *types.RootResponse
)

func TestMain(m *testing.M) {
	snippetClient = &fakeSnipClient{returnsFor: map[string]response{}, called: map[string]interface{}{}}
	runner = &fakeRunner{called: map[string]interface{}{}}
	dlg = &fakeDialogue{called: map[string]interface{}{}}
	ed = &fakeEditor{}
	writer = &fakeWWriter{called: map[string]interface{}{}}
	cxf = func() context.Context {
		return context.Background()
	}
	prefs = &out.Prefs{}

	rootPrinter := func(rr *types.RootResponse) error {
		rootCalled = rr
		return nil
	}

	snippets = NewSnippets(prefs, snippetClient, runner, ed, writer, cxf, rootPrinter, dlg)
	code := m.Run()
	os.Exit(code)
}

type fakeWWriter struct {
	called map[string]interface{}
}

func (fc *fakeWWriter) Write(p vwrite.Handler) {
	fc.called["Write"] = p
}

func (fc *fakeWWriter) EWrite(p vwrite.Handler) error {
	funcName := getFuncName(p)
	fc.called["EWrite"] = funcName
	return nil
}

func (fc *fakeWWriter) PopCalled(name string) interface{} {
	c := fc.called[name]
	delete(fc.called, name)
	return c
}

type fakeEditor struct {
}

func (*fakeEditor) Invoke(s *types.Snippet, onchange func(s types.Snippet)) error {
	panic("implement me")
}

func (*fakeEditor) Close(s *types.Snippet) (uint, error) {
	panic("implement me")
}

type fakeDialogue struct {
	called map[string]interface{}
}

func (fc *fakeDialogue) PopCalled(name string) interface{} {
	c := fc.called[name]
	delete(fc.called, name)
	return c
}

func (fc *fakeDialogue) ChooseSnippet(s []*types.Snippet) *types.Snippet {
	fc.called["ChooseSnippet"] = s
	if len(s) == 0 {
		return nil
	}
	return s[0]
}

func (fc *fakeDialogue) Modal(handler vwrite.Handler, autoYes bool) *out.DialogResponse {
	fc.called["Modal"] = getFuncName(handler)
	return &out.DialogResponse{Ok: true}
}

func (fc *fakeDialogue) FormField(field vwrite.Handler, mask bool) (*out.DialogResponse, error) {
	panic("implement me")
}

type fakeRunner struct {
	called map[string]interface{}
}

func (fc *fakeRunner) PopCalled(name string) interface{} {
	c := fc.called[name]
	delete(fc.called, name)
	return c
}

func (fc *fakeRunner) Run(s *types.Snippet, args []string) error {
	fc.called["Run"] = s.Alias.URI()
	return nil
}

type response struct {
	err error
	val interface{}
}

type fakeSnipClient struct {
	called     map[string]interface{}
	returnsFor map[string]response
}

func (fc *fakeSnipClient) PopCalled(name string) interface{} {
	c := fc.called[name]
	delete(fc.called, name)
	return c
}

func (fc *fakeSnipClient) Create(ctx context.Context, in *types.CreateRequest, opts ...grpc.CallOption) (*types.CreateResponse, error) {
	fc.called["Create"] = in
	return &types.CreateResponse{}, nil
}

func (fc *fakeSnipClient) Update(ctx context.Context, in *types.UpdateRequest, opts ...grpc.CallOption) (*types.UpdateResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) Move(ctx context.Context, in *types.MoveRequest, opts ...grpc.CallOption) (*types.MoveResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) Rename(ctx context.Context, in *types.RenameRequest, opts ...grpc.CallOption) (*types.RenameResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) Patch(ctx context.Context, in *types.PatchRequest, opts ...grpc.CallOption) (*types.PatchResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) Clone(ctx context.Context, in *types.CloneRequest, opts ...grpc.CallOption) (*types.CloneResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) Tag(ctx context.Context, in *types.TagRequest, opts ...grpc.CallOption) (*types.TagResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) UnTag(ctx context.Context, in *types.UnTagRequest, opts ...grpc.CallOption) (*types.UnTagResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) Get(ctx context.Context, in *types.GetRequest, opts ...grpc.CallOption) (*types.ListResponse, error) {
	fc.called["Get"] = in
	res, ok := fc.returnsFor["Get"]
	if ok {
		return res.val.(*types.ListResponse), res.err
	}
	return nil, nil
}

func (fc *fakeSnipClient) List(ctx context.Context, in *types.ListRequest, opts ...grpc.CallOption) (*types.ListResponse, error) {
	fc.called["List"] = in
	return &types.ListResponse{Username: "username1", Pouch: &types.Pouch{Name: "pouch1"}}, nil
}

func (fc *fakeSnipClient) Delete(ctx context.Context, in *types.DeleteRequest, opts ...grpc.CallOption) (*types.DeleteResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) GetRoot(ctx context.Context, in *types.RootRequest, opts ...grpc.CallOption) (*types.RootResponse, error) {
	fc.called["GetRoot"] = in
	res, ok := fc.returnsFor["GetRoot"]
	if ok {
		return res.val.(*types.RootResponse), res.err
	}
	return nil, nil
}

func (fc *fakeSnipClient) CreatePouch(ctx context.Context, in *types.CreatePouchRequest, opts ...grpc.CallOption) (*types.CreatePouchResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) RenamePouch(ctx context.Context, in *types.RenamePouchRequest, opts ...grpc.CallOption) (*types.RenamePouchResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) MakePouchPrivate(ctx context.Context, in *types.MakePrivateRequest, opts ...grpc.CallOption) (*types.MakePrivateResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) DeletePouch(ctx context.Context, in *types.DeletePouchRequest, opts ...grpc.CallOption) (*types.DeletePouchResponse, error) {
	panic("implement me")
}

func (fc *fakeSnipClient) Alpha(ctx context.Context, in *types.AlphaRequest, opts ...grpc.CallOption) (*types.AlphaResponse, error) {
	fc.called["Alpha"] = in
	return &types.AlphaResponse{}, nil
}

func (fc *fakeSnipClient) TypeAhead(ctx context.Context, in *types.TypeAheadRequest, opts ...grpc.CallOption) (*types.TypeAheadResponse, error) {
	fc.called["TypeAhead"] = in
	return &types.TypeAheadResponse{}, nil
}

func (fc *fakeSnipClient) LogUse(ctx context.Context, in *types.UseContext, opts ...grpc.CallOption) (*types.LogUseResponse, error) {
	fc.called["LogUse"] = in
	return &types.LogUseResponse{}, nil
}

// getFuncName gets the function name of the values pointer value of the vwrite.Handler
func getFuncName(p vwrite.Handler) string {
	tn := rt.FuncForPC(reflect.ValueOf(p).Pointer()).Name()
	prefix := "github.com/rjarmstrong/kwk/src/"
	funcName := strings.Split(strings.Replace(tn, prefix, "", -1), ".")[1]
	return funcName
}
