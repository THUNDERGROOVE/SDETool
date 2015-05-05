package sde

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/util/log"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// getsde checks that our version is downloaded and opens it.
func getsde(version string) SDE {
	defer Debug(time.Now())

	download(version)
	db, err := sql.Open("sqlite3", fmt.Sprintf("sde.%v.db", version))
	if err != nil {
		fmt.Printf("Error opening DB,%v\n", err.Error())
	}
	return SDE{
		db,
		version}
}

var (
	PrintDebug bool
)

func Debug(t time.Time) {
	defer func() {
		if r := recover(); r != nil {
			//_, path, line, _ := runtime.Caller(1)
			//_, file := filepath.Split(path)
			//duration := time.Since(t)
			//log.LogError("Debug() paniced.  Couldn't get function name", file, line, path, duration)
		}
	}()
	_, path, line, _ := runtime.Caller(1)
	_, file := filepath.Split(path)
	duration := time.Since(t)
	function := strings.Split(strings.Split(strings.Split(string(debug.Stack()), "\n")[3], ":")[0], "\t")[1]
	log.Trace(fmt.Sprintf("%v:%v:%v took %v", file, line, function, duration.String()))
}

// download attempts to download a version of the SDE provided
func download(version string) {
	defer Debug(time.Now())

	if _, err := os.Stat(fmt.Sprintf("sde.%v.db", version)); os.IsNotExist(err) {
		if Versions[version].Zipped {
			downloadZip(version)
		} else {
			downloadDb(version)
		}

	}
}

func downloadZip(version string) {
	fmt.Printf("SDE file not found.  Downloading %v", version)
	res, err := http.Get(Versions[version].URL)
	if err != nil {
		fmt.Printf("\nError downloading SDE file: %v\n", err.Error())
	}
	fmt.Printf("\rDone. Saving ...                               ")
	out, err1 := ioutil.ReadAll(res.Body)
	if err1 != nil {
		fmt.Printf("\nError reading to file\n")
		os.Exit(1)
	}
	err2 := ioutil.WriteFile(fmt.Sprintf("sde.%v.db.zip", version), out, 0777)
	if err2 != nil {
		fmt.Printf("\nError writing to disk %v\n", err2.Error())
		os.Exit(1)
	}
	fmt.Printf("\rDone. Decompressing Zip...                                ")
	z, err3 := zip.OpenReader(fmt.Sprintf("sde.%v.db.zip", version))
	if err3 != nil {
		fmt.Printf("\nError reading Zip\n")
		os.Exit(1)
	}
	fmt.Printf("\rDone. Reading DB file from zip                               ")
	r, err4 := z.File[0].Open()
	if err4 != nil {
		fmt.Printf("\nError reading file from zip\n")
		os.Exit(1)
	}
	data, err5 := ioutil.ReadAll(r)
	if err5 != nil {
		fmt.Printf("\nError reading bytes from zip\n")
		os.Exit(1)
	}
	fmt.Printf("\rDone. Copying file to disk                               \n")
	err6 := ioutil.WriteFile(fmt.Sprintf("sde.%v.db", version), data, 0777)
	if err6 != nil {
		fmt.Printf("\nError writing data to disk\n")
		os.Exit(1)
	}
}
func downloadDb(version string) {
	fmt.Printf("SDE file not found.  Downloading %v", version)
	res, err := http.Get(Versions[version].URL)
	if err != nil {
		fmt.Printf("\nError downloading SDE file: %v\n", err.Error())
	}
	fmt.Printf("\rDone. Saving ...                               ")
	out, err1 := ioutil.ReadAll(res.Body)
	if err1 != nil {
		fmt.Printf("\nError reading to file\n")
		os.Exit(1)
	}
	err2 := ioutil.WriteFile(fmt.Sprintf("sde.%v.db", version), out, 0777)
	if err2 != nil {
		fmt.Printf("\nError writing to disk %v\n", err2.Error())
		os.Exit(1)
	}
}
