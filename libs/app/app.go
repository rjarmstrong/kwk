package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/api"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/openers"
)

type KwkApp struct {
	App *cli.App
	Api api.IApi
}

func NewKwkApp(a api.IApi, s system.ISystem, w gui.IWriter, o openers.IOpen) *KwkApp {

	app := cli.NewApp()

	app.Commands = append(app.Commands, accountCommands(a)...)
	app.Commands = append(app.Commands, systemCommands(s, w)...)
	app.Commands = append(app.Commands, aliasCommands(a, s, w)...)

	app.CommandNotFound = func(c *cli.Context, fullKey string) {
		if k := a.Get(fullKey); k != nil {
			if k.Total == 1 {
				o.Open(&k.Items[0], c.Args()[1:])
			} else if k.Total > 1 {

			} else {
				w.PrintWithTemplate("notfound", fullKey)
			}
			return
		}
	}
	return &KwkApp{Api:a, App:app}
}
