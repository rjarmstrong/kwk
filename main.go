package main

import (
	"os"
	"gopkg.in/urfave/cli.v1"
	"fmt"
	"github.com/kwk-links/kwk-cli/openers"
	"github.com/kwk-links/kwk-cli/api"
	"github.com/atotto/clipboard"
	"github.com/kwk-links/kwk-cli/system"
	"github.com/olekukonko/tablewriter"
	"github.com/dustin/go-humanize"
	"strings"
	"github.com/kwk-links/kwk-cli/gui"
	"bufio"
	"time"
	"github.com/kwk-links/kwk-cli/charting"
)

func main() {

	app := cli.NewApp()
	os.Setenv("version", "v0.0.1")
	settings := system.NewSettings("leveldb")
	apiClient := api.New(settings)
	cli.HelpPrinter = system.Help
	opener := openers.NewOpener(apiClient)

	// run opener version checker

	app.CommandNotFound = func(c *cli.Context, kwklink string) {
		// if it is a fullKey then else get tag
		if k := apiClient.Get(kwklink); k != nil {
			if k.Uri != "" {
				opener.Open(k, c.Args()[1:])
			} else {
				fmt.Printf(gui.Colour(gui.Yellow, "kwklink: '%s' not found\n"), kwklink)
			}
			return
		}
		fmt.Println("Command or kwklink not found.")

	}
	app.Commands = []cli.Command{
		{
			Name:    "open",
			Aliases: []string{"o"},
			Action: func(c *cli.Context) error {
				args := c.Args()
				list := apiClient.List([]string(args))
				for _, v := range list.Items {
					fmt.Println(gui.Colour(gui.LightBlue, v.FullKey))
					opener.Open(&v, []string{})
				}
				return nil
			},
		},
		{
			Name:    "openc",
			Aliases: []string{"oc"},
			Action:  func(c *cli.Context) error {
				args := c.Args()
				list := apiClient.List([]string(args))
				for _, v := range list.Items {
					opener.OpenCovert(v.Uri)
				}
				return nil
			},
		},

		{
			Name:    "inspect",
			Aliases: []string{"i"},
			Action:  func(c *cli.Context) error {
				args := c.Args()
				link := apiClient.Get(args.First())
				if link != nil {
					system.PrettyPrint(link)
				} else {
					fmt.Println("Invalid kwklink")
				}
				return nil
			},
		},
		{
			Name:    "new",
			Aliases: []string{"create", "save"},
			Action:  func(c *cli.Context) error {
				if k := apiClient.Create(c.Args().Get(0), c.Args().Get(1)); k != nil {
					clipboard.WriteAll(k.FullKey)
					fmt.Println(k.FullKey)
				}
				return nil
			},
		},
		{
			Name:    "clone",
			Aliases: []string{"copy"},
			Action:  func(c *cli.Context) error {
				if k := apiClient.Get(c.Args().Get(0)); k != nil {
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
					clipboard.WriteAll(k.FullKey)
					fmt.Printf(gui.Colour(gui.LightBlue, "Cloned %s -> %s"), originalKey, k.FullKey)
				} else {
					fmt.Println("Invalid kwklink")
				}
				return nil
			},
		},
		{
			Name:    "edit",
			Aliases: []string{"e"},
			Action:  func(c *cli.Context) error {
				key := c.Args().First()
				if err := opener.Edit(key); err != nil {
					fmt.Println(err)
				}
				return nil
			},
		},
		{
			Name:    "cat",
			Aliases: []string{"raw", "read", "print", "get"},
			Action:  func(c *cli.Context) error {
				uri := apiClient.Get(c.Args().First())
				clipboard.WriteAll(uri.Uri)
				fmt.Println(uri.Uri)
				return nil
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"rm"},
			Action:  func(c *cli.Context) error {
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf(gui.Colour(gui.LightBlue, "Are you sure you want to delete %s y/n? "), c.Args().First())
				yesNo, _, _ := reader.ReadRune()
				if string(yesNo) == "y" {
					apiClient.Delete(c.Args().First())
					fmt.Println("Deleted")
				} else {
					messages := []string{"without a scratch", "uninjured", "intact", "unaffected", "unharmed", "unscathed", "out of danger", "safe and sound", "unblemished", "alive and well"}
					rnd := time.Now().Nanosecond() % (len(messages) - 1)
					fmt.Printf("'%s' is %s.\n", c.Args().First(), messages[rnd])
				}
				return nil
			},
		},
		{
			Name:    "tag",
			Aliases: []string{"t"},
			Action:  func(c *cli.Context) error {
				args := []string(c.Args())
				apiClient.Tag(args[0], args[1:]...)
				fmt.Println("Tagged")
				return nil
			},
		},
		{
			Name:    "untag",
			Aliases: []string{"ut"},
			Action:  func(c *cli.Context) error {
				args := []string(c.Args())
				apiClient.UnTag(args[0], args[1:]...)
				fmt.Println("UnTagged")
				return nil
			},
		},
		{
			Name:    "cd",
			Aliases: []string{},
			Action:  func(c *cli.Context) error {
				args := []string(c.Args())
				fmt.Println(gui.Colour(gui.LightBlue, "Switched to kwk.co/" + args[0] + "/"))
				return nil
			},
		},
		{
			Name:    "back",
			Aliases: []string{"b"},
			Action:  func(c *cli.Context) error {
				fmt.Print("Some text")
				//fmt.Printf("\x0c%s", "Some more text")
				fmt.Print(gui.ClearLine)
				fmt.Print(gui.MoveBack)
				fmt.Print("\u2588 ")
				fmt.Print("Some more text")
				fmt.Print(" \u2580")
				return nil
				//https://en.wikipedia.org/wiki/Block_Elements
				//https://en.wikipedia.org/wiki/Braille_Patterns#Chart
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Action:  func(c *cli.Context) error {

				args := c.Args()
				list := apiClient.List([]string(args))
				fmt.Print(gui.Colour(gui.LightBlue, "\nkwk.co/" + "rjarmstrong/"))
				fmt.Printf(gui.Build(102, " ") + "%d of %d records\n\n", len(list.Items), list.Total)

				tbl := tablewriter.NewWriter(os.Stdout)
				tbl.SetHeader([]string{"Alias", "Version", "URI", "Tags", ""})
				tbl.SetAutoWrapText(false)
				tbl.SetBorder(false)
				tbl.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
				tbl.SetCenterSeparator("")
				tbl.SetColumnSeparator("")
				tbl.SetAutoFormatHeaders(false)
				tbl.SetHeaderLine(true)
				tbl.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

				for _, v := range list.Items {
					v.Uri = strings.Replace(v.Uri, "https://", "", 1)
					v.Uri = strings.Replace(v.Uri, "http://", "", 1)
					v.Uri = strings.Replace(v.Uri, "www.", "", 1)
					v.Uri = strings.Replace(v.Uri, "\n", " ", -1)
					if len(v.Uri) >= 40 {
						v.Uri = v.Uri[0:10] + gui.Colour(gui.Subdued, "...") + v.Uri[len(v.Uri) - 30:len(v.Uri)]
					}

					var tags = []string{}
					for _, v := range v.Tags {
						if v == "error" {
							tags = append(tags, gui.Colour(gui.Pink, v))
						} else {
							tags = append(tags, v)
						}

					}

					tbl.Append([]string{
						gui.Colour(gui.LightBlue, v.Key) + gui.Colour(gui.Subdued, "." + v.Extension),
						fmt.Sprintf("%d", v.Version),
						fmt.Sprintf("%s", v.Uri),
						strings.Join(tags, ", "),
						humanize.Time(v.Created),
					})

				}
				tbl.Render()

				if len(list.Items) == 0 {
					fmt.Println(gui.Colour(gui.Yellow, "No records on this page! Use a lower page number.\n"))
				} else {
					//gui.Colour(gui.Subdued, nextcmd)
					//nextcmd := fmt.Sprintf("For next page run: kwk list %v", 2)
				}
				if list.Size != 0 {
					fmt.Printf("\n %d of %d pages", list.Page, (list.Total / list.Size) + 1)
				}
				fmt.Print("\n\n")
				return nil
			},
		},
		{
			Name:    "covert",
			Aliases: []string{"c"},
			Action:  func(c *cli.Context) error {
				k := apiClient.Get(c.Args().First())
				opener.OpenCovert(k.Uri)
				return nil
			},
		},
		{
			Name:        "upgrade",
			Action: func(c *cli.Context) error {
				system.Upgrade()
				return nil
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Action:  func(c *cli.Context) error {
				fmt.Println(os.Getenv("version"))
				return nil
			},
		},
		{
			Name:    "rename",
			Aliases: []string{"mv"},
			Action:  func(c *cli.Context) error {
				//TODO: When updating a pinned kwklink, must force to give a new name
				// (since it is technically no longer the original)
				apiClient.Rename(c.Args().Get(0), c.Args().Get(1));
				return nil
			},
		},
		{
			Name:    "patch",
			Aliases: []string{"patch"},
			Action:  func(c *cli.Context) error {
				//TODO: When updating a pinned kwklink, must force to give a new name
				// (since it is technically no longer the original)

				if k := apiClient.Get(c.Args().Get(0)); k != nil {
					uri := k.Uri
					if len(c.Args()) == 3 {
						uri = strings.Replace(uri, c.Args().Get(1), c.Args().Get(2), -1)
					} else {
						uri = c.Args().Get(1)
					}
					apiClient.Patch(k.FullKey, uri);
					clipboard.WriteAll(k.FullKey)
					fmt.Printf(gui.Colour(gui.LightBlue, "Patched %s"), k.FullKey)
				} else {
					fmt.Println("Invalid kwklink")
				}
				return nil
			},
		},
		{
			Name:    "login",
			Aliases: []string{"signin"},
			Action:  func(c *cli.Context) error {
				apiClient.Login(c.Args().Get(0), c.Args().Get(1));
				return nil
			},
		},
		{
			Name:    "logout",
			Aliases: []string{"signout"},
			Action:  func(c *cli.Context) error {
				apiClient.Logout();
				return nil
			},
		},
		{
			Name:    "signup",
			Aliases: []string{"register"},
			Action:  func(c *cli.Context) error {
				apiClient.SignUp(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2));
				return nil
			},
		},
		{
			Name:    "profile",
			Aliases: []string{"me"},
			Action:  func(c *cli.Context) error {
				apiClient.PrintProfile();
				return nil
			},
		},
		{
			Name:    "stats",
			Aliases: []string{"analytics"},
			Action:  func(c *cli.Context) error {
				list := apiClient.List(c.Args())
				charting.PrintTags(list)
				return nil
			},
		},
	}

	app.Run(os.Args)
}
