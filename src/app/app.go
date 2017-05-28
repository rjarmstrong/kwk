package app

import (
	"bufio"
	"fmt"
	gu "github.com/inconshreveable/go-update"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/vwrite"
	"github.com/urfave/cli"
	"io"
	"strings"
	"github.com/kwk-super-snippets/cli/src/updater"
	"github.com/kwk-super-snippets/cli/src/store"
)

var (
	cliInfo   = types.AppInfo{}
	principal = &UserWithToken{}
	cfg       = &CLIConfig{}
)

type KwkCLI struct {
	App      *cli.App
	Users    types.UsersClient
	Snippets types.SnippetsClient
	store.File
	store.Doc
	Updater  *updater.Runner
	Runner   Runner
	Dialogue Dialog
	vwrite.Writer
	errs.Handler
}

func NewCLI(rd io.Reader, wr io.Writer, info types.AppInfo, eh errs.Handler) *KwkCLI {
	cliInfo = info
	w := vwrite.New(wr)
	r := bufio.NewReader(rd)
	d := NewDialog(w, r)

	conn, err := GetConn(cfg.APIHost, cfg.TestMode)
	if err != nil {
		eh.Handle(errs.ApiDown)
		return nil
	}
	sc := types.NewSnippetsClient(conn)
	uc := types.NewUsersClient(conn)
	f := store.NewDiskFile()
	jsn := store.NewJson(f, cfg.DocPath)
	jsn.Get(cfg.UserDocName, principal, 0)
	runner := NewRunner(f, sc)

	out.SetColors(out.ColorsDefault())

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
	sysCli := NewSystem(w, updater.New(info.String(), &updater.S3Repo{}, gu.Apply, gu.RollbackError, jsn))
	ap.Commands = append(ap.Commands, systemRoutes(sysCli)...)
	snipCli := NewSnippets(sc, runner, d, w)
	ap.Commands = append(ap.Commands, snippetsRoutes(snipCli)...)
	ap.CommandNotFound = getDefaultCommand(snipCli)
	cli.HelpPrinter = dash.GetWriter()
	return &KwkCLI{
		App:      ap,
		File:     f,
		Doc:      jsn,
		Runner:   runner,
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
