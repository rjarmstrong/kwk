package api

import "github.com/kwk-links/kwk-cli/libs/system"

type IApi interface {
	PrintProfile()
	Login(username string, password string) *system.User
	SignUp(email string, username string, password string) *system.User
	Signout()
	//Given a fullKey will get one or many aliases that match.
	Get(fullKey string) *KwkLinkList
	Create(uri string, fullPath string) *KwkLink
}
