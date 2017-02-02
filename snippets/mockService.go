package snippets

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
)

type ServiceMock struct {
	GetCalledWith     models.Alias
	RenameCalledWith  []string
	CreateCalledWith  []string
	ReturnItemsForGet []models.Snippet
	PatchCalledWith   []string
	DeleteCalledWith  []interface{}
	CloneCalledWith   []interface{}
	TagCalledWith     map[string][]string
	UnTagCalledWith   map[string][]string
	ListCalledWith    *models.ListParams
}

func (sm *ServiceMock) Move(username string, sourcePouch string, targetPouch string, names []*models.SnipName) (string, error) {
	panic("not imp")
}

func (sm *ServiceMock) Get(a models.Alias) (*models.SnippetList, error) {
	sm.GetCalledWith = a
	return &models.SnippetList{Items: sm.ReturnItemsForGet, Total: int64(len(sm.ReturnItemsForGet))}, nil
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

func (sm *ServiceMock) List(l *models.ListParams) (*models.SnippetList, error) {
	sm.ListCalledWith = l
	return &models.SnippetList{}, nil
}
