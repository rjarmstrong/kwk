package app

import (
	"github.com/rjarmstrong/kwk-types/errs"
	"github.com/rjarmstrong/kwk-types/vwrite"
	"github.com/rjarmstrong/kwk/src/app/handlers"
	"github.com/rjarmstrong/kwk/src/app/routes"
	"github.com/rjarmstrong/kwk/src/out"
	"github.com/rjarmstrong/kwk/src/runtime"
	urf "github.com/urfave/cli"
)

func createUrfaveApp(u *handlers.Users, s *handlers.Snippets, eh errs.Handler, rg runtime.RootGetter, w vwrite.Writer) *urf.App {
	uf := urf.NewApp()
	uf.Name = out.AppName
	uf.Description = out.AppDescription
	uf.Usage = ""
	uf.UsageText = "kwk [global options] command [command options] [arguments...]"
	uf.EnableBashCompletion = true
	uf.Authors = []urf.Author{{Name: "R J Armstrong", Email: "rj@kwk.co"}}
	uf.Copyright = "© 2017 Gimanzo Systems Ltd \n"
	uf.Version = out.Version(info)

	help := urf.HelpPrinter
	uf.Commands = append(uf.Commands, urf.Command{
		Name:    "help",
		Aliases: []string{"h", "?"},
		Action: func(c *urf.Context) error {
			urf.HelpPrinter = help
			urf.ShowAppHelp(c)
			w.Write(out.FreeText("SHORTCUTS:\n"))
			w.Write(out.FreeText("     [edit|cat|run] @env    Alias for settings/<local environment>.yml\n"))
			w.Write(out.FreeText("     [edit|cat|run] @prefs  Alias for settings/prefs.yml\n"))
			w.Write(out.FreeText("\n"))
			w.Write(out.FreeText("     @today                 List snippets used today\n"))
			w.Write(out.FreeText("     @week                  List snippets used this week\n"))
			w.Write(out.FreeText("\n"))
			w.Write(out.FreeText("\n"))
			return nil
		},
	})
	uf.Commands = append(uf.Commands, routes.Users(u)...)
	uf.Commands = append(uf.Commands, routes.Snippets(*principal, prefs, s)...)
	uf.Commands = append(uf.Commands, routes.Pouches(prefs, s)...)
	uf.CommandNotFound = routes.DefaultRoute(&info, s, eh)
	routes.SetupFlags(prefs, uf)
	urf.HelpPrinter = NewDashBoard(w, eh, rg).GetWriter()
	return uf
}
