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
)

func main() {

	app := cli.NewApp()
	os.Setenv("version", "v0.0.1")
	settings := system.NewSettings("kwk.bolt.db")
	apiClient := api.New(settings)
	cli.HelpPrinter = system.Help

	app.CommandNotFound = func(context *cli.Context, kwklinkString string) {
		if k := apiClient.Decode(kwklinkString); k != nil {
			fmt.Println(k)
			openers.Open(k.Uri)
			return
		}
		fmt.Println("Command or kwklink not found.")

	}
	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"create","save"},
			Action:  func(c *cli.Context) error {
				if k := apiClient.Create(c.Args().Get(0), c.Args().Get(1)); k != nil {
					clipboard.WriteAll(k.Key)
					fmt.Println(k.Key)
				}
				return nil
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Action:  func(c *cli.Context) error {
				uri := apiClient.Decode(c.Args().First())
				fmt.Println(uri.Uri)
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Action:  func(c *cli.Context) error {
				//c.Args().First()
				list := apiClient.List(1)
				fmt.Println()
				table := tablewriter.NewWriter(os.Stdout)
				//table.SetHeader([]string{Colour(subdued, "Kwklink"), Colour(subdued, "Type"), Colour(subdued, "URI"), Colour(subdued, "Tags"), ""})
				table.SetAutoWrapText(false)
				table.SetBorder(false)
				table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
				table.SetCenterSeparator("")
				table.SetColumnSeparator("")
				table.SetAutoFormatHeaders(false)
				table.SetHeaderLine(false)
				table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

				for _, v := range list.Items {
					fmt.Printf("\n%+q %s", v.Uri, v.Uri)
					v.Uri = strings.Replace(v.Uri, "https://", "", 1)
					v.Uri = strings.Replace(v.Uri, "http://", "", 1)
					v.Uri = strings.Replace(v.Uri, "www.", "", 1)
					if len(v.Uri) >= 40 {
						v.Uri = v.Uri[0:10] + gui.Colour(gui.Subdued, "...") + v.Uri[len(v.Uri)-30:len(v.Uri)]
					}

					table.Append([]string{gui.Colour(gui.LightBlue, v.Key), "web", v.Uri, "Hot,Fake,Fresh", humanize.Time(v.Created)})

				}
				table.Render()

				nextcmd := fmt.Sprintf("For next page run: kwk list %v", 2)
				fmt.Printf("\n %v of %v pages \t #tip %s", 1, 11, gui.Colour(gui.Subdued, nextcmd))
				fmt.Print("\n\n")
				return nil
			},
		},
		{
			Name:    "covert",
			Aliases: []string{"c"},
			Action:  func(c *cli.Context) error {
				k := apiClient.Decode(c.Args().First())
				openers.OpenCovert(k.Uri)
				return nil
			},
		},
		{
			Name:	"upgrade",
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
			Name:        "template",
			Aliases:     []string{"t"},
			Usage:       "options for task templates",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Category: "zong",
					Action: func(c *cli.Context) error {
						fmt.Println("new task template: ", c.Args().First())
						return nil
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Category: "zong",
					Action: func(c *cli.Context) error {
						fmt.Println("removed task template: ", c.Args().First())
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
