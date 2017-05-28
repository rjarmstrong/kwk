package main

import (
	"github.com/kwk-super-snippets/cli/src/app"
	//"runtime/pprof"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/types"
	"os"
	"github.com/kelseyhightower/envconfig"
	"encoding/json"
	"log"
	"strconv"
	"runtime/pprof"
)


var (
	version     string = "v-.-.-"
	build       string = "0"
	releaseTime string
)

func main() {
	initVariables()
	if app.Config.Profile {
		defer profile().Close()
	}
	// Init binary info.
	var info = types.AppInfo {
		Version: version,
		Build:   build,
	}
	info.Time, _ = strconv.ParseInt(releaseTime, 10, 64)

	//  Updater
	//args := strings.Join(os.Args[1:], "+")
	jsn := app.NewJson(app.NewIO(), "settings")
	up := app.NewUpdateRunner(jsn, info.String())
	//if args == "update+silent" {
	//	up.Run()
	//	return
	//}
	//if args != "update" {
	//	app.SilentCheckAndRun()
	//}

	/// The app
	eh := out.NewErrHandler(os.Stdout)
	cli := app.NewCLI(os.Stdin, os.Stdout, info, up, eh)
	if cli == nil {
		return
	}
	eh.Handle(cli.App.Run(os.Args))
}

func initVariables() {
	err := envconfig.Process("KWK", &app.Config)
	if err != nil {
		log.Fatal(err.Error())
	}
	out.DebugEnabled = app.Config.Debug
	if app.Config.TestMode {
		app.Config.APIHost = "localhost:8000"
	}
	b, _ := json.MarshalIndent(app.Config, "", "  ")
	out.Debug("CONFIG: %s", string(b))
}

func profile() *os.File {
	var prof = "kwk_profile"
	f, err := os.Create(prof)
	if err != nil {
		panic(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	return f
}
