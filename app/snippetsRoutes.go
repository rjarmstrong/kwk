package app

import (
	"gopkg.in/urfave/cli.v1"
)

func Snippets(s *SnippetCli) []cli.Command {
	c := []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"create", "save"},
			Action: func(c *cli.Context) error {
				s.New(c.Args().Get(0), c.Args().Get(1))
				return nil
			},
		},
		{
			Name:    "describe",
			Aliases: []string{"update", "d"},
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
			Aliases: []string{"mv"},
			Action: func(c *cli.Context) error {
				s.Rename(c.Args().Get(0), c.Args().Get(1))
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
			Name:    "clone",
			Aliases: []string{"fork", "copy", "c"},
			Action: func(c *cli.Context) error {
				s.Clone(c.Args().First(), c.Args().Get(1))
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
			Action: func(c *cli.Context) error {
				s.Delete(c.Args().First())
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
			Name:    "list",
			Aliases: []string{"ls"},
			Action: func(c *cli.Context) error {
				s.List([]string(c.Args())...)
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
		{
			Name:    "lock",
			Aliases: []string{},
			Action: func(c *cli.Context) error {
				s.Rename(c.Args().First(), "."+c.Args().First())
				return nil
			},
		},
		{
			Name:    "unlock",
			Aliases: []string{},
			Action: func(c *cli.Context) error {
				s.Rename(c.Args().First(), c.Args().First())
				return nil
			},
		},
	}
	return c
}
