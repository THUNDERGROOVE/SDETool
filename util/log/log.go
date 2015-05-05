package log

import (
	"log"
	"os"
)

var Log *log.Logger

func init() {
	FileLogger()
}

func StdoutLogger() {
	Log = log.New(os.Stdout, "[SDETool]", 0)
}
func FileLogger() {
	os.Remove("log.txt")
	f, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0777)
	Log = log.New(f, "[SDETool]", 0)
}

func Println(v ...interface{}) {
	if Log != nil {
		Log.Println(v...)
	}
}

func Printf(f string, v ...interface{}) {
	if Log != nil {
		Log.Printf(f, v...)
	}
}

func Trace(v ...interface{}) {
	if Log != nil {
		Log.SetPrefix("[SDETool][Trace]")
		Log.Println(v)
		Log.SetPrefix("[SDETool]")
	}
}
