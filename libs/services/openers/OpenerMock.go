package openers

import "github.com/kwk-links/kwk-cli/libs/models"

type OpenerMock struct {
	OpenCalledWith []interface{}
	EditCalledWith *models.Alias
	EditError error
}

func (o *OpenerMock) Open(alias *models.Alias, args []string) error {
	o.OpenCalledWith = []interface{}{alias, args}
	return nil
}

func (o *OpenerMock) Edit(alias *models.Alias) error {
	o.EditCalledWith = alias
	if o.EditError != nil {
		return o.EditError
	}
	return nil
}