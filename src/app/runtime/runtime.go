package runtime

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"gopkg.in/yaml.v2"
	"github.com/kwk-super-snippets/cli/src/store"
)

type Resolver interface {
	Local() (string, error)
	Own() (string, error)
	Default() (string, error)
	Fallback() (string, error)
}

var (
	envAlias   = newRuntimeAlias("env", "yml", true)
	prefsAlias = newRuntimeAlias("prefs", "yml", false)
	EnvURI = envAlias.URI()
)

type SnippetGetter func(req *types.GetRequest) (*types.ListResponse, error)
type SnippetMaker func(req *types.CreateRequest) error
type DocGetter func() (string, error)

type Runtime struct {
	errs.Handler
	sg          SnippetGetter
	sm          SnippetMaker
	file        store.File
	snippetPath string
	loggedIn    bool
}

func Configure(env *yaml.MapSlice, prefs *out.Prefs, loggedIn bool, sg SnippetGetter, sm SnippetMaker, path string, f store.File, eh errs.Handler) {
	c := &Runtime{loggedIn: loggedIn, sg: sg, sm: sm, snippetPath: path, file: f, Handler: eh}
	env = c.getEnv()
	prefs = c.getPrefs()
}

func (cs *Runtime) getEnv() *yaml.MapSlice {
	content, err := cs.resolveDoc(envAlias, func() (string, error) { return defaultEnv, nil }, types.Role_Environment)
	if err != nil {
		out.Debug("Failed to load env settings.")
		out.LogErr(err)
		content = defaultEnv
	}
	env := &yaml.MapSlice{}
	err = yaml.Unmarshal([]byte(content), env)
	if err != nil {
		cs.Handle(errs.New(errs.CodeInvalidConfigSection,
			fmt.Sprintf("Invalid kwk *env.yml detected. `kwk edit env` to fix. %s", err)))
		return DefaultEnv()
	}
	return env
}

func (cs *Runtime) getPrefs() *out.Prefs {
	fallback := func() (string, error) {
		ph := &PrefsFile{KwkPrefs: "v1", Options: *DefaultPrefs()}
		b, err := yaml.Marshal(ph)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	c, err := cs.resolveDoc(prefsAlias, fallback, types.Role_Preferences)
	if err != nil {
		cs.Handle(err)
		return nil
	}
	prefs := &out.Prefs{}
	parse := func(p string) (*out.Prefs, error) {
		ph := &PrefsFile{}
		err := yaml.Unmarshal([]byte(p), ph)
		if err != nil {
			return nil, err
		}
		prefs = &ph.Options
		return prefs, nil
	}
	res, err := parse(c)
	if err != nil {
		cs.Handle(errs.New(errs.CodeInvalidConfigSection,
			fmt.Sprintf("Invalid kwk *prefs.yml detected. `kwk edit prefs` to fix. %s", err)))
		return DefaultPrefs()
	}
	out.Debug("SETTING PREFS: %+v", *res)
	return res
}

func (cs *Runtime) resolveDoc(a *types.Alias, fallback DocGetter, role types.Role) (string, error) {
	if !cs.loggedIn {
		out.Debug("RUNTIME: Getting Fallback %s", a.URI())
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
	conf, err = cs.getDefaultDoc(a, fallback, role)
	if err == nil {
		return conf, nil
	}
	return "", err
}

func (cs *Runtime) getLocalDoc(a *types.Alias) (string, error) {
	out.Debug("RUNTIME: Getting Local %s", a.URI())
	return cs.file.Read(cs.snippetPath, a.URI(), true, 0)
}

func (cs *Runtime) getRemoteDoc(a *types.Alias) (string, error) {
	out.Debug("RUNTIME: Getting Remote %s", a.URI())
	res, err := cs.sg(&types.GetRequest{Alias: a})
	if err != nil {
		return "", err
	}
	content := res.Items[0].Content
	_, err = cs.file.Write(cs.snippetPath, a.URI(), content, true)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (cs *Runtime) getDefaultDoc(a *types.Alias, fallback DocGetter, role types.Role) (string, error) {
	out.Debug("RUNTIME: Getting Default %s", a.URI())
	content, err := fallback()
	if err != nil {
		return "", err
	}
	err = cs.sm(&types.CreateRequest{Content: content, Alias: a, Role: role})
	if err != nil {
		return "", err
	}
	_, err = cs.file.Write(cs.snippetPath, a.URI(), content, true)
	if err != nil {
		return "", err
	}
	return content, nil
}
