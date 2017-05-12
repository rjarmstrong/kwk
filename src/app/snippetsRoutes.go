package app

import (
	"github.com/urfave/cli"
)

func snippetsRoutes(s *snippets) []cli.Command {
	c := []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"create", "save"},
			Action: func(c *cli.Context) error {
				return s.Create(c.Args())
			},
		},
		{
			Name: "ls",
			Action: func(c *cli.Context) error {
				Prefs().HorizontalLists = true
				return s.List("", c.Args().First())
			},
		},
		{
			Name:    "enchilada",
			Aliases: []string{""},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "List all snippets.",
				},
			},
			Action: func(c *cli.Context) error {
				all := c.Bool("all")
				if all {
					Prefs().ListAll = true
				}
				return s.Flatten(c.Args().First())
			},
		},
		{
			Name:    "expand",
			Aliases: []string{"x"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "List all snippets.",
				},
			},
			Action: func(c *cli.Context) error {
				all := c.Bool("all")
				if all {
					Prefs().ListAll = true
				}
				Prefs().AlwaysExpandRows = true
				// TODO: This is not quite right as it means we can't expand other users lists
				return s.List("", c.Args().First())
			},
		},
		{
			Name:    "mkdir",
			Aliases: []string{""},
			Action: func(c *cli.Context) error {
				return s.CreatePouch(c.Args().First())
			},
		},
		{
			Name: "lock",
			Action: func(c *cli.Context) error {
				return s.Lock(c.Args().First())
			},
		},
		{
			Name: "unlock",
			Action: func(c *cli.Context) error {
				return s.UnLock(c.Args().First())
			},
		},
		{
			Name:    "find",
			Aliases: []string{"f"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				return s.Search(args...)
			},
		},

		/*
			The following are actions on existing snippets (and in some cases pouches):
		*/
		{
			Name: "describe",
			Action: func(c *cli.Context) error {
				return s.Describe(c.Args().Get(0), c.Args().Get(1))

			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Action: func(c *cli.Context) error {
				covert := c.Bool("covert")
				if covert {
					Prefs().Covert = true
				}
				return s.Run(c.Args().First(), []string(c.Args())[1:])

			},
		},
		{
			Name:    "view",
			Aliases: []string{"get"},
			Action: func(c *cli.Context) error {
				return s.InspectListOrRun(c.Args().First(), true)

			},
		},
		{
			Name:    "raw",
			Aliases: []string{"cat"},
			Action: func(c *cli.Context) error {
				return s.Cat(c.Args().First())

			},
		},
		{
			Name:    "rename",
			Aliases: []string{"mv", "move"},
			Action: func(c *cli.Context) error {
				return s.Move(c.Args())

			},
		},
		{
			Name:    "clone",
			Aliases: []string{"cp", "copy"},
			Action: func(c *cli.Context) error {
				return s.Clone(c.Args().First(), c.Args().Get(1))

			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Action: func(c *cli.Context) error {
				return s.Edit(c.Args().First())

			},
		},
		{
			Name:    "patch",
			Aliases: []string{"replace"},
			Action: func(c *cli.Context) error {
				return s.Patch(c.Args().First(), c.Args().Get(1), c.Args().Get(2))

			},
		},
		{
			Name:    "delete",
			Aliases: []string{"rm"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yes, y",
					Usage: "Automatically accept yes is modal dialogs.",
				},
			},
			Action: func(c *cli.Context) error {
				autoYes := c.Bool("yes")
				if autoYes {
					Prefs().AutoYes = true
				}
				return s.Delete(c.Args())

			},
		},
		{
			Name:    "tag",
			Aliases: []string{"t"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				return s.Tag(args[0], args[1:]...)

			},
		},
		{
			Name:    "untag",
			Aliases: []string{"ut"},
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
