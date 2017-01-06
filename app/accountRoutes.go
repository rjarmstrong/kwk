package app

import (
	"gopkg.in/urfave/cli.v1"
)

func Accounts(a *AccountCli) []cli.Command {
	c := []cli.Command{
		{
			Name:    "profile",
			Aliases: []string{"me", "whoami"},
			Action: func(c *cli.Context) error {
				a.Get()
				return nil
			},
		},
		{
			Name:    "signin",
			Aliases: []string{"login", "switch", "cd"},
			Action: func(c *cli.Context) error {
				a.SignIn(c.Args().Get(0), c.Args().Get(1))
				return nil
			},
		},
		{
			Name:    "signup",
			Aliases: []string{"register"},
			Action: func(c *cli.Context) error {
				a.SignUp()
				return nil
			},
		},
		{
			Name:    "signout",
			Aliases: []string{"logout"},
			Action: func(c *cli.Context) error {
				a.SignOut()
				return nil
			},
		},
	}
	return c
}
