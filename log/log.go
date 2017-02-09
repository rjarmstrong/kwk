package log

import (
	"os"
	"io"
	"log"
	"fmt"
)

var Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
var dbg = log.New(os.Stdout, "DEBUG: ", log.LUTC|log.Lshortfile)

func SetOutput(out io.Writer) {
	dbg.SetOutput(out)
	Error.SetOutput(out)
}

func Debug(in ...interface{}){
	dbg.Output(2, fmt.Sprintf("%v", in))
}

