// Let's add this docstring thingy
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/THUNDERGROOVE/SDETool/scripting/langs"
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/THUNDERGROOVE/SDETool/sde/version"
	"github.com/THUNDERGROOVE/SDETool/util"
	_ "github.com/THUNDERGROOVE/SDETool/util/log"
	// Langs
	_ "github.com/THUNDERGROOVE/SDETool/scripting/lua"
	"github.com/d4l3k/messagediff"
)

var (
	branch     string
	tagVersion string
	commit     string
)

// Matches a TypeID for DUST
var TypeIDRegex = regexp.MustCompile(`3\d\d\d\d\d`)

/*func dumpSDE() {
	if err := dumperFlagset.Parse(os.Args[2:]); err != nil {
		fmt.Printf("[ERROR] Couldn't parse args[%v]\n", err.Error())
	}

	// @TODO: Move this to it's own function
	// Is there a better way to do this?
	cmd := exec.Command("sdedumper",
		"-i", fmt.Sprintf("%s", *dumperInFile),
		"-o", fmt.Sprintf("%s", *dumperOutFile),
		"-ver", fmt.Sprintf("%s", *dumperVersionString),
		"-official", fmt.Sprintf("%t", *dumperOfficial),
		"-v", fmt.Sprintf("%t", *dumperVerbose))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		break
	}
	if strings.Contains(err.Error(), exec.ErrNotFound.Error()) {
		fmt.Printf("You do not have sdedumper installed.  Please install\n")
	} else {
		log.Printf("Error running sdedumper [%s]", err.Error())
	}
}*/

func loadSDE(filename string) *sde.SDE {
	if filename == "" {
		SDE, err := version.LoadLatest()

		if err != nil {
			fmt.Printf("[ERROR] %v\n", err.Error())
			return nil
		}

		if SDE == nil {
			fmt.Printf("Failed to automatically load an SDE file.  Please load it manually\n")
			return nil
		}
		return SDE
	} else {
		SDE, err := sde.Load(filename)
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err.Error())
			return nil
		}
		return SDE
	}
	return nil
}

func main() {
	// If the first argument is a script then run it instead of parsing things
	if len(os.Args) > 1 {
		n := os.Args[1]
		if strings.Contains(n, ".") {
			if util.Exists(n) {
				ext := filepath.Ext(n)[1:]
				if err := langs.RunScript(ext, n); err != nil {
					fmt.Println("Error running script", err.Error())
				}
				return
			}
		}
	}

	if len(os.Args) <= 1 {
		printNoArgsText()
		return
	}

	cmd := cli.NewApp()
	cmd.Name = "SDETool"
	cmd.Author = "Nick Powell; @THUNDERGROOVE"
	cmd.Version = fmt.Sprintf("%v-%v-%v", branch, tagVersion, commit)
	cmd.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "sdefile, s",
			Value: "",
			Usage: "",
		},
	}
	cmd.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search the sde!",
			Action:  _cliSearchSDE,
		},
		{
			Name:    "lookup",
			Aliases: []string{"l"},
			Usage:   "Lookup a specific type",
			Action:  _cliLookupSDE,
		},
		{
			Name:    "diff",
			Aliases: []string{"d"},
			Usage:   "diff two types",
			Action:  _cliDiff,
		},
	}
	cmd.Run(os.Args)
}

func _cliDiff(c *cli.Context) {
	file := c.GlobalString("sdefile")
	SDE := loadSDE(file)

	t1 := c.Args().Get(0)
	t2 := c.Args().Get(1)
	if t1 == t2 {
		fmt.Printf("Why are you trying to compare two of the same types?\n")
	}

	if !TypeIDRegex.Match([]byte(t1)) || !TypeIDRegex.Match([]byte(t2)) {
		fmt.Printf("One of the given typeIDs were not valid\n")
	}

	tid1, _ := strconv.Atoi(t1)
	tid2, _ := strconv.Atoi(t2)

	tt1, err := SDE.GetType(tid1)
	if err != nil {
		fmt.Printf("[ERROR] failed to get first type: %v\n", err.Error())
	}

	tt2, err := SDE.GetType(tid2)
	if err != nil {
		fmt.Printf("[ERROR] failed to get second type: %v\n", err.Error())
	}

	diff, equal := messagediff.PrettyDiff(tt1, tt2)
	fmt.Printf("Equal?: %v\n", equal)
	fmt.Printf("\n\n%v\n\n", diff)
}

func _cliLookupSDE(c *cli.Context) {
	file := c.GlobalString("sdefile")
	SDE := loadSDE(file)

	s := c.Args().First()
	if !TypeIDRegex.Match([]byte(s)) {
		fmt.Printf("Not recognized as a typeID doing a search instead and using the first result\n")

		res, err := SDE.Search(s)
		if err != nil {
			fmt.Printf("[ERROR] failed getting type: %v\n", err.Error())
			return
		}
		t := res[0]
		fmt.Printf("[ %v | %v | %v ]\n", t.TypeID, t.TypeName, t.GetName())
		return
	}
	id, _ := strconv.Atoi(s)
	t, err := SDE.GetType(id)
	if err != nil {
		fmt.Printf("[ERROR] failed getting type: %v\n", err.Error())
		return
	}
	fmt.Printf("[ %v | %v | %v ]\n", t.TypeID, t.TypeName, t.GetName())
}

func _cliSearchSDE(c *cli.Context) {
	file := c.GlobalString("sdefile")
	SDE := loadSDE(file)

	s := c.Args().First()
	values, err := SDE.Search(s)
	if err != nil {
		fmt.Printf("[ERROR] Failed to search: %v\n", err.Error())
	}

	for _, v := range values {
		fmt.Printf("[ %v | %v | %v ]\n", v.TypeID, v.TypeName, v.GetName())
	}
}

func printNoArgsText() {
	fmt.Println("SDETool written by Nick Powell; @THUNDERGROOVE")
	fmt.Printf("Version %v@%v#%v\n", tagVersion, branch, commit)
	fmt.Println("Do 'help' for help")
}
