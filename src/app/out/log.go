package out

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
	"runtime/debug"
)

var DebugEnabled bool
var fileOut = &lumberjack.Logger{
	Filename:   path.Join(KwkPath(), "kwk.log"),
	MaxSize:    3, // megabytes
	MaxBackups: 2,
	MaxAge:     5, //days})
}
var fileLogger = log.New(fileOut, "KWK: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
var ErrorLogger = log.New(os.Stderr, "KWKCLI ERR: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
var DebugLogger = log.New(os.Stdout, "KWKCLI: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

func Debug(format string, in ...interface{}) {
	if !DebugEnabled {
		return
	}
	var mess string
	if len(in) > 0 {
	  mess = fmt.Sprintf(format, in...)
	} else {
	  mess = format
	}
	DebugLogger.Output(2, mess)
}

// LogErrM allows to log an error and specify a custom message.
func LogErrM(message string, err error) error {
	Debug(message)
	LogErr(err)
	return err
}

func LogErr(err error) error {
	fileLogger.Println(err)
	fileLogger.Output(2, string(debug.Stack()))
	if !DebugEnabled {
		return nil
	}
	ErrorLogger.Println(err)
	ErrorLogger.Output(2, string(debug.Stack()))
	return err
}
