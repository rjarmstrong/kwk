package aliases

import "github.com/kwk-links/kwk-cli/libs/models"

type IAliases interface {
	Create(uri string, fullKey string) *models.Alias
	Rename(fullKey string, newFullKey string) *models.Alias
	Patch(fullKey string, uri string) *models.Alias
	Delete(fullKey string)
	Clone(fullKey string, newKey string) *models.Alias
	Tag(fullKey string, tag ...string) *models.Alias
	UnTag(fullKey string, tag ...string) *models.Alias
	Get(fullKey string) *models.AliasList
	List(args []string) *models.AliasList
}
