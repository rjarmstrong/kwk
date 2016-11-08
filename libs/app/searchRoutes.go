package app

import (
	"gopkg.in/urfave/cli.v1"
)

func Search(ctrl *SearchController) []cli.Command {
	c := []cli.Command{
		{
			Name:    "find",
			Aliases: []string{"search"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				ctrl.Search(args...)
				return nil
			},
		},
	}
	return c
}
