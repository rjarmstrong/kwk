package main

import (
	"bitbucket.com/sharingmachine/kwkcli/snippets"
	"bitbucket.com/sharingmachine/kwkcli/ui/tmpl"
	"bitbucket.com/sharingmachine/kwkcli/ui/dlg"
	"bitbucket.com/sharingmachine/kwkcli/setup"
	"bitbucket.com/sharingmachine/kwkcli/app"
	"bitbucket.com/sharingmachine/kwkcli/rpc"
	"bitbucket.com/sharingmachine/kwkcli/cmd"
	"bitbucket.com/sharingmachine/kwkcli/models"
	"bitbucket.com/sharingmachine/kwkcli/update"
	//"runtime/pprof"
	"bufio"
	"os"
	"strings"
	"strconv"
	"bitbucket.com/sharingmachine/kwkcli/user"
	"bitbucket.com/sharingmachine/kwkcli/persist"
)

var (
	KWK_TEST_MODE = false
	f persist.IO
	j persist.Persister
)

func main() {
	_, KWK_TEST_MODE = os.LookupEnv("KWK_TEST_MODE")
	args := strings.Join(os.Args[1:], "+")
	f = persist.New()
	j = persist.NewJson(f, "settings")

	if args == "update+silent" {
		update.NewRunner(j).Run()
	} else if args == "update" {
		runKwk()
	}else {
		update.SilentCheckAndRun()
		runKwk()
	}
}

var version string = "v-.-.-"
var build string = "0"
var releaseTime string

func runKwk() {
	models.Client.Version = version
	models.Client.Build = build
	models.Client.Time, _ = strconv.ParseInt(releaseTime, 10, 64)
	//profile().Close()

	host := os.Getenv("API_HOST")
	if host == "" {
		if KWK_TEST_MODE {
			host = "localhost:8000"
		} else {
			host = "api.kwk.co:443"
		}
	}
	w := tmpl.NewWriter(os.Stdout)
	conn, err := rpc.GetConn(host, KWK_TEST_MODE)
	if err != nil {
		w.HandleErr(models.ErrOneLine(models.Code_ApiDown, " The kwk api is down, please try again."))
		return
	}
	defer conn.Close()
	h := rpc.NewHeaders()
	a := user.NewAccount(conn, j, h)
	ss := snippets.New(conn, h)
	conf := setup.NewConfigProvider(ss, f, a, w)
	conf.Load()
	o := cmd.NewStdRunner(f, ss)
	r := bufio.NewReader(os.Stdin)
	d := dlg.New(w, r)
	api := rpc.New(conn, h)
	kwkApp := app.NewApp(ss, f, j, o, a, d, w, api)
	kwkApp.App.Run(os.Args)
}

//func profile() *os.File {
//	var cpuprofile = "kwkprofile"
//	f, err := os.Create(cpuprofile)
//	if err != nil {
//		panic(err)
//	}
//	if err := pprof.StartCPUProfile(f); err != nil {
//		panic(err)
//	}
//	return f
//}
