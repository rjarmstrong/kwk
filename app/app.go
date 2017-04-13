package app

import (
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"github.com/urfave/cli"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/log"
	"strings"
	"fmt"
	"bitbucket.com/sharingmachine/kwkcli/file"
	"bitbucket.com/sharingmachine/kwkcli/user"
)

type KwkApp struct {
	App            *cli.App
	Snippets       snippets.Service
	File           file.IO
	Settings       config.Persister
	Acc            user.Account
	Runner         cmd.Runner
	Dialogue       dlg.Dialog
	TemplateWriter tmpl.Writer
	Api            rpc.Service
}

func NewApp(a snippets.Service, f file.IO, t config.Persister, r cmd.Runner, u user.Account,
	d dlg.Dialog, w tmpl.Writer, api rpc.Service) *KwkApp {

	ap := cli.NewApp()

	ap.Flags = []cli.Flag {
		cli.BoolFlag{
			Name: "covert, x",
			Usage: "Open browser in covert mode.",
			Destination:&models.Prefs().Covert,
		},
		cli.BoolFlag{
			Name: "debug, d",
			Usage: "Debug.",
			Destination:&log.EnableDebug,
			EnvVar:"DEBUG",
		},
		cli.BoolFlag {
			Name: "naked, n",
			Usage: "list without styles",
			Destination:&models.Prefs().Naked,
		},
		cli.BoolFlag {
			Name: "ansi",
			Usage: "Prints ansi escape sequences for debugging purposes",
			Destination:&models.Prefs().PrintAnsi,
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

	ap.Version = models.Client.String()
	dash := NewDashBoard(w, a)
	cli.HelpPrinter = dash.GetWriter()

	accCli := NewAccountCli(u, t, w, d, dash)
	ap.Commands = append(ap.Commands, Accounts(accCli)...)

	sysCli := NewSystemCli(f, api, w, t)
	ap.Commands = append(ap.Commands, System(sysCli)...)

	snipCli := NewSnippetCli(a, r, f, d, w, t)
	ap.Commands = append(ap.Commands, Snippets(snipCli)...)
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
		Api:            api,
	}
}

func (a *KwkApp) Run(args ...string) {
	params := []string{"[app]"}
	params = append(params, args...)
	a.App.Run(params)
}