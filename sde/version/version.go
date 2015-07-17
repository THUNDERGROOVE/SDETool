package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/THUNDERGROOVE/SDETool/sde"
)

const (
	Upstream    = "http://dl.maximumtwang.com/SDE/"
	VersionFile = "versions.json"
)

// LoadVersions checks the download server for all of the available versions
func LoadVersions() (map[string]string, error) {
	url := fmt.Sprintf("%v%v", Upstream, VersionFile)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	vers := make(map[string]string)
	if err := json.Unmarshal(data, &vers); err != nil {
		return nil, err
	}
	return vers, nil
}

// GenVersions generates a versions map given all the SDE files in the current
// working directory and marshals it to file.
func GenVersions() (map[string]string, error) {
	dirs, err := ioutil.ReadDir(".")
	if err != nil {
		return nil, err
	}
	ver := make(map[string]string)
	for _, f := range dirs {
		if !f.IsDir() {
			if filepath.Ext(f.Name()) == ".sde" {
				fmt.Printf("Found SDE file %v\n", f.Name())
				SDE, err := sde.Load(f.Name())
				if err != nil {
					fmt.Printf("Couldn't load SDE file: [%v]\n", err.Error())
					continue
				}
				ver[SDE.Version] = f.Name()
			}
		}
	}
	fmt.Println("Saving versions to file")
	if data, err := json.Marshal(ver); err != nil {
		fmt.Printf("Error marshaling data %v\n", err.Error())
	} else {
		if err := ioutil.WriteFile(VersionFile, data, 0777); err != nil {
			fmt.Printf("Error writing to file: %v\n", err.Error())
		}
	}
	return ver, nil
}

// GetVersion downloads a given version if possible
func GetVersion(v string, file string) error {
	path := GetVersionPath(file)
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	fmt.Printf("Version %v not downloaded yet\n", v)
	url := fmt.Sprintf("%v%v", Upstream, file)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("HTTP 404: Not found: %v", url)
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, data, 0777); err != nil {
		return err
	}
	return nil
}

// GetVersionPath returns where the version files are stored
func GetVersionPath(v string) string {
	if _, err := os.Stat(filepath.Join(getappdatafolder(), ".SDETool")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(getappdatafolder(), ".SDETool"), 0777)
	}
	f := filepath.Join(getappdatafolder(), ".SDETool", v)
	return f
}

func getappdatafolder() string {
	u, err := user.Current()
	if err != nil {
		fmt.Println("ERROR UNABLE TO GET INSTANCE OF USER")
	}
	return u.HomeDir
}

func loadLatestOffline() (*sde.SDE, error) {
	f, err := ioutil.ReadDir(filepath.Join(getappdatafolder(), ".SDETool"))
	if err != nil {
		return nil, err
	}
	var newest string
	n := -1
	for _, v := range f {
		if filepath.Ext(v.Name()) == ".sde" && !v.IsDir() {
			veri := parseVersionFile(v.Name())
			log.Printf("File %s was given a version of %v", v.Name(), veri)
			if veri > n {
				newest = v.Name()
				n = veri
			}
		} else {
			log.Printf("Skipping %s", v.Name())
		}
	}
	file := filepath.Join(getappdatafolder(), ".SDETool", newest)
	log.Printf("Attempting to load %v", file)
	return sde.Load(file)
}

// LoadLatest loads the latest SDE files available
func LoadLatest() (*sde.SDE, error) {
	ver, err := LoadVersions()
	if err != nil {

		return loadLatestOffline()
	}
	var newest string
	n := -1
	for k, _ := range ver {
		veri := parseVersion(k)
		if veri > n {
			newest = k
			n = veri
		}
	}
	if err := GetVersion(newest, ver[newest]); err != nil {
		return nil, err
	}
	path := GetVersionPath(ver[newest])
	return sde.Load(path)
}

func parseVersion(v string) int {
	s := strings.Split(v, " ")
	up := s[0]
	ver := s[1]
	var out int
	switch strings.ToLower(up) {
	case "warlords":
		out += 1000
	case "uprising":
	default:
		fmt.Printf("Unknown titled version: %v\n", up)
	}
	ver = strings.Replace(ver, ".", "", -1)
	i, _ := strconv.Atoi(ver)
	out += i
	return out
}

func parseVersionFile(v string) int {
	v = strings.Split(v, ".")[0]
	v = strings.Replace(v, "dust-", "", -1)
	v = strings.Replace(v, "-", " ", -1)

	return parseVersion(v)
}
