package app

import (
	"github.com/urfave/cli"
)

func systemRoutes(s *system) []cli.Command {
	c := []cli.Command{
		{
			Name: "update",
			Action: func(c *cli.Context) error {
				s.Update()
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
