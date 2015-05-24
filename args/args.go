package args

import (
	"gopkg.in/alecthomas/kingpin.v1"
)

var (
	SDETool = kingpin.New("SDETool", "A script enabled SDE lookup tool")

	SDEFile = SDETool.Flag("sde", "A provided SDETool file to open").Short('s').String()

	DoScript     = SDETool.Command("script", "Do a script.")
	DoScriptFile = DoScript.Arg("file", "A file to run").Required().String()
	DoScriptLang = DoScript.Arg("lang", "The language to use").Default("lua").String()

	Interpreter     = SDETool.Command("interpreter", "Opens an interpreter for the given language; Defaults to lua")
	InterpreterLang = Interpreter.Arg("lang", "Language to start an interpreter with").Default("lua").String()

	Lookup     = SDETool.Command("lookup", "Looks up a type")
	LookupTID  = Lookup.Arg("typeID", "A TypeID to look for").Int()
	LookupAttr = Lookup.Arg("attr", "Print attributes of type").Bool()

	Search       = SDETool.Command("search", "Used to search for a type")
	SearchString = Search.Arg("str", "The string used to search with").String()
	SearchQuick  = Search.Arg("quick", "Don't get all attributes for every type").Bool()

	ListLangs = SDETool.Command("langs", "Lists the compiled in langs")
)
