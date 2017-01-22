package setup

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"os"
)

const (
	SNIP_CACHE_PATH = "snip-cache"
)

type Provider interface {
	Env() *yaml.MapSlice
	Prefs() *Preferences
	Preload()
}

func GetHostConfigFullName(fullName string) string {
	if h, err := os.Hostname(); err != nil {
		panic(err)
	} else {
		return fmt.Sprintf(".%s_%s", h, fullName)
	}
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

