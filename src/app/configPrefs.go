package app

import (
	"github.com/kwk-super-snippets/types"
	"gopkg.in/yaml.v2"
)

type ConfigPrefs struct {
	snippets types.SnippetsClient
	file     IO
	a        *types.Alias
}

func NewPrefsResolvers(s types.SnippetsClient, f IO) ConfigResolver {
	return &ConfigPrefs{
		a:        NewSetupAlias("prefs", "yml", false),
		snippets: s,
		file:     f,
	}
}

func (p *ConfigPrefs) Anon() (string, error) {
	return p.Fallback()
}

func (p *ConfigPrefs) Local() (string, error) {
	return p.file.Read(snipCachePath, p.a.String(), true, 0)
}

func (p *ConfigPrefs) Own() (string, error) {
	l, err := p.snippets.Get(Ctx(), &types.GetRequest{Alias: p.a})
	if err != nil {
		return "", err
	}
	_, err = p.file.Write(snipCachePath, p.a.String(), l.Items[0].Content, true)
	if err != nil {
		return "", err
	}
	return l.Items[0].Content, nil
}

func (p *ConfigPrefs) Default() (string, error) {
	prefs, err := p.Fallback()
	if err != nil {
		return "", err
	}
	_, err = p.snippets.Create(Ctx(), &types.CreateRequest{Content: prefs, Alias: p.a,
		Role:                                                   types.Role_Preferences})
	if err != nil {
		return "", err
	}
	return prefs, nil
}

func (p *ConfigPrefs) Fallback() (string, error) {
	ph := &PreferencesHolder{KwkPrefs: "v1", Preferences: DefaultPrefs().PersistedPrefs}
	if b, err := yaml.Marshal(ph); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
