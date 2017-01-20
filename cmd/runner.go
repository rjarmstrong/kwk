package cmd

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/config"
)

type Runner interface {
	Run(s *models.Snippet, args []string) error
	Edit(s *models.Snippet) error
	LoadPreferences() *config.PersistedPrefs
}

// GetDefaultConf is used to resolve configuration
// that has not been downloaded for this host.
type GetDefaultConf func () (string, error)