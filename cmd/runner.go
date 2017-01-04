package cmd

import "bitbucket.com/sharingmachine/kwkcli/models"

type Runner interface {
	Run(alias *models.Snippet, args []string) error
	Edit(alias *models.Snippet) error
}
