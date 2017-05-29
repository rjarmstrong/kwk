package app

import (
	"github.com/urfave/cli"
)

func userRoutes(a *users) []cli.Command {
	cat := "Account"
	c := []cli.Command{
		{
			Category: cat,
			Usage: "Create an account with kwk",
			Name:    "signup",
			Action: func(c *cli.Context) error {
				return a.SignUp()
			},
		},
		{
			Category: cat,
			Name:    "login",
			Usage: "Login to kwk using username and password",
			Action: func(c *cli.Context) error {
				return a.SignIn(c.Args().Get(0), c.Args().Get(1))
			},
		},
		{
			Category: cat,
			Usage: "Logout from kwk and remove all locally cached data",
			Name:    "logout",
			Action: func(c *cli.Context) error {
				return a.SignOut()
			},
		},
		{
			Category: cat,
			Usage: "Show currently signed in user",
			Name:    "me",
			Action: func(c *cli.Context) error {
				return a.Profile()
			},
		},
		{
			Category: cat,
			Name: "forgot-pass",
			Usage: "Send a password reset code to the given email address",
			Action: func(c *cli.Context) error {
				return a.ResetPassword(c.Args().First())
			},
		},
		{
			Category: cat,
			Name: "change-pass",
			Usage: "Change current password to provided new password",
			Action: func(c *cli.Context) error {
				return a.ChangePassword()
			},
		},
	}
	return c
}
