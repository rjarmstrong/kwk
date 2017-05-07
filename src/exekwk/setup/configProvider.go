package setup

import (
	"github.com/kwk-super-snippets/cli/src/gokwk"
	"github.com/kwk-super-snippets/cli/src/models"
	"github.com/kwk-super-snippets/cli/src/persist"
	"bitbucket.com/sharingmachine/types/errs"
	"fmt"
	"gopkg.in/yaml.v2"
)

type ConfigProvider struct {
	u              gokwk.Users
	envResolvers   Resolvers
	prefsResolvers Resolvers
	errs.Handler
}

func NewConfig(ss gokwk.Snippets, f persist.IO, u gokwk.Users, eh errs.Handler) Provider {
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
	if models.Env() != nil {
		return models.Env()
	}
	envString, err := cs.GetConfig(cs.envResolvers)
	if err != nil {
		models.Debug("Failed to load env settings.")
		models.LogErr(err)
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
	models.SetEnv(env)
	return models.Env()
}

func (cs *ConfigProvider) loadPrefs() {
	if models.Prefs() != nil {
		return
	}
	c, err := cs.GetConfig(cs.prefsResolvers)
	if err != nil {
		cs.Handle(err)
		return
	}
	prefs := &models.Preferences{PersistedPrefs: models.PersistedPrefs{}}
	parse := func(p string) (*models.Preferences, error) {
		ph := &models.PreferencesHolder{}
		err := yaml.Unmarshal([]byte(p), ph)
		if err != nil {
			return nil, err
		}
		prefs.PersistedPrefs = ph.Preferences
		return prefs, nil
	}
	models.Debug("Loaded prefs:%+v", c)
	res, err := parse(c)
	if err != nil {
		cs.Handle(errs.New(errs.CodeInvalidConfigSection,
			fmt.Sprintf("Invalid kwk *prefs.yml detected. `kwk edit prefs` to fix. %s", err)))
		// The fallback is expected not to fail parsing
		fb, _ := cs.prefsResolvers.Fallback()
		res, _ = parse(fb)
	}
	models.Debug("SETTING PREFS: %+v", res)
	models.SetPrefs(res)
}

func (cs *ConfigProvider) GetConfig(r Resolvers) (string, error) {
	if !cs.u.HasValidCredentials() {
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
