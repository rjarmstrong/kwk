package log

import (
	lg "log"
	"os"
	"log"
	"fmt"
	"path"
	"bitbucket.com/sharingmachine/kwkcli/cache"
	"gopkg.in/natefinch/lumberjack.v2"
)

var er = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
var dbg = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

// This flag is specified here since we don't wan't to have to depend on urfave global flags.
// That flag will also work throughout the app.

var EnableDebug bool

var fileOut = &lumberjack.Logger{
	Filename:   path.Join(cache.Path(), "kwk.log"),
	MaxSize:    3, // megabytes
	MaxBackups: 2,
	MaxAge:     5, //days})
}

func setOutput() {
	if EnableDebug {
		dbg.SetOutput(os.Stdout)
		er.SetOutput(os.Stdout)
		return
	}
	dbg.SetOutput(fileOut)
	er.SetOutput(fileOut)
}

func GetDebugLog() *lg.Logger {
	return dbg
}

func Debug(f string, args ...interface{}) {
	setOutput()
	dbg.Output(2, fmt.Sprintf(f, args...))
}

func Error(desc string, err error) {
	setOutput()
	er.Output(2, fmt.Sprintf("%s: %+v", desc, err))
}
