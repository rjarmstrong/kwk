package app

import (
	"gopkg.in/urfave/cli.v1"
)

func Accounts(ctrl *AccountController) []cli.Command {
	c := []cli.Command{
		{
			Name:    "profile",
			Aliases: []string{"me"},
			Action: func(c *cli.Context) error {
				ctrl.Get()
				return nil
			},
		},
		{
			Name:    "signin",
			Aliases: []string{"login"},
			Action: func(c *cli.Context) error {
				ctrl.SignIn(c.Args().Get(0), c.Args().Get(1))
				return nil
			},
		},
		{
			Name:    "signup",
			Aliases: []string{"register"},
			Action: func(c *cli.Context) error {
				ctrl.SignUp()
				return nil
			},
		},
		{
			Name:    "signout",
			Aliases: []string{"logout"},
			Action: func(c *cli.Context) error {
				ctrl.SignOut()
				return nil
			},
		},
		{
			Name:    "cd",
			Aliases: []string{},
			Action: func(c *cli.Context) error {
				ctrl.ChangeDirectory(c.Args().Get(0))
				return nil
			},
		},
	}
	return c
}
