package openers

import "bitbucket.com/sharingmachine/kwkcli/libs/models"

type OpenerMock struct {
	OpenCalledWith []interface{}
	EditCalledWith *models.Alias
}

func (o *OpenerMock) Open(alias *models.Alias, args []string) error {
	o.OpenCalledWith = []interface{}{alias, args}
	return nil
}

func (o *OpenerMock) Edit(alias *models.Alias) error {
	o.EditCalledWith = alias
	return nil
}
