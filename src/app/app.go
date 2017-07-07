package app

import (
	"fmt"
	"github.com/rjarmstrong/kwk-types"
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/rjarmstrong/kwk/src/app/handlers"
	"github.com/rjarmstrong/kwk/src/app/routes"
	"github.com/rjarmstrong/kwk/src/cli"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/rpc"
	"github.com/rjarmstrong/kwk/src/runtime"
	"github.com/rjarmstrong/kwk/src/store"
	urf "github.com/urfave/cli"
	"io"
	"strings"
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
	d := out.NewDialog(w, r)
	f := store.NewDiskFile()
	doc := store.NewJson(f, cli.DocPath)
	srw := store.NewSnippetReadWriter(f)

	api, err := rpc.GetApi(principal, prefs, &i, cfg.APIHost, cfg.TestMode)
	if err != nil {
		out.LogErr(err)
		eh.Handle(errs.ApiDown)
		return nil
	}
	sc := types.NewSnippetsClient(api.ClientConn)
	uc := types.NewUsersClient(api.ClientConn)

	// SERVICES
	rp := rootPrinter(prefs, w, &principal.User)
	users := handlers.NewUsers(principal, uc, doc, w, d, api.Cxf, prefs, rp)
	runner := runtime.NewRunner(prefs, env, w, srw, useLogger(api.Cxf, sc))
	editor := runtime.NewEditor(env, prefs, snippetPatcher(api.Cxf, sc), srw)
	snippets := handlers.NewSnippets(prefs, sc, runner, editor, w, api.Cxf, rp, d)

	// RUNTIME
	users.LoadPrincipal(principal)
	runtime.Configure(env, prefs, principal.User.Username, snippetGetter(api.Cxf, sc), snippetMaker(api.Cxf, sc), srw, eh)
	out.Debug("PREFS: %+v", prefs)

	// LEVEL
	setProcessLevel()
	if node != nil && node.Level > 0 {
		eh.Handle(routes.RunNode(*principal, prefs, node, snippets))
		return nil
	}

	return &KwkCLI{
		uf:      createUrfaveApp(users, snippets, eh, rootGetter(api.Cxf, sc), w),
		Handler: eh,
	}
}

func (a *KwkCLI) Run(args ...string) {
	routes.ReplaceArg("@env", runtime.GetEnvURI())
	routes.ReplaceArg("@prefs", "prefs")
	a.Handle(a.uf.Run(args))
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
