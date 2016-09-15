package api

import (
	"github.com/kwk-links/kwk-cli/libs/system"
)

type ApiMock struct {
	PrintProfileCalled bool
	LoginCalledWith    []string
	SignupCalledWith    []string
	SignoutCalled bool
}

func (a *ApiMock) PrintProfile() {
	a.PrintProfileCalled = true
}

func (a *ApiMock) Login(username string, password string) *system.User {
	a.LoginCalledWith = []string{username, password}
	return &system.User{}
}

func (a *ApiMock) SignUp(email string, username string, password string) *system.User {
	a.SignupCalledWith = []string{ email, username, password }
	return &system.User{}
}

func (a *ApiMock) Signout() {
	a.SignoutCalled = true
}