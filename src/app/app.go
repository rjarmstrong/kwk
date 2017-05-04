package app

import (
	"bitbucket.com/sharingmachine/kwkcli/src/cmd"
	"bitbucket.com/sharingmachine/kwkcli/src/gokwk"
	"bitbucket.com/sharingmachine/kwkcli/src/models"
	"bitbucket.com/sharingmachine/kwkcli/src/persist"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/src/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/src/update"
	"bitbucket.com/sharingmachine/types"
	"fmt"
	"github.com/urfave/cli"
	"strings"
	"bitbucket.com/sharingmachine/kwkcli/src/ui"
)

var CLIInfo = types.AppInfo{}

type KwkApp struct {
	App            *cli.App
	Acc            gokwk.Users
	Snippets       gokwk.Snippets
	File           persist.IO
	Settings       persist.Persister
	Updater        *update.Runner
	Runner         cmd.Runner
	Dialogue       dlg.Dialog
	TemplateWriter tmpl.Writer
}

func NewApp(a gokwk.Snippets, f persist.IO, t persist.Persister, r cmd.Runner, u gokwk.Users,
	d dlg.Dialog, w tmpl.Writer, up *update.Runner) *KwkApp {

	ap := cli.NewApp()

	ap.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "covert, x",
			Usage:       "Open browser in covert mode.",
			Destination: &models.Prefs().Covert,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Debug.",
			Destination: &models.DebugEnabled,
			EnvVar:      "DEBUG",
		},
		cli.BoolFlag{
			Name:        "naked, n",
			Usage:       "list without styles",
			Destination: &models.Prefs().Naked,
		},
		cli.BoolFlag{
			Name:        "ansi",
			Usage:       "Prints ansi escape sequences for debugging purposes",
			Destination: &models.Prefs().PrintAnsi,
		},

		//cli.BoolFlag{
		//	Name: "global, g",
		//	Usage: "Searches all users public snippets plus your private ones.",
		//},
		//cli.BoolFlag{
		//	Name: "encrypt, e",
		//	Usage: "When creating a new snippet encrypt.",
		//},
		//cli.BoolFlag{
		//	Name: "decrypt, d",
		//	Usage: "When viewing a snippet decrypt if necc.",
		//},
	}

	ap.Version = CLIInfo.String()
	dash := NewDashBoard(w, a)
	help := cli.HelpPrinter
	ap.Commands = append(ap.Commands, cli.Command{
		Name:"help",
		Aliases:[]string{"h"},
		Action: func(c *cli.Context) error {
			cli.HelpPrinter = help
			cli.ShowAppHelp(c)
			return nil
		},
	})

	cli.ErrWriter = ui.ErrWriter{}
	cli.HelpPrinter = dash.GetWriter()
	accCli := NewAccountCli(u, t, w, d, dash)
	ap.Commands = append(ap.Commands, userRoutes(accCli)...)

	sysCli := NewSystemCli(w, up)
	ap.Commands = append(ap.Commands, systemRoutes(sysCli)...)

	snipCli := NewSnippet(a, r, f, d, w, t)
	ap.Commands = append(ap.Commands, snippetsRoutes(snipCli)...)
	ap.CommandNotFound = func(c *cli.Context, distinctName string) {
		i := c.Args().Get(1)
		if strings.HasPrefix(distinctName, "@") {
			fmt.Println("listing:", distinctName)
			snipCli.GetEra(distinctName)
			return
		}
		switch i {
		case "run":
			snipCli.Run(c.Args().First(), []string(c.Args())[2:])
			return
		case "r":
			snipCli.Run(c.Args().First(), []string(c.Args())[2:])
			return
		case "edit":
			snipCli.Edit(c.Args().First())
			return
		case "e":
			snipCli.Edit(c.Args().First())
			return
		}
		snipCli.InspectListOrRun(c.Args().First(), false, []string(c.Args())[1:]...)
	}

	return &KwkApp{
		App:            ap,
		File:           f,
		Settings:       t,
		Runner:         r,
		Acc:            u,
		Dialogue:       d,
		Snippets:       a,
		TemplateWriter: w,
	}
}

func (a *KwkApp) Run(args ...string) {
	params := []string{"[app]"}
	params = append(params, args...)
	a.App.Run(params)
}
