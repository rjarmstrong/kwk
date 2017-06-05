package app

import (
	"fmt"
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/kwk-super-snippets/cli/src/cli"
	"github.com/kwk-super-snippets/cli/src/rpc"
	"github.com/kwk-super-snippets/cli/src/runtime"
	"github.com/kwk-super-snippets/cli/src/store"
	"github.com/kwk-super-snippets/cli/src/style"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/vwrite"
	urf "github.com/urfave/cli"
	"io"
	"os"
	"strings"
	"time"
)

var (
	node      *runtime.ProcessNode
	info      = types.AppInfo{}
	principal = &cli.UserWithToken{User: types.User{}}
	cfg       = &cli.AppConfig{}
	prefs     = runtime.DefaultPrefs()
	env       = runtime.DefaultEnv()
)

type KwkCLI struct {
	uf *urf.App
	errs.Handler
}

func NewCLI(r io.Reader, wr io.Writer, i types.AppInfo) *KwkCLI {
	info = i
	// IO
	out.SetColors(out.ColorsDefault())
	eh := out.NewErrHandler(wr)
	w := vwrite.New(wr)
	d := NewDialog(w, r)
	f := store.NewDiskFile()
	jsn := store.NewJson(f, cfg.DocPath)
	srw := NewSnippetReadWriter(f)

	// API
	rpc, err := rpc.GetRpc(principal, &i, cfg.APIHost, cfg.TestMode)
	if err != nil {
		eh.Handle(errs.ApiDown)
		return nil
	}
	sc := types.NewSnippetsClient(rpc.ClientConn)
	uc := types.NewUsersClient(rpc.ClientConn)

	// SERVICES
	dash := NewDashBoard(w, eh, rootGetter(rpc.Cxf, sc))
	users := NewUsers(uc, jsn, w, d, dash, rpc.Cxf)
	runner := runtime.NewRunner(prefs, env, w, srw, useLogger(rpc.Cxf, sc))
	editor := runtime.NewEditor(env, prefs, snippetPatcher(rpc.Cxf, sc), srw)
	snippets := NewSnippets(sc, runner, editor, d, w, rpc.Cxf)

	// RUNTIME
	jsn.Get(cfg.UserDocName, principal, 0)
	runtime.Configure(env, prefs, principal.User.Username, snippetGetter(rpc.Cxf, sc), snippetMaker(rpc.Cxf, sc), srw, eh)
	out.Debug("PREFS: %+v", prefs)
	setProcessLevel()
	if node != nil && node.Level > 0 {
		eh.Handle(nodeRun(snippets))
		return nil
	}

	// APP
	ap := urf.NewApp()
	ap.Name = style.Fmt256(style.ColorPouchCyan, style.IconSnippet) + "  kwk super snippets"
	ap.Description = "A smart & friendly snippet manager for the CLI"
	ap.Usage = ""
	ap.UsageText = "kwk [global options] command [command options] [arguments...]"
	ap.EnableBashCompletion = true
	ap.Authors = []urf.Author{
		{
			Name:  "R J Armstrong",
			Email: "rj@kwk.co",
		},
	}
	ap.Copyright = "© 2017 Gimanzo Systems Ltd \n"

	ap = setupFlags(ap)
	ap.Version = fmt.Sprintf("\n\n%s Version : %s\n%s Revision: %s\n%s Released: %s\n",
		style.Margin, i.Version, style.Margin, i.Build, style.Margin, time.Unix(i.Time, 0).Format(time.RFC822))
	help := urf.HelpPrinter
	ap.Commands = append(ap.Commands, urf.Command{
		Name:    "help",
		Aliases: []string{"h", "?"},
		Action: func(c *urf.Context) error {
			urf.HelpPrinter = help
			urf.ShowAppHelp(c)
			return nil
		},
	})
	ap.Commands = append(ap.Commands, userRoutes(users)...)
	ap.Commands = append(ap.Commands, snippetsRoutes(snippets)...)
	ap.Commands = append(ap.Commands, pouchRoutes(snippets)...)
	ap.CommandNotFound = getDefaultCommand(snippets, eh)
	urf.HelpPrinter = dash.GetWriter()

	return &KwkCLI{
		uf:      ap,
		Handler: eh,
	}
}

func nodeRun(snippets *snippets) error {
	if len(os.Args) < 3 {
		return errs.New(errs.CodeInvalidArgument, "Invalid kwk call '%+v' in app.\n Invoke snippets as follows: kwk run <uri>", strings.Join(os.Args, " "))
	}
	if os.Args[1] != "run" {
		return errs.New(errs.CodeInvalidArgument, "'run' keyword required as first arg within an app.")
	}

	return snippets.NodeRun(os.Args[2], os.Args[3:])
}

func setProcessLevel() {
	n, err := runtime.GetCallerNode()
	if err != nil {
		out.Debug("NODE:", err)
	}
	node = n
	if node != nil {
		out.DebugLogger.SetPrefix(fmt.Sprintf("%s%d-KWK: ", strings.Repeat("--", int(node.Level)), node.Level))
	}
}

func (a *KwkCLI) Run(args ...string) {
	a.Handle(a.uf.Run(args))
}

func snippetPatcher(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.SnippetPatcher {
	return func(req *types.PatchRequest) (*types.PatchResponse, error) {
		return sc.Patch(cxf(), req)
	}
}

func rootGetter(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.RootGetter {
	return func(req *types.RootRequest) (*types.RootResponse, error) {
		return sc.GetRoot(cxf(), req)
	}
}

func useLogger(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.UseLogger {
	return func(req *types.UseContext) (*types.LogUseResponse, error) {
		return sc.LogUse(cxf(), req)
	}
}

func snippetGetter(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.SnippetGetter {
	return func(req *types.GetRequest) (*types.ListResponse, error) {
		return sc.Get(cxf(), req)
	}
}

func snippetMaker(cxf cli.ContextFunc, sc types.SnippetsClient) runtime.SnippetMaker {
	return func(req *types.CreateRequest) error {
		_, err := sc.Create(cxf(), req)
		return err
	}
}

func getDefaultCommand(snipCli *snippets, eh errs.Handler) func(*urf.Context, string) {
	return func(c *urf.Context, firstArg string) {
		i := c.Args().Get(1)
		if strings.HasPrefix(firstArg, "@") {
			fmt.Println("listing:", firstArg)
			snipCli.GetEra(firstArg)
			return
		}
		var err error
		switch i {
		case "version":
			fmt.Println(info.String())
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

func setupFlags(ap *urf.App) *urf.App {
	ap.Flags = []urf.Flag{
		urf.BoolFlag{
			Name:        "covert, x",
			Usage:       "Open browser in covert mode.",
			Destination: &prefs.Covert,
		},
		urf.BoolFlag{
			Name:        "naked, n",
			Usage:       "List without styles",
			Destination: &prefs.Naked,
		},
		urf.BoolFlag{
			Name:        "ansi",
			Usage:       "Prints ansi escape sequences for debugging purposes",
			Destination: &style.PrintAnsi,
		},
		urf.BoolFlag{
			Name:        "quiet, q",
			Usage:       "List names only",
			Destination: &prefs.Quiet,
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
