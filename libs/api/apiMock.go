package api

import (
	"github.com/kwk-links/kwk-cli/libs/system"
)

type ApiMock struct {
	PrintProfileCalled bool
	LoginCalledWith    []string
	SignupCalledWith   []string
	SignoutCalled      bool
	GetCalledWith      string
	RenameCalledWith []string
	CreateCalledWith   []string
	ReturnItemsForGet  []Alias
	PatchCalledWith []string
	DeleteCalledWith string
	CloneCalledWith []string
}

func (a *ApiMock) PrintProfile() {
	a.PrintProfileCalled = true
}

func (a *ApiMock) Login(username string, password string) *system.User {
	a.LoginCalledWith = []string{username, password}
	return &system.User{}
}

func (a *ApiMock) SignUp(email string, username string, password string) *system.User {
	a.SignupCalledWith = []string{email, username, password }
	return &system.User{}
}

func (a *ApiMock) Signout() {
	a.SignoutCalled = true
}

func (a *ApiMock) Get(fullKey string) *AliasList {
	a.GetCalledWith = fullKey
	return &AliasList{Items:a.ReturnItemsForGet, Total:len(a.ReturnItemsForGet)}
}

func (a *ApiMock) Create(uri string, fullKey string) *Alias {
	a.CreateCalledWith = []string{uri, fullKey}
	if fullKey == "" {
		fullKey = "x5hi23"
	}
	return &Alias{FullKey:fullKey}
}

func (a *ApiMock) Rename(fullKey string, newFullKey string) *Alias {
	a.RenameCalledWith = []string{fullKey, newFullKey}
	return &Alias{FullKey:newFullKey}
}

func (a *ApiMock) Patch(fullKey string, uri string) *Alias {
	a.PatchCalledWith = []string{fullKey, uri}
	return &Alias{FullKey:fullKey, Uri:uri}
}

func (a *ApiMock) Delete(fullKey string) {
	a.DeleteCalledWith = fullKey
}

func (a *ApiMock) Clone(fullKey string, newKey string) *Alias {
	a.CloneCalledWith = []string{fullKey,newKey}
	return &Alias{}
}