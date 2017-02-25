package setup

import (
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"gopkg.in/yaml.v2"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/log"
)

type ConfigProvider struct {
	u              account.Manager
	envResolvers   Resolvers
	prefsResolvers Resolvers
	w tmpl.Writer
}

func NewConfigProvider(ss snippets.Service, s sys.Manager, u account.Manager, w tmpl.Writer) Provider {
	env := NewEnvResolvers(ss, s, u)
	prefs := NewPrefsResolvers(ss, s, u)
	return &ConfigProvider{envResolvers:env, prefsResolvers:prefs, u:u, w:w}
}

func (cs *ConfigProvider) Load(){
	_, err := cs.loadEnv()
	if err != nil {
		cs.w.HandleErr(err)
	}
	cs.loadPrefs()
}

func (cs *ConfigProvider) loadEnv() (*yaml.MapSlice, error) {
	if models.Env() != nil {
		return models.Env(), nil
	}
	envString, err := cs.GetConfig(cs.envResolvers)
	if err != nil {
		return nil, err
	}
	env := &yaml.MapSlice{}
	if err := yaml.Unmarshal([]byte(envString), env); err != nil {
		return nil, err
	}
	models.SetEnv(env)
	return models.Env(), nil
}

func (cs *ConfigProvider) loadPrefs() {
	if models.Prefs() != nil {
		return
	}
	if c, err := cs.GetConfig(cs.prefsResolvers); err != nil {
		cs.w.HandleErr(err)
		return
	} else {
		prefs := &models.Preferences{PersistedPrefs:models.PersistedPrefs{}}
		parse := func(p string) (*models.Preferences, error) {
			ph := &models.PreferencesHolder{}
			if err := yaml.Unmarshal([]byte(p), ph); err != nil {
				return nil, err
			} else {
				prefs.PersistedPrefs = ph.Preferences
				return prefs, nil
			}
		}
		log.Debug("Loaded prefs:%+v", c)
		if res, err := parse(c); err != nil {
			// TODO: USE TEMPLATE WRITER
			cs.w.HandleErr(models.ErrOneLine(models.Code_InvalidConfigSection, "Invalid kwk *prefs.yml detected. `kwk edit prefs` to fix."))
			// The fallback is expected not to fail parsing
			fb, _ := cs.prefsResolvers.Fallback()
			res, _ = parse(fb)
		} else {
			log.Debug("SETTING PREFS: %+v", res)
			models.SetPrefs(res)
		}
	}
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


