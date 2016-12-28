package app

import (
	"bitbucket.com/sharingmachine/kwkcli/openers"
	"bitbucket.com/sharingmachine/kwkcli/search"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/system"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"gopkg.in/urfave/cli.v1"
	"os"
)

type KwkApp struct {
	App *cli.App
	Snippets       snippets.Service
	System         system.ISystem
	Settings       config.Settings
	AccountManage  account.Manager
	Openers        openers.IOpen
	Dialogues      dlg.Dialogue
	TemplateWriter tmpl.Writer
}

func New(a snippets.Service, s system.ISystem, t config.Settings, o openers.IOpen, u account.Manager,
	d dlg.Dialogue, w tmpl.Writer, h search.Term) *KwkApp {

	app := cli.NewApp()
	os.Setenv(system.APP_VERSION, "0.0.1")
	//cli.HelpPrinter = system.Help

	accCli := NewAccountCli(u, t, w, d)
	app.Commands = append(app.Commands, Accounts(accCli)...)

	sysCli := NewSystemCli(s, u, w)
	app.Commands = append(app.Commands, System(sysCli)...)

	snipCli := NewSnippetCli(a, o, s, d, w, t, h)
	app.Commands = append(app.Commands, Snippets(snipCli)...)
	app.CommandNotFound = func(c *cli.Context, fullKey string) {
		snipCli.Open(fullKey, []string(c.Args())[1:])
	}
	searchCli := NewSearchCli(h, w, d)
	app.Commands = append(app.Commands, Search(searchCli)...)

	return &KwkApp{App: app, System: s, Settings: t, Openers: o, AccountManage: u, Dialogues: d, Snippets: a, TemplateWriter: w}
}

func (a *KwkApp) Run(args ...string) {
	params := []string{"[app]"}
	params = append(params, args...)
	a.App.Run(params)
}
