package setup

import (
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/models"
)

type ConfigProvider struct {
	u              account.Manager
	envResolvers   Resolvers
	prefsResolvers Resolvers
	prefs          *Preferences
	env            *yaml.MapSlice
	w tmpl.Writer
}

func NewConfigProvider(ss snippets.Service, s sys.Manager, u account.Manager, w tmpl.Writer) Provider {
	env := NewEnvResolvers(ss, s, u)
	prefs := NewPrefsResolvers(ss, s, u)
	return &ConfigProvider{envResolvers:env, prefsResolvers:prefs, u:u, w:w}
}

func (cs *ConfigProvider) Preload(){
	_, err := cs.Env()
	if err != nil {
		cs.w.HandleErr(err)
	}
	cs.Prefs()
}

func (cs *ConfigProvider) Env() (*yaml.MapSlice, error) {
	if cs.env != nil {
		return cs.env, nil
	}
	env, err := cs.GetConfig(cs.envResolvers)
	if err != nil {
		return nil, err
	}
	cs.env = &yaml.MapSlice{}
	if err := yaml.Unmarshal([]byte(env), cs.env); err != nil {
		return nil, err
	}
	return cs.env, nil
}

func (cs *ConfigProvider) Prefs() *Preferences {
	if cs.prefs != nil {
		return cs.prefs
	}
	if c, err := cs.GetConfig(cs.prefsResolvers); err != nil {
		cs.w.HandleErr(err)
		return nil
	} else {
		cs.prefs = &Preferences{PersistedPrefs:PersistedPrefs{}}
		parse := func(p string) (*Preferences, error) {
			ph := &PreferencesHolder{}
			if err := yaml.Unmarshal([]byte(p), ph); err != nil {
				return nil, err
			} else {
				cs.prefs.PersistedPrefs = ph.Preferences
				return cs.prefs, nil
			}
		}
		if res, err := parse(c); err != nil {
			// TODO: USE TEMPLATE WRITER
			cs.w.HandleErr(models.ErrOneLine(models.Code_InvalidConfigSection, "Invalid kwk *prefs.yml detected. `kwk edit prefs` to fix."))
			// The fallback is expected not to fail parsing
			fb, _ := cs.prefsResolvers.Fallback()
			res, _ = parse(fb)
			return res
		} else {
			return res
		}
	}
}

func (cs *ConfigProvider) GetConfig(r Resolvers) (string, error) {
	_, ok := r.(*EnvResolvers)
	if ok && sys.KWK_TEST_MODE {
		testEnv := "./cmd/testEnv.yml"
		log.Debug("Running with:", testEnv)
		b, err := ioutil.ReadFile(testEnv)
		return string(b), nil
		if err != nil {
			return "", err
		}
	}
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
