package cli

import (
	"github.com/rjarmstrong/kwk-types"
	"golang.org/x/net/context"
)

// ContextFunc is a function which returns the Context. Suitable for late binding.
type ContextFunc func() context.Context

// RootPrinter prints the root directory of a given username. Listing pouches and
// recent snippets.
type RootPrinter func(r *types.RootResponse) error

const (
	DefaultApiHost     = "api.kwk.co:443"
	DefaultTestApiHost = "localhost:8000"
	DocPath            = "docs"
)

// AppConfig are runtime configuration options for this application mainly used for
// diagnostics.
type AppConfig struct {
	CpuProfile bool   `default:"false" json:"KWK_PROFILE"`
	Debug      bool   `default:"false" json:"KWK_DEBUG"`
	APIHost    string `default:"api.kwk.co:443" json:"KWK_APIHOST"`
	TestMode   bool   `default:"false" json:"KWK_TESTMODE"`
}

// UserWithToken is a way to store the user with the kwk access token.
type UserWithToken struct {
	AccessToken string `json:"access_token"`
	User        types.User
}

// HasAccessToken - does the current user have a non-empty access token?
func (m *UserWithToken) HasAccessToken() bool {
	return m.AccessToken != ""
}
