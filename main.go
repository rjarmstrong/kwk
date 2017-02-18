package main

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/search"
	"bitbucket.com/sharingmachine/kwkcli/config"
	"bitbucket.com/sharingmachine/kwkcli/account"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/setup"
	"bitbucket.com/sharingmachine/kwkcli/app"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/update"
	"bitbucket.com/sharingmachine/kwkcli/sys"
	"runtime/pprof"
	"bufio"
	"os"
	"strings"
)

var version string = "-"

var build string = "-"

var s sys.Manager
var j config.Persister

func main() {
	_, sys.KWK_TEST_MODE = os.LookupEnv("KWK_TEST_MODE")
	args := strings.Join(os.Args[1:], "+")
	s = sys.New()
	j = config.NewJsonSettings(s, "settings")

	if args == "update+silent" {
		update.NewRunner(j).Run()
	} else if args == "update" {
		runKwk()
	}else {
		update.SilentCheckAndRun()
		runKwk()
	}
}

func runKwk() {
	sys.Version = version
	sys.Build = build

	f, l := sys.NewLogger()
	defer f.Close()
	//profile().Close()

	host := os.Getenv("API_HOST")
	if host == "" {
		if sys.KWK_TEST_MODE {
			host = "localhost:8000"
		} else {
			host = "api.kwk.co:443"
		}
	}
	w := tmpl.NewWriter(os.Stdout)
	conn, err := rpc.GetConn(host, l, sys.KWK_TEST_MODE)
	if err != nil {
		l.Println(err)
		w.HandleErr(models.ErrOneLine(models.Code_ApiDown, " The kwk api is down, please try again."))
		return
	}
	defer conn.Close()

	v := version + "+" + build
	h := rpc.NewHeaders(j, v)
	u := account.NewStdManager(conn, j, h)
	ss := snippets.New(conn, j, h)

	su := setup.NewConfigProvider(ss, s, u, w)
	o := cmd.NewStdRunner(s, ss, su)
	r := bufio.NewReader(os.Stdin)
	d := dlg.New(w, r)
	ch := search.NewAlphaTerm(conn, j, h)
	api := rpc.New(conn, h)

	kwkApp := app.New(ss, s, j, o, u, d, w, ch, api, su)
	kwkApp.App.Version = v

	su.Preload()
	kwkApp.App.Run(os.Args)
}

func profile() *os.File {
	var cpuprofile = "kwkprofile"
	f, err := os.Create(cpuprofile)
	if err != nil {
		panic(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	return f
}
