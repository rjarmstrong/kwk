package app

import (
	"github.com/urfave/cli"
)

func System(s *SystemCli) []cli.Command {
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
		{
			Name:    "apperr",
			Aliases: []string{},
			Action: func(c *cli.Context) error {
				s.TestAppErr(c.Args().First() == "multi")
				return nil
			},
		},
		{
			Name:    "transerr",
			Aliases: []string{},
			Action: func(c *cli.Context) error {
				s.TestTransErr()
				return nil
			},
		},
	}
	return c
}
