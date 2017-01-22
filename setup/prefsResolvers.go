package setup

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"gopkg.in/yaml.v2"
)

type PrefsResolvers struct {
	snippets snippets.Service
	system         sys.Manager
	account        account.Manager
	hostConfigName string
}

func NewPrefsResolvers(s snippets.Service, sys sys.Manager, a account.Manager) Resolvers {
	return &PrefsResolvers{
		hostConfigName:GetHostConfigFullName("prefs.yml"),
		snippets:s,
		system:sys,
		account:a,
	}
}

func (p *PrefsResolvers) Anon() (string, error) {
	//fmt.Println("GETTING ANON")
	return p.Default()
}

func (p *PrefsResolvers) Local() (string, error) {
	//fmt.Println("GETTING LOCAL")
	return p.system.ReadFromFile(SNIP_CACHE_PATH, p.hostConfigName, true, 0)
}

func (p *PrefsResolvers) Own() (string, error) {
	//fmt.Println("GETTING OWN")
	if u, err := p.account.Get(); err != nil {
		return "", err
	} else {
		if l, err := p.snippets.Get(&models.Alias{FullKey: p.hostConfigName, Username: u.Username }); err != nil {
			return "", err
		} else {
			if _, err := p.system.WriteToFile(SNIP_CACHE_PATH, p.hostConfigName, l.Items[0].Snip, true); err != nil {
				return "", err
			}
			return l.Items[0].Snip, nil
		}
	}
}

func (p *PrefsResolvers) Default() (string, error) {
	//fmt.Println("GETTING DEFAULT")
	if prefs, err := p.Fallback(); err != nil {
		return "", err
	} else {
		if p.account.HasValidCredentials() {
			if _, err := p.snippets.Create(prefs, p.hostConfigName, models.RolePreferences); err != nil {
				return "", err
			}
		}
		return prefs, nil
	}
}

func (p *PrefsResolvers) Fallback() (string, error) {
	ph := &PreferencesHolder{KwkPrefs:"v1", Preferences:DefaultPrefs().PersistedPrefs }
	if b, err := yaml.Marshal(ph); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}