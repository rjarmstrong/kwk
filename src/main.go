package main

import (
	"github.com/kwk-super-snippets/cli/src/app"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/cli/src/updater"
	"github.com/kwk-super-snippets/types"
	"os"
	"runtime/pprof"
	"strconv"
)

var (
	version     string = "v-.-.-"
	build       string = "0"
	releaseTime string
)

func main() {
	cfg := app.GetConfig()
	if cfg.CpuProfile {
		defer runCpuProfile().Close()
	}
	update()
	eh := out.NewErrHandler(os.Stdout)
	info := getAppInfo()
	cli := app.NewCLI(os.Stdin, os.Stdout, info, eh)
	if cli == nil {
		return
	}
	eh.Handle(cli.App.Run(os.Args))
}

func update() {
	// If update argument supplied then we don't want to spawn an update,
	// rather actually run the update in this process.
	if len(os.Args) > 1 && os.Args[1] == "update" {
		return
	}
	// If update argument not supplied then run update in a new process.
	updater.SpawnUpdate()

}

func getAppInfo() types.AppInfo {
	var info = types.AppInfo{
		Version: version,
		Build:   build,
	}
	info.Time, _ = strconv.ParseInt(releaseTime, 10, 64)
	return info
}

func runCpuProfile() *os.File {
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
