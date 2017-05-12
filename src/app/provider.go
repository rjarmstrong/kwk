package app

const (
	snipCachePath = "snip-cache"
)

type Provider interface {
	//Load loads preferences and environments.
	Load()
}

type Resolvers interface {
	Anon() (string, error)
	Local() (string, error)
	Own() (string, error)
	Default() (string, error)
	Fallback() (string, error)
}

type ProviderMock struct {
}

func (ProviderMock) Load() {
}
