package log

import (
	lg "log"
	"os"
	"io"
	"log"
	"fmt"
	"path"
	"bitbucket.com/sharingmachine/kwkcli/cache"
	"gopkg.in/natefinch/lumberjack.v2"
)

var er = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
var dbg = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

// This flag is specified here since we don't wan't to have to depend on urfave global flags.
// That flag will also work throughout the app.

var EnableDebug bool

func init() {
	output := &lumberjack.Logger{
		Filename:   path.Join(cache.Path(), "kwk.log"),
		MaxSize:    3, // megabytes
		MaxBackups: 2,
		MaxAge:     5, //days
	}
	if !EnableDebug {
		setOutput(output)
	} else {
		Debug("--debug")
	}
	er.SetOutput(output)
}

func setOutput(out io.Writer) {
	dbg.SetOutput(out)
	er.SetOutput(out)
}

func GetDebugLog() *lg.Logger {
	return dbg
}

func Debug(f string, args ...interface{}){
	dbg.Output(2, fmt.Sprintf(f, args...))
}

func Error(desc string, err error) {
	er.Output(2, fmt.Sprintf("%s: %+v", desc, err))
}

