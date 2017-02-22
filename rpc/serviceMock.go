package rpc

import "bitbucket.com/sharingmachine/kwkcli/models"

type SysMock struct {

}

func (*SysMock) GetApiInfo() (*models.ApiInfo, error) {
	return &models.ApiInfo{}, nil
}

func (*SysMock) TestAppError(bool) (error) {
	panic("implement me")
}

func (*SysMock) TestTransportError() (error) {
	panic("implement me")
}
