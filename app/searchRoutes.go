package app

import (
	"github.com/urfave/cli"
)

func Search(s *SearchCli) []cli.Command {
	c := []cli.Command{
		{
			Name:    "find",
			Aliases: []string{"search"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				s.Search(args...)
				return nil
			},
		},
	}
	return c
}
