package openers

import "github.com/kwk-links/kwk-cli/libs/api"

type IOpen interface {
	Open(alias *api.Alias, args []string)
	Edit(alias *api.Alias) error
}
