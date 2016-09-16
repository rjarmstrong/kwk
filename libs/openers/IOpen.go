package openers

import "github.com/kwk-links/kwk-cli/libs/api"

type IOpen interface {
	Open(link *api.Alias, args []string)
	Edit(key string) error
}
