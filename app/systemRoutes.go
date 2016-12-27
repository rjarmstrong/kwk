package app

import (
	"gopkg.in/urfave/cli.v1"
)

func System(s *SystemCli) []cli.Command {
	c := []cli.Command{
		{
			Name: "upgrade",
			Action: func(c *cli.Context) error {
				s.Upgrade()
				return nil
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Action: func(c *cli.Context) error {
				s.GetVersion()
				return nil
			},
		},
	}
	return c
}
