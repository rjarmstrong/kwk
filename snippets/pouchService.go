package snippets

import (
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/rpc/src/snipsRpc"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
)

type Pouches struct {
	persister config.Persister
	c         snipsRpc.PouchesRpcClient
	headers   *rpc.Headers
}

type Root struct {
	Pouches []models.Pouch
	Snippets []*models.Snippet
}

func (rt *Root) IsPouch(name string) bool {
	return false
}