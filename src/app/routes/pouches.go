package routes

import (
	"fmt"
	"github.com/rjarmstrong/kwk/src/app/handlers"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/style"
	"github.com/urfave/cli"
)

func Pouches(prefs *out.Prefs, s *handlers.Snippets) []cli.Command {
	cat := fmt.Sprintf("\n     %s  Pouches", style.Fmt256(style.ColorPouchCyan, style.IconPouch))
	spc := "  "
	c := []cli.Command{
		{
			Category: cat,
			Name:     "mkdir",
			Usage:    spc + "Create a pouch",
			Action: func(c *cli.Context) error {
				return s.CreatePouch(c.Args().First())
			},
		},
		{
			Category: cat,
			Name:     "mv",
			Usage:    spc + "Rename a pouch",
			Action: func(c *cli.Context) error {
				return s.Move(c.Args())

			},
		},
		{
			Category: cat,
			Name:     "rm",
			Usage:    spc + "Delete a pouch and all its contained snippets",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yes, y",
					Usage: "Automatically accept yes in modal dialogs.",
				},
			},
			Action: func(c *cli.Context) error {
				autoYes := c.Bool("yes")
				if autoYes {
					prefs.AutoYes = true
				}
				return s.Delete(c.Args())

			},
		},
		{
			Category: cat,
			Name:     "ls",
			Usage:    spc + "List snippets of pouch horizontally",
			Action: func(c *cli.Context) error {
				prefs.ListHorizontal = true
				return s.List("", c.Args().First())
			},
		},
		{
			Category: cat,
			Name:     "lock",
			Usage:    spc + "Make all snippets in or created in this pouch PRIVATE",
			Action: func(c *cli.Context) error {
				return s.Lock(c.Args().First())
			},
		},
		{
			Category: cat,
			Name:     "unlock",
			Usage:    spc + "Make all the snippets in or created in this pouch PUBLIC",
			Action: func(c *cli.Context) error {
				return s.UnLock(c.Args().First())
			},
		},
		{
			Category: cat,
			Name:     "expand",
			Usage:    spc + "Fully expand all snippets when viewing a pouch\n",
			Aliases:  []string{"x"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "List all snippets.",
				},
			},
			Action: func(c *cli.Context) error {
				all := c.Bool("all")
				if all {
					prefs.PrivateView = true
				}
				prefs.ExpandedRows = true
				// TODO: This is not quite right as it means we can't expand other users lists
				return s.List("", c.Args().First())
			},
		},
	}
	return c
}
