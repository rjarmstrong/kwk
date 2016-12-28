package snippets

import "bitbucket.com/sharingmachine/kwkcli/models"

type Mock struct {
	GetCalledWith     *models.KwkKey
	RenameCalledWith  []string
	CreateCalledWith  []string
	ReturnItemsForGet []models.Snippet
	PatchCalledWith   []string
	DeleteCalledWith  string
	CloneCalledWith   []interface{}
	TagCalledWith     map[string][]string
	UnTagCalledWith   map[string][]string
	ListCalledWith    []interface{}
}

func (a *Mock) Get(k *models.KwkKey) (*models.SnippetList, error) {
	a.GetCalledWith = k
	return &models.SnippetList{Items: a.ReturnItemsForGet, Total: int64(len(a.ReturnItemsForGet))}, nil
}

func (a *Mock) Create(uri string, fullKey string) (*models.CreateSnippet, error) {
	a.CreateCalledWith = []string{uri, fullKey}
	if fullKey == "" {
		fullKey = "x5hi23"
	}
	return &models.CreateSnippet{Snippet: &models.Snippet{FullKey: fullKey}}, nil
}

func (a *Mock) Update(fullKey string, description string) (*models.Snippet, error) {
	panic("not implemented")
}

func (a *Mock) Rename(fullKey string, newFullKey string) (*models.Snippet, string, error) {
	a.RenameCalledWith = []string{fullKey, newFullKey}
	return &models.Snippet{FullKey: newFullKey}, fullKey, nil
}

func (a *Mock) Patch(fullKey string, target string, patch string) (*models.Snippet, error) {
	a.PatchCalledWith = []string{fullKey, target, patch}
	return &models.Snippet{FullKey: fullKey, Snip: patch}, nil
}

func (a *Mock) Delete(fullKey string) error {
	a.DeleteCalledWith = fullKey
	return nil
}

func (a *Mock) Clone(k *models.KwkKey, newKey string) (*models.Snippet, error) {
	a.CloneCalledWith = []interface{}{k, newKey}
	return &models.Snippet{}, nil
}

func (a *Mock) Tag(fullKey string, tag ...string) (*models.Snippet, error) {
	m := map[string][]string{}
	m[fullKey] = tag
	a.TagCalledWith = m
	return &models.Snippet{}, nil
}

func (a *Mock) UnTag(fullKey string, tag ...string) (*models.Snippet, error) {
	m := map[string][]string{}
	m[fullKey] = tag
	a.UnTagCalledWith = m
	return &models.Snippet{}, nil
}

func (a *Mock) List(username string, size int64, since int64, tags ...string) (*models.SnippetList, error) {
	a.ListCalledWith = []interface{}{username, size, since, tags}
	return &models.SnippetList{}, nil
}
