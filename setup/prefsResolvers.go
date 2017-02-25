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
	system   sys.Manager
	account  account.Manager
	a        models.Alias
}

func NewPrefsResolvers(s snippets.Service, sys sys.Manager, a account.Manager) Resolvers {
	return &PrefsResolvers{
		a:        *models.NewSetupAlias("prefs", "yml", false),
		snippets: s,
		system:   sys,
		account:  a,
	}
}

func (p *PrefsResolvers) Anon() (string, error) {
	return p.Fallback()
}

func (p *PrefsResolvers) Local() (string, error) {
	return p.system.ReadFromFile(SNIP_CACHE_PATH, p.a.String(), true, 0)
}

func (p *PrefsResolvers) Own() (string, error) {
	if l, err := p.snippets.Get(p.a); err != nil {
		return "", err
	} else {
		if _, err := p.system.WriteToFile(SNIP_CACHE_PATH, p.a.String(), l.Snippets[0].Snip, true); err != nil {
			return "", err
		}
		return l.Snippets[0].Snip, nil
	}
}

func (p *PrefsResolvers) Default() (string, error) {
	if prefs, err := p.Fallback(); err != nil {
		return "", err
	} else {
		if _, err := p.snippets.Create(prefs, p.a, models.SnipRolePreferences); err != nil {
			return "", err
		}
		return prefs, nil
	}
}

func (p *PrefsResolvers) Fallback() (string, error) {
	ph := &models.PreferencesHolder{KwkPrefs: "v1", Preferences: models.DefaultPrefs().PersistedPrefs }
	if b, err := yaml.Marshal(ph); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
