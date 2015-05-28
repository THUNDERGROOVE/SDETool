package main

/*
	@TODO:  Integrate cmd/dumper into the main tool if it is compiled
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/commands"
	"github.com/THUNDERGROOVE/SDETool/scripting/langs"
	"github.com/THUNDERGROOVE/SDETool/util"
	// Langs
	_ "github.com/THUNDERGROOVE/SDETool/scripting/lua"
)

var (
	Version string
	Branch  string
	Commit  string
)

func main() {
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

	c := strings.Join(os.Args[1:], " ")
	s := args.NewScanner(strings.NewReader(c))

	tl := s.ScanAll()
	for _, v := range tl {
		if v.Token == args.ILLEGAL {
			fmt.Printf("Token %v; Literal: '%v'\n", v.Token, v.Literal)
		}
	}

	commands.RegisterCommands(tl)

	if !args.FlagReg.Parse(tl) {
		fmt.Println("SDETool written by Nick Powell; @THUNDERGROOVE")
		fmt.Printf("Version %v branch %v commit %v\n", Version, Branch, Commit)
		fmt.Println("Do --help for help")
	}
}
