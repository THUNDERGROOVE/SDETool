package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"database/sql"

	"github.com/THUNDERGROOVE/SDETool/sde"
)

var (
	InFile        = flag.String("i", "", "The input SQLite3 database file the SDE resides in")
	OutFile       = flag.String("o", "", "The output file that is compatible with SDETool")
	VersionString = flag.String("ver", "", "The string encoded into the dump")
	Official      = flag.Bool("official", false, "Set this flag if the SQLite db is from CCP")
	Verbose       = flag.Bool("v", false, "Print more stuffs")
)

var (
	DB *sql.DB
)

func init() {
	flag.Parse()
}

func main() {
	if *InFile == "" {
		fmt.Println("You must supply an input file with the -i flag")
		return
	}

	if *VersionString == "" {
		fmt.Println("You must supply a version string with -ver")
		return
	}

	fmt.Printf("dumper using %v threads\n", runtime.GOMAXPROCS(4))

	offstr := "Unofficial"
	if *Official {
		offstr = "Official"
	}
	fmt.Println("Starting dump proccess for:")
	fmt.Printf(" => %v:%v %v\n", *InFile, *VersionString, offstr)

	db, err := sql.Open("sqlite3", *InFile)
	DB = db
	if err != nil {
		ReportError("couldn't open the sqlite file %v [%v]", *InFile, err.Error)
		return
	}

	out := TypeAttributeGetter(TypeBaseGetter())

	SDE := &sde.SDE{
		Version:  *VersionString,
		Official: *Official,
		Types:    make(map[int]*sde.SDEType),
	}

	var i int
	count := TableSize("CatmaTypes")
	countWidth := len(strconv.Itoa(count))
	for tt := range out {
		t := tt
		i += 1
		iwidth := len(strconv.Itoa(i))
		fmt.Printf("\r                                                                                 ")
		fmt.Printf("\r[%v/%v] @%p [%v] %v",
			strings.Repeat(" ", countWidth-iwidth)+strconv.Itoa(i),
			count,
			&t,
			t.TypeID,
			t.GetName(),
		)
		SDE.Types[t.TypeID] = &t
	}

	if err := sde.Save(*OutFile, SDE); err != nil {
		ReportError("Couldn't save SDE dump %v", err.Error())
	}
}

func TypeBaseGetter() chan sde.SDEType {
	out := make(chan sde.SDEType, 0)
	go func() {
		fmt.Println("Begning CatmaType lookup")
		rows, err := DB.Query("SELECT * FROM CatmaTypes;")
		if err != nil {
			ReportError("SQLite3 error: %v", err.Error())
		}
		for rows.Next() {
			var typeID int
			var typeName string
			rows.Scan(&typeID, &typeName)
			t := sde.SDEType{
				TypeID:     typeID,
				TypeName:   typeName,
				Attributes: make(map[string]interface{}),
			}
			VerbosePrint("Attribute getting sending type")
			out <- t
		}
		close(out)
	}()
	return out
}

func TypeAttributeGetter(in chan sde.SDEType) chan sde.SDEType {
	out := make(chan sde.SDEType)
	go func() {
		fmt.Println("Attribute lookup goroutine started")
		for t := range in {
			VerbosePrint("Attribute Lookup received type [%v][%v]", t.TypeID, t.TypeName)
			rows, err := DB.Query(
				fmt.Sprintf(
					"SELECT catmaAttributeName, catmaValueInt, catmaValueReal,catmaValueText FROM CatmaAttributes where TypeID == '%v';",
					t.TypeID))
			if err != nil {
				ReportError("SQLlite3 error: %v", err.Error())
				continue
			}
			for rows.Next() {
				var (
					cattr  string
					cvint  string
					cvreal string
					cvtext string
				)
				if err := rows.Scan(&cattr, &cvint, &cvreal, &cvtext); err != nil {
					ReportError("Couldn't scan db row [%v]", err.Error())
					continue
				}

				if cvint != "None" {
					v, err := strconv.Atoi(cvint)
					if err != nil {
						ReportError("Couldn't parse int from catmaValue field")
						continue
					}
					t.Attributes[cattr] = v
				}
				if cvreal != "None" {
					v, err := strconv.ParseFloat(cvreal, 64)
					if err != nil {
						ReportError("Couldn't parse int from catmaValue field")
						continue
					}
					t.Attributes[cattr] = v
				}
				if cvtext != "None" {
					t.Attributes[cattr] = cvtext
				}
			}
			out <- t
		}
		close(out)
	}()
	return out
}

func Consumer(chans ...chan sde.SDEType) chan sde.SDEType {
	out := make(chan sde.SDEType)

	for _, v := range chans {
		go func() {
			for t := range v {
				out <- t // Yeah yeah shutup go vet
			}
		}()
	}
	return out
}

func TableSize(table string) int {
	if DB == nil {
		ReportError("TableSize called with nil DB")
		return 0
	}
	rows, err := DB.Query(fmt.Sprintf("SELECT Count(*) FROM %v;", table))
	if err != nil {
		ReportError("SQLite error: %v", err.Error())
	}
	rows.Next()
	var count int
	rows.Scan(&count)
	return count
}

func ReportError(f string, a ...interface{}) {
	_, path, line, _ := runtime.Caller(1)
	_, file := filepath.Split(path)
	f = fmt.Sprintf("[ERROR]%v:%v", file, line) + f + "\n"
	fmt.Fprintf(os.Stderr, f, a...)
}

func VerbosePrint(f string, a ...interface{}) {
	if *Verbose {
		_, path, line, _ := runtime.Caller(1)
		_, file := filepath.Split(path)
		f = fmt.Sprintf("[INFO]%v:%v", file, line) + f + "\n"
		fmt.Printf(f, a...)
	}
}
