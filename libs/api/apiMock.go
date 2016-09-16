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
	CreateCalledWith   []string
	ReturnItemsForGet  []KwkLink
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

func (a *ApiMock) Get(fullKey string) *KwkLinkList {
	a.GetCalledWith = fullKey
	return &KwkLinkList{Items:a.ReturnItemsForGet, Total:len(a.ReturnItemsForGet)}
}

func (a *ApiMock) Create(uri string, fullKey string) *KwkLink {
	a.CreateCalledWith = []string{uri, fullKey}
	if fullKey == "" {
		fullKey = "x5hi23"
	}
	return &KwkLink{FullKey:fullKey}
}