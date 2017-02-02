package cmd

import
(
	"bitbucket.com/sharingmachine/kwkcli/models"
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
