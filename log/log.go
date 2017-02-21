package log

import (
	lg "log"
	"os"
	"io"
	"log"
	"fmt"
	"flag"
	"path"
	"bitbucket.com/sharingmachine/kwkcli/cache"
)

var er = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
var dbg = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)

// This flag is specified here since we don't wan't to have to depend on urfave global flags.
// That flag will also work throughout the app.
var dbgFlag = flag.Bool("debug", false, "Show debug trace.")

func init() {
	flag.Parse()
	f, err := os.OpenFile(path.Join(cache.Path(), "kwk.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	if !*dbgFlag {
		SetOutput(f)
	} else {
		Debug("--debug")
	}
	er.SetOutput(f)
}

func SetOutput(out io.Writer) {
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

