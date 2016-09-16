package openers

import "github.com/kwk-links/kwk-cli/libs/api"

type OpenerMock struct {
	OpenCalledWith []interface{}
	EditCalledWith string
	EditError error
}

func (o *OpenerMock) Open(alias *api.Alias, args []string) {
	o.OpenCalledWith = []interface{}{alias, args}
}

func (o *OpenerMock) Edit(key string) error {
	o.EditCalledWith = key
	if o.EditError != nil {
		return o.EditError
	}
	return nil
}