package app

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/urfave/cli"
)

func snippetsRoutes(s *snippets) []cli.Command {
	cat := fmt.Sprintf("  %s  Snippets", style.Fmt256(style.ColorPouchCyan, style.IconSnippet))
	spc := "  "
	c := []cli.Command{
		{
			Category: cat,
			Name:     "new",
			Usage:    spc + "Create a new snippet",
			Aliases:  []string{"n"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "pipe, p",
					Usage: "Pipe output creates content of new snippet",
				},
			},
			Action: func(c *cli.Context) error {
				pipe := c.Bool("pipe")
				return s.Create(c.Args(), pipe)
			},
		},
		{
			Category: cat,
			Name:     "view",
			Aliases:  []string{"v"},
			Usage:    spc + "View details of snippet",
			Action: func(c *cli.Context) error {
				return s.InspectListOrRun(c.Args().First(), true)

			},
		},
		{
			Category: cat,
			Name:     "edit",
			Usage:    spc + "Edit a snippet in the configured editor",
			Aliases:  []string{"e"},
			Action: func(c *cli.Context) error {
				return s.Edit(c.Args().First())

			},
		},
		{
			Category: cat,
			Name:     "find",
			Usage:    spc + "Find a snippet by keyword",
			Aliases:  []string{"f"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				return s.Search(args...)
			},
		},
		{
			Category: cat,
			Name:     "tag",
			Usage:    spc + "Add n number of tags",
			Aliases:  []string{"t"},
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				return s.Tag(args[0], args[1:]...)

			},
		},
		{
			Category: cat,
			Name:     "untag",
			Usage:    spc + "Remove n number of tags",
			Action: func(c *cli.Context) error {
				args := []string(c.Args())
				return s.UnTag(args[0], args[1:]...)

			},
		},
		{
			Category: cat,
			Name:     "describe",
			Usage:    spc + "Add a short description to a snippet",
			Action: func(c *cli.Context) error {
				return s.Describe(c.Args().Get(0), c.Args().Get(1))

			},
		},
		{
			Category: cat,
			Name:     "cat",
			Usage:    spc + "Print the raw content of the snippet",
			Action: func(c *cli.Context) error {
				return s.Cat(c.Args().First())

			},
		},
		{
			Category: cat,
			Name:     "mv",
			Usage:    spc + "Move snippets between pouches",
			Action: func(c *cli.Context) error {
				return s.Move(c.Args())

			},
		},
		{
			Category: cat,
			Name:     "rm",
			Usage:    spc + "Delete snippets",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yes, y",
					Usage: "Automatically accept yes is modal dialogs",
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
			Name:     "clone",
			Usage:    spc + "Clone another snippet given any snippet URI",
			Action: func(c *cli.Context) error {
				return s.Clone(c.Args().First(), c.Args().Get(1))

			},
		},
		{
			Category: cat,
			Name:     "enchilada",
			Usage:    spc + "List all snippets 'un-pouched', useful for exporting or for other bulk management",
			Action: func(c *cli.Context) error {
				return s.Flatten(c.Args().First())
			},
		},
		{
			Category: cat,
			Name:     "run",
			Usage:    spc + "When prefs.yml has RequireRunKeyword=true, run is required to execute a snippet",
			Aliases:  []string{"r"},
			Action: func(c *cli.Context) error {
				return s.Run(c.Args().First(), []string(c.Args())[1:])

			},
		},
		{
			Category: cat,
			Name:     "patch",
			Usage:    spc + "Replace matching string within a snippet\n",
			Action: func(c *cli.Context) error {
				return s.Patch(c.Args().First(), c.Args().Get(1), c.Args().Get(2))

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
