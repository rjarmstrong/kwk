package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/system"
)

func systemCommands(s system.ISystem, w gui.IInteraction) []cli.Command {
	c := []cli.Command{
		{
			Name:    "upgrade",
			Action: func(c *cli.Context) error {
				s.Upgrade()
				return nil
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Action:  func(c *cli.Context) error {
				w.Respond("version", s.GetVersion())
				return nil
			},
		},
		{
			Name:    "cd",
			Aliases: []string{},
			Action:  func(c *cli.Context) error {
				s.ChangeDirectory(c.Args().Get(0))
				w.Respond("cd", c.Args().Get(0))
				return nil
			},
		},
	}
	return c
}
