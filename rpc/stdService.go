package rpc

import (
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/rpc/src/sysRpc"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"google.golang.org/grpc"
)

type stdService struct {
	Settings config.Persister
	client   sysRpc.SysRpcClient
	headers  *Headers
}


func New(conn *grpc.ClientConn, h *Headers) *stdService {
	return &stdService{client: sysRpc.NewSysRpcClient(conn), headers: h}
}

func (s *stdService) GetApiInfo() (*models.InfoResponse, error) {
	if r, err := s.client.GetApiInfo(s.headers.Context(), &sysRpc.InfoRequest{}); err != nil {
		return nil, err
	} else {
		return &models.InfoResponse{Build:r.Build, Version:r.Version, Revision: r.Revision}, nil
	}
}

func (s *stdService) TestAppError(multi bool) (error) {
	request := &sysRpc.ErrorRequest{}
	request.Multi = multi
	if _, err := s.client.TestAppError(s.headers.Context(), request); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *stdService) TestTransportError() (error) {
	if _, err := s.client.TestTransportError(s.headers.Context(), &sysRpc.ErrorRequest{}); err != nil {
		return err
	} else {
		return nil
	}
}