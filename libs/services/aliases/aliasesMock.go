package aliases

import "bitbucket.com/sharingmachine/kwkcli/libs/models"

type AliasesMock struct {
	GetCalledWith      string
	RenameCalledWith []string
	CreateCalledWith   []string
	ReturnItemsForGet  []models.Alias
	PatchCalledWith []string
	DeleteCalledWith string
	CloneCalledWith []string
	TagCalledWith map[string][]string
	UnTagCalledWith map[string][]string
	ListCalledWith []interface{}
}

func (a *AliasesMock) Get(fullKey string) (*models.AliasList, error) {
	a.GetCalledWith = fullKey
	return &models.AliasList{Items:a.ReturnItemsForGet, Total:int32(len(a.ReturnItemsForGet))}, nil
}

func (a *AliasesMock) Create(uri string, fullKey string) (*models.CreateAlias, error) {
	a.CreateCalledWith = []string{uri, fullKey}
	if fullKey == "" {
		fullKey = "x5hi23"
	}
	return &models.CreateAlias{Alias: &models.Alias{FullKey:fullKey}}, nil
}

func (a *AliasesMock) Rename(fullKey string, newFullKey string) (*models.Alias, error) {
	a.RenameCalledWith = []string{fullKey, newFullKey}
	return &models.Alias{FullKey:newFullKey}, nil
}

func (a *AliasesMock) Patch(fullKey string, uri string) (*models.Alias, error){
	a.PatchCalledWith = []string{fullKey, uri}
	return &models.Alias{FullKey:fullKey, Uri:uri}, nil
}

func (a *AliasesMock) Delete(fullKey string) error {
	a.DeleteCalledWith = fullKey
	return nil
}

func (a *AliasesMock) Clone(fullKey string, newKey string) (*models.Alias, error) {
	a.CloneCalledWith = []string{fullKey,newKey}
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

func (a *AliasesMock) List(username string, page int32, size int32, tags ...string) (*models.AliasList, error) {
	a.ListCalledWith = []interface{}{username, page, size, tags}
	return &models.AliasList{}, nil
}