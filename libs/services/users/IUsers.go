package users

import "github.com/kwk-links/kwk-cli/libs/models"

type IUsers interface {
	Login(username string, password string) *models.User
	SignUp(email string, username string, password string) *models.User
	Get() *models.User
	Signout()
}
