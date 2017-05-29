package app

import (
	"github.com/urfave/cli"
)

func systemRoutes(s *system) []cli.Command {
	cat := "System"
	c := []cli.Command{
		{
			Category: cat,
			Name: "update",
			Action: func(c *cli.Context) error {
				return s.Update()
			},
		},
		{
			Category: cat,
			Name:    "version",
			Aliases: []string{"v"},
			Action: func(c *cli.Context) error {
				return s.GetVersion()
			},
		},
	}
	return c
}
