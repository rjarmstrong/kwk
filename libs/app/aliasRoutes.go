package app

import (
	"gopkg.in/urfave/cli.v1"
)

func Alias(ctrl *AliasController) []cli.Command {
	c := []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"create", "save"},
			Action:  func(c *cli.Context) error {
				ctrl.New(c.Args().Get(0), c.Args().Get(1))
				return nil
			},
		},
		{
			Name:    "inspect",
			Aliases: []string{"i"},
			Action: func(c *cli.Context) error {
				ctrl.Inspect(c.Args().First())
				return nil
			},
		},
		{
			Name:    "cat",
			Aliases: []string{"raw", "read", "print", "get"},
			Action:  func(c *cli.Context) error {
				ctrl.Cat(c.Args().First())
				return nil
			},
		},
		{
			Name:    "rename",
			Aliases: []string{"mv"},
			Action:  func(c *cli.Context) error {
				ctrl.Rename(c.Args().Get(0), c.Args().Get(1))
				return nil
			},
		},
		{
			Name:    "clone",
			Aliases: []string{},
			Action:  func(c *cli.Context) error {
				ctrl.Clone(c.Args().First(), c.Args().Get(1))
				return nil
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Action:  func(c *cli.Context) error {
				ctrl.Edit(c.Args().First())
				return nil
			},
		},
		{
			Name:    "patch",
			Aliases: []string{"replace"},
			Action:  func(c *cli.Context) error {
				ctrl.Patch(c.Args().First(), c.Args().Get(1), c.Args().Get(2))
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"rm"},
			Action:  func(c *cli.Context) error {
				ctrl.Delete(c.Args().First())
				return nil
			},
		},
		{
			Name:    "tag",
			Aliases: []string{"t"},
			Action:  func(c *cli.Context) error {
				args := []string(c.Args())
				ctrl.Tag(args[0], args[1:]...)
				return nil
			},
		},
		{
			Name:    "untag",
			Aliases: []string{"ut"},
			Action:  func(c *cli.Context) error {
				args := []string(c.Args())
				ctrl.UnTag(args[0], args[1:]...)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Action:  func(c *cli.Context) error {
				ctrl.List([]string(c.Args())...)
				return nil
			},
		},

	}
	return c
}