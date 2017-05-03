package gokwk

import "bitbucket.com/sharingmachine/kwkcli/src/models"

type UsersMock struct {
	GetCalled        bool
	LoginCalledWith  []string
	SignupCalledWith []string
	SignoutCalled    bool
	GetCalledWith    string
	SignInResponse   *models.User
}

func (a *UsersMock) Get() (*models.User, error) {
	a.GetCalled = true
	return &models.User{}, nil
}

func (a *UsersMock) SignIn(username string, password string) (*models.User, error) {
	a.LoginCalledWith = []string{username, password}
	if a.SignInResponse == nil {
		a.SignInResponse = &models.User{}
	}
	return a.SignInResponse, nil
}

func (a *UsersMock) SignUp(email string, username string, password string, inviteCode string) (*models.User, error) {
	a.SignupCalledWith = []string{email, username, password}
	return &models.User{}, nil
}

func (a *UsersMock) Signout() error {
	a.SignoutCalled = true
	return nil
}

func (u *UsersMock) HasValidCredentials() bool {
	return false
}

func (u *UsersMock) ResetPassword(email string) (bool, error){
	panic("not imple")
}

func (u *UsersMock) ChangePassword(p models.ChangePasswordParams) (bool, error) {
	panic("not imple")
}
