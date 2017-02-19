package log

import (
	"os"
	"io"
	"log"
	"fmt"
	"io/ioutil"
	"flag"
)

var er = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
var dbg = log.New(os.Stdout, "DEBUG: ", log.LUTC|log.Lshortfile)

// This flag is specified here since we don't wan't to have to depend on urfave global flags.
// That flag will also work throughout the app.
var dbgFlag = flag.Bool("debug", false, "Show debug trace.")

func init() {
	flag.Parse()
	if !*dbgFlag {
		SetOutput(ioutil.Discard)
	} else {
		Debug("--debug")
	}
}

func SetOutput(out io.Writer) {
	dbg.SetOutput(out)
	er.SetOutput(out)
}

func Debug(f string, args ...interface{}){
	dbg.Output(2, fmt.Sprintf(f, args...))
}

func Error(desc string, err error) {
	er.Output(2, fmt.Sprintf("%s: %+v", desc, err))
}

