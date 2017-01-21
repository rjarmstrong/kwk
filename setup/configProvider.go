package setup

import (
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"gopkg.in/yaml.v2"
)

type ConfigProvider struct {
	u              account.Manager
	envResolvers   Resolvers
	prefsResolvers Resolvers
	prefs          *Preferences
	env            *yaml.MapSlice

}

func NewConfigProvider(ss snippets.Service, s sys.Manager, u account.Manager) Provider {
	env := NewEnvResolvers(ss, s, u)
	prefs := NewPrefsResolvers(ss, s, u)
	return &ConfigProvider{envResolvers:env, prefsResolvers:prefs, u:u}
}

func (cs *ConfigProvider) Preload(){
	cs.Env()
	cs.Prefs()
}

func (cs *ConfigProvider) Env() *yaml.MapSlice {
	if cs.env != nil {
		return cs.env
	}
	env, err := cs.GetConfig(cs.envResolvers)
	if err != nil {
		panic(err)
	}
	cs.env = &yaml.MapSlice{}
	if err := yaml.Unmarshal([]byte(env), cs.env); err != nil {
		panic(err)
	}
	return cs.env
}

func (cs *ConfigProvider) Prefs() *Preferences {
	if cs.prefs != nil {
		return cs.prefs
	}
	if c, err := cs.GetConfig(cs.prefsResolvers); err != nil {
		panic(err)

	} else {
		cs.prefs = &Preferences{PersistedPrefs:PersistedPrefs{}}
		if err := yaml.Unmarshal([]byte(c), cs.prefs); err != nil {
			panic(err)
		} else {
			return cs.prefs
		}

	}
}

func (cs *ConfigProvider) GetConfig(r Resolvers) (string, error) {
	//if os.Getenv(sys.KWK_TESTMODE) != "" && fullName == "env.yml" {
	//	testEnv := "./cmd/testEnv.yml"
	//	// TODO: use log
	//	//fmt.Println(">> Running with:", testEnv, " <<")
	//	b, err := ioutil.ReadFile(testEnv)
	//	return string(b), nil
	//	if err != nil {
	//		return "", err
	//	}
	//}
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
