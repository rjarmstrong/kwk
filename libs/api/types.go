package api

import "github.com/kwk-links/kwk-cli/libs/system"

type IApi interface {
	PrintProfile()
	Login(username string, password string) *system.User
	SignUp(email string, username string, password string) *system.User
	Signout()
}
