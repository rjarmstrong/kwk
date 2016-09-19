package openers

import "github.com/kwk-links/kwk-cli/libs/api"

type OpenerMock struct {
	OpenCalledWith []interface{}
	EditCalledWith *api.Alias
	EditError error
}

func (o *OpenerMock) Open(alias *api.Alias, args []string) {
	o.OpenCalledWith = []interface{}{alias, args}
}

func (o *OpenerMock) Edit(alias *api.Alias) error {
	o.EditCalledWith = alias
	if o.EditError != nil {
		return o.EditError
	}
	return nil
}