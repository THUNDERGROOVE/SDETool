package version

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/THUNDERGROOVE/SDETool/sde"
)

const (
	Upstream    = "http://dl.maximumtwang.com/SDE/"
	VersionFile = "versions.json"
)

func LoadVersions() (map[string]string, error) {
	url := fmt.Sprintf("%v/%v", Upstream, VersionFile)
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
	return ver, nil
}
