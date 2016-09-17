package app

import (
	"gopkg.in/urfave/cli.v1"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"github.com/kwk-links/kwk-cli/libs/api"
	"fmt"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/kwk-links/kwk-cli/libs/openers"
	"bufio"
	"os"
	"time"
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

var aliasTemplates = map[string]gui.Template{
	"inspect" : func(input interface{}) interface{} {
		if input != nil {
			system.PrettyPrint(input)
		} else {
			fmt.Println("Invalid kwklink")
		}
		return nil
	},
	"new" : func(input interface{}) interface{}{
		k := input.(*api.Alias)
		fmt.Println(k.FullKey)
		return nil
	},
	"cat" : func(input interface{}) interface{}{
		k := input.(*api.AliasList)
		fmt.Println(k)
		return nil
	},
	"notfound" : func(input interface{}) interface{}{
		fmt.Printf(gui.Colour(gui.Yellow, "kwklink: '%s' not found\n"), input)
		return nil
	},
	"patch" : func(input interface{}) interface{}{
		if input != nil {
			k := input.(*api.Alias)
			fmt.Printf(gui.Colour(gui.LightBlue, "Patched %s"), k.FullKey)
		} else {
			fmt.Println("Invalid kwklink")
		}
		return nil
	},
	// delete returns a boolean indicating whether the user agreed to delete or not.
	"delete" : func(input interface{}) interface{} {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf(gui.Colour(gui.LightBlue, "Are you sure you want to delete %s y/n? "), input)
		yesNo, _, _ := reader.ReadRune()
		return string(yesNo) == "y"
	},
	"deleted" : func(input interface{}) interface{}{
		fmt.Println("Deleted")
		return nil
	},
	"notdeleted": func(input interface{}) interface{}{
		messages := []string{"without a scratch", "uninjured", "intact", "unaffected", "unharmed",
			"unscathed", "out of danger", "safe and sound", "unblemished", "alive and well"}
		rnd := time.Now().Nanosecond() % (len(messages) - 1)
		fmt.Printf("'%s' is %s.\n", input, messages[rnd])
		return nil
	},
	/*
	Move to serverside
		originalKey := k.FullKey
			uri := k.Uri
			if c.Args().Get(1) != "" && c.Args().Get(2) != "" {
				uri = strings.Replace(uri, c.Args().Get(1), c.Args().Get(2), -1)
			}
			kwklink := ""
			if c.Args().Get(3) != "" {
				kwklink = c.Args().Get(3)
			}
			k = apiClient.Create(uri, kwklink)
	 */

	"clone": func(input interface{}) interface{}{
		k := input.(*api.Alias)
		if input != nil {
			fmt.Printf(gui.Colour(gui.LightBlue, "Cloned as %s"), k.FullKey)
		} else {
			fmt.Println("Invalid kwklink")
		}
		return nil
	},
}