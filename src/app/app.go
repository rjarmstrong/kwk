package app

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/vwrite"
	"github.com/urfave/cli"
	"strings"
)

var (
	CLIInfo   = types.AppInfo{}
	principal = &UserWithToken{}
)

const profileFileName = "profile.json"

type KwkApp struct {
	App      *cli.App
	Users    types.UsersClient
	Snippets types.SnippetsClient
	File     IO
	Settings DocStore
	Updater  *Updater
	Runner   Runner
	Dialogue Dialog
	vwrite.Writer
	errs.Handler
}

func NewApp(a types.SnippetsClient, f IO, docs DocStore, r Runner, u types.UsersClient,
	d Dialog, w vwrite.Writer, up Updater, eh errs.Handler) *KwkApp {
	out.SetColors(out.ColorsDefault())
	docs.Get(profileFileName, principal, 0)
	NewConfig(a, f, u, eh)
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
	accCli := NewUsers(u, docs, w, d, dash)
	ap.Commands = append(ap.Commands, userRoutes(accCli)...)
	sysCli := NewSystem(w, up)
	ap.Commands = append(ap.Commands, systemRoutes(sysCli)...)
	snipCli := NewSnippets(a, r, d, w, docs)
	ap.Commands = append(ap.Commands, snippetsRoutes(snipCli)...)
	ap.CommandNotFound = getDefaultCommand(snipCli)
	return &KwkApp{
		App:      ap,
		File:     f,
		Settings: docs,
		Runner:   r,
		Users:    u,
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
			Destination: &Prefs().Covert,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Debug.",
			Destination: &out.DebugEnabled,
			EnvVar:      "DEBUG",
		},
		cli.BoolFlag{
			Name:        "naked, n",
			Usage:       "list without styles",
			Destination: &Prefs().Naked,
		},
		cli.BoolFlag{
			Name:        "ansi",
			Usage:       "Prints ansi escape sequences for debugging purposes",
			Destination: &style.PrintAnsi,
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

//func GetPrincipal() (*types.User, error) {
//	if principal != nil && principal.Id != "" {
//		return principal, nil
//	}
//	principal = &types.User{}
//	if err := u.settings.Get(profileFileName, principal, 0); err != nil {
//		return nil, err
//	} else {
//		return principal, nil
//	}
//	u, err := c.client.
//	if err != nil {
//		return err
//	}
//	return c.EWrite(out.UserProfile(u))
//}
