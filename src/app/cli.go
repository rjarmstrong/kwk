package app

import (
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
	"github.com/kwk-super-snippets/cli/src/app/runtime"
)

var (
	cliInfo   = types.AppInfo{}
	principal = &UserWithToken{}
	cfg       = &CLIConfig{}
	prefs     = runtime.DefaultPrefs()
	env       = runtime.DefaultEnv()
)

type KwkCLI struct {
	app *cli.App
	errs.Handler
}

func NewCLI(r io.Reader, wr io.Writer, info types.AppInfo) *KwkCLI {
	cliInfo = info

	// IO
	out.SetColors(out.ColorsDefault())
	eh := out.NewErrHandler(wr)
	w := vwrite.New(wr)
	d := NewDialog(w, r)
	f := store.NewDiskFile()
	jsn := store.NewJson(f, cfg.DocPath)

	// API
	conn, err := GetConn(cfg.APIHost, cfg.TestMode)
	if err != nil {
		eh.Handle(errs.ApiDown)
		return nil
	}
	sc := types.NewSnippetsClient(conn)
	uc := types.NewUsersClient(conn)
	runner := NewRunner(f, sc, cfg.SnippetPath)

	// SERVICES
	dash := NewDashBoard(w, eh, sc)
	users := NewUsers(uc, jsn, w, d, dash)
	snippets := NewSnippets(sc, runner, d, w)
	system := NewSystem(w, updater.New(info.String(), &updater.S3Repo{}, gu.Apply, gu.RollbackError, jsn))

	// APP
	jsn.Get(cfg.UserDocName, principal, 0)
	runtime.Configure(env, prefs, principal.HasAccessToken(), snippetGetter(sc), snippetMaker(sc), cfg.SnippetPath, f, eh)
	out.Debug("PREFS: %+v", prefs)
	ap := cli.NewApp()
	ap.Name = "kwk super snippets"
	ap.Description = "A super snippet manager for the CLI"
	ap.Usage = "kwk COMMAND"
	ap.UsageText = "kwk [global options] command [command options] [arguments...]"
	ap = setupFlags(ap)
	ap.Version = cliInfo.String()
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
	ap.Commands = append(ap.Commands, userRoutes(users)...)
	ap.Commands = append(ap.Commands, snippetsRoutes(snippets)...)
	ap.Commands = append(ap.Commands, pouchRoutes(snippets)...)
	ap.Commands = append(ap.Commands, systemRoutes(system)...)
	ap.CommandNotFound = getDefaultCommand(snippets)
	cli.HelpPrinter = dash.GetWriter()

	return &KwkCLI{
		app:     ap,
		Handler: eh,
	}
}

func (a *KwkCLI) Run(args ...string) {
	a.Handle(a.app.Run(args))
}

func snippetGetter(sc types.SnippetsClient) runtime.SnippetGetter {
	return func(req *types.GetRequest) (*types.ListResponse, error) {
		return sc.Get(Ctx(), req)
	}
}

func snippetMaker(sc types.SnippetsClient) runtime.SnippetMaker {
	return func(req *types.CreateRequest) error {
		_, err := sc.Create(Ctx(), req)
		return err
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
			Destination: &prefs.Covert,
		},
		cli.BoolFlag{
			Name:        "naked, n",
			Usage:       "list without styles",
			Destination: &prefs.Naked,
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
