package cli

import (
	"github.com/kwk-super-snippets/types"
	"golang.org/x/net/context"
)

type ContextFunc func() context.Context

type RootPrinter func(r *types.RootResponse) error

const (
	DefaultApiHost     = "api.kwk.co:443"
	DefaultTestApiHost = "localhost:8000"
)

type AppConfig struct {
	CpuProfile bool   `default:"false" json:"KWK_PROFILE"`
	Debug      bool   `default:"false" json:"KWK_DEBUG"`
	APIHost    string `default:"api.kwk.co:443" json:"KWK_APIHOST"`
	TestMode   bool   `default:"false" json:"KWK_TESTMODE"`
	DocPath    string `default:"docs" json:"KWK_DOCSPATH"` // Any other non-snippet files.
}

type UserWithToken struct {
	AccessToken string `json:"access_token"`
	User        types.User
}

func (m *UserWithToken) HasAccessToken() bool {
	return m.AccessToken != ""
}
