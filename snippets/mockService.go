package snippets

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
)

type ServiceMock struct {
	GetCalledWith         models.Alias
	RenameCalledWith      []string
	CreateCalledWith      []string
	ReturnItemsForGet     []*models.Snippet
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

func (sm *ServiceMock) LogRun(a models.Alias, s models.RunStatus) {
	panic("not impl")
}

func (sm *ServiceMock) AlphaSearch(term string) (*models.SearchTermResponse, error) {
	panic("not impl")
}

func (sm *ServiceMock) Move(username string, sourcePouch string, targetPouch string, names []*models.SnipName) (string, error) {
	panic("not imp")
}

func (sm *ServiceMock) Get(a models.Alias) (*models.ListView, error) {
	sm.GetCalledWith = a
	return &models.ListView{Snippets: sm.ReturnItemsForGet, Total: int64(len(sm.ReturnItemsForGet))}, nil
}

func (sm *ServiceMock) Create(snip string, a models.Alias, role models.SnipRole) (*models.CreateSnippetResponse, error) {
	sm.CreateCalledWith = []string{snip, a.String()}
	if a.Name == "" {
		a.Name = "x5hi23"
	}
	return &models.CreateSnippetResponse{Snippet: &models.Snippet{Alias: a}}, nil
}

func (sm *ServiceMock) Update(a models.Alias, description string) (*models.Snippet, error) {
	panic("not implemented")
}

func (sm *ServiceMock) Rename(a models.Alias, new models.SnipName) (*models.Snippet, *models.SnipName, error) {
	sm.RenameCalledWith = []string{a.String(), new.String()}
	s := models.NewSnippet("")
	s.Alias.SnipName = new
	return s, &new, nil
}

func (sm *ServiceMock) Patch(a models.Alias, target string, patch string) (*models.Snippet, error) {
	sm.PatchCalledWith = []string{a.String(), target, patch}
	s := models.NewSnippet("")
	s.Alias = a
	return s, nil
}

func (sm *ServiceMock) Delete(username string, pouch string, names []*models.SnipName) error {
	sm.DeleteCalledWith = []interface{}{username, pouch, names}
	return nil
}

func (sm *ServiceMock) Clone(a models.Alias, new models.Alias) (*models.Snippet, error) {
	sm.CloneCalledWith = []interface{}{a.String(), new.String()}
	return &models.Snippet{}, nil
}

func (sm *ServiceMock) Tag(a models.Alias, tag ...string) (*models.Snippet, error) {
	m := map[string][]string{}
	m[a.String()] = tag
	sm.TagCalledWith = m
	return &models.Snippet{}, nil
}

func (sm *ServiceMock) UnTag(a models.Alias, tag ...string) (*models.Snippet, error) {
	m := map[string][]string{}
	m[a.String()] = tag
	sm.UnTagCalledWith = m
	return &models.Snippet{}, nil
}

func (sm *ServiceMock) List(l *models.ListParams) (*models.ListView, error) {
	sm.ListCalledWith = l
	return &models.ListView{}, nil
}

func (sm *ServiceMock) CreatePouch(name string) (string, error) {
	sm.CreatePouchCalledWith = name
	return name, nil
}

func (sm *ServiceMock) DeletePouch(name string) (bool, error) {
	sm.DeletePouchCalledWith = name
	return true, nil
}

func (sm *ServiceMock) GetRoot (username string, all bool) (*models.ListView, error){
	sm.GetRootCalledWith = []interface{}{username, all}
	return &models.ListView{}, nil
}
func (sm *ServiceMock) RenamePouch (pouch string, newPouch string) (string, error){
	sm.RenamePouchCalledWith = []string{pouch, newPouch}
	return pouch, nil
}
func (sm *ServiceMock) MakePrivate (pouch string, private bool) (bool, error){
	sm.MakePrivateCalledWith = []interface{}{pouch, private}
	return true, nil
}
