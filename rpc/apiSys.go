package rpc

import (
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/rpc/src/sysRpc"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"google.golang.org/grpc"
)

type ApiSys struct {
	Settings config.Settings
	client   sysRpc.SysRpcClient
	headers  *Headers
}


func New(conn *grpc.ClientConn, h *Headers) *ApiSys {
	return &ApiSys{client: sysRpc.NewSysRpcClient(conn), headers: h}
}

func (s *ApiSys) GetApiInfo() (*models.InfoResponse, error) {
	if _, err := s.client.GetApiInfo(s.headers.GetContext(), &sysRpc.InfoRequest{}); err != nil {
		return nil, err
	} else {
		return &models.InfoResponse{}, nil
	}
}

func (s *ApiSys) TestAppError(multi bool) (error) {
	request := &sysRpc.ErrorRequest{}
	request.Multi = multi
	if _, err := s.client.TestAppError(s.headers.GetContext(), request); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *ApiSys) TestTransportError() (error) {
	if _, err := s.client.TestTransportError(s.headers.GetContext(), &sysRpc.ErrorRequest{}); err != nil {
		return err
	} else {
		return nil
	}
}