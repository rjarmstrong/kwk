package aliases

import "bitbucket.com/sharingmachine/kwkcli/libs/models"

type IAliases interface {
	Create(uri string, fullKey string) (*models.CreateAlias, error)
	Update(fullKey string, description string) (*models.Alias, error)
	Rename(fullKey string, newFullKey string) (*models.Alias, string, error)
	Patch(fullKey string, target string, patch string) (*models.Alias, error)
	Delete(fullKey string) error
	Clone(k *models.KwkKey, newKey string) (*models.Alias, error)
	Tag(fullKey string, tag ...string) (*models.Alias, error)
	UnTag(fullKey string, tag ...string) (*models.Alias, error)
	Get(k *models.KwkKey) (*models.AliasList, error)
	List(username string, page int64, size int64, tags ...string) (*models.AliasList, error)
}
