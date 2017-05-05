package app

import (
	"github.com/urfave/cli"
)

func systemRoutes(s *system) []cli.Command {
	c := []cli.Command{
		{
			Name: "update",
			Action: func(c *cli.Context) error {
				return s.Update()
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Action: func(c *cli.Context) error {
				return s.GetVersion()
			},
		},
	}
	return c
}
