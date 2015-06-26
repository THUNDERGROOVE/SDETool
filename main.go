package main

/*
	@TODO:  Integrate cmd/dumper into the main tool if it is compiled and in the
	path; I don't wanna ;_;
*/

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/THUNDERGROOVE/SDETool/scripting/langs"
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/THUNDERGROOVE/SDETool/sde/version"
	"github.com/THUNDERGROOVE/SDETool/util"
	// Langs
	_ "github.com/THUNDERGROOVE/SDETool/scripting/lua"
)

var (
	// Branch is current branch of the git repo
	Branch string
	// Version is the latest git tag
	Version string
	// Commit is the "short" hash of the latest git commit
	Commit string
)

func main() {
	// If the first argument is a script then run it instead of parsing things
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

	if len(os.Args) <= 1 {
		printNoArgsText()
		return
	}

	var SDE *sde.SDE
	var Type *sde.SDEType
	var t []*sde.SDEType
	var MultiTypes bool
	var err error
	SDE, err = version.LoadLatest()

	if err != nil {
		fmt.Printf("[ERROR] %v\n", err.Error())
	}

	if SDE == nil {
		fmt.Printf("Failed to automatically load an SDE file.  Please load it manually\n")
		return
	}

	switch os.Args[1] {
	case "lookup":
		if err := lookupFlagset.Parse(os.Args[2:]); err != nil {
			fmt.Printf("[ERROR] Couldn't parse args [%v]\n", err.Error())
		}
		var err error

		switch {
		case *lookupSDE != "":
			SDE, err = sde.Load(*lookupSDE)
			fallthrough
		case *lookupTID != 0:
			Type, err = SDE.GetType(*lookupTID)
		case *lookupTN != "":
			t, err = SDE.Search(*lookupTN)
			if len(t) != 0 {
				Type = t[0]
			}
		case *lookupTD != "":
			t, err = SDE.Search(*lookupTD)
			if len(t) != 0 {
				Type = t[0]
			}
		}
		if Type != nil && *lookupAttr {
			for k, v := range Type.Attributes {
				fmt.Printf(" %v | %v\n", k, v)
			}
		}
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err.Error())
		}
	case "search":
		if err := searchFlagset.Parse(os.Args[2:]); err != nil {
			fmt.Printf("[ERROR] Couldn't parse args[%v]\n", err.Error())
		}
		var err error
		switch {
		case *searchSDE != "":
			SDE, err = sde.Load(*searchSDE)
			fallthrough
		case *searchName != "":
			t, err = SDE.Search(*searchName)
			if len(t) != 0 {
				Type = t[0]
			}
			MultiTypes = true
			fallthrough
		case *searchAttr == true:
			if len(t) == 1 {
				for k, v := range Type.Attributes {
					fmt.Printf(" %v| %v\n", k, v)
				}
				MultiTypes = false
			}
			if err != nil {
				fmt.Printf("[ERROR], %v\n", err.Error())
			}
		}
	case "help":
		fmt.Printf(HelpText)
	default:
		printNoArgsText()
	}
	if Type != nil && !MultiTypes {
		fmt.Printf("%v | %v | %v\n", Type.TypeID, Type.TypeName, Type.GetName())
	} else if MultiTypes {
		for _, v := range t {
			fmt.Printf("%v | %v | %v\n", v.TypeID, v.TypeName, v.GetName())
		}
	} else {
		fmt.Printf("No type resolved\n")
	}
}

func printNoArgsText() {
	fmt.Println("SDETool written by Nick Powell; @THUNDERGROOVE")
	fmt.Printf("Version %v branch %v commit %v\n", Version, Branch, Commit)
	fmt.Println("Do 'help' for help")
}
