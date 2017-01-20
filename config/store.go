package config

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"google.golang.org/grpc/codes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"runtime"
	"fmt"
	"os"
	"bitbucket.com/sharingmachine/kwkcli/models"
)

const (
	PREFS_SUFFIX = "prefs.yml"
	ENV_SUFFIX   = "env.yml"
	SNIP_CACHE_PATH = "snip-cache"
)

type Store interface {
	Env() *yaml.MapSlice
	Prefs() *Preferences
	GetSubSection(yml *yaml.MapSlice, name string) (yaml.MapSlice, []string)
}

type StdConfigStore struct {
	snippets snippets.Service
	system sys.Manager
	account account.Manager
	clone models.SnippetCloner
}

func GetHostConfigFullName(fullName string) string {
	if h, err := os.Hostname(); err != nil {
		panic(err)
	} else {
		return fmt.Sprintf(".%s_%s", h, fullName)
	}
}

func (cs *StdConfigStore) Env() *yaml.MapSlice {
	getDefault := func() (string, error) {
		defaultEnv := fmt.Sprintf("%s-%s.yml", runtime.GOOS, runtime.GOARCH)
		defaultAlias := &Alias{FullKey:defaultEnv, Username:"env"}
		if snip, err := cs.snippets.Clone(defaultAlias, GetHostConfigFullName(ENV_SUFFIX)); err != nil {
			return "", err
		} else {
			return snip.Snip, nil
		}
	}

	env, err := cs.GetConfig(ENV_SUFFIX, getDefault)
	if err != nil {
		panic(err)
	}
	c := &yaml.MapSlice{}
	if err := yaml.Unmarshal([]byte(env), c); err != nil {
		panic(err)
	}
	return c
}

func (cs *StdConfigStore) Prefs() *Preferences {
	getDefault := func ()(string, error) {
		dp := DefaultPrefs()
		b, err := yaml.Marshal(dp.PersistedPrefs)
		if err != nil {
			panic(err)
		}
		if cs.account.HasValidCredentials() {
			if _, err := cs.snippets.Create(string(b), GetHostConfigFullName(PREFS_SUFFIX), models.RolePreferences); err != nil {
				return "", err
			}
		}
		return string(b), nil
	}
	if c, err := cs.GetConfig(PREFS_SUFFIX, getDefault); err != nil {
		panic(err)
	} else {
		pp := &PersistedPrefs{}
		if err := yaml.Unmarshal([]byte(c), pp); err != nil {
			panic(err)
		} else {
			return pp
		}

	}
}


func (cs *StdConfigStore) GetConfig(fullName string, ownConf models.GetConf, localConf models.GetConf, anonConf models.GetConf, defaultConf models.GetConf) (string, error) {
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
	hostConfigName := GetHostConfigFullName(fullName)
	if !cs.account.HasValidCredentials() {
		return anonConf()
	}
	if conf, err := localConf(); err == nil {
		return conf, nil
	} else if conf, err = ownConf(); err == nil {
		return conf, nil
	} else if conf, err = defaultConf(); err == nil {
		return conf, nil
	} else {
		return "", err
	}

	// TODO: check yml version is compatible with this build else force upgrade.

	var u *User
	if u, err := cs.account.Get(); err != nil {

	}
	if ok, _ := cs.system.FileExists(SNIP_CACHE_PATH, hostConfigName, true); !ok {
		if l, err := cs.snippets.Get(&Alias{FullKey: hostConfigName, Username: u.Username }); err != nil {
			if err.(*models.ClientErr).TransportCode == codes.NotFound {
				if conf, err := getDefault(); err != nil {
					return "", err
				} else {
					_, err := cs.system.WriteToFile(SNIP_CACHE_PATH, hostConfigName, conf, true)
					return conf, err
				}
			} else {
				return "", err
			}
		} else {
			//r.system.WriteToFile(FILE_CACHE_PATH, hostConfigName, l.Items[0].Snip, true)
			return l.Items[0].Snip, nil
		}
	} else {
		if e, err := cs.system.ReadFromFile(SNIP_CACHE_PATH, hostConfigName, true, 0); err != nil {
			return "", err
		} else {
			return e, nil
		}
	}
}

func (*StdConfigStore) GetSubSection(yml *yaml.MapSlice, name string) (yaml.MapSlice, []string) {
	f := func(yml *yaml.MapSlice, name string) (yaml.MapSlice, []string) {
		for _, v := range *yml {
			if v.Key == name {
				if slice, ok := v.Value.(yaml.MapSlice); ok {
					return slice, nil
				}
				if _, ok := v.Value.([]interface{}); ok {
					items := []string{}
					for _, v2 := range v.Value.([]interface{}) {
						items = append(items, v2.(string))
					}
					return nil, items
				}
				return nil, []string{v.Value.(string)}
			}
		}
		return nil, nil
	}
	if sub, bottom := f(yml, name); sub == nil && bottom == nil {
		return f(yml, "default")
	} else {
		return sub, bottom
	}
}

