package main


import (
"os"
"gopkg.in/urfave/cli.v1"
	"fmt"
	//"io"
)

func main() {
	app := cli.NewApp()
	app.CommandNotFound = func(context *cli.Context, name string) {

		fmt.Println("opening url %!", name)
	}
	//app.EnableBashCompletion = true
	//cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
	//	fmt.Println("Ha HA.  I pwnd the help!!1")
	//}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action:  func(c *cli.Context) error {
				fmt.Println("added task: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action:  func(c *cli.Context) error {
				fmt.Println("completed task: ", c.Args().First())
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
		{
			Name:    "*",
			Aliases: []string{"*"},
			Usage:   "fallthrough",
			HideHelp: true,
			Action:  func(c *cli.Context) error {
				fmt.Println("fall-through: ", c.Args().First())
				return nil
			},
		},
	}

	app.Run(os.Args)
}
