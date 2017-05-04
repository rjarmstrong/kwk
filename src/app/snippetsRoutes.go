package app

import (
	"github.com/urfave/cli"
	"bitbucket.com/sharingmachine/kwkcli/src/models"
)

func snippetsRoutes(s *snippets) []cli.Command {
	c := []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"create", "save"},
			Action: func(c *cli.Context) error {
				s.Create(c.Args())
				return nil
			},
		},
		{
			Name:    "ls",
			Action: func(c *cli.Context) error {
				models.Prefs().HorizontalLists = true;
				s.List("", c.Args().First())
				return nil
			},
		},
		{
			Name:    "enchilada",
			Aliases: []string{""},
			Flags: []cli.Flag {
				cli.BoolFlag{
					Name: "all, a",
					Usage: "List all snippets.",
				},
			},
			Action: func(c *cli.Context) error {
				all := c.Bool("all")
				if all {
					models.Prefs().ListAll = true
				}
				s.Flatten(c.Args().First())
				return nil
			},
		},
		{
			Name:    "expand",
			Aliases: []string{"x"},
			Flags: []cli.Flag {
				cli.BoolFlag{
					Name: "all, a",
					Usage: "List all snippets.",
				},
			},
			Action: func(c *cli.Context) error {
				all := c.Bool("all")
				if all {
					models.Prefs().ListAll = true
				}
				models.Prefs().AlwaysExpandRows = true
				// TODO: This is not quite right as it means we can't expand other users lists
				s.List("", c.Args().First())
				return nil
			},
		},
		{
			Name:    "mkdir",
			Aliases: []string{""},
			Action: func(c *cli.Context) error {
				s.CreatePouch(c.Args().First())
				return nil
			},
		},
		{
			Name:    "lock",
			Action: func(c *cli.Context) error {
				s.Lock(c.Args().First())
				return nil
			},
		},
		{
			Name:    "unlock",
			Action: func(c *cli.Context) error {
				s.UnLock(c.Args().First())
				return nil
			},
		},
		{
			Name:    "find",
			Aliases: []string{"f"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				s.Search(args...)
				return nil
			},
		},


		/*
		The following are actions on existing snippets (and in some cases pouches):
		 */
		{
			Name:    "describe",
			Action: func(c *cli.Context) error {
				s.Describe(c.Args().Get(0), c.Args().Get(1))
				return nil
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Action: func(c *cli.Context) error {
				covert := c.Bool("covert")
				if covert {
					models.Prefs().Covert = true
				}
				s.Run(c.Args().First(), []string(c.Args())[1:])
				return nil
			},
		},
		{
			Name:    "view",
			Aliases: []string{"get"},
			Action: func(c *cli.Context) error {
				s.InspectListOrRun(c.Args().First(), true)
				return nil
			},
		},
		{
			Name:    "raw",
			Aliases: []string{"cat"},
			Action: func(c *cli.Context) error {
				s.Cat(c.Args().First())
				return nil
			},
		},
		{
			Name:    "rename",
			Aliases: []string{"mv", "move"},
			Action: func(c *cli.Context) error {
				s.Move(c.Args())
				return nil
			},
		},
		{
			Name:    "clone",
			Aliases: []string{"cp", "copy"},
			Action: func(c *cli.Context) error {
				s.Clone(c.Args().First(), c.Args().Get(1))
				return nil
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Action: func(c *cli.Context) error {
				s.Edit(c.Args().First())
				return nil
			},
		},
		{
			Name:    "patch",
			Aliases: []string{"replace"},
			Action: func(c *cli.Context) error {
				s.Patch(c.Args().First(), c.Args().Get(1), c.Args().Get(2))
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"rm"},
			Flags: []cli.Flag {
				cli.BoolFlag{
					Name: "yes, y",
					Usage: "Automatically accept yes is modal dialogs.",
				},
			},
			Action: func(c *cli.Context) error {
				autoYes := c.Bool("yes")
				if autoYes {
					models.Prefs().AutoYes = true
				}
				s.Delete(c.Args())
				return nil
			},
		},
		{
			Name:    "tag",
			Aliases: []string{"t"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				s.Tag(args[0], args[1:]...)
				return nil
			},
		},
		{
			Name:    "untag",
			Aliases: []string{"ut"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				s.UnTag(args[0], args[1:]...)
				return nil
			},
		},
		{
			Name:    "share",
			Aliases: []string{"send"},
			Action: func(c *cli.Context) error {
				s.Share(c.Args().First(), c.Args().Get(2))
				return nil
			},
		},
	}
	return c
}
