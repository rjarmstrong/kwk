package app

import (
	"github.com/urfave/cli"
)

func snippetsRoutes(s *snippets) []cli.Command {
	cat := "Snipppets"
	c := []cli.Command{
		{
			Category: cat,
			Name:     "new",
			Action: func(c *cli.Context) error {
				return s.Create(c.Args())
			},
		},
		{
			Category: cat,
			Name:     "view",
			Action: func(c *cli.Context) error {
				return s.InspectListOrRun(c.Args().First(), true)

			},
		},
		{
			Category: cat,
			Name:     "cat",
			Action: func(c *cli.Context) error {
				return s.Cat(c.Args().First())

			},
		},
		{
			Category: cat,
			Name:     "rm",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yes, y",
					Usage: "Automatically accept yes is modal dialogs.",
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
			Name:     "mv",
			Action: func(c *cli.Context) error {
				return s.Move(c.Args())

			},
		},
		{
			Category: cat,
			Name:     "edit",
			Aliases:  []string{"e"},
			Action: func(c *cli.Context) error {
				return s.Edit(c.Args().First())

			},
		},
		{
			Category: cat,
			Name:     "find",
			Flags: []cli.Flag {
				cli.BoolFlag{
					Name: "global, g",
					Usage: "Search everyone's public snippets, including yours.",
				},
			},
			Aliases:  []string{"f"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				return s.Search(args...)
			},
		},
		{
			Category: cat,
			Name:     "clone",
			Action: func(c *cli.Context) error {
				return s.Clone(c.Args().First(), c.Args().Get(1))

			},
		},

		{
			Category: cat,
			Name:     "enchilada",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "List all snippets.",
				},
			},
			Action: func(c *cli.Context) error {
				all := c.Bool("all")
				if all {
					prefs.ListAll = true
				}
				return s.Flatten(c.Args().First())
			},
		},
		{
			Category: cat,
			Name:     "describe",
			Action: func(c *cli.Context) error {
				return s.Describe(c.Args().Get(0), c.Args().Get(1))

			},
		},
		{
			Category: cat,
			Name:     "run",
			Aliases:  []string{"r"},
			Action: func(c *cli.Context) error {
				covert := c.Bool("covert")
				if covert {
					prefs.Covert = true
				}
				return s.Run(c.Args().First(), []string(c.Args())[1:])

			},
		},
		{
			Category: cat,
			Name:     "patch",
			Action: func(c *cli.Context) error {
				return s.Patch(c.Args().First(), c.Args().Get(1), c.Args().Get(2))

			},
		},
		{
			Category: cat,
			Name:     "tag",
			Aliases:  []string{"t"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				return s.Tag(args[0], args[1:]...)

			},
		},
		{
			Category: cat,
			Name:     "untag",
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				return s.UnTag(args[0], args[1:]...)

			},
		},
		//{
		//	Name:    "share",
		//	Aliases: []string{"send"},
		//	Action: func(c *cli.Context) error {
		//		return s.Share(c.Args().First(), c.Args().Get(2))
		//
		//	},
		//},
	}
	return c
}
