package app

import (
	"bitbucket.com/sharingmachine/kwkcli/search"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"github.com/urfave/cli"
	"bitbucket.com/sharingmachine/kwkcli/setup"
)

type KwkApp struct {
	App            *cli.App
	Snippets       snippets.Service
	System         sys.Manager
	Settings       config.Persister
	AccountManage  account.Manager
	Runner         cmd.Runner
	Dialogue       dlg.Dialog
	TemplateWriter tmpl.Writer
	Search         search.Term
	Api 	       rpc.Service
}

func New(a snippets.Service, s sys.Manager, t config.Persister, r cmd.Runner, u account.Manager,
	d dlg.Dialog, w tmpl.Writer, h search.Term, api rpc.Service, su setup.Provider) *KwkApp {

	app := cli.NewApp()
	dash := NewDashBoard(w, a)
	cli.HelpPrinter = dash.GetWriter()
	app.Flags = []cli.Flag {
		cli.BoolFlag{
			Name: "covert, x",
			Usage: "Open browser in covert mode.",
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


	accCli := NewAccountCli(u, t, w, d)
	app.Commands = append(app.Commands, Accounts(accCli)...)

	sysCli := NewSystemCli(s, api, u, w, t)
	app.Commands = append(app.Commands, System(sysCli)...)

	snipCli := NewSnippetCli(a, r, s, d, w, t, h, su)
	app.Commands = append(app.Commands, Snippets(snipCli)...)
	app.CommandNotFound = func(c *cli.Context, fullKey string) {
		covert := c.Bool("covert")
		if covert {
			su.Prefs().Covert = true
		}
		snipCli.Run(fullKey, []string(c.Args())[1:])
	}
	searchCli := NewSearchCli(h, w, d)
	app.Commands = append(app.Commands, Search(searchCli)...)

	return &KwkApp{
		App: app,
		System: s,
		Settings: t,
		Runner: r,
		AccountManage: u,
		Dialogue: d,
		Snippets: a,
		TemplateWriter: w,
		Search:h,
		Api:api,
	}
}

func (a *KwkApp) Run(args ...string) {
	params := []string{"[app]"}
	params = append(params, args...)
	a.App.Run(params)
}