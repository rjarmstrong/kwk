package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/api"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/openers"
)

func aliasCommands(a api.IApi, s system.ISystem, i gui.IInteraction, o openers.IOpen) []cli.Command {
	c := []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"create", "save"},
			Action:  func(c *cli.Context) error {
				if k := a.Create(c.Args().Get(0), c.Args().Get(1)); k != nil {
					s.CopyToClipboard(k.FullKey)
					i.Respond("new", k)
				}
				return nil
			},
		},
		{
			Name:    "inspect",
			Aliases: []string{"i"},
			Action: func(c *cli.Context) error {
				i.Respond("inspect", a.Get(c.Args().First()))
				return nil
			},
		},
		{
			Name:    "cat",
			Aliases: []string{"raw", "read", "print", "get"},
			Action:  func(c *cli.Context) error {
				i.Respond("cat", a.Get(c.Args().First()))
				return nil
			},
		},
		{
			Name:    "rename",
			Aliases: []string{"mv"},
			Action:  func(c *cli.Context) error {
				i.Respond("rename", a.Rename(c.Args().Get(0), c.Args().Get(1)))
				return nil
			},
		},
		{
			Name:    "clone",
			Aliases: []string{},
			Action:  func(c *cli.Context) error {
				i.Respond("clone", a.Clone(c.Args().First(), c.Args().Get(1)))
				return nil
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Action:  func(c *cli.Context) error {
				key := c.Args().First()
				i.Respond("edit", o.Edit(key))
				return nil
			},
		},
		{
			Name:    "patch",
			Aliases: []string{"patch"},
			Action:  func(c *cli.Context) error {
				i.Respond("patch", a.Patch(c.Args().First(), c.Args()[1]))
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"rm"},
			Action:  func(c *cli.Context) error {
				fullKey := c.Args().First()
				if i.Respond("delete", fullKey).(bool) {
					a.Delete(fullKey)
					i.Respond("deleted", fullKey)
				} else {
					i.Respond("notdeleted", fullKey)
				}
				return nil
			},
		},

	}
	return c
}