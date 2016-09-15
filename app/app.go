package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/api"
)

type KwkApp struct {
	App *cli.App
	Api api.IApi
}

func NewKwkApp(a api.IApi) *KwkApp {
	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "profile",
			Aliases: []string{"me"},
			Action: func(c *cli.Context) error {
				a.PrintProfile()
				return nil
			},
		},
		{
			Name:    "signin",
			Aliases: []string{"login"},
			Action:  func(c *cli.Context) error {
				return nil
			},
		},
	}

	return &KwkApp{Api:a, App:app}
}
