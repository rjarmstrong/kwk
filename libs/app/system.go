package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/controllers"
)

func System(ctrl *controllers.SystemController) []cli.Command {
	c := []cli.Command{
		{
			Name:    "upgrade",
			Action: func(c *cli.Context) error {
				ctrl.Upgrade()
				return nil
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Action:  func(c *cli.Context) error {
				ctrl.GetVersion()
				return nil
			},
		},
		{
			Name:    "cd",
			Aliases: []string{},
			Action:  func(c *cli.Context) error {
				ctrl.ChangeDirectory(c.Args().Get(0))
				return nil
			},
		},
	}
	return c
}
