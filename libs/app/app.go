package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"github.com/kwk-links/kwk-cli/libs/services/openers"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"github.com/kwk-links/kwk-cli/libs/services/aliases"
	"github.com/kwk-links/kwk-cli/libs/services/users"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
)

type KwkApp struct {
	App     *cli.App

	Aliases aliases.IAliases
	System system.ISystem
	Settings settings.ISettings
	Users users.IUsers
	Openers openers.IOpen
	Dialogues gui.IDialogues
	TemplateWriter gui.ITemplateWriter
}

func NewKwkApp(a aliases.IAliases, s system.ISystem, t settings.ISettings, o openers.IOpen, u users.IUsers, d gui.IDialogues, w gui.ITemplateWriter) *KwkApp {

	app := cli.NewApp()

	//cli.HelpPrinter = system.Help

	accCtrl := NewAccountController(u, t, w, d)
	app.Commands = append(app.Commands, Accounts(accCtrl)...)

	sysCtrl := NewSystemController(s, u, w)
	app.Commands = append(app.Commands, System(sysCtrl)...)

	aliasCtrl := NewAliasController(a, o, s, d, w)
	app.Commands = append(app.Commands, Alias(aliasCtrl)...)
	app.CommandNotFound = func(c *cli.Context, fullKey string) {
		aliasCtrl.Open(fullKey, []string(c.Args())[1:])
	}
	return &KwkApp{App:app, System:s, Settings:t, Openers:o, Users:u, Dialogues:d, Aliases:a, TemplateWriter:w}
}
