package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/api"
	"fmt"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/openers"
)

func aliasCommands(a api.IApi, s system.ISystem, i gui.IInteraction, o openers.IOpen) []cli.Command {
	m := NewMultiResultPrompt(o, i)

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
				alias := a.Get(c.Args().First())
				i.Respond("cat", alias)
				return nil
			},
		},
		{
			Name:    "rename",
			Aliases: []string{"mv"},
			Action:  func(c *cli.Context) error {
				a.Rename(c.Args().Get(0), c.Args().Get(1));
				i.Respond("rename", &api.Alias{})
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
			Name:    "covert",
			Aliases: []string{"c"},
			Action:  func(c *cli.Context) error {
				fullKey := c.Args().First()
				m.CheckAndPrompt(fullKey, a.Get(fullKey), c.Args())
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

	}
	return c
}

var aliasTemplates = map[string]gui.Template{
	"inspect" : func(input interface{}) {
		if input != nil {
			system.PrettyPrint(input)
		} else {
			fmt.Println("Invalid kwklink")
		}
	},
	"new" : func(input interface{}) {
		k := input.(*api.Alias)
		fmt.Println(k.FullKey)
	},
	"cat" : func(input interface{}) {
		k := input.(*api.AliasList)
		fmt.Println(k)
	},
	"notfound" : func(input interface{}) {
		fmt.Printf(gui.Colour(gui.Yellow, "kwklink: '%s' not found\n"), input)
	},
	"patch" : func(input interface{}) {
		if input != nil {
			k := input.(*api.Alias)
			fmt.Printf(gui.Colour(gui.LightBlue, "Patched %s"), k.FullKey)
		} else {
			fmt.Println("Invalid kwklink")
		}
	},

}