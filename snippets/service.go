package snippets

import "bitbucket.com/sharingmachine/kwkcli/models"

type Service interface {
	Create(uri string, fullKey string) (*models.CreateSnippet, error)
	Update(fullKey string, description string) (*models.Snippet, error)
	Rename(fullKey string, newFullKey string) (*models.Snippet, string, error)
	Patch(fullKey string, target string, patch string) (*models.Snippet, error)
	Delete(fullKey string) error
	Clone(k *models.KwkKey, newKey string) (*models.Snippet, error)
	Tag(fullKey string, tag ...string) (*models.Snippet, error)
	UnTag(fullKey string, tag ...string) (*models.Snippet, error)
	Get(k *models.KwkKey) (*models.SnippetList, error)
	List(username string, size int64, since int64, tags ...string) (*models.SnippetList, error)
}
