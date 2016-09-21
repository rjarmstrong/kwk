package openers

import (
	"github.com/kwk-links/kwk-cli/libs/models"
)

type IOpen interface {
	Open(alias *models.Alias, args []string) error
	Edit(alias *models.Alias) error
}
