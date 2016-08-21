package main


import (
"os"
"gopkg.in/urfave/cli.v1"
	"fmt"
	"io"
	"github.com/fatih/color"
)

func main() {
	app := cli.NewApp()
	app.CommandNotFound = func(context *cli.Context, name string) {

		fmt.Printf("opening url %s!\n", name)
	}
	//app.EnableBashCompletion = true

	c := color.New(color.FgCyan).Add(color.Underline)
	cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
		c.Printf("=== Welcome to kwk ===")
	}

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
	}

	app.Run(os.Args)
}
