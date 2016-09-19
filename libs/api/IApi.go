package api

import (
	"github.com/kwk-links/kwk-cli/libs/models"
)

type IApi interface {
	//migrate to integration
	PrintProfile()
	Login(username string, password string) *models.User
	SignUp(email string, username string, password string) *models.User
	Signout()
	Create(uri string, fullPath string) *Alias
	Rename(fullKey string, newFullKey string) *Alias
	Patch(fullKey string, uri string) *Alias
	Delete(fullKey string)
	Clone(fullKey string, newKey string) *Alias
	Tag(fullKey string, tag ...string) *Alias
	UnTag(fullKey string, tag ...string) *Alias
	Get(fullKey string) *AliasList
	List(args []string) *AliasList
}
