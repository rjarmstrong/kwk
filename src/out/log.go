package out

import (
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path"
	"runtime/debug"
)

// DebugEnabled is a setting allowing turning Debug to stdout on and off.
var DebugEnabled bool

var fileOut = &lumberjack.Logger{
	Filename:   path.Join(KwkPath(), "main.log"),
	MaxSize:    3, // megabytes
	MaxBackups: 2,
	MaxAge:     5, //days})
}
var fileLogger = log.New(fileOut, "KWK: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

// ErrorLogger formats errors for printing to stdout when DebugEnabled, or to main.log when not.
var ErrorLogger = log.New(os.Stderr, "KWK:ERR: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

// DebugLogger formats debug output to stdout when DebugEnabled.
var DebugLogger = log.New(os.Stdout, "KWK: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

// Debug used to log all non-errors.
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

// LogErr logs errors as is.
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
