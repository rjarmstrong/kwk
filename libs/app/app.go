package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"github.com/kwk-links/kwk-cli/libs/services/openers"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"github.com/kwk-links/kwk-cli/libs/services/aliases"
	"github.com/kwk-links/kwk-cli/libs/services/users"
	"github.com/kwk-links/kwk-cli/libs/controllers"
)

type KwkApp struct {
	App     *cli.App
}

func NewKwkApp(a aliases.IAliases, s system.ISystem, settings settings.ISettings, o openers.IOpen, u users.IUsers) *KwkApp {

	app := cli.NewApp()

	//cli.HelpPrinter = system.Help

	accCtrl := controllers.NewAccountController(u, settings)
	app.Commands = append(app.Commands, Accounts(accCtrl)...)

	sysCtrl := controllers.NewSystemController(s, u)
	app.Commands = append(app.Commands, System(sysCtrl)...)

	aliasCtrl := controllers.NewAliasController(a, o)
	app.Commands = append(app.Commands, Alias(aliasCtrl)...)
	app.CommandNotFound = func(c *cli.Context, fullKey string) {
		aliasCtrl.Get(fullKey)
	}
	return &KwkApp{App:app}
}
