package rpc

import "bitbucket.com/sharingmachine/kwkcli/models"

type Sys interface {
	GetApiInfo() (*models.InfoResponse, error)
	TestAppError(multi bool) (error)
	TestTransportError() (error)
}
