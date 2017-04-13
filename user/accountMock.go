package user

import "bitbucket.com/sharingmachine/kwkcli/models"

type AccountMock struct {
	GetCalled        bool
	LoginCalledWith  []string
	SignupCalledWith []string
	SignoutCalled    bool
	GetCalledWith    string
	SignInResponse   *models.User
}

func (a *AccountMock) Get() (*models.User, error) {
	a.GetCalled = true
	return &models.User{}, nil
}

func (a *AccountMock) SignIn(username string, password string) (*models.User, error) {
	a.LoginCalledWith = []string{username, password}
	if a.SignInResponse == nil {
		a.SignInResponse = &models.User{}
	}
	return a.SignInResponse, nil
}

func (a *AccountMock) SignUp(email string, username string, password string, inviteCode string) (*models.User, error) {
	a.SignupCalledWith = []string{email, username, password}
	return &models.User{}, nil
}

func (a *AccountMock) Signout() error {
	a.SignoutCalled = true
	return nil
}

func (u *AccountMock) HasValidCredentials() bool {
	return false
}

func (u *AccountMock) ResetPassword(email string) (bool, error){
	panic("not imple")
}

func (u *AccountMock) ChangePassword(p models.ChangePasswordParams) (bool, error) {
	panic("not imple")
}
