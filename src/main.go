package main

import (
	gu "github.com/inconshreveable/go-update"
	"github.com/kwk-super-snippets/cli/src/app"
	"github.com/kwk-super-snippets/cli/src/app/out"
	"github.com/kwk-super-snippets/cli/src/store"
	"github.com/kwk-super-snippets/cli/src/updater"
	"github.com/kwk-super-snippets/types"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
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
	if hasUpdateFlag() {
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
	var fileOut = &lumberjack.Logger{
		Filename:   path.Join(out.KwkPath(), "update.log"),
		MaxSize:    3, // megabytes
		MaxBackups: 2,
		MaxAge:     5, //days})
	}
	eh := out.NewErrHandler(fileOut)
	f := store.NewDiskFile()
	jsn := store.NewJson(f, cfg.DocPath)
	up := updater.New(info.String(), &updater.S3Repo{}, gu.Apply, gu.RollbackError, jsn)
	eh.Handle(up.Run())
}

func hasUpdateFlag() bool {
	out.Debug("UPDATE MODE: %+v", os.Args)
	return len(os.Args) > 1 && os.Args[1] == updater.UpdateFlag
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
