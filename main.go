package main

import (
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/scripting"
	"github.com/THUNDERGROOVE/SDETool/scripting/langs"
	"gopkg.in/alecthomas/kingpin.v1"
	"os"
	// Langs
	_ "github.com/THUNDERGROOVE/SDETool/scripting/lua"
)

func main() {
	switch kingpin.MustParse(args.SDETool.Parse(os.Args[1:])) {
	case args.ListLangs.FullCommand():
		fmt.Println("Compiled in languages: ")
		for k, _ := range langs.ScriptingLangs {
			fmt.Printf("  %v\n", k)
		}
	case args.DoScript.FullCommand():
		if err := langs.RunScript(*args.DoScriptLang, *args.DoScriptFile); err != nil {
			if err == langs.NoSuchLang {
				fmt.Printf("The language provided '%v' is not compiled in. \n", *args.DoScriptLang)
				return
			} else {
				fmt.Printf("The following error occured while running the script '%v': %v", *args.DoScriptFile, err.Error())
			}
		}
	case args.Interpreter.FullCommand():
		if err := langs.Interpreter(*args.InterpreterLang); err != nil {
			if err == scripting.InterpreterNotImplemented {
				fmt.Printf("The language %v does not implement an interpreter.\n", *args.InterpreterLang)
			} else {
				fmt.Printf("The following error occured while in the interpreter %v\n", err.Error())
			}
		}
	default:
		fmt.Println("SDETool written by Nick Powell; @THUNDERGROOVE")
		fmt.Println("Do --help for help")
	}
}
