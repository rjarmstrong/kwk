package main


import (
"os"
"gopkg.in/urfave/cli.v1"
	"fmt"
	"io"
	"github.com/fatih/color"
	"github.com/kwk-links/kwk-cli/openers"
	"github.com/kwk-links/kwk-cli/api"
	"github.com/atotto/clipboard"
	"github.com/kwk-links/kwk-cli/system"
)

func main() {

	app := cli.NewApp()
	os.Setenv("version", "v0.0.1")
	settings := system.NewSettings("kwk.bolt.db")
	apiClient := api.New(settings)
	//app.EnableBashCompletion = true

	c := color.New(color.FgCyan).Add(color.Bold)
	cli.HelpPrinter = func(w io.Writer, template string, data interface{}) {
		c.Printf("\n ===================================================================== ")
		c.Printf("\n ~~~~~~~~~~~~~~~~~~~~~~~~   KWK Power Links.  ~~~~~~~~~~~~~~~~~~~~~~~~ \n\n")
		c.Printf(" The ultimate URI manager. Create short and memorable codes called\n")
		c.Printf(" `kwklinks` to store URLs, computer paths, AppLinks etc.\n\n")
		c.Printf(" Usage: kwk [kwklink|cmd] [subcmd] [args]\n")
		fmt.Print("\n e.g.: `kwk open got-spoilers` to open all G.O.T. spoiler websites.\n")

		c.Printf("\n Commands:\n")
		fmt.Print("    <kwklink,..>                      - Open and navigate to uris in default browser etc.\n")
		fmt.Print("    new        <uri> [name]           - Create a new kwklink, optionally provide a memorable name\n")
		fmt.Print("    list       [tag,..] [and|or|not]  - List kwklinks, filter by tags\n")
		fmt.Print("    search     [term] [tag]           - Search kwklinks and their metadata by keyword, filter by tags\n")
		fmt.Print("    suggest    <uri>                  - List suggested kwklinks or tags for the given uri\n")
		fmt.Print("    tag        <kwklink> [tag,..]     - Add tags to a kwklink\n")
		fmt.Print("    open       <tag>,.. [page]        - Open links for given tags, 5 at a time\n")
		fmt.Print("    untag      <kwklink> [tag,..]     - Remove tags from a kwklink\n")

		fmt.Print("    update\n")
		fmt.Print("      kwklink  <kwklink> <kwklink>    - Update kwklink <old> <new>\n")
		fmt.Print("      uri      <kwklink> <uri>        - Update uri\n")
		fmt.Print("    delete     <kwklink>              - Deletes kwklink with warning prompt. Will give 404.\n")
		fmt.Print("    detail     <kwklink>              - Get details and info\n")
		fmt.Print("    covert     <kwklink>              - Open in covert (incognito mode)\n")
		fmt.Print("    get        <kwklink> [page]       - Gets URIs without navigating. (Copies first to clipboard)\n")

		c.Printf("\n Analytics:\n")
		fmt.Print("    stats      [kwklink][tag]         - Get statistics and filter by kwklink or tag\n")

		c.Printf("\n Account:\n")
		fmt.Print("    login      <secret_key>           - Login with secret key.\n")
		fmt.Print("    logout                            - Clears locally cached secret key.\n")
		fmt.Print("    signup     <email> <password> <username>  - Sign-up with a username.\n")

		fmt.Print("\n\n  * Filter only Tags: today yesterday thisweek lastweek thismonth lastmonth thisyear lastyear")
		fmt.Print("\n ** kwklinks are case sensitive")

		fmt.Print("\n\n More Commands: `kwk [admin|device] help`")

		//Day II: fmt.Print("	lock       <kwklink> <pin>          - Lock a kwklink with a pin\n")
		//Day II: fmt.Print("	subscribe  <domain>	            - Subscribe with custom domain. Free for 30 days.\n")
		//Day II: fmt.Print("	rate <kwklink> 9	            - Subscribe with custom domain. Free for 30 days.\n")
		//Day II: fmt.Print("	note <kwklink> "I like this one	    - Subscribe with custom domain. Free for 30 days.\n")

		//fmt.Printf("\n Admin:\n")
		//fmt.Printf("	cache       ls                  - List locally cached kwklinks.\n")
		//fmt.Printf("	cache       clear               - Clears any locally cached data.\n")
		//fmt.Printf("	upgrade                    	- Downloads and upgrades kwk cli client.\n")
		//fmt.Printf("	config      warn  [on|off]      - Warns if attempting to open dodgy kwklink.\n")
		//fmt.Printf("	config      quiet [on|off]      - Prevents links from being printed to console.\n")
		//fmt.Printf("	version                    	\n")
		c.Printf("\n ===================================================================== \n\n")
	}

	app.CommandNotFound = func(context *cli.Context, kwklinkString string) {
		if k := apiClient.Decode(kwklinkString); k != nil {
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
				fmt.Println(uri)
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
