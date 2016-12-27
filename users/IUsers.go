package users

import "bitbucket.com/sharingmachine/kwkcli/models"

type IUsers interface {
	SignIn(username string, password string) (*models.User, error)
	SignUp(email string, username string, password string) (*models.User, error)
	Get() (*models.User, error)
	Signout() error
}
