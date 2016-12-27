package openers

import "bitbucket.com/sharingmachine/kwkcli/models"

type OpenerMock struct {
	OpenCalledWith []interface{}
	EditCalledWith *models.Snippet
}

func (o *OpenerMock) Open(alias *models.Snippet, args []string) error {
	o.OpenCalledWith = []interface{}{alias, args}
	return nil
}

func (o *OpenerMock) Edit(alias *models.Snippet) error {
	o.EditCalledWith = alias
	return nil
}
