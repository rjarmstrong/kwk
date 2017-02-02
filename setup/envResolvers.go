package setup

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"runtime"
	"fmt"
)

// TODO: check yml version is compatible with this build else force upgrade.

type EnvResolvers struct {
	snippets snippets.Service
	system   sys.Manager
	account  account.Manager
	alias    models.Alias
}

func NewEnvResolvers(s snippets.Service, sys sys.Manager, a account.Manager) Resolvers {
	return &EnvResolvers{
		alias:*models.NewSetupAlias("env", "yml"),
		snippets:      s,
		system:        sys,
		account:       a,
	}
}

func (e *EnvResolvers) Anon() (string, error) {
	return "", nil
}

func (e *EnvResolvers) Local() (string, error) {
	return e.system.ReadFromFile(SNIP_CACHE_PATH, e.alias.Path(), true, 0)
}

func (e *EnvResolvers) Own() (string, error) {
	if u, err := e.account.Get(); err != nil {
		return "", err
	} else {
		if l, err := e.snippets.Get(*models.NewAlias(u.Username, e.alias.Pouch, e.alias.Name, e.alias.Ext)); err != nil {
			return "", err
		} else {
			if _, err := e.system.WriteToFile(SNIP_CACHE_PATH, e.alias.Path(), l.Items[0].Snip, true); err != nil {
				return "", err
			}
			return l.Items[0].Snip, nil
		}
	}
}

func (e *EnvResolvers) Default() (string, error) {
	if env, err := e.Fallback(); err != nil {
		return "", err
	} else {
		if _, err := e.system.WriteToFile(SNIP_CACHE_PATH, e.alias.Path(), env, true); err != nil {
			return "", err
		} else {
			return env, nil
		}
	}
}

func (e *EnvResolvers) Fallback() (string, error) {
	n := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	defaultA := models.NewAlias("env", models.ROOT_POUCH, n, "yml")
	local := models.NewAlias("", models.SETUP_POUCH, n, "yml")
	if snip, err := e.snippets.Clone(*defaultA, *local); err != nil {
		return "", err
	} else {
		return snip.Snip, nil
	}
}
