package app

import (
	"gopkg.in/urfave/cli.v1"
)

func System(ctrl *SystemController) []cli.Command {
	c := []cli.Command{
		{
			Name: "upgrade",
			Action: func(c *cli.Context) error {
				ctrl.Upgrade()
				return nil
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Action: func(c *cli.Context) error {
				ctrl.GetVersion()
				return nil
			},
		},
	}
	return c
}
