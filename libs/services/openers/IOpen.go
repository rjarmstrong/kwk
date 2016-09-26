package openers

import "bitbucket.com/sharingmachine/kwkcli/libs/models"

type IOpen interface {
	Open(alias *models.Alias, args []string) error
	Edit(alias *models.Alias) error
}
