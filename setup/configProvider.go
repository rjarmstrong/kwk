package setup

import (
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"gopkg.in/yaml.v2"
	_ "io/ioutil"
	_ "bitbucket.com/sharingmachine/kwkcli/log"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/models"
)

type ConfigProvider struct {
	u              account.Manager
	envResolvers   Resolvers
	prefsResolvers Resolvers
	w tmpl.Writer
}

var prefs  *Preferences
var env   *yaml.MapSlice

func NewConfigProvider(ss snippets.Service, s sys.Manager, u account.Manager, w tmpl.Writer) Provider {
	env := NewEnvResolvers(ss, s, u)
	prefs := NewPrefsResolvers(ss, s, u)
	return &ConfigProvider{envResolvers:env, prefsResolvers:prefs, u:u, w:w}
}

func Prefs() *Preferences {
	return prefs
}

func Env() *yaml.MapSlice {
	return env
}

func (cs *ConfigProvider) Load(){
	_, err := cs.localEnv()
	if err != nil {
		cs.w.HandleErr(err)
	}
	cs.loadPrefs()
}

func (cs *ConfigProvider) localEnv() (*yaml.MapSlice, error) {
	if env != nil {
		return env, nil
	}
	envString, err := cs.GetConfig(cs.envResolvers)
	if err != nil {
		return nil, err
	}
	env = &yaml.MapSlice{}
	if err := yaml.Unmarshal([]byte(envString), env); err != nil {
		return nil, err
	}
	return env, nil
}

func (cs *ConfigProvider) loadPrefs() {
	if prefs != nil {
		return
	}
	if c, err := cs.GetConfig(cs.prefsResolvers); err != nil {
		cs.w.HandleErr(err)
		return
	} else {
		prefs = &Preferences{PersistedPrefs:PersistedPrefs{}}
		parse := func(p string) (*Preferences, error) {
			ph := &PreferencesHolder{}
			if err := yaml.Unmarshal([]byte(p), ph); err != nil {
				return nil, err
			} else {
				prefs.PersistedPrefs = ph.Preferences
				return prefs, nil
			}
		}
		if res, err := parse(c); err != nil {
			// TODO: USE TEMPLATE WRITER
			cs.w.HandleErr(models.ErrOneLine(models.Code_InvalidConfigSection, "Invalid kwk *prefs.yml detected. `kwk edit prefs` to fix."))
			// The fallback is expected not to fail parsing
			fb, _ := cs.prefsResolvers.Fallback()
			res, _ = parse(fb)
		} else {
			prefs = res
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


