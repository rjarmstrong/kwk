package setup

import (
	"github.com/kwk-super-snippets/cli/src/gokwk"
	"github.com/kwk-super-snippets/cli/src/models"
	"github.com/kwk-super-snippets/cli/src/persist"
	"bitbucket.com/sharingmachine/types"
	"gopkg.in/yaml.v2"
)

type PrefsResolvers struct {
	snippets gokwk.Snippets
	file     persist.IO
	a        types.Alias
}

func NewPrefsResolvers(s gokwk.Snippets, f persist.IO) Resolvers {
	return &PrefsResolvers{
		a:        *models.NewSetupAlias("prefs", "yml", false),
		snippets: s,
		file:     f,
	}
}

func (p *PrefsResolvers) Anon() (string, error) {
	return p.Fallback()
}

func (p *PrefsResolvers) Local() (string, error) {
	return p.file.Read(SNIP_CACHE_PATH, p.a.String(), true, 0)
}

func (p *PrefsResolvers) Own() (string, error) {
	if l, err := p.snippets.Get(p.a); err != nil {
		return "", err
	} else {
		if _, err := p.file.Write(SNIP_CACHE_PATH, p.a.String(), l.Snippets[0].Snip, true); err != nil {
			return "", err
		}
		return l.Snippets[0].Snip, nil
	}
}

func (p *PrefsResolvers) Default() (string, error) {
	if prefs, err := p.Fallback(); err != nil {
		return "", err
	} else {
		if _, err := p.snippets.Create(prefs, p.a, types.RolePreferences); err != nil {
			return "", err
		}
		return prefs, nil
	}
}

func (p *PrefsResolvers) Fallback() (string, error) {
	ph := &models.PreferencesHolder{KwkPrefs: "v1", Preferences: models.DefaultPrefs().PersistedPrefs}
	if b, err := yaml.Marshal(ph); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
