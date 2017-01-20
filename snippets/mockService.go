package snippets

import "bitbucket.com/sharingmachine/kwkcli/models"

type ServiceMock struct {
	GetCalledWith     *models.Alias
	RenameCalledWith  []string
	CreateCalledWith  []string
	ReturnItemsForGet []models.Snippet
	PatchCalledWith   []string
	DeleteCalledWith  string
	CloneCalledWith   []interface{}
	TagCalledWith     map[string][]string
	UnTagCalledWith   map[string][]string
	ListCalledWith    *models.ListParams
}

func (a *ServiceMock) Get(k *models.Alias) (*models.SnippetList, error) {
	a.GetCalledWith = k
	return &models.SnippetList{Items: a.ReturnItemsForGet, Total: int64(len(a.ReturnItemsForGet))}, nil
}

func (a *ServiceMock) Create(uri string, fullKey string, role models.SnipRole) (*models.CreateSnippetResponse, error) {
	a.CreateCalledWith = []string{uri, fullKey}
	if fullKey == "" {
		fullKey = "x5hi23"
	}
	return &models.CreateSnippetResponse{Snippet: &models.Snippet{FullName: fullKey}}, nil
}

func (a *ServiceMock) Update(fullKey string, description string) (*models.Snippet, error) {
	panic("not implemented")
}

func (a *ServiceMock) Rename(fullKey string, newFullKey string) (*models.Snippet, string, error) {
	a.RenameCalledWith = []string{fullKey, newFullKey}
	return &models.Snippet{FullName: newFullKey}, fullKey, nil
}

func (a *ServiceMock) Patch(fullKey string, target string, patch string) (*models.Snippet, error) {
	a.PatchCalledWith = []string{fullKey, target, patch}
	return &models.Snippet{FullName: fullKey, Snip: patch}, nil
}

func (a *ServiceMock) Delete(fullKey string) error {
	a.DeleteCalledWith = fullKey
	return nil
}

func (a *ServiceMock) Clone(k *models.Alias, newKey string) (*models.Snippet, error) {
	a.CloneCalledWith = []interface{}{k, newKey}
	return &models.Snippet{}, nil
}

func (a *ServiceMock) Tag(fullKey string, tag ...string) (*models.Snippet, error) {
	m := map[string][]string{}
	m[fullKey] = tag
	a.TagCalledWith = m
	return &models.Snippet{}, nil
}

func (a *ServiceMock) UnTag(fullKey string, tag ...string) (*models.Snippet, error) {
	m := map[string][]string{}
	m[fullKey] = tag
	a.UnTagCalledWith = m
	return &models.Snippet{}, nil
}

func (a *ServiceMock) List(l *models.ListParams) (*models.SnippetList, error) {
	a.ListCalledWith = l
	return &models.SnippetList{}, nil
}
