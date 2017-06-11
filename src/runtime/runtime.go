package runtime

import (
	"fmt"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/store"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/errs"
	"gopkg.in/yaml.v2"
)

type Runtime struct {
	errs.Handler
	sg   SnippetGetter
	sm   SnippetMaker
	file store.SnippetReadWriter
}

var (
	envAlias   *types.Alias
	prefsAlias *types.Alias
)

func GetEnvURI() string {
	if envAlias == nil {
		return ""
	}
	return envAlias.URI()
}

func Configure(env *yaml.MapSlice, prefs *out.Prefs, username string, sg SnippetGetter, sm SnippetMaker, f store.SnippetReadWriter, eh errs.Handler) {
	c := &Runtime{sg: sg, sm: sm, file: f, Handler: eh}
	envAlias = newRuntimeAlias(username, "env", "yml", true)
	prefsAlias = newRuntimeAlias(username, "prefs", "yml", false)
	c.resolvePrefs(prefs)
	c.resolveEnv(env)
}

func (cs *Runtime) resolveEnv(env *yaml.MapSlice) {
	content, err := cs.resolveDoc(envAlias, func() (string, error) { return defaultEnvString, nil }, types.Role_Environment)
	err = yaml.Unmarshal([]byte(content), env)
	if err != nil {
		cs.Handle(errs.New(errs.CodeInvalidArgument, fmt.Sprintf("Invalid kwk *env.yml detected. `kwk edit env` to fix. %s", err)))
		*env = *DefaultEnv()
	}
	//out.Debug("ENV: %+v", *env)
}

func (cs *Runtime) resolvePrefs(prefs *out.Prefs) {
	fallback := func() (string, error) {
		return getPrefsAsString(*DefaultPrefs()), nil
	}
	content, err := cs.resolveDoc(prefsAlias, fallback, types.Role_Preferences)
	parse := func(p string) error {
		ph := &PrefsFile{}
		err := yaml.Unmarshal([]byte(p), ph)
		if err != nil {
			return err
		}
		*prefs = ph.Options
		return nil
	}
	err = parse(content)
	if err != nil {
		cs.Handle(errs.New(errs.CodeInvalidArgument,
			fmt.Sprintf("Invalid kwk *prefs.yml detected. `kwk edit prefs` to fix. %s", err)))
		*prefs = *DefaultPrefs()
	}
	//out.Debug("PREFS: %+v", *prefs)
}

func (cs *Runtime) resolveDoc(a *types.Alias, fallback DocGetter, role types.Role) (string, error) {
	if a.Username == "" {
		out.Debug("RUNTIME: No username available. Getting Fallback %s", a.URI())
		return fallback()
	}
	conf, err := cs.getLocalDoc(a)
	if err == nil {
		return conf, nil
	}
	conf, err = cs.getRemoteDoc(a)
	if err == nil {
		return conf, nil
	}
	conf, err = cs.createSnippetWithDefault(a, fallback, role)
	if err == nil {
		return conf, nil
	}
	out.Debug("Failed to resolve any %s. Using fallback.", a.URI())
	return fallback()
}

func (cs *Runtime) getLocalDoc(a *types.Alias) (string, error) {
	out.Debug("RUNTIME: Getting Local %s", a.URI())
	return cs.file.Read(a.URI())
}

func (cs *Runtime) getRemoteDoc(a *types.Alias) (string, error) {
	out.Debug("RUNTIME: Getting Remote %s", a.URI())
	res, err := cs.sg(&types.GetRequest{Alias: a})
	if err != nil {
		return "", err
	}
	content := res.Items[0].Content
	_, err = cs.file.Write(a.URI(), content)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (cs *Runtime) createSnippetWithDefault(a *types.Alias, fallback DocGetter, role types.Role) (string, error) {
	out.Debug("RUNTIME: Getting Default %s", a.URI())
	content, err := fallback()
	if err != nil {
		return "", err
	}
	err = cs.sm(&types.CreateRequest{Content: content, Alias: a, Role: role})
	if err != nil {
		return "", err
	}
	_, err = cs.file.Write(a.URI(), content)
	if err != nil {
		return "", err
	}
	return content, nil
}
