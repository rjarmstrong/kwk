package setup

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"gopkg.in/yaml.v2"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"fmt"
	"bitbucket.com/sharingmachine/kwkcli/user"
	"bitbucket.com/sharingmachine/kwkcli/persist"
)

type ConfigProvider struct {
	u              user.Account
	envResolvers   Resolvers
	prefsResolvers Resolvers
	w              tmpl.Writer
}

func NewConfigProvider(ss snippets.Service, f persist.IO, u user.Account, w tmpl.Writer) Provider {
	env := NewEnvResolvers(ss, f)
	prefs := NewPrefsResolvers(ss, f)
	return &ConfigProvider{envResolvers: env, prefsResolvers: prefs, u: u, w: w}
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
		log.Error("Failed to load env settings.", err)
		envString, _ = cs.envResolvers.Fallback()
	}
	env := &yaml.MapSlice{}
	err = yaml.Unmarshal([]byte(envString), env)
	if err != nil {
		cs.w.HandleErr(models.ErrOneLine(models.Code_InvalidConfigSection,
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
	c, err := cs.GetConfig(cs.prefsResolvers);
	if err != nil {
		cs.w.HandleErr(err)
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
	log.Debug("Loaded prefs:%+v", c)
	res, err := parse(c)
	if err != nil {
		// TODO: USE TEMPLATE WRITER
		cs.w.HandleErr(models.ErrOneLine(models.Code_InvalidConfigSection,
			fmt.Sprintf("Invalid kwk *prefs.yml detected. `kwk edit prefs` to fix. %s", err)))
		// The fallback is expected not to fail parsing
		fb, _ := cs.prefsResolvers.Fallback()
		res, _ = parse(fb)
	}
	log.Debug("SETTING PREFS: %+v", res)
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
