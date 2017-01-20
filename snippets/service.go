package snippets

import "bitbucket.com/sharingmachine/kwkcli/models"

type Service interface {
	Create(snip string, fullName string, role models.SnipRole) (*models.CreateSnippetResponse, error)
	Update(fullName string, description string) (*models.Snippet, error)
	Rename(fullName string, newFullName string) (*models.Snippet, string, error)
	Patch(fullName string, target string, patch string) (*models.Snippet, error)
	Delete(fullName string) error
	Clone(k *models.Alias, newName string) (*models.Snippet, error)
	Tag(fullName string, tag ...string) (*models.Snippet, error)
	UnTag(fullName string, tag ...string) (*models.Snippet, error)
	Get(k *models.Alias) (*models.SnippetList, error)
	List(l *models.ListParams) (*models.SnippetList, error)
}
