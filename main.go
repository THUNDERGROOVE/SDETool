package main

/*
	@TODO:  Integrate cmd/dumper into the main tool if it is compiled
*/

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/scripting/langs"
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/THUNDERGROOVE/SDETool/util"
	"github.com/nsf/termbox-go"
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
	var SDE *sde.SDE
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

	// @TODO:
	//	Finish implementing all the old commands and flags

	c := strings.Join(os.Args[1:], " ")
	s := args.NewScanner(strings.NewReader(c))

	tl := s.ScanAll()
	for _, v := range tl {
		if v.Token == args.ILLEGAL {
			fmt.Printf("Token %v; Literal: '%v'\n", v.Token, v.Literal)
		}
	}

	args.FlagReg.RegisterCmd("search", func(tok args.TokLit, index int) {
		if SDE == nil {
			fmt.Printf("NO SDE file was loaded\n")
			return
		}
		if arg := tl.Next(index); arg != nil {
			switch arg.Token {
			case args.STRING:
				fmt.Printf("Searching using string '%v'\n", arg.Literal)
				out, err := SDE.Search(arg.Literal)
				if err != nil {
					fmt.Printf("Error searching the SDE: %v\n", err.Error())
				}
				for _, v := range out {
					fmt.Printf("  [%v] [%v]\n", v.TypeID, v.GetName())
				}
			default:
				fmt.Printf("search doesn't take a token type of %v\n", arg.Token)
			}
		} else {
			fmt.Printf("No argument was supplied for search\n")
		}
	})

	args.FlagReg.RegisterCmd("lookup", func(tok args.TokLit, index int) {
		if SDE == nil {
			fmt.Printf("No SDE file was loaded\n")
			return
		}
		if arg := tl.Next(index); arg != nil {
			switch arg.Token {
			case args.INT: // Lookup by TypeID
				fmt.Printf("Looking up by typeID\n")
				var err error
				i, _ := strconv.Atoi(arg.Literal)
				Type, err = SDE.GetType(i)
				if err != nil {
					fmt.Printf("Error getting type: [%v]\n", err.Error())
				}
			case args.STRING: // Lookup by mDisplayName followed by TypeName if that fails
				fmt.Printf("Looking up by display name or typename: %v\n", arg.Literal)
			default:
				fmt.Printf("Lookup doesn't know how to handle token %v\n", arg.Token)
			}
		} else {
			fmt.Printf("No argument supplied for lookup\n")
		}
		if Type == nil {
			fmt.Printf("Type returned was nil?\n")
			return
		}
		fmt.Printf("%v | %v | %v\n", Type.TypeID, Type.TypeName, Type.GetName())
		if subcmd := tl.Next(index + 2); subcmd != nil {
			if subcmd.Token == args.STRING {
				switch subcmd.Literal {
				case "attrs":
					cleanPrintAttrs(Type)
				default:
					fmt.Printf("Unknown sub command: '%v'\n", subcmd.Literal)
				}
			} else {
				fmt.Println("Subcommand had type", subcmd.Token)
			}
		} else {
			fmt.Printf("No sub command provided\n")
		}
	})

	args.FlagReg.Register("--sde", "-s", func(tok args.TokLit, index int) {
		if arg := tl.Next(index); arg != nil {
			var err error
			fmt.Printf("Loading SDE: %v\n", arg.Literal)
			SDE, err = sde.Load(arg.Literal)
			if err != nil {
				fmt.Printf("Error encoutnered while loading the SDE file:  [%v]\n", err.Error())
			}
		} else {
			fmt.Printf("No argument supplied for loading an SDE file index: %v\n", index)
		}
	})

	args.FlagReg.Parse(tl)

	fmt.Println("SDETool written by Nick Powell; @THUNDERGROOVE")
	fmt.Printf("Version %v branch %v commit %v\n", Version, Branch, Commit)
	fmt.Println("Do --help for help")
}

func cleanPrintAttrs(t *sde.SDEType) {
	width := 0
	for k, _ := range t.Attributes {
		if len(k) > width {
			width = len(k)
		}
	}
	w, _ := GetSize()
	fmt.Println("Using width of", width)
	// @TODO: Seriously need to refactor this piece of garbage LOL; But hey it looks pretty
	// someAttribute---------> Really really really really really really really |
	//                     |-> really really really really really really really |
	//                     \-> really really really long text                   |
	//
	for k, v := range t.Attributes {
		switch str := v.(type) {
		case string:
			v = interface{}(strings.Replace(str, "\n", " ", -1))
		}
		if len(fmt.Sprintf("%v", v)) > w-width-2 {
			s := splitEvery(fmt.Sprintf("%v", v), w-width-2-5)
			var once sync.Once
			for ki, v := range s {
				var skip bool
				once.Do(func() {
					fmt.Printf("  %v> %v\n",
						k+strings.Repeat("-", width-len(k)),
						v,
					)

					skip = true
				})
				if skip {
					continue
				}
				if ki == len(s)-1 {
					fmt.Printf("  %v\\--> %v\n",
						strings.Repeat(" ", width-3),
						v,
					)
				} else {
					fmt.Printf("  %v|--> %v\n",
						strings.Repeat(" ", width-3),
						v,
					)
				}
			}
		} else {
			fmt.Printf("  %v> %v\n",
				k+strings.Repeat("-", width-len(k)),
				v,
			)
		}
	}
}

func GetSize() (width int, height int) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	termbox.Close()
	return w, h
}

func splitEvery(s string, i int) []string {
	o := make([]string, 0)
	c := int(float64(len(s)) / math.Floor(float64(i)))
	for off := 0; off < c; off++ {
		var b int
		if off == c {
			b = len(s) - (c * i)
		} else {
			b = i
		}
		o = append(o, s[off*i:off*i+b])
	}
	return o
}
