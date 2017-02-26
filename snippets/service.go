package snippets

import "bitbucket.com/sharingmachine/kwkcli/models"

type Service interface {
	Create(snip string, a models.Alias, role models.SnipRole) (*models.CreateSnippetResponse, error)
	Update(a models.Alias, description string) (*models.Snippet, error)
	Rename(a models.Alias, new models.SnipName) (*models.Snippet, *models.SnipName, error)
	Patch(a models.Alias, target string, patch string) (*models.Snippet, error)
	Delete(username string, pouch string, names []*models.SnipName) error
	Move(username string, sourcePouch string, targetPouch string, names []*models.SnipName) (string, error)
	Clone(a models.Alias, new models.Alias) (*models.Snippet, error)
	Tag(a models.Alias, tag ...string) (*models.Snippet, error)
	UnTag(a models.Alias, tag ...string) (*models.Snippet, error)
	Get(a models.Alias) (*models.ListView, error)
	List(l *models.ListParams) (*models.ListView, error)
	AlphaSearch(term string) (*models.SearchTermResponse, error)
	LogRun(a models.Alias, s models.RunStatus)
	SetPreview(a models.Alias, p string) error

	GetRoot (username string, all bool) (*models.ListView, error)
	CreatePouch (pouch string) (string, error)
	RenamePouch (pouch string, newPouch string) (string, error)
	MakePrivate (pouch string, private bool) (bool, error)
	DeletePouch (pouch string) (bool, error)
}
