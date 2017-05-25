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
	"os"
	"bufio"
)

var (
	cliInfo   = types.AppInfo{}
	principal = &UserWithToken{}
	Config    CLIConfig
)

const profileFileName = "profile.json"

type KwkCLI struct {
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

func NewCLI(info types.AppInfo, up Updater, eh errs.Handler) *KwkCLI {

	cliInfo = info

	conn, err := GetConn(Config.APIHost, Config.TestMode)
	if err != nil {
		eh.Handle(errs.ApiDown)
		return nil
	}

	w := vwrite.New(os.Stdout)
	defer conn.Close()
	r := bufio.NewReader(os.Stdin)
	d := NewDialog(w, r)
	sc := types.NewSnippetsClient(conn)
	uc := types.NewUsersClient(conn)
	f := NewIO()
	jsn := NewJson(f, "settings")
	run := NewRunner(f, sc)

	out.SetColors(out.ColorsDefault())

	jsn.Get(profileFileName, principal, 0)

	InitConfig(sc, f, uc, eh)

	ap := cli.NewApp()
	ap = setupFlags(ap)
	ap.Version = cliInfo.String()
	dash := NewDashBoard(w, eh, sc)
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
	accCli := NewUsers(uc, jsn, w, d, dash)
	ap.Commands = append(ap.Commands, userRoutes(accCli)...)
	sysCli := NewSystem(w, up)
	ap.Commands = append(ap.Commands, systemRoutes(sysCli)...)
	snipCli := NewSnippets(sc, run, d, w, jsn)
	ap.Commands = append(ap.Commands, snippetsRoutes(snipCli)...)
	ap.CommandNotFound = getDefaultCommand(snipCli)
	cli.HelpPrinter = dash.GetWriter()
	return &KwkCLI{
		App:      ap,
		File:     f,
		Settings: jsn,
		Runner:   run,
		Users:    uc,
		Dialogue: d,
		Snippets: sc,
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

func (a *KwkCLI) Run(args ...string) {
	params := []string{"[app]"}
	params = append(params, args...)
	a.Handle(a.App.Run(params))
}
