package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/services/system"
	"github.com/kwk-links/kwk-cli/libs/services/gui"
	"github.com/kwk-links/kwk-cli/libs/services/openers"
	"github.com/kwk-links/kwk-cli/libs/services/settings"
	"github.com/kwk-links/kwk-cli/libs/services/aliases"
)

type KwkApp struct {
	App     *cli.App
}

func NewKwkApp(a aliases.IAliases, s system.ISystem, settings settings.ISettings, i gui.IInteraction, o openers.IOpen) *KwkApp {
	app := cli.NewApp()

	//cli.HelpPrinter = system.Help

	app.Commands = append(app.Commands, accountCommands(a)...)
	app.Commands = append(app.Commands, systemCommands(s, settings, i)...)
	app.Commands = append(app.Commands, aliasCommands(a, s, i, o)...)

	m := NewMultiResultPrompt(o, i)
	app.CommandNotFound = func(c *cli.Context, fullKey string) {
		m.CheckAndPrompt(fullKey, a.Get(fullKey), c.Args())
	}
	return &KwkApp{App:app}
}
