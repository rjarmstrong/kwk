package users

import "github.com/kwk-links/kwk-cli/libs/models"

type UsersMock struct {
	GetCalled        bool
	LoginCalledWith  []string
	SignupCalledWith []string
	SignoutCalled    bool
	GetCalledWith    string
}

func (a *UsersMock) Get() (*models.User, error){
	a.GetCalled = true
	return &models.User{}, nil
}

func (a *UsersMock) SignIn(username string, password string) (*models.User, error) {
	a.LoginCalledWith = []string{username, password}
	return &models.User{}, nil
}

func (a *UsersMock) SignUp(email string, username string, password string) (*models.User, error) {
	a.SignupCalledWith = []string{email, username, password }
	return &models.User{}, nil
}

func (a *UsersMock) Signout() {
	a.SignoutCalled = true
}