package app

import (
	"github.com/urfave/cli"
	"bitbucket.com/sharingmachine/kwkcli/models"
)

func Snippets(s *SnippetCli) []cli.Command {
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
			Name:    "list",
			Aliases: []string{"ls"},
			Flags: []cli.Flag {
				cli.BoolFlag{
					Name: "all, a",
					Usage: "List all snippets.",
				},
				cli.IntFlag{
					Name: "expand, x",
					Usage: "Expand list item to N rows.",
				},
			},
			Action: func(c *cli.Context) error {
				all := c.Bool("all")
				if all {
					models.Prefs().ListAll = true
				}
				expand := c.Int("expand")
				if expand > 0 {
					models.Prefs().ExpandLines = expand
				}
				s.List([]string(c.Args())...)
				return nil
			},
		},
		{
			Name:    "mkpouch",
			Aliases: []string{"mkdir"},
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
			Name:    "suggest",
			Action: func(c *cli.Context) error {
				s.Suggest(c.Args().First())
				return nil
			},
		},
		{
			Name:    "find",
			Aliases: []string{"search"},
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
			Name:    "inspect",
			Aliases: []string{"i"},
			Action: func(c *cli.Context) error {
				s.Inspect(c.Args().First())
				return nil
			},
		},
		{
			Name:    "cat",
			Aliases: []string{"raw", "read", "print", "get"},
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
