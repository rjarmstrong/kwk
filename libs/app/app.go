package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/api"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/gui"
)

type KwkApp struct {
	App *cli.App
	Api api.IApi
}

func NewKwkApp(a api.IApi, s system.ISystem, w gui.IWriter) *KwkApp {
	app := cli.NewApp()
	app.Commands  = append(app.Commands, accountCommands(a)...)
	app.Commands  = append(app.Commands, systemCommands(s, w)...)
	return &KwkApp{Api:a, App:app}
}
