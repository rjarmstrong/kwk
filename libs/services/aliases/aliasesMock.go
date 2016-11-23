package aliases

import "bitbucket.com/sharingmachine/kwkcli/libs/models"

type AliasesMock struct {
	GetCalledWith    *models.KwkKey
	RenameCalledWith []string
	CreateCalledWith   []string
	ReturnItemsForGet  []models.Alias
	PatchCalledWith []string
	DeleteCalledWith string
	CloneCalledWith []interface{}
	TagCalledWith map[string][]string
	UnTagCalledWith map[string][]string
	ListCalledWith []interface{}
}

func (a *AliasesMock) Get(k *models.KwkKey) (*models.AliasList, error) {
	a.GetCalledWith = k
	return &models.AliasList{Items:a.ReturnItemsForGet, Total:int64(len(a.ReturnItemsForGet))}, nil
}

func (a *AliasesMock) Create(uri string, fullKey string) (*models.CreateAlias, error) {
	a.CreateCalledWith = []string{uri, fullKey}
	if fullKey == "" {
		fullKey = "x5hi23"
	}
	return &models.CreateAlias{Alias: &models.Alias{FullKey:fullKey}}, nil
}

func (a *AliasesMock) Update(fullKey string, description string) (*models.Alias, error){
	panic("not implemented")
}

func (a *AliasesMock) Rename(fullKey string, newFullKey string) (*models.Alias, string, error) {
	a.RenameCalledWith = []string{fullKey, newFullKey}
	return &models.Alias{FullKey:newFullKey}, fullKey, nil
}

func (a *AliasesMock) Patch(fullKey string, target string, patch string) (*models.Alias, error){
	a.PatchCalledWith = []string{fullKey, target, patch}
	return &models.Alias{FullKey:fullKey, Uri:patch}, nil
}

func (a *AliasesMock) Delete(fullKey string) error {
	a.DeleteCalledWith = fullKey
	return nil
}

func (a *AliasesMock) Clone(k *models.KwkKey, newKey string) (*models.Alias, error) {
	a.CloneCalledWith = []interface{}{k,newKey}
	return &models.Alias{},nil
}

func (a *AliasesMock) Tag(fullKey string, tag ...string) (*models.Alias, error){
	m := map[string][]string{}
	m[fullKey] = tag
	a.TagCalledWith = m
	return &models.Alias{},nil
}

func (a *AliasesMock) UnTag(fullKey string, tag ...string) (*models.Alias, error){
	m := map[string][]string{}
	m[fullKey] = tag
	a.UnTagCalledWith = m
	return &models.Alias{},nil
}

func (a *AliasesMock) List(username string, page int64, size int64, tags ...string) (*models.AliasList, error) {
	a.ListCalledWith = []interface{}{username, page, size, tags}
	return &models.AliasList{}, nil
}