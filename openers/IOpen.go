package openers

import "bitbucket.com/sharingmachine/kwkcli/models"

type IOpen interface {
	Open(alias *models.Snippet, args []string) error
	Edit(alias *models.Snippet) error
}
