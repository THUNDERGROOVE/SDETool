package log

import (
	"log"
	"os"
)

var Log *log.Logger

func init() {
	Log = log.New(os.Stdout, "[SDETool]", 0)
}

func Println(v ...interface{}) {
	Log.Println(v...)
}

func Printf(f string, v ...interface{}) {
	Log.Printf(f, v...)
}

func Trace(v ...interface{}) {
	Log.SetPrefix("[SDETool][Trace]")
	Log.Println(v)
	Log.SetPrefix("[SDETool]")
}
