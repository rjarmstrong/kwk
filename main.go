package main


import (
"os"
"gopkg.in/urfave/cli.v1"
	"fmt"
	"io"
	"github.com/fatih/color"
	"github.com/kwk-links/cli/openers"
	"github.com/kwk-links/cli/api"
	"github.com/atotto/clipboard"
)

func main() {
	app := cli.NewApp()
	//app.EnableBashCompletion = true

	c := color.New(color.FgCyan).Add(color.Bold)
	cli.HelpPrinter = func(w io.Writer, template string, data interface{}) {
		c.Printf("\n ===================================================================== ")
		c.Printf("\n ~~~~~~~~~~~~~~~~~~~~~~~~   KWK Power Links.  ~~~~~~~~~~~~~~~~~~~~~~~~ \n\n")
		c.Printf(" Manage any kind of link. \n\n")
		c.Printf(" Usage: kwk [cmd/kwklink] [subcmd] ... [args]\n")

		c.Printf("\n Commands:\n")
		fmt.Print("	<kwklink,..>               	- Open and navigate.\n")
		fmt.Print("  	create 	   <uri> [kwklink]     	- Create a kwklink and optionally provide a preferred kwklink.\n")
		fmt.Print("	suggest    <uri>              	- Returns a list of suggested kwklinks for the given uri.\n")
		fmt.Print("  	open-tag   [page_number]     	- Open links for given tag, 5 links per page.\n")
		fmt.Print("	tag        <kwklink> [tag1,..]  - Add tags to a kwklink.\n")
		fmt.Print("	untag      <kwklink> [tag1,..] 	- Remove tags from a kwklink.\n")
		fmt.Print("	lock       <kwklink> <pin>      - Lock a kwklink with a pin.\n")
		fmt.Print("	update     <kwklink> <new kwklink> - Made a mistake? Update a kwklink.\n")
		fmt.Print("	detail     <kwklink>            - Get details and info.\n")
		fmt.Print("	covert     <kwklink>            - Open in covert (incognito mode).\n")
		fmt.Print("	get        <kwklink>            - Get the URI without navigating and copies it.\n")
		fmt.Print("	search     <term> 	 	- Search kwklinks.\n")

		c.Printf("\n Analytics:\n")
		fmt.Print("	stats      [kwklink][range]   	- Get stats summary for all keys\n")
		fmt.Print("	stats      <kwklink>            - Get stats for a key plus sub keys.\n")

		c.Printf("\n Account:\n")
		fmt.Print("	login      <secret_key>  	- Login with secret key.\n")
		fmt.Print("	logout                     	- Clears locally cached secret key.\n")
		fmt.Print("	signup     <kwklink> <password> - Sign-up with a personal root kwklink.\n")

		//fmt.Printf("\n Admin:\n")
		//fmt.Printf("	cache       ls                  - List locally cached kwklinks.\n")
		//fmt.Printf("	cache       clear               - Clears any locally cached data.\n")
		//fmt.Printf("	upgrade                    	- Downloads and upgrades kwk cli client.\n")
		//fmt.Printf("	config      warn  [on|off]      - Warns if attempting to open dodgy kwklink.\n")
		//fmt.Printf("	config      quiet [on|off]      - Prevents links from being printed to console.\n")
		//fmt.Printf("	version                    	- Enough said.\n")
		c.Printf("\n ===================================================================== \n\n")
	}

	app.CommandNotFound = func(context *cli.Context, kwklink string) {
		c := &api.ApiClient{}
		uri := c.Decode(kwklink)
		openers.Open(uri)
	}
	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Action:  func(c *cli.Context) error {
				a := api.ApiClient{}
				kwklink := a.Create(c.Args().Get(0), c.Args().Get(1))
				clipboard.WriteAll(kwklink.Key)
				fmt.Println(kwklink.Key)
				return nil
			},
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Action:  func(c *cli.Context) error {
				a := &api.ApiClient{}
				uri := a.Decode(c.Args().First())
				fmt.Println(uri)
				return nil
			},
		},
		{
			Name:    "covert",
			Aliases: []string{"c"},
			Action:  func(c *cli.Context) error {
				a := &api.ApiClient{}
				uri := a.Decode(c.Args().First())
				openers.OpenCovert(uri)
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
