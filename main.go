package main

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/search"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/app"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"runtime/pprof"
	"bufio"
	"os"
)

var version string = "-"

var build string = "-"

func main() {
	f, l := sys.NewLogger()
	defer f.Close()
	//profile().Close()

	host := os.Getenv("API_HOST")
	if host == "" {
		host = "api.kwk.co:443"
	}
	conn := rpc.GetConn(host, l)
	defer conn.Close()

	s := sys.New()
	t := config.NewFileSettings(s, "settings")
	h := rpc.NewHeaders(t)
	u := account.NewStdManager(conn, t, h)
	a := snippets.New(conn, t, h)
	w := tmpl.NewWriter(os.Stdout)
	o := cmd.NewStdRunner(s, a, w, u, t)
	r := bufio.NewReader(os.Stdin)
	d := dlg.New(w, r)
	ch := search.NewAlphaTerm(conn, t, h)
	api := rpc.New(conn, h)

	kwkApp := app.New(a, s, t, o, u, d, w, ch, api)
	kwkApp.App.Version = version + "+" + build

	configure(o, t)
	kwkApp.App.Run(os.Args)
}


func configure(r *cmd.StdRunner, t config.Settings) {
	t.SetPersistedPrefs(r.LoadPreferences())
	t.GetEnv()
}

func profile() *os.File  {
	var cpuprofile = "kwkprofile"
	f, err := os.Create(cpuprofile)
	if err != nil {
		panic( err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	return f
}
