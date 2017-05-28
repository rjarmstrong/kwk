package app

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"runtime"
	"strings"
	"github.com/kwk-super-snippets/cli/src/store"
)

// TODO: check yml version is compatible with this build else force upgrade.

type ConfigEnv struct {
	snippets types.SnippetsClient
	file     store.File
	alias    *types.Alias
	runtime  string
}

func NewEnvResolvers(s types.SnippetsClient, f store.File) ConfigResolver {
	r := strings.ToLower(fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
	return &ConfigEnv{
		runtime:  r,
		alias:    NewSetupAlias("env", "yml", true),
		snippets: s,
		file:     f,
	}
}

func (e *ConfigEnv) Anon() (string, error) {
	return e.Fallback()
}

func (e *ConfigEnv) Local() (string, error) {
	return e.file.Read(cfg.SnippetPath, e.alias.URI(), true, 0)
}

func (e *ConfigEnv) Own() (string, error) {
	out.Debug("GETTING ENV: %s", e.alias.URI())
	res, err := e.snippets.Get(Ctx(), &types.GetRequest{Alias: e.alias})
	if err != nil {
		return "", err
	}
	content := res.Items[0].Content
	_, err = e.file.Write(cfg.SnippetPath, e.alias.URI(), content, true)
	if err != nil {
		return "", err
	}
	return content, nil
}

func (e *ConfigEnv) Default() (string, error) {
	env, err := e.Fallback()
	if err != nil {
		return "", err
	}
	_, err = e.snippets.Create(Ctx(), &types.CreateRequest{Content: env, Alias: e.alias, Role: types.Role_Environment})
	if err != nil {
		return "", err
	}
	_, err = e.file.Write(cfg.SnippetPath, e.alias.String(), env, true)
	if err != nil {
		return "", err
	}
	return env, nil
}

func (e *ConfigEnv) Fallback() (string, error) {
	fb := fallbackMap[e.runtime]
	if fb == "" {
		return "", errs.New(errs.CodeEnvironmentNotSupported, fmt.Sprintf("No default environment configuration for your system: %s", e.runtime))
	}
	return fb, nil
}

var fallbackMap = map[string]string{
	"darwin-amd64": darwinAmd64,
	"linux-amd64":  linuxAmd64,
}

const darwinAmd64 = `kwkenv: "1"
editors:
#  Specify one app for each file type to edit.
#  sh: [vim]
#  go: [gogland]
#  py: [vscode]
#  url: [textedit]
  default: ["textedit"]
apps:
  webstorm: ["open", "-a", "webstorm", "$DIR"]
  textedit: ["open", "-e", "$FULL_NAME"]
  vscode: ["open", "-a", "Visual Studio Code", "$DIR"]
  vim: ["vi", "$FULL_NAME" ]
  emacs: ["emacs", "$FULL_NAME" ]
  nano: ["nano", "$FULL_NAME" ]
  default: ["open", "-t", "$FULL_NAME"]
runners:
  sh: ["/bin/bash", "-c", "$SNIP"]
  url: ["open", "$SNIP"]
  url-covert: ["/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--incognito", "$SNIP"]
  js: ["node", "-e", "$SNIP"] #nodejs
  py: ["python", "-c", "$SNIP"] #python
  php: ["php", "-r", "$SNIP"] #php
  scpt: ["osascript", "-e", "$SNIP"] #applescript
  applescript: ["osascript", "-e", "$SNIP"] #applescript
  rb: ["ruby", "-e", "$SNIP"] #ruby
  pl: ["perl", "-E", "$SNIP" ] #perl
  exs: ["elixir", "-e", "$SNIP"] # elixir
  java:
    compile: ["javac", "$FULL_NAME"]
    run: ["java", "$CLASS_NAME"]
  scala:
    compile: ["scalac", "-d", "$DIR", "$FULL_NAME"]
    run: ["scala", "$NAME"]
  cs: #c sharp (dotnet core) Under development
    compile: ["dotnet", "restore", "/Volumes/development/go/src/github.com/kwk-super-snippets/cli/src/dotnet/project.json"]
    run: ["dotnet", "run", "--project", "/Volumes/development/go/src/github.com/kwk-super-snippets/cli/src/dotnet/project.json", "$FULL_NAME",]
  go: #golang
    run: ["go", "run", "$FULL_NAME"]
  rs: #rust
    compile: ["rustc", "-o", "$NAME", "$FULL_NAME"]
    run: ["$NAME"]
  cpp: # c++
    compile: ["g++", "$FULL_NAME", "-o", "$FULL_NAME.out" ]
    run: ["$FULL_NAME.out"]
  path: ["echo", "$SNIP" ]
  xml: ["echo", "$SNIP"]
  json: ["echo", "$SNIP"]
  yml: ["echo", "$SNIP"]
  md:
    run: ["mdless", "$FULL_NAME"]
  default: ["echo", "$SNIP"]
security: #https://gist.github.com/pmarreck/5388643
  encrypt: []
  decrypt: []
  sign: []
  verify: []`

const linuxAmd64 = `kwkenv: "1"
editors:
#  Specify one app for each file type to edit.
#  sh: [vim]
#  go: [emacs]
#  py: [nano]
#  url: [vim]
  default: ["vim"]
apps:
  vim: ["vi", "$FULL_NAME"]
  emacs: ["emacs", "$FULL_NAME" ]
  nano: ["nano", "$FULL_NAME" ]
  default: ["vi", "$FULL_NAME"]
runners:
  jl: ["julia", "-e", "$SNIP"]
  sh: ["/bin/bash", "-c", "$SNIP"]
  url: ["firefox", "--new-tab", "$SNIP"]
  url-covert: ["firefox", "--private-window", "$SNIP"]
  js: ["node", "-e", "$SNIP"] #nodejs
  py: ["python", "-c", "$SNIP"] #python
  php: ["php", "-r", "$SNIP"] #php
  scpt: ["osascript", "-e", "$SNIP"] #applescript
  applescript: ["osascript", "-e", "$SNIP"] #applescript
  rb: ["ruby", "-e", "$SNIP"] #ruby
  pl: ["perl", "-E", "$SNIP" ] #perl
  exs: ["elixir", "-e", "$SNIP"] # elixir
  java:
    compile: ["javac", "$FULL_NAME"]
    run: ["java", "$CLASS_NAME"]
  scala:
    compile: ["scalac", "-d", "$DIR", "$FULL_NAME"]
    run: ["scala", "$NAME"]
  go: #golang
    run: ["go", "run", "$FULL_NAME"]
  rs: #rust
    compile: ["rustc", "-o", "$NAME", "$FULL_NAME"]
    run: ["$NAME"]
  cpp: # c++
    compile: ["g++", "$FULL_NAME", "-o", "$FULL_NAME.out" ]
    run: ["$FULL_NAME.out"]
  path: ["echo", "$SNIP" ]
  xml: ["echo", "$SNIP"]
  json: ["echo", "$SNIP"]
  yml: ["echo", "$SNIP"]
  default: ["echo", "$SNIP"]
security: #https://gist.github.com/pmarreck/5388643
  encrypt: []
  decrypt: []
  sign: []
  verify: []`
