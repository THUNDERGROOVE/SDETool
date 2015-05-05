package args

import (
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	SDETool = kingpin.New("SDETool", "A script enabled SDE lookup tool")

	DoScript     = SDETool.Command("script", "Do a script.")
	DoScriptFile = DoScript.Arg("file", "A file to run").Required().String()
	DoScriptLang = DoScript.Arg("lang", "The language to use").Default("lua").String()

	Interpreter     = SDETool.Command("interpreter", "Opens an interpreter for the given language; Defaults to lua")
	InterpreterLang = Interpreter.Arg("lang", "Language to start an interpreter with").Default("lua").String()

	ListLangs = SDETool.Command("langs", "Lists the compiled in langs")
)
