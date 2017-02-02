package setup

import (
	"gopkg.in/yaml.v2"
)

const (
	SNIP_CACHE_PATH = "snip-cache"
)

type Provider interface {
	Env() *yaml.MapSlice
	Prefs() *Preferences
	Preload()
}


type Resolvers interface{
	Anon() (string, error)
	Local() (string, error)
	Own() (string, error)
	Default() (string, error)
	Fallback() (string, error)
}

type ProviderMock struct {

}

func (ProviderMock) Env() *yaml.MapSlice {
	return &yaml.MapSlice{}
}

func (ProviderMock) Prefs() *Preferences {
	return DefaultPrefs()
}

func (ProviderMock) Preload() {
}

