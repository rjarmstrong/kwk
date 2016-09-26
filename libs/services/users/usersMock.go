package users

import "bitbucket.com/sharingmachine/kwkcli/libs/models"

type UsersMock struct {
	GetCalled        bool
	LoginCalledWith  []string
	SignupCalledWith []string
	SignoutCalled    bool
	GetCalledWith    string
	SignInResponse *models.User
}

func (a *UsersMock) Get() (*models.User, error){
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

func (a *UsersMock) SignUp(email string, username string, password string) (*models.User, error) {
	a.SignupCalledWith = []string{email, username, password }
	return &models.User{}, nil
}

func (a *UsersMock) Signout() error {
	a.SignoutCalled = true
	return nil
}