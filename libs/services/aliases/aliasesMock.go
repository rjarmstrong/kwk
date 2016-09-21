package aliases

import "github.com/kwk-links/kwk-cli/libs/models"

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
	ListCalledWith []string
}

func (a *AliasesMock) Get(fullKey string) *models.AliasList {
	a.GetCalledWith = fullKey
	return &models.AliasList{Items:a.ReturnItemsForGet, Total:len(a.ReturnItemsForGet)}
}

func (a *AliasesMock) Create(uri string, fullKey string) *models.Alias {
	a.CreateCalledWith = []string{uri, fullKey}
	if fullKey == "" {
		fullKey = "x5hi23"
	}
	return &models.Alias{FullKey:fullKey}
}

func (a *AliasesMock) Rename(fullKey string, newFullKey string) *models.Alias {
	a.RenameCalledWith = []string{fullKey, newFullKey}
	return &models.Alias{FullKey:newFullKey}
}

func (a *AliasesMock) Patch(fullKey string, uri string) *models.Alias {
	a.PatchCalledWith = []string{fullKey, uri}
	return &models.Alias{FullKey:fullKey, Uri:uri}
}

func (a *AliasesMock) Delete(fullKey string) {
	a.DeleteCalledWith = fullKey
}

func (a *AliasesMock) Clone(fullKey string, newKey string) *models.Alias {
	a.CloneCalledWith = []string{fullKey,newKey}
	return &models.Alias{}
}

func (a *AliasesMock) Tag(fullKey string, tag ...string) *models.Alias {
	m := map[string][]string{}
	m[fullKey] = tag
	a.TagCalledWith = m
	return &models.Alias{}
}

func (a *AliasesMock) UnTag(fullKey string, tag ...string) *models.Alias {
	m := map[string][]string{}
	m[fullKey] = tag
	a.UnTagCalledWith = m
	return &models.Alias{}
}

func (a *AliasesMock) List(args []string) *models.AliasList {
	a.ListCalledWith = args
	return &models.AliasList{}
}