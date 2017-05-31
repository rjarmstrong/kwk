package app

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/cli/src/app/runtime"
	"github.com/kwk-super-snippets/cli/src/store"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/vwrite"
	"github.com/urfave/cli"
	"io"
	"strings"
)

var (
	cliInfo   = types.AppInfo{}
	principal = &UserWithToken{User: types.User{}}
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
	srw := NewSnippetReadWriter(f)

	// API
	conn, err := GetConn(cfg.APIHost, cfg.TestMode)
	if err != nil {
		eh.Handle(errs.ApiDown)
		return nil
	}
	sc := types.NewSnippetsClient(conn)
	uc := types.NewUsersClient(conn)

	// SERVICES
	dash := NewDashBoard(w, eh, sc)
	users := NewUsers(uc, jsn, w, d, dash)
	runner := NewRunner(srw, sc)
	setProcessLevel(err, eh)
	snippets := NewSnippets(sc, runner, d, w)

	// APP
	jsn.Get(cfg.UserDocName, principal, 0)
	runtime.Configure(env, prefs, principal.User.Username, snippetGetter(sc), snippetMaker(sc), srw, eh)
	out.Debug("PREFS: %+v", prefs)
	ap := cli.NewApp()

	ap.Name = style.Fmt256(style.ColorPouchCyan, style.IconSnippet) + "  kwk super snippets"
	ap.Description = "A smart & friendly snippet manager for the CLI"
	ap.Usage = ""
	ap.UsageText = "kwk [global options] command [command options] [arguments...]"
	ap.EnableBashCompletion = true
	ap.Authors = []cli.Author{
		{
			Name:  "Richard J Armstrong",
			Email: "richard.armstrong@gimanzo.com",
		},
	}
	ap.Copyright = "© 2017 Gimanzo Systems Ltd \n"

	ap = setupFlags(ap)
	ap.Version = cliInfo.String()
	help := cli.HelpPrinter
	ap.Commands = append(ap.Commands, cli.Command{
		Name:    "help",
		Aliases: []string{"h", "?"},
		Action: func(c *cli.Context) error {
			cli.HelpPrinter = help
			cli.ShowAppHelp(c)
			return nil
		},
	})
	ap.Commands = append(ap.Commands, userRoutes(users)...)
	ap.Commands = append(ap.Commands, snippetsRoutes(snippets)...)
	ap.Commands = append(ap.Commands, pouchRoutes(snippets)...)
	ap.CommandNotFound = getDefaultCommand(snippets, eh)
	cli.HelpPrinter = dash.GetWriter()

	return &KwkCLI{
		app:     ap,
		Handler: eh,
	}
}
func setProcessLevel(err error, eh errs.Handler) {
	node, err := GetCallerNode()
	if err != nil {
		out.Debug("NODE:", err)
	}
	if node != nil {
		out.DebugLogger.SetPrefix(fmt.Sprintf("%d%sKWK: ", node.Level, strings.Repeat("--", int(node.Level))))
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

func getDefaultCommand(snipCli *snippets, eh errs.Handler) func(*cli.Context, string) {
	return func(c *cli.Context, firstArg string) {
		i := c.Args().Get(1)
		if strings.HasPrefix(firstArg, "@") {
			fmt.Println("listing:", firstArg)
			snipCli.GetEra(firstArg)
			return
		}
		var err error
		switch i {
		case "run":
			err = snipCli.Run(c.Args().First(), []string(c.Args())[2:])
		case "r":
			err = snipCli.Run(c.Args().First(), []string(c.Args())[2:])
		case "edit":
			err = snipCli.Edit(c.Args().First())
		case "e":
			err = snipCli.Edit(c.Args().First())
		default:
			err = snipCli.InspectListOrRun(c.Args().First(), false, []string(c.Args())[1:]...)
		}
		if err != nil {
			eh.Handle(err)
		}
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
