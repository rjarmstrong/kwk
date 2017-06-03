package main

import (
	gu "github.com/inconshreveable/go-update"
	"github.com/kwk-super-snippets/cli/src/app"
	"github.com/kwk-super-snippets/cli/src/out"
	"github.com/kwk-super-snippets/cli/src/store"
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
	info := getAppInfo()

	// If update argument supplied then we don't want to run the app
	// rather actually run the update.
	if isUpdateMode() {
		runUpdate(cfg, info)
		return
	}

	cli := app.NewCLI(os.Stdin, os.Stdout, info)
	if cli == nil {
		return
	}

	cli.Run(os.Args...)
	// If update argument not supplied then run update in a new process.
	updater.SpawnUpdate()
}

func runUpdate(cfg *app.CLIConfig, info types.AppInfo) {
	out.DebugEnabled = false
	f := store.NewDiskFile()
	jsn := store.NewJson(f, cfg.DocPath)
	up := updater.New(info.String(), &updater.S3Repo{}, gu.Apply, gu.RollbackError, jsn)
	err := up.Run()
	if err != nil {
		out.LogErrM("Update exited with err:", err)
	}
}

func isUpdateMode() bool {
	ok := len(os.Args) == 2 && os.Args[1] == updater.UpdateFlag
	if ok {
		out.DebugLogger.SetPrefix("KWK:UM: ")
	}
	return ok
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
