package setup

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"fmt"
	"runtime"
	"strings"
)

// TODO: check yml version is compatible with this build else force upgrade.

type EnvResolvers struct {
	snippets snippets.Service
	system   sys.Manager
	account  account.Manager
	alias    models.Alias
	runtime string
}

func NewEnvResolvers(s snippets.Service, sys sys.Manager, a account.Manager) Resolvers {
	r := strings.ToLower(fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH))
	return &EnvResolvers{
		runtime: 	r,
		alias:    *models.NewSetupAlias("env", "yml"),
		snippets: s,
		system:   sys,
		account:  a,
	}
}

func (e *EnvResolvers) Anon() (string, error) {
	return e.Fallback()
}

func (e *EnvResolvers) Local() (string, error) {
	return e.system.ReadFromFile(SNIP_CACHE_PATH, e.alias.String(), true, 0)
}

func (e *EnvResolvers) Own() (string, error) {
	if l, err := e.snippets.Get(e.alias); err != nil {
		return "", err
	} else {
		if _, err := e.system.WriteToFile(SNIP_CACHE_PATH, e.alias.String(), l.Items[0].Snip, true); err != nil {
			return "", err
		}
		return l.Items[0].Snip, nil
	}
}

func (e *EnvResolvers) Default() (string, error) {
	if env, err := e.Fallback(); err != nil {
		return "", err
	} else {
		if snip, err := e.snippets.Create(env, e.alias, models.RoleEnvironment); err != nil {
			return "", err
		} else {
			return snip.Snippet.Snip, nil
		}
		if _, err := e.system.WriteToFile(SNIP_CACHE_PATH, e.alias.String(), env, true); err != nil {
			return "", err
		} else {
			return env, nil
		}
	}
}

func (e *EnvResolvers) Fallback() (string, error) {
	fb := fallbackMap[e.runtime]
	if fb == "" {
		return "", models.ErrOneLine(models.Code_EnvironmentNotSupported, fmt.Sprintf("No default environment configuration for your system: %s", e.runtime))
	}
	return fb, nil
}

var fallbackMap = map[string]string{
	"darwin-amd64": darwinAmd64,
	"linux-amd64":  linuxAmd64,
}

const darwinAmd64 = `kwkenv: "1"
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
    compile: ["dotnet", "restore", "/Volumes/development/go/src/bitbucket.com/sharingmachine/kwkcli/dotnet/project.json"]
    run: ["dotnet", "run", "--project", "/Volumes/development/go/src/bitbucket.com/sharingmachine/kwkcli/dotnet/project.json", "$FULL_NAME",]
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
apps:
  webstorm: ["open", "-a", "webstorm", "$DIR"]
  textedit: ["open", "-e", "$FULL_NAME"]
  vscode: ["open", "-a", "Visual Studio Code", "$DIR"]
  vim: ["vi", "$FULL_NAME" ]
  emacs: ["emacs", "$FULL_NAME" ]
  nano: ["nano", "$FULL_NAME" ]
  default: ["open", "-t", "$FULL_NAME"]
editors:
#  Specify one app for each file type to edit.
#  sh: [vim]
#  go: [gogland]
#  py: [vscode]
#  url: [textedit]
  default: ["vim"]
security: #https://gist.github.com/pmarreck/5388643
  encrypt: []
  decrypt: []
  sign: []
  verify: []`

const linuxAmd64 = `kwkenv: "1"
runners:
  sh: ["/bin/bash", "-c", "$SNIP"]
  url: ["open", "$SNIP"]
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
apps:
  vim: ["vi", "$FULL_NAME"]
  emacs: ["emacs", "$FULL_NAME" ]
  nano: ["nano", "$FULL_NAME" ]
  default: ["vi", "$FULL_NAME"]
editors:
#  Specify one app for each file type to edit.
#  sh: [vim]
#  go: [emacs]
#  py: [nano]
#  url: [vim]
  default: ["vim"]
security: #https://gist.github.com/pmarreck/5388643
  encrypt: []
  decrypt: []
  sign: []
  verify: []`
