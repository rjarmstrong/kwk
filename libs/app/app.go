package app

import (
	"bitbucket.com/sharingmachine/kwkcli/libs/services/aliases"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/openers"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/search"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/settings"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/system"
	"bitbucket.com/sharingmachine/kwkcli/libs/services/users"
	"gopkg.in/urfave/cli.v1"
	"os"
	"bitbucket.com/sharingmachine/kwkcli/libs/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/libs/ui/tmpl"
)

type KwkApp struct {
	App *cli.App

	Aliases        aliases.IAliases
	System         system.ISystem
	Settings       settings.ISettings
	Users          users.IUsers
	Openers        openers.IOpen
	Dialogues      dlg.Dialogue
	TemplateWriter tmpl.Writer
}

func NewKwkApp(a aliases.IAliases, s system.ISystem, t settings.ISettings, o openers.IOpen, u users.IUsers,
	d dlg.Dialogue, w tmpl.Writer, h search.ISearch) *KwkApp {

	app := cli.NewApp()
	os.Setenv(system.APP_VERSION, "0.0.1")
	//cli.HelpPrinter = system.Help

	accCtrl := NewAccountController(u, t, w, d)
	app.Commands = append(app.Commands, Accounts(accCtrl)...)

	sysCtrl := NewSystemController(s, u, w)
	app.Commands = append(app.Commands, System(sysCtrl)...)

	aliasCtrl := NewAliasController(a, o, s, d, w, t, h)
	app.Commands = append(app.Commands, Alias(aliasCtrl)...)
	app.CommandNotFound = func(c *cli.Context, fullKey string) {
		aliasCtrl.Open(fullKey, []string(c.Args())[1:])
	}
	searchCtrl := NewSearchController(h, w, d)
	app.Commands = append(app.Commands, Search(searchCtrl)...)

	return &KwkApp{App: app, System: s, Settings: t, Openers: o, Users: u, Dialogues: d, Aliases: a, TemplateWriter: w}
}

func (a *KwkApp) Run(args ...string) {
	params := []string{"[app]"}
	params = append(params, args...)
	a.App.Run(params)
}
