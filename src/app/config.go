package app

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"gopkg.in/yaml.v2"
)

const (
	snipCachePath = "snip-cache"
)

type Provider interface {
	//Load loads preferences and environments.
	Load()
}

type ConfigResolver interface {
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


type ConfigProvider struct {
	u              types.UsersClient
	envResolvers   ConfigResolver
	prefsResolvers ConfigResolver
	errs.Handler
}

func NewConfig(ss types.SnippetsClient, f IO, u types.UsersClient, eh errs.Handler) Provider {
	env := NewEnvResolvers(ss, f)
	prefs := NewPrefsResolvers(ss, f)
	c := &ConfigProvider{envResolvers: env, prefsResolvers: prefs, u: u, Handler: eh}
	c.Load()
	return c
}

func (cs *ConfigProvider) Load() {
	cs.loadEnv()
	cs.loadPrefs()
}

func (cs *ConfigProvider) loadEnv() *yaml.MapSlice {
	if Env() != nil {
		return Env()
	}
	envString, err := cs.GetConfig(cs.envResolvers)
	if err != nil {
		out.Debug("Failed to load env settings.")
		out.LogErr(err)
		envString, _ = cs.envResolvers.Fallback()
	}
	env := &yaml.MapSlice{}
	err = yaml.Unmarshal([]byte(envString), env)
	if err != nil {
		cs.Handle(errs.New(errs.CodeInvalidConfigSection,
			fmt.Sprintf("Invalid kwk *env.yml detected. `kwk edit env` to fix. %s", err)))
		envString, _ = cs.envResolvers.Fallback()
		yaml.Unmarshal([]byte(envString), env)
	}
	SetEnv(env)
	return Env()
}

func (cs *ConfigProvider) loadPrefs() {
	if Prefs() != nil {
		return
	}
	c, err := cs.GetConfig(cs.prefsResolvers)
	if err != nil {
		cs.Handle(err)
		return
	}
	prefs := &Preferences{PersistedPrefs: PersistedPrefs{}}
	parse := func(p string) (*Preferences, error) {
		ph := &PreferencesHolder{}
		err := yaml.Unmarshal([]byte(p), ph)
		if err != nil {
			return nil, err
		}
		prefs.PersistedPrefs = ph.Preferences
		return prefs, nil
	}
	out.Debug("Loaded prefs:%+v", c)
	res, err := parse(c)
	if err != nil {
		cs.Handle(errs.New(errs.CodeInvalidConfigSection,
			fmt.Sprintf("Invalid kwk *prefs.yml detected. `kwk edit prefs` to fix. %s", err)))
		// The fallback is expected not to fail parsing
		fb, _ := cs.prefsResolvers.Fallback()
		res, _ = parse(fb)
	}
	out.Debug("SETTING PREFS: %+v", res)
	SetPrefs(res)
}

func (cs *ConfigProvider) GetConfig(r ConfigResolver) (string, error) {
	if principal == nil {
		return r.Anon()
	}
	if conf, err := r.Local(); err == nil {
		return conf, nil
	} else if conf, err = r.Own(); err == nil {
		return conf, nil
	} else if conf, err = r.Default(); err == nil {
		return conf, nil
	} else {
		return "", err
	}
}
