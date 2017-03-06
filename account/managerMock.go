package account

import "bitbucket.com/sharingmachine/kwkcli/models"

type ManagerMock struct {
	GetCalled        bool
	LoginCalledWith  []string
	SignupCalledWith []string
	SignoutCalled    bool
	GetCalledWith    string
	SignInResponse   *models.User
}

func (a *ManagerMock) Get() (*models.User, error) {
	a.GetCalled = true
	return &models.User{}, nil
}

func (a *ManagerMock) SignIn(username string, password string) (*models.User, error) {
	a.LoginCalledWith = []string{username, password}
	if a.SignInResponse == nil {
		a.SignInResponse = &models.User{}
	}
	return a.SignInResponse, nil
}

func (a *ManagerMock) SignUp(email string, username string, password string, inviteCode string) (*models.User, error) {
	a.SignupCalledWith = []string{email, username, password}
	return &models.User{}, nil
}

func (a *ManagerMock) Signout() error {
	a.SignoutCalled = true
	return nil
}

func (u *ManagerMock) HasValidCredentials() bool {
	return false
}

func (u *ManagerMock) ResetPassword(email string) (bool, error){
	panic("not imple")
}

func (u *ManagerMock) ChangePassword(p models.ChangePasswordParams) (bool, error) {
	panic("not imple")
}
