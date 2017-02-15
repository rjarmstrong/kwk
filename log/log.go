package log

import (
	"os"
	"io"
	"log"
	"fmt"
	"io/ioutil"
)

var Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
var dbg = log.New(ioutil.Discard, "DEBUG: ", log.LUTC|log.Lshortfile)

func SetOutput(out io.Writer) {
	dbg.SetOutput(out)
	Error.SetOutput(out)
}

func Debug(in ...interface{}){
	dbg.Output(2, fmt.Sprintf("%v", in))
}

