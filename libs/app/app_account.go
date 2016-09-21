package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/services/users"
)

func accountCommands(u users.IUsers) []cli.Command {
	c := []cli.Command{
		{
			Name:    "profile",
			Aliases: []string{"me"},
			Action: func(c *cli.Context) error {
				u.Get()
				return nil
			},
		},
		{
			Name:    "signin",
			Aliases: []string{"login"},
			Action:  func(c *cli.Context) error {
				u.Login(c.Args().Get(0), c.Args().Get(1))
				return nil
			},
		},
		{
			Name:    "signup",
			Aliases: []string{"register"},
			Action:  func(c *cli.Context) error {
				u.SignUp(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2));
				return nil
			},
		},
		{
			Name:    "signout",
			Aliases: []string{"logout"},
			Action:  func(c *cli.Context) error {
				u.Signout();
				return nil
			},
		},
	}
	return c
}
