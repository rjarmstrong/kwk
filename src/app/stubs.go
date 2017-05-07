package app

import (
	"bytes"
	"github.com/kwk-super-snippets/cli/src/gokwk"
	"github.com/kwk-super-snippets/cli/src/models"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/vwrite"
	"reflect"
	"runtime/debug"
	"testing"
)

func ErrIf(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err, string(debug.Stack()))
	}
}

type RunnerMock struct {
	RunCalledWith  []interface{}
	EditCalledWith *types.Snippet
}

func (o *RunnerMock) Run(alias *types.Snippet, args []string) error {
	o.RunCalledWith = []interface{}{alias, args}
	return nil
}

func (o *RunnerMock) Edit(alias *types.Snippet) error {
	o.EditCalledWith = alias
	return nil
}

type UpdaterMock struct {
}

func (*UpdaterMock) Run() error {
	panic("implement me")
}

type ErrorHandlerMock struct {
	HandleCalledW error
}

func (e *ErrorHandlerMock) Handle(err error) {
	e.HandleCalledW = err
}

type WriterMock struct {
	*bytes.Buffer
}

func (w *WriterMock) Write(p vwrite.Handler) {
	p.Write(w.Buffer)
}

func (w *WriterMock) EWrite(p vwrite.Handler) error {
	p.Write(w.Buffer)
	return nil
}

type DialogMock struct {
}

func (*DialogMock) Modal(handler vwrite.Handler, autoYes bool) *DialogResponse {
	panic("implement me")
}

func (*DialogMock) FormField(field vwrite.Handler, mask bool) (*DialogResponse, error) {
	panic("implement me")
}

func (*DialogMock) MultiChoice(vwrite.Handler, []*types.Snippet) (*types.Snippet, error) {
	panic("implement me")
}

type IoMock struct {
}

func (*IoMock) Delete(subDirName string, suffixPath string) error {
	panic("implement me")
}

func (*IoMock) Write(subDirName string, suffixPath string, snippet string, incHoldingDir bool) (string, error) {
	panic("implement me")
}

func (*IoMock) Read(subDirName string, suffixPath string, incHoldingDir bool, fresherThan int64) (string, error) {
	panic("implement me")
}

func (*IoMock) DeleteAll() error {
	panic("implement me")
}

type SnippetsMock struct {
	GetCalledWith         types.Alias
	RenameCalledWith      []string
	CreateCalledWith      []string
	ReturnItemsForGet     []*types.Snippet
	PatchCalledWith       []string
	DeleteCalledWith      []interface{}
	CloneCalledWith       []interface{}
	TagCalledWith         map[string][]string
	UnTagCalledWith       map[string][]string
	ListCalledWith        *models.ListParams
	CreatePouchCalledWith string
	DeletePouchCalledWith string
	GetRootCalledWith     []interface{}
	RenamePouchCalledWith []string
	MakePrivateCalledWith []interface{}
}

func (sm *SnippetsMock) LogUse(types.Alias, types.UseStatus, types.UseType, *gokwk.UseContext) {
	panic("not impl")
}

func (sm *SnippetsMock) AlphaSearch(term string) (*models.SearchTermResponse, error) {
	panic("not impl")
}

func (sm *SnippetsMock) Move(username string, sourcePouch string, targetPouch string, names []*types.SnipName) (string, error) {
	panic("not imp")
}

func (sm *SnippetsMock) Get(a types.Alias) (*models.ListView, error) {
	sm.GetCalledWith = a
	return &models.ListView{Snippets: sm.ReturnItemsForGet, Total: int64(len(sm.ReturnItemsForGet))}, nil
}

func (sm *SnippetsMock) Create(snip string, a types.Alias, role types.SnipRole) (*models.CreateSnippetResponse, error) {
	sm.CreateCalledWith = []string{snip, a.String()}
	if a.Name == "" {
		a.Name = "x5hi23"
	}
	return &models.CreateSnippetResponse{Snippet: &types.Snippet{Alias: a}}, nil
}

func (sm *SnippetsMock) Update(a types.Alias, description string) (*types.Snippet, error) {
	panic("not implemented")
}

func (sm *SnippetsMock) Rename(a types.Alias, new types.SnipName) (*types.Snippet, *types.SnipName, error) {
	sm.RenameCalledWith = []string{a.String(), new.String()}
	s := &types.Snippet{}
	s.Alias.SnipName = new
	return s, &new, nil
}

func (sm *SnippetsMock) Patch(a types.Alias, target string, patch string) (*types.Snippet, error) {
	sm.PatchCalledWith = []string{a.String(), target, patch}
	s := &types.Snippet{}
	s.Alias = a
	return s, nil
}

func (sm *SnippetsMock) Delete(username string, pouch string, names []*types.SnipName) error {
	sm.DeleteCalledWith = []interface{}{username, pouch, names}
	return nil
}

func (sm *SnippetsMock) Clone(a types.Alias, new types.Alias) (*types.Snippet, error) {
	sm.CloneCalledWith = []interface{}{a.String(), new.String()}
	return &types.Snippet{}, nil
}

func (sm *SnippetsMock) Tag(a types.Alias, tag ...string) (*types.Snippet, error) {
	m := map[string][]string{}
	m[a.String()] = tag
	sm.TagCalledWith = m
	return &types.Snippet{}, nil
}

func (sm *SnippetsMock) UnTag(a types.Alias, tag ...string) (*types.Snippet, error) {
	m := map[string][]string{}
	m[a.String()] = tag
	sm.UnTagCalledWith = m
	return &types.Snippet{}, nil
}

func (sm *SnippetsMock) List(l *models.ListParams) (*models.ListView, error) {
	sm.ListCalledWith = l
	return &models.ListView{}, nil
}

func (sm *SnippetsMock) CreatePouch(name string) (string, error) {
	sm.CreatePouchCalledWith = name
	return name, nil
}

func (sm *SnippetsMock) DeletePouch(name string) (bool, error) {
	sm.DeletePouchCalledWith = name
	return true, nil
}

func (sm *SnippetsMock) GetRoot(username string, all bool) (*models.ListView, error) {
	sm.GetRootCalledWith = []interface{}{username, all}
	return &models.ListView{}, nil
}
func (sm *SnippetsMock) RenamePouch(pouch string, newPouch string) (string, error) {
	sm.RenamePouchCalledWith = []string{pouch, newPouch}
	return pouch, nil
}
func (sm *SnippetsMock) MakePrivate(pouch string, private bool) (bool, error) {
	sm.MakePrivateCalledWith = []interface{}{pouch, private}
	return true, nil
}

type UsersMock struct {
	GetCalled        bool
	LoginCalledWith  []string
	SignupCalledWith []string
	SignoutCalled    bool
	GetCalledWith    string
	SignInResponse   *models.User
}

func (a *UsersMock) Get() (*models.User, error) {
	a.GetCalled = true
	return &models.User{}, nil
}

func (a *UsersMock) SignIn(username string, password string) (*models.User, error) {
	a.LoginCalledWith = []string{username, password}
	if a.SignInResponse == nil {
		a.SignInResponse = &models.User{}
	}
	return a.SignInResponse, nil
}

func (a *UsersMock) SignUp(email string, username string, password string, inviteCode string) (*models.User, error) {
	a.SignupCalledWith = []string{email, username, password}
	return &models.User{}, nil
}

func (a *UsersMock) Signout() error {
	a.SignoutCalled = true
	return nil
}

func (u *UsersMock) HasValidCredentials() bool {
	return false
}

func (u *UsersMock) ResetPassword(email string) (bool, error) {
	panic("not imple")
}

func (u *UsersMock) ChangePassword(p models.ChangePasswordParams) (bool, error) {
	panic("not imple")
}

type PersisterMock struct {
	GetCalledWith             []interface{}
	ChangeDirectoryCalledWith string
	UpsertCalledWith          []interface{}
	DeleteCalledWith          string
	GetHydrates               interface{}
	GetReturns                error
}

func (s *PersisterMock) DeleteAll() error {
	panic("ni")
}

func (s *PersisterMock) Delete(fullKey string) error {
	s.DeleteCalledWith = fullKey
	return nil
}

func (s *PersisterMock) Get(fullKey string, input interface{}, fresherThan int64) error {
	s.GetCalledWith = []interface{}{fullKey, input}
	if s.GetReturns != nil {
		return s.GetReturns
	}
	if s.GetHydrates != nil {
		v1 := reflect.ValueOf(input).Elem()
		v2 := reflect.ValueOf(s.GetHydrates).Elem()
		v1.Set(v2)
	}
	return nil
}

func (s *PersisterMock) Upsert(dir string, data interface{}) error {
	s.UpsertCalledWith = []interface{}{dir, data}
	return nil
}
