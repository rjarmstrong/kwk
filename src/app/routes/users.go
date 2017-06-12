package routes

import (
	"fmt"
	"github.com/rjarmstrong/kwk/src/app/handlers"
	"github.com/rjarmstrong/kwk/src/style"
	"github.com/urfave/cli"
)

func Users(a *handlers.Users) []cli.Command {
	cat := fmt.Sprintf("  %s  Account", style.Fmt256(style.ColorPouchCyan, style.IconAccount))
	c := []cli.Command{
		{
			Category: cat,
			Usage:    "Create an account with kwk",
			Name:     "signup",
			Action: func(c *cli.Context) error {
				return a.SignUp()
			},
		},
		{
			Category: cat,
			Name:     "login",
			Usage:    "Login to kwk using username and password",
			Aliases:  []string{"signin"},
			Action: func(c *cli.Context) error {
				return a.SignIn(c.Args().Get(0), c.Args().Get(1))
			},
		},
		{
			Category: cat,
			Usage:    "Logout from kwk and remove all locally cached data",
			Name:     "logout",
			Aliases:  []string{"signout"},
			Action: func(c *cli.Context) error {
				return a.SignOut()
			},
		},
		{
			Category: cat,
			Usage:    "Show currently signed in user",
			Name:     "me",
			Action: func(c *cli.Context) error {
				return a.Profile()
			},
		},
		{
			Category: cat,
			Name:     "forgot-pass",
			Usage:    "Send a password reset code to the given email address",
			Action: func(c *cli.Context) error {
				return a.ForgotPassword(c.Args().First())
			},
		},
		{
			Category: cat,
			Name:     "change-pass",
			Usage:    "Change current password to provided new password\n",
			Action: func(c *cli.Context) error {
				return a.ChangePassword()
			},
		},
	}
	return c
}
