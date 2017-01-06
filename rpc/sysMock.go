package rpc

import "bitbucket.com/sharingmachine/kwkcli/models"

type SysMock struct {

}

func (*SysMock) GetApiInfo() (*models.InfoResponse, error) {
	panic("implement me")
}

func (*SysMock) TestAppError() (error) {
	panic("implement me")
}

func (*SysMock) TestTransportError() (error) {
	panic("implement me")
}
