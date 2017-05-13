package app

import (
	"github.com/kwk-super-snippets/types"
	"gopkg.in/yaml.v2"
)

type PrefsResolvers struct {
	snippets types.SnippetsClient
	file     IO
	a        *types.Alias
}

func NewPrefsResolvers(s types.SnippetsClient, f IO) Resolvers {
	return &PrefsResolvers{
		a:        NewSetupAlias("prefs", "yml", false),
		snippets: s,
		file:     f,
	}
}

func (p *PrefsResolvers) Anon() (string, error) {
	return p.Fallback()
}

func (p *PrefsResolvers) Local() (string, error) {
	return p.file.Read(snipCachePath, p.a.String(), true, 0)
}

func (p *PrefsResolvers) Own() (string, error) {
	l, err := p.snippets.Get(GetCtx(), &types.GetRequest{Alias: p.a})
	if err != nil {
		return "", err
	}
	_, err = p.file.Write(snipCachePath, p.a.String(), l.Items[0].Content, true)
	if err != nil {
		return "", err
	}
	return l.Items[0].Content, nil
}

func (p *PrefsResolvers) Default() (string, error) {
	prefs, err := p.Fallback()
	if err != nil {
		return "", err
	}
	_, err = p.snippets.Create(GetCtx(), &types.CreateRequest{Content: prefs, Alias: p.a,
		Role: types.Role_Preferences})
	if err != nil {
		return "", err
	}
	return prefs, nil
}

func (p *PrefsResolvers) Fallback() (string, error) {
	ph := &PreferencesHolder{KwkPrefs: "v1", Preferences: DefaultPrefs().PersistedPrefs}
	if b, err := yaml.Marshal(ph); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}
