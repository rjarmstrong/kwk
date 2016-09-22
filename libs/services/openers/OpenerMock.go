package openers

import "github.com/kwk-links/kwk-cli/libs/models"

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