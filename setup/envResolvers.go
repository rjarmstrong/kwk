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
	snippets       snippets.Service
	system         sys.Manager
	account        account.Manager
	hostConfigName string
}

func NewEnvResolvers(s snippets.Service, sys sys.Manager, a account.Manager) Resolvers {
	return &EnvResolvers{
		hostConfigName:GetHostConfigFullName("env.yml"),
		snippets:      s,
		system:        sys,
		account:       a,
	}
}

func (e *EnvResolvers) Anon() (string, error) {
	return "", nil
}

func (e *EnvResolvers) Local() (string, error) {
	return e.system.ReadFromFile(SNIP_CACHE_PATH, e.hostConfigName, true, 0)
}

func (e *EnvResolvers) Own() (string, error) {
	if u, err := e.account.Get(); err != nil {
		return "", err
	} else {
		if l, err := e.snippets.Get(&models.Alias{FullKey: e.hostConfigName, Username: u.Username }); err != nil {
			return "", err
		} else {
			if _, err := e.system.WriteToFile(SNIP_CACHE_PATH, e.hostConfigName, l.Items[0].Snip, true); err != nil {
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
		if _, err := e.system.WriteToFile(SNIP_CACHE_PATH, e.hostConfigName, env, true); err != nil {
			return "", err
		} else {
			return env, nil
		}
	}
}

func (e *EnvResolvers) Fallback() (string, error) {
	defaultEnv := fmt.Sprintf("%s-%s.yml", runtime.GOOS, runtime.GOARCH)
	defaultAlias := &models.Alias{FullKey:defaultEnv, Username:"env"}
	if snip, err := e.snippets.Clone(defaultAlias, e.hostConfigName); err != nil {
		return "", err
	} else {
		return snip.Snip, nil
	}
}
