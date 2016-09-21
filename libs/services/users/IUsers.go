package users

import "github.com/kwk-links/kwk-cli/libs/models"

type IUsers interface {
	SignIn(username string, password string) (*models.User, error)
	SignUp(email string, username string, password string) (*models.User, error)
	Get() (*models.User, error)
	Signout()
}
