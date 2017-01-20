package cmd

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/config"
)

type RunnerMock struct {
	OpenCalledWith []interface{}
	EditCalledWith *models.Snippet
}

func (o *RunnerMock) Run(alias *models.Snippet, args []string) error {
	o.OpenCalledWith = []interface{}{alias, args}
	return nil
}

func (o *RunnerMock) Edit(alias *models.Snippet) error {
	o.EditCalledWith = alias
	return nil
}

func (o *RunnerMock) LoadPreferences() *config.PersistedPrefs {
	return nil
}
