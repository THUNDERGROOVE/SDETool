package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/scripting"
	"github.com/THUNDERGROOVE/SDETool/scripting/langs"
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/THUNDERGROOVE/SDETool/util"
	"gopkg.in/alecthomas/kingpin.v1"
	// Langs
	_ "github.com/THUNDERGROOVE/SDETool/scripting/lua"
)

var (
	Version string
	Branch  string
	Commit  string
)

func main() {
	var Type *sde.SDEType
	// Attempt to figure out what the fuck to do before kingpin gets involved.
	if len(os.Args) > 1 {
		n := os.Args[1]
		if util.Exists(n) {
			ext := filepath.Ext(n)[1:]
			if err := langs.RunScript(ext, n); err != nil {
				fmt.Println("Error running script", err.Error())
			}
			return
		}
	}

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
	case args.Lookup.FullCommand():
		if *args.SDEFile == "" {
			fmt.Println("You must supply an SDEFile using the sde flag")
			return
		}
		SDE, err := sde.Load(*args.SDEFile)
		if err != nil {
			fmt.Println("Error while opening the SDE file")
			return
		}
		if *args.LookupTID != 0 {
			t, err := SDE.GetType(*args.LookupTID)
			if err == nil {
				Type = t
			} else {
				fmt.Println("Couldn't find the type:", *args.LookupTID)
				os.Exit(1)
			}
		}
		if Type != nil {
			if *args.LookupAttr {
				fmt.Printf("Attributes for:\n  %v | %v | %v\n", Type.GetName(), Type.TypeID, Type.TypeName)
				for k, v := range Type.Attributes {
					fmt.Printf("  %v |  %v\n", k, v)
				}
			}
		} else {
			fmt.Println("Failed to resolve a type.")
		}
	case args.Search.FullCommand():
		if *args.SDEFile == "" {
			fmt.Println("You must supply an SDEFile using the sde flag")
			return
		}
		SDE, err := sde.Load(*args.SDEFile)
		if err != nil {
			fmt.Println("Error while opening the SDE file")
			return
		}
		if *args.SearchString != "" {
			if *args.SearchQuick {
				fmt.Println("Quick search not implemented yet")
			}
			fmt.Println("Searching for:", *args.SearchString)
			types, err := SDE.Search(*args.SearchString)
			if err != nil {
				fmt.Println("Error searching for types", err.Error())
				os.Exit(1)
			}
			for _, v := range types {
				fmt.Printf("  %v | %v\n", v.TypeID, v.GetName())
			}
		}
	default:
		fmt.Println("SDETool written by Nick Powell; @THUNDERGROOVE")
		fmt.Printf("Version %v branch %v commit %v\n", Version, Branch, Commit)
		fmt.Println("Do --help for help")
	}
}
