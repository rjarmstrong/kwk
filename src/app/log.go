package app

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
	Filename:   path.Join(kwkPath(), "kwk.log"),
	MaxSize:    3, // megabytes
	MaxBackups: 2,
	MaxAge:     5, //days})
}
var fileLogger = log.New(fileOut, "KWK: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
var ErrorLogger = log.New(os.Stderr, "KWKCLI ERR: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
var DebugLogger = log.New(os.Stdout, "KWKCLI: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

func Debug(in ...interface{}) {
	if !DebugEnabled {
		return
	}
	DebugLogger.Output(2, fmt.Sprintf("%v", in))
}

func LogErrM(message string, err error) error {
	Debug(message)
	LogErr(err)
	return err
}

func LogErr(err error) error {
	fileLogger.Output(2, string(debug.Stack()))
	if !DebugEnabled {
		return nil
	}
	ErrorLogger.Println(err)
	ErrorLogger.Output(2, string(debug.Stack()))
	return err
}
