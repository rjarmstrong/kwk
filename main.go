package main

import (
	"os"
	"github.com/kwk-links/kwk-cli/libs/charting"
	"gopkg.in/urfave/cli.v1"
	"fmt"
	"github.com/kwk-links/kwk-cli/libs/openers"
	"github.com/kwk-links/kwk-cli/libs/api"
	"github.com/atotto/clipboard"
	"github.com/kwk-links/kwk-cli/libs/system"
	"github.com/olekukonko/tablewriter"
	"github.com/dustin/go-humanize"
	"strings"
	"github.com/kwk-links/kwk-cli/libs/gui"
	"bufio"
	"time"
)

func main() {

	app := cli.NewApp()
	os.Setenv("version", "v0.0.1")
	settings := system.NewSettings("leveldb")
	apiClient := api.New(settings)
	cli.HelpPrinter = system.Help

	// run opener version checker

	app.Commands = []cli.Command{
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
