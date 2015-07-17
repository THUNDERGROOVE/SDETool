package log

import (
	"log"
	"os"
)

func init() {
	FileLogger()
}

func StdoutLogger() {
	log.SetOutput(os.Stdout)
}
func FileLogger() {
	os.Remove("log.txt")
	f, _ := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY, 0777)
	log.SetOutput(f)
	log.SetPrefix("[SDETool]")
	log.SetFlags(log.Lshortfile)
}

func Trace(v ...interface{}) {
	log.SetPrefix("[SDETool][Trace]")
	log.Println(v)
	log.SetPrefix("[SDETool]")
}
