package cmd

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
)

type Runner interface {
	Run(s *models.Snippet, args []string) error
	Edit(s *models.Snippet) error
}