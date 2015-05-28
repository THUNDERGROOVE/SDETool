package main

/*
	@TODO:  Integrate cmd/dumper into the main tool if it is compiled
*/

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/THUNDERGROOVE/SDETool/args"
	"github.com/THUNDERGROOVE/SDETool/commands"
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
