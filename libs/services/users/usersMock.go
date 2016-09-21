package users

import "github.com/kwk-links/kwk-cli/libs/models"

type UsersMock struct {
	GetCalled        bool
	LoginCalledWith  []string
	SignupCalledWith []string
	SignoutCalled    bool
	GetCalledWith    string
}

func (a *UsersMock) Get() {
	a.GetCalled = true
}

func (a *UsersMock) SignIn(username string, password string) *models.User {
	a.LoginCalledWith = []string{username, password}
	return &models.User{}
}

func (a *UsersMock) SignUp(email string, username string, password string) *models.User {
	a.SignupCalledWith = []string{email, username, password }
	return &models.User{}
}

func (a *UsersMock) Signout() {
	a.SignoutCalled = true
}