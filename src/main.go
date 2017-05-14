package main

import (
	"github.com/kwk-super-snippets/cli/src/app"
	//"runtime/pprof"
	"bufio"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types"
	"github.com/kwk-super-snippets/types/errs"
	"github.com/kwk-super-snippets/types/vwrite"
	"os"
	"strconv"
	"strings"
)

var (
	KWK_TEST_MODE = false
	f             app.IO
	j             app.Persister
)

func main() {
	_, KWK_TEST_MODE = os.LookupEnv("KWK_TEST_MODE")
	args := strings.Join(os.Args[1:], "+")
	f = app.NewIO()
	j = app.NewJson(f, "settings")

	up := app.NewUpdateRunner(j, app.CLIInfo.String())
	if args == "update+silent" {
		up.Run()
	} else if args == "update" {
		runKwk(up)
	} else {
		app.SilentCheckAndRun()
		runKwk(up)
	}
}

var version string = "v-.-.-"
var build string = "0"
var releaseTime string

func runKwk(up app.Updater) {
	app.CLIInfo.Version = version
	app.CLIInfo.Build = build
	app.CLIInfo.Time, _ = strconv.ParseInt(releaseTime, 10, 64)

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
	conn, err := app.GetConn(host, KWK_TEST_MODE)
	if err != nil {
		eh.Handle(errs.ApiDown)
		return
	}
	defer conn.Close()
	r := bufio.NewReader(os.Stdin)
	d := app.NewDialog(w, r)
	sc := types.NewSnippetsClient(conn)
	uc := types.NewUsersClient(conn)
	o := app.NewRunner(f, sc)
	kwkApp := app.NewApp(sc, f, j, o, uc, d, w, up, eh)
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
