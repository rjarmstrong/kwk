package rpc

import "bitbucket.com/sharingmachine/kwkcli/models"

type Sys interface {
	GetApiInfo() (*models.InfoResponse, error)
	TestAppError() (error)
	TestTransportError() (error)
}
