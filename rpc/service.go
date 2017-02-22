package rpc

import "bitbucket.com/sharingmachine/kwkcli/models"

type Service interface {
	GetApiInfo() (*models.ApiInfo, error)
	TestAppError(multi bool) (error)
	TestTransportError() (error)
}
