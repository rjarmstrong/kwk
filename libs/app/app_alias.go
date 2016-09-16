package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/api"
	"fmt"
	"github.com/kwk-links/kwk-cli/libs/system"
)

func aliasCommands(a api.IApi, s system.ISystem, w gui.IWriter) []cli.Command {
	c := []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"create", "save"},
			Action:  func(c *cli.Context) error {
				if k := a.Create(c.Args().Get(0), c.Args().Get(1)); k != nil {
					s.CopyToClipboard(k.FullKey)
					w.PrintWithTemplate("new", k)
				}
				return nil
			},
		},
		{
			Name:    "inspect",
			Aliases: []string{"i"},
			Action: func(c *cli.Context) error {
				w.PrintWithTemplate("inspect", a.Get(c.Args().First()))
				return nil
			},
		},
		{
			Name:    "cat",
			Aliases: []string{"raw", "read", "print", "get"},
			Action:  func(c *cli.Context) error {
				alias := a.Get(c.Args().First())
				w.PrintWithTemplate("cat", alias)
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
		k := input.(api.KwkLink)
		fmt.Println(k.FullKey)
	},
	"cat" : func(input interface{}) {
		k := input.(api.KwkLinkList)
		fmt.Println(k)
	},
	"notfound" : func(input interface{}) {
		fmt.Printf(gui.Colour(gui.Yellow, "kwklink: '%s' not found\n"), input)
	},
}
