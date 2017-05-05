package app

import (
	"github.com/urfave/cli"
)

func userRoutes(a *users) []cli.Command {
	c := []cli.Command{
		{
			Name:    "profile",
			Aliases: []string{"me", "whoami"},
			Action: func(c *cli.Context) error {
				return a.Get()
			},
		},
		{
			Name:    "signin",
			Aliases: []string{"login", "switch", "cd"},
			Action: func(c *cli.Context) error {
				return a.SignIn(c.Args().Get(0), c.Args().Get(1))
			},
		},
		{
			Name:    "signup",
			Aliases: []string{"register"},
			Action: func(c *cli.Context) error {
				return a.SignUp()
			},
		},
		{
			Name:    "signout",
			Aliases: []string{"logout"},
			Action: func(c *cli.Context) error {
				return a.SignOut()
			},
		},
		{
			Name: "reset-password",
			Action: func(c *cli.Context) error {
				return a.ResetPassword(c.Args().First())
			},
		},
		{
			Name: "change-password",
			Action: func(c *cli.Context) error {
				return a.ChangePassword()
			},
		},
	}
	return c
}
