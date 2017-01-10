package rpc

import "bitbucket.com/sharingmachine/kwkcli/models"

type SysMock struct {

}

func (*SysMock) GetApiInfo() (*models.InfoResponse, error) {
	return &models.InfoResponse{}, nil
}

func (*SysMock) TestAppError(bool) (error) {
	panic("implement me")
}

func (*SysMock) TestTransportError() (error) {
	panic("implement me")
}
