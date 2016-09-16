package openers

import "github.com/kwk-links/kwk-cli/libs/api"

type OpenerMock struct {
	OpenCalledWith []interface{}
}

func (o *OpenerMock) Open(link *api.KwkLink, args []string) {
	o.OpenCalledWith = []interface{}{link, args}
}