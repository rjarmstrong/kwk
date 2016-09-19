package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/api"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/openers"
	"github.com/kwk-links/kwk-cli/libs/settings"
)

type KwkApp struct {
	App *cli.App
	Api api.IApi
}

func NewKwkApp(a api.IApi, s system.ISystem, settings settings.ISettings, i gui.IInteraction, o openers.IOpen) *KwkApp {
	app := cli.NewApp()

	//cli.HelpPrinter = system.Help

	app.Commands = append(app.Commands, accountCommands(a)...)
	app.Commands = append(app.Commands, systemCommands(s, settings, i)...)
	app.Commands = append(app.Commands, aliasCommands(a, s, i, o)...)

	m := NewMultiResultPrompt(o, i)
	app.CommandNotFound = func(c *cli.Context, fullKey string) {
		m.CheckAndPrompt(fullKey, a.Get(fullKey), c.Args())
	}
	return &KwkApp{Api:a, App:app}
}
