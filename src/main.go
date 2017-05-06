package main

import (
	"bitbucket.com/sharingmachine/kwkcli/src/app"
	//"runtime/pprof"
	"bitbucket.com/sharingmachine/kwkcli/src/app/out"
	"bitbucket.com/sharingmachine/kwkcli/src/exekwk/cmd"
	"bitbucket.com/sharingmachine/kwkcli/src/exekwk/update"
	"bitbucket.com/sharingmachine/kwkcli/src/gokwk"
	"bitbucket.com/sharingmachine/kwkcli/src/persist"
	"bitbucket.com/sharingmachine/types/errs"
	"bitbucket.com/sharingmachine/types/vwrite"
	"bufio"
	"os"
	"strconv"
	"strings"
)

var (
	KWK_TEST_MODE = false
	f             persist.IO
	j             persist.Persister
)

func main() {
	_, KWK_TEST_MODE = os.LookupEnv("KWK_TEST_MODE")
	args := strings.Join(os.Args[1:], "+")
	f = persist.New()
	j = persist.NewJson(f, "settings")

	up := update.NewRunner(j, app.CLIInfo.String())
	if args == "update+silent" {
		up.Run()
	} else if args == "update" {
		runKwk(up)
	} else {
		update.SilentCheckAndRun()
		runKwk(up)
	}
}

var version string = "v-.-.-"
var build string = "0"
var releaseTime string

func runKwk(up update.Updater) {
	app.CLIInfo.Version = version
	app.CLIInfo.Build = build
	app.CLIInfo.Time, _ = strconv.ParseInt(releaseTime, 10, 64)
	//profile().Close()

	host := os.Getenv("API_HOST")
	if host == "" {
		if KWK_TEST_MODE {
			host = "localhost:8000"
		} else {
			host = "api.kwk.co:443"
		}
	}
	eh := out.NewErrHandler(os.Stdout)
	w := vwrite.New(os.Stdout)
	conn, err := gokwk.GetConn(host, KWK_TEST_MODE)
	if err != nil {
		eh.Handle(errs.ApiDown)
		return
	}
	defer conn.Close()
	u := gokwk.NewUsers(conn, j, app.CLIInfo)
	ss := gokwk.New(conn, app.CLIInfo)
	o := cmd.NewStdRunner(f, ss)
	r := bufio.NewReader(os.Stdin)
	d := app.NewDialog(w, r)
	kwkApp := app.NewApp(ss, f, j, o, u, d, w, up, eh)
	eh.Handle(kwkApp.App.Run(os.Args))
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
