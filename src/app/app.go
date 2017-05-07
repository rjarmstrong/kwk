package app

import (
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/cli/src/exekwk/cmd"
	"github.com/kwk-super-snippets/cli/src/exekwk/setup"
	"github.com/kwk-super-snippets/cli/src/exekwk/update"
	"github.com/kwk-super-snippets/cli/src/gokwk"
	"github.com/kwk-super-snippets/cli/src/models"
	"github.com/kwk-super-snippets/cli/src/persist"
	"bitbucket.com/sharingmachine/types"
	"bitbucket.com/sharingmachine/types/errs"
	"bitbucket.com/sharingmachine/types/vwrite"
	"fmt"
	"github.com/urfave/cli"
	"strings"
)

var CLIInfo = types.AppInfo{}

type KwkApp struct {
	App      *cli.App
	Acc      gokwk.Users
	Snippets gokwk.Snippets
	File     persist.IO
	Settings persist.Persister
	Updater  *update.Updater
	Runner   cmd.Runner
	Dialogue Dialog
	vwrite.Writer
	errs.Handler
}

func NewApp(a gokwk.Snippets, f persist.IO, t persist.Persister, r cmd.Runner, u gokwk.Users,
	d Dialog, w vwrite.Writer, up update.Updater, eh errs.Handler) *KwkApp {
	out.SetColors(out.ColorsDefault())
	setup.NewConfig(a, f, u, eh)
	ap := cli.NewApp()
	ap = setupFlags(ap)
	ap.Version = CLIInfo.String()
	dash := NewDashBoard(w, eh, a)
	help := cli.HelpPrinter
	ap.Commands = append(ap.Commands, cli.Command{
		Name:    "help",
		Aliases: []string{"h"},
		Action: func(c *cli.Context) error {
			cli.HelpPrinter = help
			cli.ShowAppHelp(c)
			return nil
		},
	})
	cli.HelpPrinter = dash.GetWriter()
	accCli := NewAccount(u, t, w, d, dash)
	ap.Commands = append(ap.Commands, userRoutes(accCli)...)
	sysCli := NewSystem(w, up)
	ap.Commands = append(ap.Commands, systemRoutes(sysCli)...)
	snipCli := NewSnippet(a, r, d, w, t)
	ap.Commands = append(ap.Commands, snippetsRoutes(snipCli)...)
	ap.CommandNotFound = getDefaultCommand(snipCli)
	return &KwkApp{
		App:      ap,
		File:     f,
		Settings: t,
		Runner:   r,
		Acc:      u,
		Dialogue: d,
		Snippets: a,
		Writer:   w,
		Handler:  eh,
	}
}

func getDefaultCommand(snipCli *snippets) func(*cli.Context, string) {
	return func(c *cli.Context, firstArg string) {
		i := c.Args().Get(1)
		if strings.HasPrefix(firstArg, "@") {
			fmt.Println("listing:", firstArg)
			snipCli.GetEra(firstArg)
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
}

func setupFlags(ap *cli.App) *cli.App {
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
	return ap
}

func (a *KwkApp) Run(args ...string) {
	params := []string{"[app]"}
	params = append(params, args...)
	a.Handle(a.App.Run(params))
}
