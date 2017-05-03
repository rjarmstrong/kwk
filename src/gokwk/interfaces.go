package gokwk

import (
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	t "bitbucket.com/sharingmachine/types"
)

type Snippets interface {
	Create(snip string, a t.Alias, role t.SnipRole) (*models.CreateSnippetResponse, error)
	Update(a t.Alias, description string) (*t.Snippet, error)
	Rename(a t.Alias, new t.SnipName) (*t.Snippet, *t.SnipName, error)
	Patch(a t.Alias, target string, patch string) (*t.Snippet, error)
	Delete(username string, pouch string, names []*t.SnipName) error
	Move(username string, sourcePouch string, targetPouch string, names []*t.SnipName) (string, error)
	Clone(a t.Alias, new t.Alias) (*t.Snippet, error)
	Tag(a t.Alias, tag ...string) (*t.Snippet, error)
	UnTag(a t.Alias, tag ...string) (*t.Snippet, error)
	Get(a t.Alias) (*models.ListView, error)
	List(l *models.ListParams) (*models.ListView, error)
	AlphaSearch(term string) (*models.SearchTermResponse, error)
	LogUse(a t.Alias, s t.UseStatus, u t.UseType, ctx *UseContext)

	GetRoot (username string, all bool) (*models.ListView, error)
	CreatePouch (pouch string) (string, error)
	RenamePouch (pouch string, newPouch string) (string, error)
	MakePrivate (pouch string, private bool) (bool, error)
	DeletePouch (pouch string) (bool, error)
}

type Users interface {
	SignIn(username string, password string) (*models.User, error)
	SignUp(email string, username string, password string, inviteCode string) (*models.User, error)
	Get() (*models.User, error)
	Signout() error
	HasValidCredentials() bool
	ResetPassword(email string) (bool, error)
	ChangePassword(p models.ChangePasswordParams) (bool, error)
}
